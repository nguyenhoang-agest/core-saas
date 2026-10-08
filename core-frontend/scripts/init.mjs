#!/usr/bin/env node

/**
 * Starter Kit Initialization Script
 * Usage:
 *   node scripts/init.mjs --name <project-name>
 *   pnpm init-kit --name <project-name>
 *   npm run init-kit -- --name <project-name>
 */

import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, '..');

function parseArgs() {
  const args = process.argv.slice(2);
  let name = '';

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--name' || args[i] === '-n') {
      name = args[i + 1] || '';
      i++;
    } else if (args[i].startsWith('--name=')) {
      name = args[i].split('=')[1] || '';
    }
  }

  return { name: name.trim() };
}

function toTitleCase(str) {
  return str
    .replace(/[-_]+/g, ' ')
    .split(' ')
    .filter(Boolean)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
}

async function main() {
  const { name } = parseArgs();

  if (!name) {
    console.error('❌ Error: Missing required argument --name');
    console.log('Usage: pnpm init-kit --name <project-name>');
    console.log('       npm run init-kit -- --name <project-name>');
    process.exit(1);
  }

  const pkgPath = path.join(rootDir, 'package.json');
  const appConfigPath = path.join(rootDir, 'config', 'app.config.ts');
  const envExamplePath = path.join(rootDir, '.env.example');
  const envPath = path.join(rootDir, '.env');

  const titleName = toTitleCase(name);

  console.log(`🚀 Initializing SaaS Starter Kit for project: "${name}" (${titleName})...\n`);

  // 1. Update package.json
  if (fs.existsSync(pkgPath)) {
    const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
    pkg.name = name.toLowerCase().replace(/[^a-z0-9-_]/g, '-');
    pkg.version = '0.1.0';
    fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2) + '\n', 'utf8');
    console.log(`✅ Updated package.json name to "${pkg.name}"`);
  }

  // 2. Update config/app.config.ts
  if (fs.existsSync(appConfigPath)) {
    let configContent = fs.readFileSync(appConfigPath, 'utf8');
    configContent = configContent.replace(
      /appName:\s*['"][^'"]*['"]/,
      `appName: '${titleName}'`
    );
    fs.writeFileSync(appConfigPath, configContent, 'utf8');
    console.log(`✅ Updated config/app.config.ts branding to "${titleName}"`);
  }

  // 3. Update .env.example & .env
  for (const file of [envExamplePath, envPath]) {
    if (fs.existsSync(file)) {
      let content = fs.readFileSync(file, 'utf8');
      content = content.replace(
        /NEXT_PUBLIC_APP_NAME=.*/,
        `NEXT_PUBLIC_APP_NAME="${titleName}"`
      );
      fs.writeFileSync(file, content, 'utf8');
      console.log(`✅ Updated ${path.basename(file)} NEXT_PUBLIC_APP_NAME`);
    }
  }

  console.log(`\n🎉 Kit initialized successfully!`);
  console.log(`   Next steps:`);
  console.log(`   1. Add modules: pnpm new-module <name>`);
  console.log(`   2. Run dev:     pnpm dev / npm run dev\n`);
}

main().catch((err) => {
  console.error('Fatal initialization error:', err);
  process.exit(1);
});
