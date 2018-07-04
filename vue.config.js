path = require('path')

module.exports = {
  // where to output the built files
  outputDir: 'public/dist',

  chainWebpack: config => {
    config.entry('app')
      .clear()
      .add('./public/src/main.js')

    config.resolve.alias
      .set('@', path.join(__dirname, './public/src'))

    config.plugin('copy')
      .tap(args => {
        args[0][0].from = path.join(__dirname, './public/public')
        return args
      })
  }
}
