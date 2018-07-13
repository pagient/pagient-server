import Vuex from "vuex";
import createLogger from "vuex/dist/logger";
import * as getters from "./getters";
import * as actions from "./actions";
import mutations from "./mutations";

const state = {
  isLoggedIn: !!localStorage.getItem("token"),
  authToken: localStorage.getItem("token"),
  currentClientId: null,
  clients: {
    /*
    id: {
      id,
      name
    }
    */
  },
  pagers: {
    /*
    id: {
      id,
      name
    }
    */
  },
  patients: {
    /*
    id: {
      id,
      name,
      ssn,
      clientID,
      pagerID,
      status,
      active
    }
    */
  }
};

const plugins = [];
if (process.env.NODE_ENV !== "production") {
  plugins.push(createLogger());
}

export default new Vuex.Store({
  state,
  getters,
  actions,
  mutations,
  plugins: plugins,
  strict: process.env.NODE_ENV !== "production"
});
