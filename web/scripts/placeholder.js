#!/usr/bin/env node

// placeholder.js replaces the placeholders in index.html with real values
// It is needed to load assets whose paths are not fixed because they change
// every time they are generated.

const fs = require('fs');
const path = require('path');

const baseURL = process.env.BASE_URL;
const assetBaseURL = process.env.ASSET_BASE_URL;
const publicPath = process.env.PUBLIC_PATH;
const compiledPath = process.env.COMPILED_PATH;
if (!publicPath) {
  throw new Error('No PUBLIC_PATH environment variable found');
}
if (!compiledPath) {
  throw new Error('No COMPILED_PATH environment variable found');
}

const indexHtmlPath = `${publicPath}/index.html`;

const assetManifestPath = path.resolve(compiledPath, 'webpack-manifest.json');
let manifest;

try {
  // eslint-disable-next-line import/no-dynamic-require,global-require
  manifest = require(assetManifestPath);
  // eslint-disable-next-line no-empty
} catch (e) {}

const isProduction = process.env.NODE_ENV === 'PRODUCTION';

function getJSBundleTag() {
  let jsBundleUrl;
  if (isProduction) {
    jsBundleUrl = `${baseURL}${manifest['app.js']}`;
  } else {
    jsBundleUrl = `${baseURL}/dist/app.js`;
  }

  return `<script src="${jsBundleUrl}"></script>`;
}

// Replace the placeholders with real values
fs.readFile(indexHtmlPath, 'utf8', (err, data) => {
  if (err) {
    console.log('Error while reading index.html');
    console.log(err);
    process.exit(1);
  }

  const jsBundleTag = getJSBundleTag();
  let result = data.replace(/<!--JS_BUNDLE_PLACEHOLDER-->/g, jsBundleTag);
  result = result.replace(/<!--ASSET_BASE_PLACEHOLDER-->/g, assetBaseURL);

  if (isProduction) {
    const cssBundleUrl = `${baseURL}${manifest['app.css']}`;
    const cssBundleTag = `<link rel="stylesheet" href="${cssBundleUrl}" />`;

    result = result.replace(/<!--CSS_BUNDLE_PLACEHOLDER-->/g, cssBundleTag);
  }

  fs.writeFile(indexHtmlPath, result, 'utf8', writeErr => {
    if (writeErr) {
      console.log('Error while writing index.html');
      console.log(writeErr);
      process.exit(1);
    }
  });
});
