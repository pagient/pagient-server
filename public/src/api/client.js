import axios from "axios";

const base =
  process.env.NODE_ENV === "production"
    ? "/api/clients"
    : `${process.env.VUE_APP_API_ROOT}/api/clients`;

export function getAllClients() {
  return axios.get(base);
}
