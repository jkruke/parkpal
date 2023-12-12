<template>
  <div>
    <Map/>
  </div>
</template>

<script>
import Map from "@/components/Map.vue";

export default {
  name: 'App',
  components: {
    Map
}
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--dark-text);
}
</style>
