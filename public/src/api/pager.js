import axios from "axios";

const base =
  process.env.NODE_ENV === "production"
    ? "/api/pagers"
    : `${process.env.VUE_APP_API_ROOT}/api/pagers`;

export function getAllPagers() {
  return axios.get(base);
}
