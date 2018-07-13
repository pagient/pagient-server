<template>
  <v-app>
    <v-content>
      <router-view></router-view>
    </v-content>
  </v-app>
</template>

<script>
import createWebSocketPlugin from "@/store/plugins/websocket";

export default {
  name: "app",
  mounted() {
    const subscribeStore = () => {
      const unsubscribe = this.$store.subscribe((mutation, state) => {
        if (mutation.type === "login" || state.isLoggedIn) {
          connectWebsocket(state.authToken);
          unsubscribe();
        }
      });
    };

    const connectWebsocket = token => {
      const websocketUrl =
        process.env.NODE_ENV === "production"
          ? `ws://${location.host}/ws`
          : process.env.VUE_APP_WEBSOCKET_ROOT;
      const socket = new WebSocket(`${websocketUrl}?jwt=${token}`);

      socket.onclose = evt => {
        if (evt.code === 1006) {
          subscribeStore();
          this.$store.commit("logout");

          if (this.$router.currentRoute.path !== "/login") {
            this.$router.push({
              path: "/login",
              query: { redirect: this.$router.currentRoute.fullPath }
            });
          }
          return;
        }
      };

      createWebSocketPlugin(socket)(this.$store);
    };

    if (this.$store.getters.isLoggedIn) {
      connectWebsocket(this.$store.getters.authToken);
      return;
    }

    subscribeStore();
  }
};
</script>

<style lang="stylus">
</style>
