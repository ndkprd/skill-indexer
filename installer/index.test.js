'use strict';

const test = require('node:test');
const assert = require('node:assert/strict');
const http = require('node:http');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const AdmZip = require('adm-zip');

const { installFromUrl, parseArgs } = require('./index.js');

function buildFixtureZipBuffer() {
  const zip = new AdmZip();
  zip.addFile(
    'myskill/SKILL.md',
    Buffer.from('---\nname: myskill\ndescription: A fixture skill.\n---\n\nHello from myskill.\n'),
  );
  zip.addFile('myskill/references/note.md', Buffer.from('Some reference notes.\n'));
  return zip.toBuffer();
}

function startFixtureServer(zipBuffer) {
  const server = http.createServer((req, res) => {
    if (req.url === '/fixture.zip') {
      res.writeHead(200, { 'Content-Type': 'application/zip' });
      res.end(zipBuffer);
      return;
    }
    if (req.url === '/not-a-zip.zip') {
      res.writeHead(200, { 'Content-Type': 'application/zip' });
      res.end('this is definitely not a zip file');
      return;
    }
    res.writeHead(404);
    res.end('not found');
  });

  return new Promise((resolve) => {
    server.listen(0, '127.0.0.1', () => resolve(server));
  });
}

function makeTempDest() {
  return fs.mkdtempSync(path.join(os.tmpdir(), 'skillstore-install-test-'));
}

test('installer end-to-end', async (t) => {
  const zipBuffer = buildFixtureZipBuffer();
  const server = await startFixtureServer(zipBuffer);
  const port = server.address().port;
  const baseUrl = `http://127.0.0.1:${port}`;

  t.after(() => {
    server.close();
  });

  await t.test('extracts the zip contents as-is into dest', async () => {
    const dest = makeTempDest();
    try {
      const extractedPaths = await installFromUrl(`${baseUrl}/fixture.zip`, dest);

      const skillMdPath = path.join(dest, 'myskill', 'SKILL.md');
      const noteMdPath = path.join(dest, 'myskill', 'references', 'note.md');

      assert.ok(fs.existsSync(skillMdPath), `${skillMdPath} should exist`);
      assert.ok(fs.existsSync(noteMdPath), `${noteMdPath} should exist`);
      assert.match(fs.readFileSync(skillMdPath, 'utf8'), /name: myskill/);

      assert.deepEqual(extractedPaths, [path.join(dest, 'myskill')]);
    } finally {
      fs.rmSync(dest, { recursive: true, force: true });
    }
  });

  await t.test('creates dest when it does not already exist', async () => {
    const parent = makeTempDest();
    const dest = path.join(parent, 'nested', 'skills');
    try {
      await installFromUrl(`${baseUrl}/fixture.zip`, dest);
      assert.ok(fs.existsSync(path.join(dest, 'myskill', 'SKILL.md')));
    } finally {
      fs.rmSync(parent, { recursive: true, force: true });
    }
  });

  await t.test('rejects with a clear message on a non-200 response', async () => {
    const dest = makeTempDest();
    try {
      await assert.rejects(
        () => installFromUrl(`${baseUrl}/does-not-exist.zip`, dest),
        (err) => {
          assert.match(err.message, /404/);
          return true;
        },
      );
    } finally {
      fs.rmSync(dest, { recursive: true, force: true });
    }
  });

  await t.test('rejects with a clear message on a corrupt zip', async () => {
    const dest = makeTempDest();
    try {
      await assert.rejects(
        () => installFromUrl(`${baseUrl}/not-a-zip.zip`, dest),
        (err) => {
          assert.match(err.message, /not a valid zip archive/);
          return true;
        },
      );
    } finally {
      fs.rmSync(dest, { recursive: true, force: true });
    }
  });

  await t.test('rejects with a clear message on an unreachable host', async () => {
    const dest = makeTempDest();
    try {
      await assert.rejects(
        () => installFromUrl('http://127.0.0.1:1/fixture.zip', dest),
        (err) => {
          assert.match(err.message, /could not reach/);
          return true;
        },
      );
    } finally {
      fs.rmSync(dest, { recursive: true, force: true });
    }
  });
});

test('parseArgs', async (t) => {
  await t.test('defaults dest to ./skills', () => {
    const { url, dest } = parseArgs(['http://example.com/a.zip']);
    assert.equal(url, 'http://example.com/a.zip');
    assert.equal(dest, './skills');
  });

  await t.test('accepts --dest as a separate argument', () => {
    const { url, dest } = parseArgs(['http://example.com/a.zip', '--dest', './out']);
    assert.equal(url, 'http://example.com/a.zip');
    assert.equal(dest, './out');
  });

  await t.test('accepts --dest=value form', () => {
    const { url, dest } = parseArgs(['http://example.com/a.zip', '--dest=./out']);
    assert.equal(url, 'http://example.com/a.zip');
    assert.equal(dest, './out');
  });
});
