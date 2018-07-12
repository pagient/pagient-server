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
    const unsubscribe = this.$store.subscribe((mutation, state) => {
      if (mutation.type === "login" || state.isLoggedIn) {
        connectWebsocket(state.authToken);
        unsubscribe();
      }
    });

    const connectWebsocket = token => {
      const socket = new WebSocket(
        process.env.VUE_APP_WEBSOCKET_ROOT + `?jwt=${token}`
      );

      createWebSocketPlugin(socket)(this.$store);
    };
  }
};
</script>

<style lang="stylus">
</style>
