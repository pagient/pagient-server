<template>
  <v-layout row wrap>
    <v-flex v-for="pager in pagers" :key="pager.id" class="custom-flex">
      <v-card @click.native="assignPager(currentPatient, pager)" class="flex-card flex-column" height="100%" :color="isPagerOverdue(pager) ? 'error' : undefined" :dark="isPagerOverdue(pager)" hover ripple>
        <v-card-title>
          <div>
            <h3 class="title font-weight-light">{{ pager.name }}</h3>
          </div>
        </v-card-title>

        <v-card-text class="grow py-0">
          <div v-if="pager.patient">
              {{ pager.patient.name }}<br>
              <strong>{{ pager.patient.ssn }}</strong>
          </div>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click.stop="callPatient(pager.patient)" :disabled="!pager.patient" :color="isPagerCalled(pager) ? 'primary' : undefined" :dark="isPagerCalled(pager)" icon large>
            <v-icon medium>vibration</v-icon>
          </v-btn>
          <v-btn @click.stop="assignPager(pager.patient, null)" :disabled="!pager.patient" icon large>
            <v-icon medium>check</v-icon>
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-flex>
    <v-snackbar
      v-model="pagerNotAssignableError"
      color="error"
      timeout.int="3000"
    >
      Pager has already been assigned!
      <v-btn
        dark
        flat
        @click="pagerNotAssignableError = false"
      >
        Close
      </v-btn>
    </v-snackbar>
  </v-layout>
</template>

<script>
import { mapActions, mapGetters } from "vuex";

export default {
  data: () => ({
    pagerNotAssignableError: false
  }),
  methods: {
    ...mapActions(["callPatient"]),
    assignPager(patient, pager) {
      // prevent overwriting of a patient if pager has already been assigned
      if (pager && pager.patient) {
        this.pagerNotAssignableError = true;
        return;
      }
      this.$store.dispatch("assignPager", { patient: patient, pager: pager });
    },
    isPagerCalled(pager) {
      return pager.patient && pager.patient.status === "called";
    },
    isPagerOverdue(pager) {
      return pager.patient && pager.patient.status === "finished";
    }
  },
  computed: mapGetters(["pagers", "currentPatient"])
};
</script>

<style scoped lang="stylus">
.flex-card
  display: flex

.custom-flex
  flex-basis: 25%

.flex-column
  flex-direction: column

.text-nobreak
  white-space: nowrap
</style>
