package httpapi

import (
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/adapter/http/openapi"
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
	"github.com/google/uuid"
)

// toUUID parses a canonical UUID string; an invalid value yields the nil UUID.
func toUUID(s string) uuid.UUID {
	u, _ := uuid.Parse(s)
	return u
}

// f32 narrows a *float64 to the *float32 the contract uses for ratings.
func f32(p *float64) *float32 {
	if p == nil {
		return nil
	}
	v := float32(*p)
	return &v
}

func currencyDTO(c domain.Currency) openapi.Currency {
	return openapi.Currency{
		Id:      toUUID(c.ID),
		Code:    c.Code,
		Name:    c.Name,
		Kind:    openapi.CurrencyKind(c.Kind),
		Network: c.Network,
		IconUrl: c.IconURL,
	}
}

func directionDTO(d domain.Direction) openapi.Direction {
	return openapi.Direction{
		Id:             toUUID(d.ID),
		Slug:           d.Slug,
		FromCurrencyId: toUUID(d.FromID),
		ToCurrencyId:   toUUID(d.ToID),
		FromCode:       d.FromCode,
		FromName:       d.FromName,
		ToCode:         d.ToCode,
		ToName:         d.ToName,
		IsPopular:      d.IsPopular,
	}
}

func exchangerDTO(e domain.Exchanger) openapi.Exchanger {
	assets := e.Assets
	if assets == nil {
		assets = []string{}
	}
	return openapi.Exchanger{
		Id:              toUUID(e.ID),
		Slug:            e.Slug,
		Name:            e.Name,
		Status:          openapi.ExchangerStatus(e.Status),
		Partner:         e.Partner,
		IsVerified:      e.IsVerified,
		ReviewsCount:    e.ReviewsCount,
		RatingAvg:       f32(e.RatingAvg),
		WebsiteUrl:      e.WebsiteURL,
		LogoUrl:         e.LogoURL,
		Description:     e.Description,
		ReserveTotal:    e.ReserveTotal,
		DirectionsCount: e.DirectionsCount,
		Assets:          assets,
		OnSince:         e.OnSinceYear,
	}
}

func rateRowDTO(r domain.RateRow) openapi.RateRow {
	return openapi.RateRow{
		ExchangerId:   toUUID(r.ExchangerID),
		ExchangerSlug: r.ExchangerSlug,
		ExchangerName: r.ExchangerName,
		Partner:       r.Partner,
		RatingAvg:     f32(r.Rating),
		ReviewsCount:  r.ReviewsCount,
		Rate:          r.Rate,
		Reserve:       r.Reserve,
		MinAmount:     r.MinAmount,
		MaxAmount:     r.MaxAmount,
		FetchedAt:     r.FetchedAt,
	}
}

func rateRowDTOs(rows []domain.RateRow) []openapi.RateRow {
	out := make([]openapi.RateRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, rateRowDTO(r))
	}
	return out
}
