export const isLoggedIn = state => state.isLoggedIn;
export const authToken = state => state.authToken;

export const clients = (state, getters) => {
  const clientIds = Object.keys(state.clients);
  // map to array and load patients
  if (clientIds.length > 0) {
    return clientIds.map(id => {
      const client = state.clients[id];
      client.patient = getters.patientByClient(client);
      return client;
    });
  }
  return [];
};

export const pagers = (state, getters) => {
  const pagerIds = Object.keys(state.pagers);
  // map to array and load patients
  if (pagerIds.length > 0) {
    return pagerIds.map(id => {
      const pager = state.pagers[id];
      pager.patient = getters.patientByPager(pager);
      return pager;
    });
  }
  return [];
};

export const patients = state => {
  const patientIds = Object.keys(state.patients);
  // map to array
  if (patientIds.length > 0) {
    return patientIds.map(id => state.patients[id]);
  }
  return [];
};

export const activeClient = state => {
  return state.activeClientId ? state.clients[state.activeClientId] : {};
};

export const activePatient = (state, getters) => {
  const client = activeClient(state);
  return getters.patientByClient(client);
};

export const patientByClient = (_, getters) => client => {
  return getters.patients.find(
    patient => patient.clientId === client.id && patient.active
  );
};

export const patientByPager = (_, getters) => pager => {
  return getters.patients.find(patient => patient.pagerId === pager.id);
};
