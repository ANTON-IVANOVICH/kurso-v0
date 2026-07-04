// Monorepo pre-commit routing.
//
// There is a single git repo (and therefore a single husky hooksPath), so this
// root config dispatches staged files to the right sub-project's tooling.
// Frontend commands `cd` into their project so ESLint/Prettier resolve the
// project-local flat config and node_modules.

import path from 'node:path'

const ROOT = process.cwd()

// Build `sh -c 'cd <dir> && <cmd> <relative files>'`. lint-staged splits the
// returned string with string-argv (no shell), so the single-quoted script is
// passed to sh as one argument and sh interprets the `&&`.
const scoped = (dir, cmd) => (files) => {
  const rel = files.map((f) => path.relative(path.join(ROOT, dir), f)).join(' ')
  return `sh -c 'cd ${dir} && ${cmd} ${rel}'`
}

const frontend = (dir) => ({
  [`${dir}/**/*.{ts,tsx,vue,js,mjs,cjs}`]: [
    scoped(dir, 'pnpm exec eslint --fix'),
    scoped(dir, 'pnpm exec prettier --write'),
  ],
  [`${dir}/**/*.{json,css,scss,md,yml,yaml,html}`]: [scoped(dir, 'pnpm exec prettier --write')],
})

export default {
  // Go: format the staged files, then vet the whole module.
  'kurso-api/**/*.go': (files) => [
    `gofmt -w ${files.join(' ')}`,
    "sh -c 'cd kurso-api && go vet ./...'",
  ],
  ...frontend('kurso-web'),
  ...frontend('kurso-admin'),
  ...frontend('kurso-partner'),
}
