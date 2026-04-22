#!/usr/bin/env node

const { spawnSync } = require('child_process');
const path = require('path');
const fs = require('fs');

const platform = process.platform;
const exeExt = platform === 'win32' ? '.exe' : '';
const binPath = path.join(__dirname, 'bin', `studex-cli${exeExt}`);

if (!fs.existsSync(binPath)) {
    console.error(`Error: studex-cli binary not found at ${binPath}.`);
    console.error('Please ensure the installation completed successfully.');
    process.exit(1);
}

const args = process.argv.slice(2);
const result = spawnSync(binPath, args, { stdio: 'inherit' });

if (result.error) {
    console.error('Error executing studex-cli:', result.error.message);
    process.exit(1);
}

process.exit(result.status || 0);
