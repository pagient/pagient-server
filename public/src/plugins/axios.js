import axios from "axios";
import store from "@/store";
import router from "@/plugins/vue-router";

axios.interceptors.request.use(
  config => {
    const token = store.getters.authToken;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => Promise.reject(error)
);

axios.interceptors.response.use(
  response => response,
  error => {
    if (
      error.response &&
      error.response.status === 401 &&
      router.currentRoute.path !== "/login"
    ) {
      store.commit("logout");

      router.push({
        path: "/login",
        query: { redirect: router.currentRoute.fullPath }
      });
    }
    return Promise.reject(error);
  }
);
