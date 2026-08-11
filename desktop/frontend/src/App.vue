<script setup lang="ts">

import {
  onMounted
} from "vue"

import MixerBoard from "./components/mixer/MixerBoard.vue"

import {
  mixerStore
} from "./stores/mixer"

import {
  useMixerSocket
} from "./composables/useMixerSocket"

import {
  useSerial
} from "./composables/useSerial"


const mixerSocket =
  useMixerSocket()



onMounted(
  async () => {


    await mixerStore.loadChannels()


    mixerSocket.start()


  })



useSerial(
  command => {


    console.log(
      "[APP] ESP32 COMMAND:",
      command
    )


    mixerStore.handleCommand(
      command
    )


  })

</script>


<template>

  <MixerBoard />

</template>