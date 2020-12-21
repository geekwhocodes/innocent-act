module.exports = {
  publicPath: '/',
  outputDir: 'dist',

  // This is to make all static file requests generated by Vue to go to
  // /frontend/*. However, this also ends up creating a `dist/frontend`
  // directory and moves all the static files in it. The physical directory
  // and the URI for assets are tightly coupled. This is handled in the Go app
  // by using stuffbin aliases.
  assetsDir: 'web',

  // Move the index.html file from dist/index.html to dist/frontend/index.html
  indexPath: './web/index.html',

  productionSourceMap: false,
  filenameHashing: true,

  devServer: {
    port: process.env.LOWKEY_FRONTEND_PORT || 8080,
    proxy: {
      '^/api': {
        target: process.env.LOWKEY_API_URL || 'http://127.0.0.1:8083'
      },
    },
  },
};
