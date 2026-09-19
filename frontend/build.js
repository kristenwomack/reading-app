// Assembles a static, deployable site into frontend/dist/.
// Copies the frontend assets plus the repo-root data files (books.json,
// goals.json) so the site can be served from a single directory.
import { existsSync, mkdirSync, rmSync, cpSync, copyFileSync, writeFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const frontendDir = dirname(fileURLToPath(import.meta.url));
const repoRoot = join(frontendDir, '..');
const distDir = join(frontendDir, 'dist');

// Frontend assets to include (relative to frontend/).
const assets = [
    'index.html',
    'favicon.svg',
    'favicon.ico',
    'favicon-16x16.png',
    'favicon-32x32.png',
    'src',
    'styles',
    'lib',
];

// Repo-root data files that must sit next to index.html.
const dataFiles = ['books.json', 'goals.json'];

rmSync(distDir, { recursive: true, force: true });
mkdirSync(distDir, { recursive: true });

for (const asset of assets) {
    const from = join(frontendDir, asset);
    if (!existsSync(from)) {
        console.warn(`build: skipping missing asset ${asset}`);
        continue;
    }
    cpSync(from, join(distDir, asset), { recursive: true });
}

for (const file of dataFiles) {
    const from = join(repoRoot, file);
    if (!existsSync(from)) {
        throw new Error(`build: required data file not found: ${file}`);
    }
    copyFileSync(from, join(distDir, file));
}

// GitHub Pages: prevent Jekyll processing of the static assets.
writeFileSync(join(distDir, '.nojekyll'), '');

console.log(`build: wrote static site to ${distDir}`);
