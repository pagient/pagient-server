import axios from "axios";

let base = process.env.VUE_APP_API_ROOT + "api/pagers";

export function getAllPagers() {
  return axios.get(base);
}
