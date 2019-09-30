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

const webpackConf = require('./webpack/rules/javascript');

module.exports = config => {
  config.set({
    frameworks: ['mocha'],
    reporters: ['mocha'],
    browsers: ['ChromeHeadlessNoSandbox'],
    customLaunchers: {
      ChromeHeadlessNoSandbox: {
        base: 'ChromeHeadless',
        // specified because the default user of gitlab ci docker image is root, and chrome does not
        // support running as root without no-sandbox
        flags: ['--no-sandbox']
      }
    },
    files: [
      'node_modules/regenerator-runtime/runtime.js',
      './src/**/*_test.ts'
    ],
    preprocessors: {
      './src/**/*.ts': ['webpack']
    },
    webpack: {
      mode: 'none',
      resolve: {
        extensions: ['.js', '.ts']
      },
      module: {
        rules: [...webpackConf({ produciton: false, browser: true })]
      }
    },
    mochaReporter: {
      showDiff: true
    }
  });
};
