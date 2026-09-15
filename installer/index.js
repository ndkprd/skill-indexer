'use strict';

const fs = require('fs');
const path = require('path');
const AdmZip = require('adm-zip');

const DEFAULT_DEST = './skills';

/**
 * Downloads the zip at `url` and extracts it into `dest`.
 *
 * The zip is expected to already carry its own top-level skill directory in
 * every entry path (e.g. `vue/SKILL.md`), so extraction writes directly into
 * `dest` with no additional wrapping folder.
 *
 * @param {string} url - HTTP(S) URL of the zip to install.
 * @param {string} dest - Directory to extract into (created if missing).
 * @returns {Promise<string[]>} Absolute paths of the top-level entries written under `dest`.
 */
async function installFromUrl(url, dest) {
  const zipBuffer = await downloadZip(url);
  const zip = openZip(zipBuffer, url);
  extractZip(zip, dest);
  return listTopLevelPaths(zip, dest);
}

async function downloadZip(url) {
  let response;
  try {
    response = await fetch(url);
  } catch (err) {
    throw new Error(`could not reach ${url}: ${err.message}`);
  }

  if (!response.ok) {
    throw new Error(`request to ${url} failed with status ${response.status} ${response.statusText}`.trim());
  }

  const arrayBuffer = await response.arrayBuffer();
  return Buffer.from(arrayBuffer);
}

function openZip(zipBuffer, url) {
  try {
    return new AdmZip(zipBuffer);
  } catch (err) {
    throw new Error(`downloaded content from ${url} is not a valid zip archive: ${err.message}`);
  }
}

function extractZip(zip, dest) {
  fs.mkdirSync(dest, { recursive: true });
  try {
    zip.extractAllTo(dest, true);
  } catch (err) {
    throw new Error(`failed to extract zip into ${dest}: ${err.message}`);
  }
}

function listTopLevelPaths(zip, dest) {
  const topLevelNames = new Set();
  for (const entry of zip.getEntries()) {
    const [firstSegment] = entry.entryName.split('/');
    if (firstSegment) {
      topLevelNames.add(firstSegment);
    }
  }
  return [...topLevelNames].map((name) => path.join(dest, name));
}

function parseArgs(argv) {
  let url;
  let dest = DEFAULT_DEST;

  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (arg === '--dest') {
      dest = argv[i + 1];
      i++;
    } else if (arg.startsWith('--dest=')) {
      dest = arg.slice('--dest='.length);
    } else if (!url && !arg.startsWith('-')) {
      url = arg;
    }
  }

  return { url, dest };
}

async function runCli(argv) {
  const { url, dest } = parseArgs(argv);

  if (!url) {
    console.error('Usage: skill-repo-store-install <zip-url> [--dest ./skills]');
    process.exitCode = 1;
    return;
  }

  try {
    const extractedPaths = await installFromUrl(url, dest);
    const resolvedDest = path.resolve(dest);
    console.log(`Installed ${extractedPaths.length} item(s) from ${url} into ${resolvedDest}`);
    for (const extractedPath of extractedPaths) {
      console.log(`  ${extractedPath}`);
    }
  } catch (err) {
    console.error(`skill-repo-store-install: ${err.message}`);
    process.exitCode = 1;
  }
}

if (require.main === module) {
  runCli(process.argv.slice(2));
}

module.exports = { installFromUrl, parseArgs };
