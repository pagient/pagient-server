import clone from "lodash.clone";
import * as api from "@/api";

export const login = ({ commit }, credentials) => {
  return api.login(credentials).then(response => {
    commit("login", response.data.token);
  });
};

export const logout = ({ commit }) => {
  return api.logout().then(() => {
    commit("logout");
  });
};

export const getAllClients = ({ commit }) => {
  return api.getAllClients().then(response => {
    commit("receiveClients", response.data);
  });
};

export const getAllPagers = ({ commit }) => {
  return api.getAllPagers().then(response => {
    commit("receivePagers", response.data);
  });
};

export const getAllPatients = ({ commit }) => {
  return api.getAllPatients().then(response => {
    commit("receivePatients", response.data);
  });
};

export const callPatient = (_, patient) => {
  if (!patient) return;
  // Copy patient to prevent direct state mutation
  patient = clone(patient);
  patient.status = "call";
  return api.updatePatient(patient);
};

export const assignPager = (_, { patient, pager }) => {
  if (!patient) return;
  // Copy patient to prevent direct state mutation
  patient = clone(patient);
  patient.pagerId = pager ? pager.id : null;
  return api.updatePatient(patient).then(() => {
    if (!patient.active && !patient.pagerId) {
      return api.deletePatient(patient);
    }
  });
};

export const selectClient = ({ commit }, client) => {
  commit("switchClient", client);
};
