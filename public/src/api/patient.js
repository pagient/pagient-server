import axios from "axios";

const base =
  process.env.NODE_ENV === "production"
    ? "/api/patients"
    : `${process.env.VUE_APP_API_ROOT} /api/patients`;

export function getAllPatients() {
  return axios.get(base);
}

export function getPatient(id) {
  return axios.get(base + "/" + id);
}

export function addPatient(patient) {
  return axios.post(base, patient);
}

export function updatePatient(patient) {
  return axios.post(base + "/" + patient.id, patient);
}

export function deletePatient(patient) {
  return axios.delete(base + "/" + patient.id);
}
