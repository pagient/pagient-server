<template>
  <v-jumbotron>
    <v-container fill-height text-xs-center>
      <v-layout align-center>
        <v-flex>
          <h3 class="display-1 mb-4">We are sorry, but we could not find the page</h3>
          <span class="subheading">{{ requestedPath }} not available</span>,
          <v-divider class="my-3"></v-divider>
          <div class="title mb-2">
            Go back to <router-link to="/">Home</router-link>!
          </div>
        </v-flex>
      </v-layout>
    </v-container>
  </v-jumbotron>
</template>

<script>
export default {
  computed: {
    requestedPath() {
      return this.$route.path;
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="stylus">
</style>
