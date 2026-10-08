#!/usr/bin/env node

/**
 * Module Scaffolding Script
 * Usage:
 *   node scripts/new-module.mjs <module-name>
 *   pnpm new-module <module-name>
 *   npm run new-module <module-name>
 */

import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, '..');

function toCamelCase(str) {
  return str
    .toLowerCase()
    .replace(/[-_]+(.)/g, (_, c) => c.toUpperCase());
}

function toTitleCase(str) {
  return str
    .replace(/[-_]+/g, ' ')
    .split(' ')
    .filter(Boolean)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
}

function copyAndTransformDir(srcDir, destDir, replacements) {
  if (!fs.existsSync(destDir)) {
    fs.mkdirSync(destDir, { recursive: true });
  }

  const entries = fs.readdirSync(srcDir, { withFileTypes: true });

  for (const entry of entries) {
    const srcPath = path.join(srcDir, entry.name);
    // Replace 'template' with module name in filename
    let destName = entry.name.replace(/template/g, replacements.name);
    const destPath = path.join(destDir, destName);

    if (entry.isDirectory()) {
      copyAndTransformDir(srcPath, destPath, replacements);
    } else {
      let content = fs.readFileSync(srcPath, 'utf8');

      content = content
        .replaceAll('__MODULE_NAME__', replacements.name)
        .replaceAll('__MODULE_TITLE__', replacements.title)
        .replaceAll('__MODULE_DESCRIPTION__', replacements.description)
        .replaceAll('__CAMEL_NAME__', replacements.camelName)
        .replaceAll('templateModule', `${replacements.camelName}Module`)
        .replaceAll('TemplatePage', `${toTitleCase(replacements.name).replace(/\s+/g, '')}Page`)
        .replaceAll('TemplateStore', `${toTitleCase(replacements.name).replace(/\s+/g, '')}Store`)
        .replaceAll('TemplateState', `${toTitleCase(replacements.name).replace(/\s+/g, '')}State`)
        .replaceAll('createTemplateStore', `create${toTitleCase(replacements.name).replace(/\s+/g, '')}Store`)
        .replaceAll('TemplateStoreProvider', `${toTitleCase(replacements.name).replace(/\s+/g, '')}StoreProvider`)
        .replaceAll('useTemplateStore', `use${toTitleCase(replacements.name).replace(/\s+/g, '')}Store`)
        .replaceAll('fetchTemplateData', `fetch${toTitleCase(replacements.name).replace(/\s+/g, '')}Data`)
        .replaceAll('useTemplateData', `use${toTitleCase(replacements.name).replace(/\s+/g, '')}Data`)
        .replaceAll('TemplateData', `${toTitleCase(replacements.name).replace(/\s+/g, '')}Data`)
        .replaceAll('template-page', `${replacements.name}-page`)
        .replaceAll('template-api', `${replacements.name}-api`)
        .replaceAll('template-store', `${replacements.name}-store`);

      fs.writeFileSync(destPath, content, 'utf8');
    }
  }
}

async function main() {
  const rawName = process.argv[2];

  if (!rawName) {
    console.error('❌ Error: Missing module name.');
    console.log('Usage: pnpm new-module <name>');
    console.log('       npm run new-module <name>');
    process.exit(1);
  }

  const name = rawName.toLowerCase().replace(/[^a-z0-9-_]/g, '-');
  const camelName = toCamelCase(name);
  const title = toTitleCase(name);
  const description = `Manage and view ${title}`;
  const componentName = `${title.replace(/\s+/g, '')}Page`;
  const moduleVarName = `${camelName}Module`;

  const templateDir = path.join(rootDir, 'src', 'modules', '_template');
  const targetDir = path.join(rootDir, 'src', 'modules', name);
  const appRouteDir = path.join(rootDir, 'src', 'app', name);
  const appConfigPath = path.join(rootDir, 'config', 'app.config.ts');

  if (fs.existsSync(targetDir)) {
    console.error(`❌ Error: Module directory "src/modules/${name}" already exists.`);
    process.exit(1);
  }

  console.log(`📦 Scaffolding new module: "${name}" (${title})...\n`);

  // 1. Copy and transform template
  copyAndTransformDir(templateDir, targetDir, {
    name,
    camelName,
    title,
    description,
  });
  console.log(`✅ Created src/modules/${name}/`);

  // 2. Create test file in module
  const testContent = `import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import React from 'react';
import { ${moduleVarName} } from './index';
import { ${componentName} } from './pages/${name}-page';
import { TenantProvider } from '@/core/tenant/tenant-context';

describe('${title} Module', () => {
  it('has a valid ModuleDefinition structure', () => {
    expect(${moduleVarName}.id).toBe('${name}');
    expect(${moduleVarName}.routes).toHaveLength(1);
    expect(${moduleVarName}.menuItems).toHaveLength(1);
  });

  it('renders ${componentName} smoke test without crashing', () => {
    render(
      <TenantProvider>
        <${componentName} />
      </TenantProvider>
    );
    expect(screen.getByText('${title}')).toBeInTheDocument();
  });
});
`;
  fs.writeFileSync(path.join(targetDir, `${name}.test.tsx`), testContent, 'utf8');
  console.log(`✅ Created src/modules/${name}/${name}.test.tsx smoke test`);

  // 3. Create thin route in src/app/<name>/page.tsx
  if (!fs.existsSync(appRouteDir)) {
    fs.mkdirSync(appRouteDir, { recursive: true });
  }
  const appPageContent = `import { ${componentName} } from '@/modules/${name}';

export const metadata = {
  title: '${title} | Core SaaS Starter Kit',
};

export default function Page() {
  return <${componentName} />;
}
`;
  fs.writeFileSync(path.join(appRouteDir, 'page.tsx'), appPageContent, 'utf8');
  console.log(`✅ Created thin route src/app/${name}/page.tsx`);

  // 4. Register in config/app.config.ts
  if (fs.existsSync(appConfigPath)) {
    let config = fs.readFileSync(appConfigPath, 'utf8');

    // Add import statement at top
    const importStatement = `import { ${moduleVarName} } from '@/modules/${name}';\n`;
    if (!config.includes(importStatement)) {
      config = importStatement + config;
    }

    // Add module to enabledModules array
    if (!config.includes(moduleVarName)) {
      config = config.replace(
        /enabledModules:\s*\[([\s\S]*?)\]/,
        (match, inner) => {
          const trimmed = inner.trim();
          const separator = trimmed ? '\n    ' : '';
          return `enabledModules: [${inner}${separator}${moduleVarName},\n  ]`;
        }
      );
      fs.writeFileSync(appConfigPath, config, 'utf8');
      console.log(`✅ Registered ${moduleVarName} in config/app.config.ts`);
    }
  }

  console.log(`\n🎉 Module "${name}" scaffolded and registered successfully!`);
  console.log(`   Route available at: /${name}`);
  console.log(`   Definition at:      src/modules/${name}/index.ts\n`);
}

main().catch((err) => {
  console.error('Fatal scaffolding error:', err);
  process.exit(1);
});
