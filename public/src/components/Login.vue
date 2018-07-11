<template>
  <v-container fluid fill-height>
    <v-layout align-center justify-center>
      <v-flex xs4 class="text-xs-center">
        <v-avatar tile size="100">
          <img src="/logo.png" alt="Pagient Brand">
        </v-avatar>
        <v-card class="mt-2 text-xs-left">
          <v-card-title>
            <h3 class="headline mb-0">Login</h3>
          </v-card-title>
          <v-card-text>
            <v-form v-model="valid" @submit.prevent="login">
              <v-alert class="mb-5" :value="requestError" color="error" icon="warning" outline>
                {{ requestError }}
              </v-alert>
              <v-text-field
                v-model="credentials.username"
                :error-messages="usernameErrors"
                label="Username"
                required
                @input="$v.credentials.username.$touch()"
                @blur="$v.credentials.username.$touch()"
              ></v-text-field>
              <v-text-field
                v-model="credentials.password"
                :append-icon="showPassword ? 'visibility_off' : 'visibility'"
                :error-messages="passwordErrors"
                :type="showPassword ? 'test' : 'password'"
                label="Password"
                required
                @click:append="showPassword = !showPassword"
                @input="$v.credentials.password.$touch()"
                @blur="$v.credentials.password.$touch()"
              ></v-text-field>
              <v-btn type="submit" class="primary" :disabled="!valid">Login</v-btn>
            </v-form>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { validationMixin } from "vuelidate";
import { required } from "vuelidate/lib/validators";
import createWebSocketPlugin from "@/store/plugins/websocket";

export default {
  mixins: [validationMixin],
  validations: {
    credentials: {
      username: { required },
      password: { required }
    }
  },
  data: () => {
    return {
      valid: true,
      requestError: "",
      credentials: {
        username: "",
        password: ""
      },
      showPassword: false
    };
  },
  computed: {
    usernameErrors() {
      const errors = [];
      if (!this.$v.credentials.username.$dirty) return errors;
      !this.$v.credentials.username.required &&
        errors.push("Username is required");
      return errors;
    },
    passwordErrors() {
      const errors = [];
      if (!this.$v.credentials.password.$dirty) return errors;
      !this.$v.credentials.password.required &&
        errors.push("Password is required");
      return errors;
    }
  },
  methods: {
    login() {
      this.$v.$touch();

      // reset request error
      this.requestError = "";

      // dispatch login action
      this.$store
        .dispatch("login", this.credentials)
        .then(() => {
          // success
          createWebSocketPlugin()(this.$store);

          const redirect = this.$route.query.redirect
            ? this.$route.query.redirect
            : "/";
          this.$router.push(redirect);
        })
        .catch(error => {
          // failure
          if (!error.response) {
            this.requestError = error.message;
            return;
          }

          if (error.response.status === 401) {
            this.requestError = "Username or password wrong!";
            return;
          }
          this.requestError = error.response.statusText;
        });
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="stylus">
</style>
