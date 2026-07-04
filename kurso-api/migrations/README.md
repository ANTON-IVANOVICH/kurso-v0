# Миграции (goose)

SQL-миграции схемы БД Kurso. Управляются [goose](https://github.com/pressly/goose),
подключённым как Go-tool в `kurso-api/go.mod` (не требует отдельной установки).

## Команды

Из корня репозитория (через `Makefile`, использует `DATABASE_URL` из окружения):

```bash
make migrate            # применить все миграции (goose up)
make migrate-down       # откатить одну миграцию
make migrate-status     # статус миграций
make migrate-create name=add_something   # новая миграция
```

Напрямую через goose:

```bash
cd kurso-api
go tool goose -dir migrations postgres "$DATABASE_URL" up
go tool goose -dir migrations postgres "$DATABASE_URL" status
```

## Порядок

Миграции применяются по номеру в имени файла. Порядок учитывает внешние ключи:

| #     | Файл                        | Таблицы |
|-------|-----------------------------|---------|
| 00001 | extensions_and_helpers      | pgcrypto/citext/pg_trgm, `set_updated_at()` |
| 00002 | currencies                  | `currencies`, `currency_aliases` |
| 00003 | exchangers                  | `exchangers` |
| 00004 | directions                  | `directions` |
| 00005 | identities                  | `users`, `admins`, `exchanger_users` |
| 00006 | rates                       | `rates`, `rates_history` (партиционирована) |
| 00007 | parser_configs              | `parser_configs` |
| 00008 | alerts                      | `alerts`, `triggered_alerts`, `push_subscriptions` |
| 00009 | reviews                     | `reviews`, `review_replies`, `review_reports` |
| 00010 | clickouts_and_payouts       | `clickouts`, `partner_payouts` |
| 00011 | audit_and_security          | `admin_audit_log`, `security_events` |

## Партиционирование `rates_history`

`rates_history` партиционирована по диапазону `recorded_at`. На этапе 0 создана
только DEFAULT-партиция, чтобы вставки всегда проходили. Помесячные партиции и их
автосоздание (pg_partman / cron) добавляются на более поздних этапах.
