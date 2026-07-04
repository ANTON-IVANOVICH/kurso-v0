# kurso-partner

Кабинет обменника Kurso (`partner.kurso.io`) — SPA на **Vite + Vue 3 + TypeScript +
PrimeVue**. Плотный десктопный кабинет: главная с метриками, курсы и состояние
фида, ответы на отзывы, статистика трафика, профиль, жалобы и биллинг.

Часть монорепозитория Kurso. Единый источник контракта с API — `../api/openapi.yaml`,
из которого генерируются типы в `src/types/api.d.ts` (доменные ответы кабинета
описаны в `src/types/models.ts`).

## Авторизация

Как в админке: access-токен живёт только в памяти, refresh — в httpOnly-cookie
(`kurso_partner_rt`), которую читает сервер. Данные кабинета читаются через
Pinia Colada. Тестовый доступ: `partner@kurso.io` / `partner12345` (обменник
CryptoBridge).

По умолчанию SPA обращается к API на `http://localhost:8080`
(переопределяется через `VITE_API_BASE`).

## Требования

- **Node.js** 24+
- **pnpm** 10+

## Установка

```bash
pnpm install
```

## Команды

- `pnpm dev` — запуск dev-сервера Vite на порту **5175** (`http://localhost:5175`)
- `pnpm build` — проверка типов (`vue-tsc`) и production-сборка в `dist/`
- `pnpm preview` — локальный предпросмотр собранного проекта
- `pnpm openapi` — генерация типов из `../api/openapi.yaml` в `src/types/api.d.ts`
- `pnpm lint` — проверка кода через ESLint
- `pnpm lint:fix` — проверка ESLint с автоисправлением
- `pnpm format` — форматирование кода через Prettier
