# kurso-admin

Админка Kurso (`admin.kurso.io`) — SPA на **Vite + Vue 3 + TypeScript**.

Этап 0: чистый базовый каркас проекта (без Pinia, vue-router и PrimeVue — они появятся на следующих этапах). Настроены линтинг (ESLint + Prettier) и генерация типов из OpenAPI-контракта.

## Требования

- **Node.js 24+**
- **pnpm 10+**

## Установка

```bash
pnpm install
```

## Команды

| Команда         | Описание                                                                   |
| --------------- | -------------------------------------------------------------------------- |
| `pnpm dev`      | Запуск dev-сервера Vite на порту **5174** (`http://localhost:5174`)        |
| `pnpm build`    | Проверка типов (`vue-tsc`) и продакшн-сборка в `dist/`                     |
| `pnpm preview`  | Локальный предпросмотр собранного билда                                    |
| `pnpm openapi`  | Генерация TypeScript-типов из `../api/openapi.yaml` в `src/types/api.d.ts` |
| `pnpm lint`     | Проверка кода ESLint                                                       |
| `pnpm lint:fix` | Проверка ESLint с автоисправлением                                         |
| `pnpm format`   | Форматирование кода Prettier                                               |

## Типы из API

Контракт API — единый источник правды в `../api/openapi.yaml`. После его изменения перегенерируйте типы:

```bash
pnpm openapi
```

Файл `src/types/api.d.ts` генерируется автоматически и не редактируется вручную (исключён из линтинга и форматирования).
