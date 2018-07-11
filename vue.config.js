var path = require("path");

module.exports = {
  // where to output the built files
  outputDir: "public/dist",

  chainWebpack: config => {
    config
      .entry("app")
      .clear()
      .add("./public/src/main.js");

    config.resolve.alias.set("@", path.resolve(__dirname, "./public/src"));

    config.plugin("html").tap(args => {
      args[0].template = path.resolve(__dirname, "./public/public/index.html");
      return args;
    });

    config.plugin("copy").tap(args => {
      args[0][0].from = path.resolve(__dirname, "./public/public");
      return args;
    });
  }
};
