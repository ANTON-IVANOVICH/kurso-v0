import pluginVue from 'eslint-plugin-vue'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import prettierConfig from '@vue/eslint-config-prettier'

export default defineConfigWithVueTs(
  {
    name: 'kurso-partner/files-to-lint',
    files: ['**/*.{ts,mts,tsx,vue}'],
  },
  {
    name: 'kurso-partner/ignores',
    ignores: ['dist', 'node_modules', 'src/types/api.d.ts'],
  },
  pluginVue.configs['flat/essential'],
  vueTsConfigs.recommended,
  prettierConfig,
  {
    name: 'kurso-partner/rules',
    rules: {
      // Allow intentionally-unused args/vars when prefixed with `_`.
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_', caughtErrors: 'none' },
      ],
    },
  },
)
