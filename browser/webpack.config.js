/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

const path = require('path');
const webpack = require('webpack');
const packageJson = require('./package.json');

const ENV = process.env.NODE_ENV;
const TARGET = process.env.TARGET;
const isProduction = ENV === 'production';

console.log(`Running webpack in ${ENV} mode`);

const webUrl = isProduction
  ? 'https://app.getdnote.com'
  : 'http://127.0.0.1:3000';
const apiUrl = isProduction
  ? 'https://api.getdnote.com'
  : 'http://127.0.0.1:5000';

const plugins = [
  new webpack.DefinePlugin({
    __API_ENDPOINT__: JSON.stringify(apiUrl),
    __WEB_URL__: JSON.stringify(webUrl),
    __VERSION__: JSON.stringify(packageJson.version)
  })
];

const moduleRules = [
  {
    test: /\.ts(x?)$/,
    exclude: /node_modules|_test\.ts(x)$/,
    loaders: ['ts-loader'],
    exclude: path.resolve(__dirname, 'node_modules')
  }
];

module.exports = env => {
  return {
    // run in production mode because of Content Security Policy error encountered
    // when running a JavaScript bundle produced in a development mode
    mode: 'production',
    entry: { popup: ['./src/scripts/popup.tsx'] },
    output: {
      filename: '[name].js',
      path: path.resolve(__dirname, 'dist', TARGET, 'scripts')
    },
    resolve: {
      extensions: ['.ts', '.tsx', '.js'],
      alias: {
        jslib: path.join(__dirname, '../jslib/src')
      },
      modules: [path.resolve('node_modules')]
    },
    module: { rules: moduleRules },
    plugins: plugins,
    optimization: {
      minimize: isProduction
    }
  };
};
