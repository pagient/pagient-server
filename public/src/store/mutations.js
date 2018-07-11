import Vue from "vue";

export default {
  login(state, token) {
    localStorage.setItem("token", token);
    state.isLoggedIn = true;
    state.authToken = token;
  },
  logout(state) {
    localStorage.removeItem("token");
    state.isLoggedIn = false;
    state.authToken = null;
  },
  receiveClients(state, clients) {
    clients.forEach(client => {
      if (!state.clients[client.id]) {
        createClient(state, client.id, client.name);
      }
    });
  },
  receivePagers(state, pagers) {
    pagers.forEach(pager => {
      if (!state.pagers[pager.id]) {
        createPager(state, pager.id, pager.name);
      }
    });
  },
  receivePatients(state, patients) {
    patients.forEach(patient => {
      if (!state.patients[patient.id]) {
        createPatient(state, patient.id, patient);
      }
    });
  },
  receivePatient(state, patient) {
    createPatient(state, patient.id, patient);

    if (patient.active) {
      setCurrentClient(state, patient.clientId);
    }
  },
  switchClient(state, client) {
    setCurrentClient(state, client.id);
  }
};

function createClient(state, id, name) {
  Vue.set(state.clients, id, {
    id,
    name
  });
}

function createPager(state, id, name) {
  Vue.set(state.pagers, id, {
    id,
    name
  });
}

function createPatient(state, id, patient) {
  Vue.set(state.patients, id, patient);
}

function setCurrentClient(state, id) {
  state.currentClientId = id;
}
