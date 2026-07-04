package httpapi

import (
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// clickout logs the outbound click (the core monetisation event) and redirects
// to the exchanger's referral URL, falling back to its website.
func (a *api) clickout(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	directionSlug := r.URL.Query().Get("direction")

	e, ok, err := a.deps.Svc.ExchangerBySlug(r.Context(), slug)
	if err != nil || !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var directionID string
	if directionSlug != "" {
		if d, ok, err := a.deps.Svc.DirectionBySlug(r.Context(), directionSlug); err == nil && ok {
			directionID = d.ID
		}
	}

	// Affiliate attribution: an explicit ?ref= wins, else the kurso_ref cookie the
	// visitor picked up when they landed via a partner link.
	refCode := strings.TrimSpace(r.URL.Query().Get("ref"))
	if refCode == "" {
		if c, err := r.Cookie("kurso_ref"); err == nil {
			refCode = strings.TrimSpace(c.Value)
		}
	}

	// Best-effort clickout log — never blocks the redirect on a logging failure.
	if _, err := a.deps.DB.Exec(r.Context(), `
		INSERT INTO clickouts (exchanger_id, direction_id, ref_code, ip, user_agent, referer)
		VALUES ($1::uuid, NULLIF($2,'')::uuid, NULLIF($3,''), NULLIF($4,'')::inet, NULLIF($5,''), NULLIF($6,''))`,
		e.ID, directionID, refCode, clientIP(r), r.UserAgent(), r.Referer()); err != nil {
		a.deps.Log.Warn("clickout log failed", "exchanger", slug, "err", err)
	}

	target := ""
	switch {
	case e.ReferralTmpl != nil && *e.ReferralTmpl != "":
		target = strings.ReplaceAll(*e.ReferralTmpl, "{direction}", directionSlug)
	case e.WebsiteURL != nil:
		target = *e.WebsiteURL
	}
	if target == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}

// clientIP extracts a valid IP from RemoteAddr, or "" when it can't be parsed
// (so the inet column receives NULL rather than a bad value).
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	if net.ParseIP(host) == nil {
		return ""
	}
	return host
}
