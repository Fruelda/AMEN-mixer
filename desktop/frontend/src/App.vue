<script setup lang="ts">

import {
  onMounted,
  ref
} from "vue"

import MixerBoard from "./components/mixer/MixerBoard.vue"
import DevicesPage from "./components/mixer/DevicesPage.vue"

import {
  mixerStore
} from "./stores/mixer"

import {
  useMixerSocket
} from "./composables/useMixerSocket"

import {
  useSerial
} from "./composables/useSerial"


// ============================================================
// PAGE
// ============================================================

type AppPage =
  | "mixer"
  | "devices"


const activePage =
  ref<AppPage>(
    "mixer"
  )


// ============================================================
// REALTIME
// ============================================================

const mixerSocket =
  useMixerSocket()


// ============================================================
// START APP
// ============================================================

onMounted(
  async () => {

    // ========================================================
    // LOAD CHANNEL
    // ========================================================

    await mixerStore.loadChannels()


    // ========================================================
    // START WEBSOCKET
    // ========================================================

    mixerSocket.start()

  }
)


// ============================================================
// SERIAL
// ============================================================

useSerial(
  command => {

    console.log(
      "[APP] ESP32 COMMAND:",
      command
    )


    mixerStore.handleCommand(
      command
    )

  }
)


// ============================================================
// NAVIGATION
// ============================================================

function openDevices() {

  activePage.value =
    "devices"

}


function backToMixer() {

  activePage.value =
    "mixer"

}

</script>


<template>

  <!-- ====================================================== -->
  <!-- MIXER PAGE -->
  <!-- ====================================================== -->

  <MixerBoard v-if="activePage === 'mixer'" @open-devices="openDevices" />


  <!-- ====================================================== -->
  <!-- DEVICES PAGE -->
  <!-- ====================================================== -->

  <DevicesPage v-else @back="backToMixer" />

</template>