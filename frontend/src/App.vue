<template>
  <div>
    <div class="container-fluid">
      <div class="row g-0">

        <nav class="col-2 bg-light pe-3 border-right">

          <h1 class="h4 py-4 text-center text-primary">
            <i class="fa-solid fa-square-parking"></i>
            <span class="d-none d-lg-inline">ParkPal</span>
          </h1>

          <div class="list-group text-center text-lg-left">
            <span class="list-group-item disabled d-none d-lg-block">
              <small>CONTROLS</small>
            </span>
            <a href="#" class="list-group-item list-group-item-action active">
              <i class="fas fa-home"></i>
              <span class="d-none d-lg-inline">Dashboard</span>
            </a>

            <a href="#" class="list-group-item list-group-item-action">
              <i class="fas fa-flag"></i>
              <span class="d-none d-lg-inline">Reports</span>
            </a>
          </div>

        </nav>

        <main class="col-10 bg-secondary">

          <nav class="navbar navbar-expand-lg navbar-light bg-light">
            <div class="flex-fill"></div>
            <div class="navbar nav">

              <li class="nav-item">
                <div class="input-group mb-3">
                  <input type="text" class="form-control" placeholder="Enter license plate number ">
                  <button class="btn btn-secondary">Search</button>
                </div>
              </li>

              <li class="nav-item">
                <a href="#" class="nav-link">
                  <i class="fas fa-cog"></i>
                </a>
              </li>
            </div>
          </nav>

          <div class="container-fluid mt-3 p-4">
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

          <div id="app"></div>

        </main>

      </div>

      <footer class="text-center py-2 text-muted">
        &copy; IOTeam
      </footer>
    </div>
    <Map/>
  </div>
</template>

<script>
import Map from "@/components/Map.vue";
import axios from "axios";

export default {
  name: 'App',
  components: {
    Map
  },
  data() {
    return {
      totalBikeCount: 0,
      totalParkingSlots: 0
    }
  },
  mounted() {
    this.requestAPI()
    setInterval(this.requestAPI, 2000)
  },
  computed: {
    freeParkingSlots() {
      return this.totalParkingSlots - this.totalBikeCount
    }
  },
  methods: {
    requestAPI() {
      axios.get("http://172.20.10.12:9091/parking-lots")
          .then(response => {
            const data = response.data
            console.log(data)
            this.totalBikeCount = data.reduce((a, b) => a + b.bikeCount, 0)
            this.totalParkingSlots = data.reduce((a, b) => a + b.totalSpace, 0)
          })
          .catch(() => {
            this.totalBikeCount = -1
            this.totalParkingSlots = -1
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
}
</style>
