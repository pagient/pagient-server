import Vue from "vue";
import VueRouter from "vue-router";
import store from "@/store";

import HomeComponent from "@/components/Home";
import LoginComponent from "@/components/Login";
import LogoutComponent from "@/components/Logout";
import NotFoundComponent from "@/components/NotFound";

Vue.use(VueRouter);

const routes = [
  { path: "/", component: HomeComponent, meta: { requiresAuth: true } },
  { path: "/login", component: LoginComponent, meta: { guestOnly: true } },
  { path: "/logout", component: LogoutComponent, meta: { requiresAuth: true } },
  { path: "*", component: NotFoundComponent }
];

const router = new VueRouter({
  routes // short for `routes: routes`
});

// requiresAuth checker
router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // this route requires auth, check if logged in
    // if not, redirect to login page
    if (!store.getters.isLoggedIn) {
      next({
        path: "/login",
        query: { redirect: to.fullPath !== "/logout" ? to.fullPath : "/" }
      });
    }
  }
  next();
});

// guestOnly checker
router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.guestOnly)) {
    // this route requires to be not authenticated, check if logged in
    // if logged in, redirect to logout page.
    if (store.getters.isLoggedIn) {
      next({
        path: "/logout",
        query: { redirect: to.fullPath }
      });
    }
  }
  next();
});

export default router;
