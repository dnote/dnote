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

const MiniCssExtractPlugin = require('mini-css-extract-plugin');

module.exports = ({ production = false } = {}) => {
  const createScssLoaders = (sourceMap = true, useModules = false) => {
    let localIdentName;
    if (production) {
      localIdentName = '[local]_[hash:base64:5]';
    } else {
      localIdentName = '[name]__[local]___[hash:base64:5]';
    }

    let modules;
    if (useModules) {
      modules = {
        localIdentName
      };
    } else {
      modules = false;
    }

    return [
      {
        loader: 'css-loader',
        options: {
          sourceMap,
          modules,
          importLoaders: 2
        }
      },
      {
        loader: 'postcss-loader',
        options: {
          sourceMap,
          ident: 'postcss',
          plugins: () => [require('autoprefixer')(), require('cssnano')()]
        }
      },
      {
        loader: 'sass-loader',
        options: {
          sourceMap
        }
      }
    ];
  };

  const createBrowserLoaders = extractCssToFile => loaders => {
    if (extractCssToFile) {
      return [MiniCssExtractPlugin.loader, ...loaders];
    }
    return [{ loader: 'style-loader' }, ...loaders];
  };

  const scssLoaders = createBrowserLoaders(production)(
    createScssLoaders(true, false)
  );
  const scssModuleLoaders = createBrowserLoaders(production)(
    createScssLoaders(true, true)
  );

  return [
    {
      test: /\.global\.scss$/,
      use: scssLoaders
    },
    {
      test: /\.scss$/,
      exclude: /\.global\.scss$/,
      use: scssModuleLoaders
    }
  ];
};
