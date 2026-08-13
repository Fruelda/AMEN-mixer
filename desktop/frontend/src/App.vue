<script setup lang="ts">
import { onMounted, ref } from "vue"

import MixerBoard from "./components/mixer/MixerBoard.vue"
import DevicesPage from "./components/mixer/DevicePage.vue"

import { mixerStore } from "./stores/mixer"
import { useMixerSocket } from "./composables/useMixerSocket"
import { useSerial } from "./composables/useSerial"


type AppPage = "mixer" | "devices"

const activePage = ref<AppPage>("mixer")
const mixerSocket = useMixerSocket()


onMounted(async () => {
  await mixerStore.loadChannels()
  mixerSocket.start()
})


useSerial(command => {
  console.log("[APP] ESP32 COMMAND:", command)
  mixerStore.handleCommand(command)
})


function openDevices() {
  activePage.value = "devices"
}


function backToMixer() {
  activePage.value = "mixer"
}
</script>


<template>
  <MixerBoard v-if="activePage === 'mixer'" @open-devices="openDevices" />

  <DevicesPage v-else @back="backToMixer" />
</template>