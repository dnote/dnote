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

const webpack = require('webpack');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const ManifestPlugin = require('webpack-manifest-plugin');

module.exports = ({
  production = false,
  test = false,
  standalone = false
} = {}) => {
  let domain;
  if (production) {
    domain = 'dnote.io';
  } else {
    domain = '127.0.0.1';
  }

  let basename;
  if (standalone) {
    basename = '/';
  } else {
    basename = '/app';
  }

  let baseURL;
  if (production) {
    baseURL = `https://${domain}${basename}`;
  } else {
    baseURL = `http://${domain}:3000${basename}`;
  }

  let stripePublicKey;
  if (test) {
    stripePublicKey = 'pk_test_5926f65DQoIilZeNOiKydfoN';
  } else {
    stripePublicKey = 'pk_live_xvouPZFPDDBSIyMUSLZwkXfR';
  }

  const compileTimeConstantForMinification = {
    __PRODUCTION__: production,
    __DEVELOPMENT__: !production,
    __DOMAIN__: JSON.stringify(domain),
    __BASE_URL__: JSON.stringify(baseURL),
    __BASE_NAME__: JSON.stringify(basename),
    __STRIPE_PUBLIC_KEY__: JSON.stringify(stripePublicKey)
  };

  if (!production) {
    return [
      new webpack.DefinePlugin(compileTimeConstantForMinification),
      new webpack.HotModuleReplacementPlugin(),
      new webpack.NoEmitOnErrorsPlugin()
    ];
  }

  return [
    new webpack.DefinePlugin(compileTimeConstantForMinification),
    new MiniCssExtractPlugin({
      filename: '[contenthash].css',
      allChunks: true
    }),
    new ManifestPlugin({
      fileName: 'webpack-manifest.json'
    })
  ];
};
