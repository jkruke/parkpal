<template>
  <table ref="mapMarker" class="marker">
    <tr>
      <td colspan="2" class="headline"><strong>{{ parkingLot.name }}</strong></td>
    </tr>
    <tr>
      <td>Occupied</td>
      <td>{{ occupied }}%</td>
    </tr>
    <tr>
      <td>Trend</td>
      <td>{{ trend }}</td>
    </tr>
  </table>
</template>

<script>
import {Map, Marker, Popup} from "maplibre-gl";

export default {
  name: "MapMarker",
  props: {
    map: Map,
    parkingLot: Object
  },
  mounted() {
    const parkingLot = this.parkingLot
    const latlong = [parseFloat(parkingLot.latitude), parseFloat(parkingLot.longitude)];
    new Marker({color: this.getColor()})
        .setLngLat(latlong)
        .setPopup(new Popup({offset: 25, closeButton: false})
            .setHTML(this.$refs.mapMarker.outerHTML))
        .addTo(this.map)
  },
  computed: {
    trend() {
      return this.parkingLot.congestionRate > 0 ? "+" + this.parkingLot.congestionRate : this.parkingLot.congestionRate
    },
    occupied() {
      return Math.round(100 * this.parkingLot.bikeCount / this.parkingLot.totalSpace)
    }
  },
  methods: {
    getColor() {
      const value = this.occupied / 100
      const hue = ((1 - value) * 120).toString(10);
      return ["hsl(", hue, ",100%,50%)"].join("");
    }
  }
}

</script>

<style scoped>
.marker {
  font-size: 1.5em;
  line-height: 1.25em;
}
.headline {
  text-align: center;
}
td {
  padding: 0.25em;
}
</style>