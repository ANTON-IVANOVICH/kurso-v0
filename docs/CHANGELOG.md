# Changelog

Все заметные изменения проекта фиксируются здесь. Формат основан на
[Keep a Changelog](https://keepachangelog.com/ru/1.1.0/), проект придерживается
[семантического версионирования](https://semver.org/lang/ru/).

## [Unreleased]

_Пока пусто._

## [0.1.0] — 2026-07-04 — Этап 0: Фундамент

Инфраструктурный этап: рабочее окружение, схема БД и контракт API. Кода фич нет.

### Общее / монорепо

- Монорепо из четырёх проектов: `kurso-api` (Go), `kurso-web` (Nuxt), `kurso-admin`
  и `kurso-partner` (Vue).
- Корневой `README.md` и README в каждом проекте.
- `.gitignore` (в т.ч. `docs/PLAN.md` не версионируется).

### Backend (`kurso-api`)

- Go-модуль, точка входа `cmd/kurso/main.go`, конфигурация через env, `slog`-логгер
  (text в dev / json в prod), graceful shutdown.
- Гексагональная структура `internal/`: `domain` (ядро и порты), `service`
  (сценарии), `adapter/http` (входящий HTTP-адаптер на chi), `platform`
  (logger, Postgres через pgx/pgxpool, Redis через go-redis, HTTP-сервер).
- Эндпоинты: `GET /healthz`, `GET /readyz` (проверка Postgres + Redis),
  `GET /api/v1/currencies` (заглушка).
- Multi-stage `Dockerfile` (distroless runtime).
- Unit-тесты: `config`, health-handler.

### База данных

- 11 миграций goose со всеми таблицами этапа 0 (20 шт.): `users`, `admins`,
  `exchanger_users`, `exchangers`, `currencies`, `currency_aliases`, `directions`,
  `rates`, `rates_history` (партиционирована по `recorded_at`), `parser_configs`,
  `alerts`, `triggered_alerts`, `push_subscriptions`, `reviews`, `review_replies`,
  `review_reports`, `clickouts`, `partner_payouts`, `admin_audit_log`,
  `security_events`.
- goose подключён как Go-tool (без отдельной установки).

### Контракт и кодогенерация

- `api/openapi.yaml` — единый источник правды (OpenAPI 3.0).
- oapi-codegen генерирует Go-типы; openapi-typescript — TS-типы во всех трёх
  фронтендах. Регенерация — `make generate`.

### Фронтенды

- `kurso-web`: Nuxt 3 + TypeScript + Tailwind.
- `kurso-admin`, `kurso-partner`: Vite + Vue 3 + TypeScript.
- ESLint + Prettier в каждом; сгенерированные типы API из общего контракта.

### Инфраструктура и DX

- Локальная разработка идёт на **хостовых Postgres/Redis** (`localhost:5432` /
  `6379`), без Docker. `make db-create` заводит роль и БД `kurso` в локальном
  Postgres.
- `docker-compose.yml` — **прод-стек** (self-host / прод-паритет) целиком:
  Postgres 16 + Redis 7 + `api` + три фронтенда (прод-образы: `kurso-web` —
  Nuxt SSR на node, `kurso-admin` / `kurso-partner` — статика в nginx со
  SPA-fallback). Управление — `make stack-up` / `stack-down`; порты хоста
  переопределяемы через `POSTGRES_PORT` / `REDIS_PORT`.
- `Makefile`: `dev` (API + 3 фронта разом), `run`, `migrate`, `generate`, `test`,
  `lint`, `build`, `db-create`, `stack-up`/`stack-down` и др.
- Pre-commit: husky + lint-staged в корне монорепо (один `core.hooksPath`);
  `gofmt`/`go vet` для Go, `eslint --fix`/`prettier --write` для фронтендов.

### CI/CD

- GitHub Actions `ci.yml`: бэкенд (gofmt, vet, тесты, сборка, проверка
  актуальности кодогена) + фронтенды матрицей (lint, build).
- `docker.yml`: сборка Docker-образа бэкенда (без пуша). Деплой — вручную.
