const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');

const version = require('./package.json').version;
const platform = process.platform;
const arch = process.arch;

const GOOS = platform === 'win32' ? 'windows' : platform === 'darwin' ? 'darwin' : 'linux';
let GOARCH = arch === 'x64' ? 'x86_64' : arch === 'arm64' ? 'arm64' : arch === 'ia32' ? 'i386' : arch;

const ext = GOOS === 'windows' ? 'zip' : 'tar.gz';
const filename = `studex-cli_${GOOS}_${GOARCH}.${ext}`;
const url = `https://github.com/itshivams/Studex-CLI/releases/download/v${version}/${filename}`;

const binDir = path.join(__dirname, 'bin');
if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
}

const exeExt = GOOS === 'windows' ? '.exe' : '';
const binPath = path.join(binDir, `studex-cli${exeExt}`);

console.log(`Downloading studex-cli from ${url}...`);

const download = (url, dest) => {
    return new Promise((resolve, reject) => {
        https.get(url, (response) => {
            if (response.statusCode === 301 || response.statusCode === 302) {
                return download(response.headers.location, dest).then(resolve, reject);
            }
            if (response.statusCode !== 200) {
                reject(new Error(`Failed to download: ${response.statusCode}`));
                return;
            }
            const file = fs.createWriteStream(dest);
            response.pipe(file);
            file.on('finish', () => {
                file.close(resolve);
            });
        }).on('error', reject);
    });
};

const extract = (archivePath, destDir) => {
    if (ext === 'zip') {
        if (process.platform === 'win32') {
            execSync(`powershell -command "Expand-Archive -Path '${archivePath}' -DestinationPath '${destDir}' -Force"`);
        } else {
            execSync(`unzip -o -q "${archivePath}" -d "${destDir}"`);
        }
    } else {
        execSync(`tar -xzf "${archivePath}" -C "${destDir}"`);
    }
    try {
        fs.unlinkSync(archivePath);
    } catch (e) { }
}

async function main() {
    const archivePath = path.join(__dirname, filename);
    try {
        await download(url, archivePath);
        console.log("Extracting...");
        extract(archivePath, binDir);
        if (GOOS !== 'windows') {
            fs.chmodSync(binPath, 0o755);
        }
        console.log("Installation completed successfully.");
    } catch (err) {
        console.error("Install failed. Please build manually or check your internet connection.");
        console.error(err);
        process.exit(1);
    }
}

main();
