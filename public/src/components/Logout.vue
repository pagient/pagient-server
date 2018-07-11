<template>
  <v-jumbotron>
    <v-container fill-height text-xs-center>
      <v-layout align-center>
        <v-flex>
          <h3 class="display-1 mb-4">It seems that you tried to enter a page that requires you to be logged out!</h3>
          <span class="subheading">Either logout or go back to <router-link to="/">Home</router-link> if you didn't intend to do this.</span>
          <v-divider class="my-3"></v-divider>
          <div class="title mb-2">
            <v-btn @click="logout" class="primary">Logout</v-btn>
          </div>
        </v-flex>
      </v-layout>
    </v-container>
  </v-jumbotron>
</template>

<script>
export default {
  methods: {
    logout() {
      this.$store.dispatch("logout").then(() => {
        const redirect = this.$route.query.redirect
          ? this.$route.query.redirect
          : "/login";
        this.$router.push(redirect);
      });
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="stylus">
</style>
