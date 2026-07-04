# Kurso

Агрегатор курсов криптообменников. Монорепозиторий с одним Go-бэкендом и тремя
фронтендами.

## Структура

| Каталог | Что это | Стек | Порт |
| --- | --- | --- | --- |
| `kurso-api/` | Бэкенд-монолит (public/admin/partner API) | Go 1.26, chi, pgx, Redis | 8080 |
| `kurso-web/` | Пользовательский сайт kurso.io | Nuxt 3, TypeScript, Tailwind | 3000 |
| `kurso-admin/` | Админка admin.kurso.io | Vite, Vue 3, TypeScript | 5174 |
| `kurso-partner/` | Кабинет обменников partner.kurso.io | Vite, Vue 3, TypeScript | 5175 |
| `api/` | `openapi.yaml` — единый контракт | OpenAPI 3.0 | — |
| `kurso-api/migrations/` | Миграции схемы БД | goose (SQL) | — |

Бэкенд организован по гексагональной архитектуре: `internal/domain` (ядро и
порты), `internal/service` (сценарии), `internal/adapter` (HTTP-адаптер),
`internal/platform` (инфраструктура: логгер, Postgres, Redis, HTTP-сервер).

## Требования

- Go **1.26+**
- Node **20+** (протестировано на 24) и **pnpm 9+**
- Локальные **Postgres** и **Redis** для разработки (по умолчанию `localhost:5432` / `6379`)
- GNU Make
- Docker + Docker Compose — только для прод-стека (`make stack-up`), для дева не нужен

## Разработка (локально, без Docker)

Дев использует **твои хостовые Postgres и Redis** (`localhost:5432` / `6379`).
Docker для разработки не нужен.

```bash
# 1. Зависимости (husky-хуки, Go-модули, все фронтенды)
make install

# 2. Один раз: роль + БД kurso в локальном Postgres
make db-create
#   Другие логин/пароль/имя БД? Задай свой DATABASE_URL, напр.:
#   export DATABASE_URL=postgres://me@localhost:5432/kurso?sslmode=disable

# 3. Миграции в локальную БД
make migrate

# 4. Всё разом — API + 3 фронта (Ctrl+C гасит всё)
make dev
#   api :8080 · web :3000 · admin :5174 · partner :5175
#   curl localhost:8080/readyz → статус зависимостей
```

Или по отдельности: `make run` (только API) и `pnpm dev` в нужном фронте.

## Полный стек в Docker (прод-паритет)

Docker-стек — это прод-развёртывание, для дева он не нужен. Одной командой:

```bash
make stack-up      # соберёт и поднимет pg + redis + api + 3 фронта
make stack-logs
make stack-down    # остановить (stack-down-v — со сносом данных)
```

Нет локальных Postgres/Redis? Можно поднять только их в Docker: `make db-up`
(при занятом 5432/6379 — `POSTGRES_PORT=5433 REDIS_PORT=6380 make db-up`).

## Полезные команды

```bash
make help              # список всех команд
make build             # собрать бэкенд и все фронтенды
make test              # тесты бэкенда
make lint              # go vet + gofmt + eslint по всем фронтендам
make generate          # регенерация типов из api/openapi.yaml (Go + TS)
make migrate-status    # статус миграций
make migrate-create name=add_something
make stack-up / make stack-down   # прод-стек в Docker
```

## Контракт и кодогенерация

`api/openapi.yaml` — единый источник правды. `make generate` перегенерирует:

- Go-типы → `kurso-api/internal/adapter/http/openapi/types.gen.go` (oapi-codegen)
- TS-типы → `<frontend>/**/types/api.d.ts` (openapi-typescript)

После правок спецификации запусти `make generate` и закоммить сгенерированные файлы.

## Pre-commit

В корне настроены husky + lint-staged (`make install` их активирует). На
`git commit` запускаются `gofmt`/`go vet` для изменённых Go-файлов и
`eslint --fix`/`prettier --write` для изменённых файлов фронтендов.

## CI

GitHub Actions (`.github/workflows/`): на каждый PR — линт и тесты бэкенда,
линт/сборка фронтендов, проверка актуальности кодогена и сборка Docker-образа
бэкенда. Деплой пока выполняется вручную.
