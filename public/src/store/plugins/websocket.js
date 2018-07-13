export default function createWebSocketPlugin(socket) {
  return store => {
    socket.addEventListener("open", () => {
      store
        .dispatch("getAllClients")
        .then(() => store.dispatch("getAllPagers"))
        .then(() => store.dispatch("getAllPatients"))
        .then(() => {
          socket.addEventListener("message", ({ data }) => {
            const message = JSON.parse(data);
            switch (message.type) {
              case "patient_add":
              case "patient_update":
                store.commit("receivePatient", message.data);
                break;
              case "patient_delete":
                store.commit("deletePatient", message.data);
                break;
            }
          });
        });
    });
  };
}
