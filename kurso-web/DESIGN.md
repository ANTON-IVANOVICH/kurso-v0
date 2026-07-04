# Kurso Web — Design System (Stage 1)

Dark, data-dense fintech aesthetic. Everything below is already wired into the
project (Tailwind tokens, global CSS, fonts, layout, base components). Build new
pages/components **on top of these tokens and components** — do not hardcode the
raw palette when a token exists.

## Fonts (loaded globally via `nuxt.config.ts`)

- **Onest** — default sans (`font-sans`). Weights 400–800.
- **Geist Mono** — tabular figures for prices/amounts/IDs. Apply the `.tnum`
  utility class to any numeric span (`<span class="tnum">81 200 ₽</span>`).
- **JetBrains Mono** — uppercase eyebrow labels (`font-label`, usually with
  `text-[11px] uppercase tracking-[0.16em] text-ink-faint`).

## Color tokens (Tailwind — use these, not raw hex)

Backgrounds: `bg-canvas` #0A0B0D (page) · `bg-panel` #0C0D10 · `bg-well` #0E1013
(inputs) · `bg-surface` #15171B (cards) · `bg-surface-hi` #15181C (dropdowns,
modals, toasts) · `bg-surface-raised` #1F2329 (secondary buttons, pills) ·
`bg-surface-nested` #121417 · `bg-surface-nav` #1A1D22 (mobile nav) ·
`bg-surface-chip` #262B32.

Borders: `border-line` #20242A (default) · `border-line-strong` #2E343B ·
`border-line-subtle` #1B1E22 (dividers) · `border-line-faint` #16191D.

Text: `text-ink` #F4F5F7 (primary) · `text-ink-bright` #C9CED4 ·
`text-ink-muted` #A8AEB6 (secondary) · `text-ink-dim` #9BA1A8 ·
`text-ink-faint` #6E757E (muted labels) · `text-ink-ghost` #4A5159.

Brand/semantic: `brand` #2E7DF2 (`brand-hover`, `brand-light`, `brand-bright`
#6BA6FF) · `success` #2BC58C (`success-bright`) · `warning` #E0A954
(`warning-deep`, `warning-light`) · `danger` #EC5B5B (`danger-soft` #E08C7A for
"worse-than-best" deltas) · `violet` #8B5CF6.

Use opacity modifiers for tints, e.g. `bg-brand/[0.12]`, `border-success/30`,
`bg-warning/[0.12]`.

## Radius & shadow tokens

Radius: `rounded` 9 · `rounded-md` 10 (inputs) · `rounded-lg` 11 (buttons) ·
`rounded-xl` 12 (cards/large buttons) · `rounded-2xl` 16 (panels) ·
`rounded-3xl` 20 (big panels) · `rounded-full` pills.

Shadow: `shadow-glow` (brand glow) · `shadow-card` · `shadow-panel` ·
`shadow-pop` (dropdowns) · `shadow-modal` · `shadow-toast`.

Animations: `animate-kpulse` (live dot), `animate-kspin` (spinners).
Utilities: `.tnum`, `.scrollx` (hidden horizontal scrollbar for chip rails),
`.scrolly` (thin styled vertical scrollbar).

Brand gradient (logo, FABs): `bg-[linear-gradient(150deg,#4A90F5,#2E7DF2)]`.

## Layout

`layouts/default.vue` already renders, around the page's `<slot />`:
`GlowBackdrop` → `SiteHeader` (responsive: desktop nav + currency pill + Войти;
mobile: logo + search + bell) → `<main>` → `SiteFooter` → `MobileBottomNav`
(mobile only, fixed). **Regular pages should NOT re-add header/footer** — just
render their content; it lands inside `<main>` (which already has bottom padding
for the mobile nav). Wrap page content in a centered container, e.g.
`<div class="mx-auto max-w-[1200px] px-4 md:px-6">`.

Pages that must be full-bleed / chrome-less (error, maintenance) set their own
layout (`definePageMeta({ layout: false })` or a custom layout) — see each task.

## Available components (auto-imported by filename, no dir prefix)

- `<AppLogo :size="40" :wordmark="true" />` — brand mark (gradient square +
  "Kurso"). `size` is the square px; wordmark scales with it.
- `<KButton variant size pill block disabled>` — `variant`:
  primary|secondary|ghost|success|danger (default primary); `size`: sm|md|lg
  (default md); `pill` (full radius), `block` (full width). Content via slot
  (icon + text supported, gap already set).
- `<KBadge tone>` — small pill label. `tone`:
  brand|success|warning|danger|violet|neutral (default neutral).
- `<KStatusDot tone pulse :size />` — status dot. `tone`:
  success|warning|danger|neutral; `pulse` toggles the live animation.
- `<KChip active pill>` — filter chip. `active` = selected (brand fill);
  `pill` = fully rounded (mobile style). Slot = label (+ optional icon).

Inline SVG icons (stroke `currentColor`, `stroke-width="2"`, round caps) as in
the designs — there is no icon library; copy the paths from the reference HTML.

## Conventions

- `<script setup lang="ts">`, Vue 3 + Nuxt 3, Tailwind utility classes.
- All UI copy in **Russian** (match the reference text exactly where given).
- Numeric values (prices, amounts, counts, rates, IDs, times) → `.tnum`.
- Internal links → `<NuxtLink to="…">`.
- No external images; use CSS/SVG and the letter/emoji avatars from the design
  (currency circles, exchanger squares with initials + brand-ish bg colors).
- Keep pages responsive: the designs give explicit desktop and mobile variants —
  implement one responsive component, not two separate DOM trees, using Tailwind
  breakpoints (`md:`), unless a task says otherwise.
