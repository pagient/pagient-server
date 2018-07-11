<template>
  <v-container fluid grid-list-xl>
    <da-clients></da-clients>
    <v-divider class="my-3"></v-divider>
    <da-pagers></da-pagers>
  </v-container>
</template>

<script>
import ClientsComponent from "./Clients.vue";
import PagersComponent from "./Pagers.vue";

export default {
  components: {
    "da-clients": ClientsComponent,
    "da-pagers": PagersComponent
  },
  mounted() {
    this.$store.dispatch("getAllPatients");
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="stylus">
.active
  box-shadow: 0 3px 3px -2px #ffce3b, 0 3px 4px 0 #ffce3b, 0 1px 8px 0 #ffce3b !important
</style>
