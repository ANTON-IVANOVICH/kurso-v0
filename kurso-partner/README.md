# kurso-partner

Кабинет обменников Kurso (`partner.kurso.io`) — SPA на **Vite + Vue 3 + TypeScript**.

Часть монорепозитория Kurso. Единый источник контракта с API — `../api/openapi.yaml`,
из которого генерируются типы в `src/types/api.d.ts`.

> Этап 0: чистая база проекта. Pinia, vue-router и PrimeVue пока не подключены.

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
