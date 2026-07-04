import pluginVue from 'eslint-plugin-vue'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import prettierConfig from '@vue/eslint-config-prettier'

export default defineConfigWithVueTs(
  {
    name: 'kurso-admin/files-to-lint',
    files: ['**/*.{ts,mts,tsx,vue}'],
  },
  {
    name: 'kurso-admin/ignores',
    ignores: ['dist', 'node_modules', 'src/types/api.d.ts'],
  },
  pluginVue.configs['flat/essential'],
  vueTsConfigs.recommended,
  prettierConfig,
  {
    name: 'kurso-admin/rules',
    rules: {
      // Allow intentionally-unused args/vars when prefixed with `_`
      // (e.g. contract-ready function signatures ahead of the backend).
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_', caughtErrors: 'none' },
      ],
    },
  },
)
