module.exports = {
  presets: ["@vue/cli-plugin-babel/preset"],
  env: {
    test: {
      plugins: ["transform-es2015-modules-commonjs"]
    }
  }
};
