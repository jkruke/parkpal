<template>
  <div id="header">
    <span class="left">IoT - ParkPal</span>
  </div>
  <div id="content">
    <ParkingLotPanel :id="id" v-for="id in parkingLots" :key="id" />
    <div>Parking recommendation: ...</div>
  </div>
</template>

<script>
import ParkingLotPanel from "@/components/ParkingLotPanel.vue";

export default {
  name: 'LandingPage',
  components: {ParkingLotPanel},
  data() {
    return {
      parkingLots: [0] // IDs of all parking lots for which we want to display a panel
    }
  }
}
</script>

<style scoped>
#header {
  background: var(--blue);
  padding: 1em;
  color: var(--white);
  font-weight: bold;
  text-align: left;
}
</style>
