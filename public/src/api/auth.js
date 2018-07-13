import axios from "axios";

const base =
  process.env.NODE_ENV === "production"
    ? "/oauth/token"
    : `${process.env.VUE_APP_API_ROOT}/oauth/token`;

export function login(credentials) {
  // credentials have to be an object { username, password }
  return axios.post(base, credentials);
}

export function logout() {
  return axios.delete(base);
}
