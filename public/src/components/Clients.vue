<template>
  <v-layout row>
    <v-flex v-for="client in clients" :key="client.id">
      <v-card @click.native="selectClient(client)" :color="activeClient.id === client.id ? 'primary' : undefined" :dark="activeClient.id === client.id" height="100%" hover ripple>
        <v-card-title>
          <div>
            <h3 class="headline">{{ client.name }}</h3>
            <span v-if="client.patient">{{ client.patient.name }} <strong>{{ client.patient.ssn }}</strong></span>
          </div>
        </v-card-title>
      </v-card>
    </v-flex>
  </v-layout>
</template>

<script>
import { mapActions, mapGetters } from "vuex";

export default {
  methods: mapActions(["selectClient"]),
  computed: mapGetters(["clients", "activeClient"])
};
</script>
