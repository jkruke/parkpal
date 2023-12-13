<template>
  <div class="map-wrap">
    <div class="map" ref="mapContainer"></div>
  </div>
</template>

<script>
import { Map, Marker, Popup } from 'maplibre-gl';

export default {
  name: "Map",
  props: {
    parkingLots: [Object]
  },
  data() {
    return {
      apiKey: 'X6Psx4JFJdoaMCgMD6q8',
      initialState: { lng: 105.84414, lat: 21.00444, zoom: 17 },
      map: null
    }
  },
  mounted() {
    const mapContainer = this.$refs.mapContainer
    this.map = new Map({
      container: mapContainer,
      style: `https://api.maptiler.com/maps/streets-v2/style.json?key=${this.apiKey}`,
      center: [this.initialState.lng, this.initialState.lat],
      zoom: this.initialState.zoom
    });
  },
  watch: {
    parkingLots() {
      for (let parkingLot of this.parkingLots) {
        console.log("ParkingLot", parkingLot.name)
        const latlong = [parseFloat(parkingLot.latitude), parseFloat(parkingLot.longitude)];
        const occupied = Math.round(100 * parkingLot.bikeCount / parkingLot.totalSpace)
        const trend = parkingLot.congestionRate > 0 ? "+" + parkingLot.congestionRate : parkingLot.congestionRate
        new Marker({color: this.getColor(occupied / 100)})
            .setLngLat(latlong)
            .setPopup(new Popup({offset: 25})
                .setHTML(`<strong>${parkingLot.name}</strong><p>Occupied: ${occupied}%</p><p>Trend: ${trend}</p>`))
            .addTo(this.map)
      }
    }
  },
  methods: {
    getColor(value) {
      const hue = ((1 - value) * 120).toString(10);
      return ["hsl(", hue, ",100%,50%)"].join("");
    }
  },
  unmounted() {
    this.map?.remove();
  }
};
</script>


<style scoped>
@import '~maplibre-gl/dist/maplibre-gl.css';

.map-wrap {
  position: relative;
  width: 100%;
  height: calc(100vh - 77px); /* calculate height of the screen minus the heading */
}

.map {
  position: absolute;
  width: 100%;
  height: 100%;
}

.watermark {
  position: absolute;
  left: 10px;
  bottom: 10px;
  z-index: 999;
}
</style>