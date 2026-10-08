import nextConfig from 'eslint-config-next';

const config = [
  ...nextConfig,
  {
    // Rule: core/ must not import from domain modules (src/modules/*)
    files: ['src/core/**/*.{ts,tsx}'],
    rules: {
      'no-restricted-imports': [
        'error',
        {
          patterns: [
            {
              group: ['@/modules/**', '**/src/modules/**'],
              message:
                'Architecture boundary violation: core/ must not import from domain modules (src/modules/*).',
            },
          ],
        },
      ],
    },
  },
  {
    // Rule: domain modules must not import other domain modules' internal paths
    files: ['src/modules/**/*.{ts,tsx}'],
    rules: {
      'no-restricted-imports': [
        'error',
        {
          patterns: [
            {
              group: ['@/modules/*/**'],
              message:
                "Architecture boundary violation: modules must not import other modules' internals. Import only from module public entrypoint (@/modules/<name>).",
            },
          ],
        },
      ],
    },
  },
];

export default config;
