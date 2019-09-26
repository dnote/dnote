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

module.exports = ({ production = false, test = false } = {}) => {
  const rootUrl = process.env.ROOT_URL;

  let stripePublicKey;
  if (test) {
    stripePublicKey = 'pk_test_5926f65DQoIilZeNOiKydfoN';
  } else {
    stripePublicKey = 'pk_live_xvouPZFPDDBSIyMUSLZwkXfR';
  }

  const cdnUrl = 'https://cdn.getdnote.com';

  const compileTimeConstantForMinification = {
    __ROOT_URL__: JSON.stringify(rootUrl),
    __STRIPE_PUBLIC_KEY__: JSON.stringify(stripePublicKey),
    __CDN_URL__: JSON.stringify(cdnUrl)
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
