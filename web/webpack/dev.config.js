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

const PATHS = require('./paths');
const rules = require('./rules');
const plugins = require('./plugins');
const resolve = require('./resolve');
const externals = require('./externals');

module.exports = env => {
  const isTest = env.isTest === 'true';

  return {
    mode: 'development',
    devtool: 'eval',
    entry: { app: ['./src/client'] },
    output: {
      path: PATHS.output,
      filename: '[name].js',
      publicPath: PATHS.public
    },
    devServer: {
      publicPath: PATHS.public,
      port: 8080
    },
    module: { rules: rules({ production: false }) },
    resolve,
    plugins: plugins({
      production: false,
      test: isTest
    }),
    externals
  };
};
