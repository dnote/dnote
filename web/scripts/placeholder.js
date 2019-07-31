/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

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
