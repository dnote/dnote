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

const PATHS = require('../paths');

const createPlugins = production => {
  const ret = [
    '@babel/plugin-proposal-class-properties',
    '@babel/plugin-transform-react-constant-elements',
    'react-hot-loader/babel'
  ];

  if (production) {
    ret.push('transform-react-remove-prop-types');
  }

  return ret;
};

module.exports = ({ production = false } = {}) => {
  const presets = [
    [
      '@babel/preset-env',
      {
        useBuiltIns: 'entry',
        corejs: '3',
        targets: '> 0.25%, not dead'
      }
    ],
    '@babel/preset-react'
  ];
  const plugins = createPlugins(production);

  return [
    {
      test: /\.js$|\.jsx$/,
      loader: 'babel-loader',
      options: {
        presets,
        plugins
      },
      exclude: PATHS.modules
    },
    {
      test: /\.ts(x?)$/,
      exclude: /node_modules/,
      use: [
        {
          loader: 'ts-loader'
        }
      ]
    },
    // All output '.js' files will have any sourcemaps re-processed by 'source-map-loader'.
    {
      enforce: 'pre',
      test: /\.js$/,
      loader: 'source-map-loader'
    }
  ];
};
