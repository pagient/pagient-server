import axios from "axios";

let base = process.env.VUE_APP_API_ROOT + "api/clients";

export function getAllClients() {
  return axios.get(base);
}
