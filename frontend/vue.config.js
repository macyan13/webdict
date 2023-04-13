module.exports = {
  outputDir: 'target/dist',
  assetsDir: 'static',

  // proxy all webpack dev-server requests starting with /api
  // to our backend (localhost:4000) using http-proxy-middleware
  // see https://cli.vuejs.org/config/#devserver-proxy
  devServer: {
    proxy: {
      '/v1/api': {
        target: 'http://localhost:4000',
        ws: true,
        changeOrigin: true
      }
    }
  },
}