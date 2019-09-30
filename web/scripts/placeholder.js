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

// bundleBaseURL is the base URL from which the application javascript bundle
// In production, it should be the same as assetBaseURL. It is used for development
// environment, in which it is configured to be the webpack development server.
const bundleBaseURL = process.env.BUNDLE_BASE_URL;
// assetBaseURL is the base URL from which all assets excluding the application
// bundle is served.
const assetBaseURL = process.env.ASSET_BASE_URL;
const publicPath = process.env.PUBLIC_PATH;
const compiledPath = process.env.COMPILED_PATH;

if (bundleBaseURL === undefined) {
  throw new Error('No BUNDLE_BASE_URL environment variable found');
}
if (assetBaseURL === undefined) {
  throw new Error('No ASSET_BASE_URL environment variable found');
}
if (publicPath === undefined) {
  throw new Error('No PUBLIC_PATH environment variable found');
}
if (compiledPath === undefined) {
  throw new Error('No COMPILED_PATH environment variable found');
}

const isProduction = process.env.NODE_ENV === 'PRODUCTION';
const indexHtmlPath = `${publicPath}/index.html`;
const assetManifestPath = path.resolve(compiledPath, 'webpack-manifest.json');

let manifest;
try {
  // eslint-disable-next-line import/no-dynamic-require,global-require
  manifest = require(assetManifestPath);
  // eslint-disable-next-line no-empty
} catch (e) {
  if (isProduction) {
    throw new Error('asset manifest not found');
  }
}

function getJSBundleTag() {
  let jsFilename;
  if (isProduction) {
    jsFilename = manifest['app.js'];
  } else {
    jsFilename = '/app.js';
  }

  const jsBundleUrl = `${bundleBaseURL}${jsFilename}`;
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
    const cssBundleUrl = `${assetBaseURL}${manifest['app.css']}`;
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
