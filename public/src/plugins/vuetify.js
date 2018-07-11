import Vue from "vue";
import Vuetify from "vuetify";
import colors from "vuetify/es5/util/colors";

Vue.use(Vuetify, {
  theme: {
    primary: colors.orange.base, // #FF9800
    secondary: colors.brown.darken1, // #6D4C41
    // accent: colors.orange.darken3, // #EF6C00
    accent: colors.orange.darken4, // #EF6C00
    error: colors.deepOrange.base, // #FF5722
    warning: colors.yellow.base, // #FFEB3B
    info: colors.blue.base, // #2196F3
    success: colors.green.base // #4CAF50
  }
});
