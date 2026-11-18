import path from 'path';
import { fileURLToPath } from 'url';
import { createRequire } from 'module';
import HtmlWebpackPlugin from 'html-webpack-plugin';
import webpack from 'webpack';
import CopyPlugin from 'copy-webpack-plugin';

const filename = fileURLToPath(import.meta.url);
const dirname = path.dirname(filename);
const require = createRequire(import.meta.url);
const pkg = require('./package.json');

// Dev server only: in standalone mode there is no shell that injects the Smart Client session
// token, so ACCESS_TOKEN (a bearer token for the XProtect IDP) is attached to proxied requests.
// BACKEND_URL is the Python backend (see ../app/src/main.py).
const accessToken = process.env.ACCESS_TOKEN;
const backendUrl = process.env.BACKEND_URL || "http://localhost:8000";

export default (env = {}, argv = {}) => {
  const isStandalone = env.hasOwnProperty('standalone');

  return {
    entry: './src/index.tsx',
    // The scripts pass --mode; webpack then also defines process.env.NODE_ENV accordingly
    // ("production" for build/build:federation), so it must not be defined again below.
    mode: argv.mode || process.env.NODE_ENV || 'development',
    devServer: {
      port: 3003,
      historyApiFallback: true,
      hot: true,
      headers: {
        'Access-Control-Allow-Origin': '*',
      },
      static: {
        directory: path.join(dirname, 'public'),
      },
      proxy: [
        {
          context: ['/api'],
          // Origin only: the proxy forwards the full `/api/...` path, so a target ending in
          // `/api` would reach the upstream as `/api/api/...`.
          target: backendUrl,
          changeOrigin: true,
          secure: false,
          ...(accessToken ? { headers: { Authorization: `Bearer ${accessToken}` } } : {}),
        },
      ],
    },
    output: {
      path: path.resolve(dirname, 'dist'),
      clean: true,
      publicPath: 'auto',
    },
    module: {
      rules: [
        {
          test: /\.tsx?$/,
          exclude: /node_modules/,
          use: {
            loader: 'ts-loader',
            options: {
              transpileOnly: true,
            },
          },
        },
        {
          test: /\.jsx?$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader',
            options: {
              presets: ['@babel/preset-env', '@babel/preset-react', '@babel/preset-typescript'],
            },
          },
        },
        {
          test: /\.css$/,
          use: ['style-loader', 'css-loader'],
        },
      ],
    },
    resolve: {
      extensions: ['.ts', '.tsx', '.js', '.jsx', '.mjs'],
      alias: {
        'process/browser': require.resolve('process/browser.js'),
      },
      fallback: {
        http: false,
        https: false,
        fs: false,
        net: false,
        tls: false,
        process: require.resolve('process/browser.js'),
      },
    },
    plugins: [
      new HtmlWebpackPlugin({
        template: './public/index.html',
        inject: 'body',
      }),
      new webpack.DefinePlugin({
        '__STANDALONE__': JSON.stringify(isStandalone),
      }),
      new CopyPlugin({
        patterns: [
            {
            from: 'public',
            to: '',
            globOptions: {
                ignore: ['**/index.html']
            }
          },
        ],
      }),
      new webpack.ProvidePlugin({
        process: 'process/browser',
      }),
      ...(!isStandalone ? [new webpack.container.ModuleFederationPlugin({
        name: 'restful_config_api',
        filename: 'remote-entry.js',
        exposes: {
          './App': './src/App.tsx',
        },
        shared: {
          react: {
            singleton: true,
            requiredVersion: pkg.dependencies.react,
          },
          'react-dom': {
            singleton: true,
            requiredVersion: pkg.dependencies['react-dom'],
          },
          '@mui/material': {
            singleton: true,
            requiredVersion: pkg.dependencies['@mui/material'],
          },
          '@mui/system': {
            singleton: true,
            requiredVersion: pkg.dependencies['@mui/system'],
          },
          '@milestone/cross-mfe-communication': {
            singleton: true,
            // Provided by the shell; use whatever version it offers instead of the one in package.json
            requiredVersion: false,
          },
          '@milestone/fe-typed-events': {
            singleton: true,
            // Provided by the shell; use whatever version it offers instead of the one in package.json
            requiredVersion: false,
          },
          '@milestone/one-shell-mfe-contract': {
            singleton: true,
            // Provided by the shell; use whatever version it offers instead of the one in package.json
            requiredVersion: false,
          },
          '@milestone/milestone-theme': {
            singleton: true,
            // Provided by the shell; use whatever version it offers instead of the one in package.json
            requiredVersion: false,
          },
          '@milestone/design-tokens': {
            singleton: true,
            // Provided by the shell; use whatever version it offers instead of the one in package.json
            requiredVersion: false,
          },
          '@milestone/ui-assets': {
            singleton: true,
            // Provided by the shell; use whatever version it offers instead of the one in package.json
            requiredVersion: false,
          },
        },
      })] : []),
    ],
  };
};
