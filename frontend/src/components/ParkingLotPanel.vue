<template>
  <div>
    <h1>Parking Lot ID: {{ id}} ({{ name }})</h1>
    <table>
      <tr>
        <td>Total count of motorbikes at parking lot</td>
        <td><strong>{{ totalCount }}</strong></td>
      </tr>
      <tr>
        <td>Current load of the gate</td>
        <td><strong>{{ currentLoad }}</strong></td>
      </tr>
    </table>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "ParkingLotPanel",
  props: {
    id: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      totalCount: 0,
      currentLoad: 0,
      name: ""
    }
  },
  mounted() {
    this.doExampleRequest()
  },
  methods: {
    doExampleRequest() {
      // This is an example how to request a REST API and use the response to change values of variables which will be
      // reactively changed in the DOM:
      axios.get("http://jsonplaceholder.typicode.com/posts")
          .then(response => {
            const data = response.data
            this.totalCount = data.length
            this.currentLoad = data[0].userId
            this.name = "Parking Lot Name"
          })
    },
    updateData() {
      // This is an example how to request a REST API and use the response to change values of variables which will be
      // reactively changed in the DOM:
      const host = "http://192.168.1.196:9091"
      axios.get(host + "/parkingLots/" + this.id)
          .then(response => {
            const data = response.data
            this.totalCount = data.bikeCount
            this.currentLoad = data.congestionRate
            this.name = data.name
          })
    }
  }
}
</script>

<style scoped>

</style>