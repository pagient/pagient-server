import "@babel/polyfill";
import Vue from "vue";
import "@/plugins/vuelidate";
import "@/plugins/vuetify";
import "@/plugins/vuex";
import "@/plugins/axios";
import router from "@/plugins/vue-router";
import store from "@/store";

import App from "@/App.vue";

import "./assets/stylus/main.styl";

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
