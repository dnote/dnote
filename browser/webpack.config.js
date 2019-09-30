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
