# kurso-web

Пользовательский сайт [kurso.io](https://kurso.io) — агрегатор курсов обменников.
Построен на **Nuxt 3** + **TypeScript** + **Tailwind CSS**.

## Требования

- **Node.js 24**
- **pnpm** (v10)

## Установка

```bash
pnpm install
```

## Команды

| Команда         | Описание                                                              |
| --------------- | -------------------------------------------------------------------- |
| `pnpm dev`      | Запуск дев-сервера на `http://localhost:3000`                        |
| `pnpm build`    | Продакшн-сборка приложения                                            |
| `pnpm preview`  | Локальный предпросмотр продакшн-сборки                               |
| `pnpm generate` | Статическая генерация сайта (Nuxt)                                    |
| `pnpm openapi`  | Генерация TypeScript-типов из `../api/openapi.yaml` в `types/api.d.ts` |
| `pnpm lint`     | Проверка кода ESLint                                                  |
| `pnpm lint:fix` | Проверка и авто-исправление ESLint                                    |
| `pnpm format`   | Форматирование кода Prettier                                          |
| `pnpm typecheck`| Проверка типов через `vue-tsc`                                        |

## Типы API

Контракт с бэкендом (`kurso-api`) описан в `../api/openapi.yaml`.
Сгенерировать типы:

```bash
pnpm openapi
```

Результат попадает в `types/api.d.ts` (файл сгенерирован автоматически, руками не редактируется).
