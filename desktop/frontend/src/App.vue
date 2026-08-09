<script setup lang="ts">

import {
  onMounted,
} from "vue"

import MixerBoard from "./components/mixer/MixerBoard.vue"

import {
  mixerStore,
} from "./stores/mixer"

import {
  useSerial,
} from "./composables/useSerial"

import {
  useRealtime,
} from "./composables/useRealtime"


/*
|--------------------------------------------------------------------------
| INITIAL LOAD
|--------------------------------------------------------------------------
*/

onMounted(
  async () => {

    await mixerStore.loadChannels()

  }
)


/*
|--------------------------------------------------------------------------
| ESP32 SERIAL
|--------------------------------------------------------------------------
*/

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


/*
|--------------------------------------------------------------------------
| WEBSOCKET
|--------------------------------------------------------------------------
*/

useRealtime(
  message => {

    console.log(
      "[APP] REALTIME:",
      message
    )


    mixerStore.applyRemoteUpdate(
      message
    )

  }
)

</script>


<template>

  <MixerBoard />

</template>