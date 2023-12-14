<template>
  <div class="map-wrap">
    <div class="map" ref="mapContainer"></div>
    <template style="display: none" v-for="parkingLot in parkingLots" :key="parkingLot.name">
      <MapMarker :map="map" :parking-lot="parkingLot"/>
    </template>
  </div>
</template>

<script>
import { Map } from 'maplibre-gl';
import MapMarker from "@/components/MapMarker.vue";

export default {
  name: "Map",
  components: {MapMarker},
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