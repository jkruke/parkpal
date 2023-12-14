<template>
  <div id="app">
    <div class="row g-0">

      <main class="col-12 bg-secondary">

        <nav class="navbar navbar-expand-lg navbar-light bg-light">
          <h1 class="h4 py-2 px-2 text-center text-primary">
            <i class="fa-solid fa-square-parking"></i>
            <span class="d-none d-lg-inline">ParkPal</span>
          </h1>
          <div class="flex-fill"></div>
          <div class="px-2">
            <div class="input-group">
              <input type="text" class="form-control" placeholder="Enter license plate number ">
              <button class="btn btn-secondary">Search</button>
            </div>
          </div>
        </nav>

        <div class="container-fluid mt-1 p-1">
          <div class="row flex-column flex-lg-row">
            <h2 class="h6 text-white-50">QUICK STATS</h2>
            <div class="col">
              <div class="card mb-3">
                <div class="card-body">
                  <h3 class="card-title h2">{{ totalParkingSlots }}</h3>
                  <span class="text-success">
                      <i class="fas fa-chart-line"></i>
                      Total Parking Slots
                    </span>
                </div>
              </div>
            </div>
            <div class="col">
              <div class="card mb-3">
                <div class="card-body">
                  <h3 class="card-title h2">{{ totalBikeCount }}</h3>
                  <span class="text-success">
                      <i class="fa-solid fa-person-biking"></i>
                      Bikes
                    </span>
                </div>
              </div>
            </div>
            <div class="col">
              <div class="card mb-3">
                <div class="card-body">
                  <h3 class="card-title h2">{{ freeParkingSlots }}</h3>
                  <span class="text-success">
                      <i class="fas fa-chart-line"></i>
                      Free Parking Slots
                    </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>

    <Map :parking-lots="parkingLots"/>

    <footer class="text-center py-2 text-muted">
      &copy; IOTeam
    </footer>
  </div>
</template>

<script>
import Map from "@/components/Map.vue";
import axios from "axios";

const dummyData = [
  {
    "id": 1,
    "name": "D8",
    "bikeCount": 100,
    "congestionRate": 80,
    "totalSpace": 1000,
    "latlong": [
      "105.84414",
      "21.00405"
    ]
  },
  {
    "id": 2,
    "name": "D9",
    "bikeCount": 323,
    "congestionRate": 60,
    "totalSpace": 500,
    "latlong": [
      "105.84414",
      "21.00605"
    ]
  }
];

export default {
  name: 'App',
  components: {
    Map
  },
  data() {
    return {
      parkingLots: []
    }
  },
  mounted() {
    setInterval(this.updateParkingLots, 2000)
  },
  computed: {
    freeParkingSlots() {
      return this.totalParkingSlots - this.totalBikeCount
    },
    totalBikeCount() {
      return this.parkingLots.reduce((a, b) => a + b.bikeCount, 0)
    },
    totalParkingSlots() {
      return this.parkingLots.reduce((a, b) => a + b.totalSpace, 0)
    }
  },
  methods: {
    updateParkingLots() {
      axios.get("http://localhost:9091/parking-lots")
          .then(response => {
            this.parkingLots = response.data
          })
          .catch(() => {
            console.info("Using dummy data...")
            this.parkingLots = dummyData
          })
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--dark-text);
  max-height: 100vh;
  overflow: hidden;
}

footer {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  background: white;
}
</style>
