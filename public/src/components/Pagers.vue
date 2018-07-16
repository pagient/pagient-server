<template>
  <v-layout row wrap>
    <v-flex v-for="pager in pagers" :key="pager.id" class="custom-flex">
      <v-card @click.native="assignPager({ patient: currentPatient, pager: pager })" class="flex-card flex-column" height="100%" :color="isPagerOverdue(pager) ? 'error' : undefined" :dark="isPagerOverdue(pager)" hover ripple>
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
          <v-btn @click.stop="assignPager({ patient: pager.patient, pager: null })" :disabled="!pager.patient" icon large>
            <v-icon medium>check</v-icon>
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-flex>
  </v-layout>
</template>

<script>
import { mapActions, mapGetters } from "vuex";

export default {
  methods: {
    ...mapActions(["assignPager", "callPatient"]),
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
