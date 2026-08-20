<script setup lang="ts">
import { computed } from "vue"

import ChannelInfo from "./ChannelInfo.vue"
import VolumeMeter from "./VolumeMeter.vue"
import MuteButton from "./MuteButton.vue"

import type {
  MixerChannel
} from "../../types/mixer"


type MeterColor =
  | "cyan"
  | "green"
  | "purple"
  | "yellow"
  | "white"


const props = defineProps<{
  channel: MixerChannel
}>()


const emit = defineEmits<{
  (
    event: "set-volume",
    id: number,
    volume: number
  ): void

  (
    event: "toggle-mute",
    id: number
  ): void
}>()


// ============================================================
// METER COLOR
// ============================================================

const meterColor = computed<MeterColor>(() => {
  switch (props.channel.app) {
    case "spotify":
      return "green"

    case "discord":
      return "purple"

    case "valeton":
      return "yellow"

    case "master":
      return "white"

    default:
      return "cyan"
  }
})


// ============================================================
// EVENTS
// ============================================================

function setVolume(volume: number) {
  emit(
    "set-volume",
    props.channel.id,
    volume
  )
}


function toggleMute() {
  emit(
    "toggle-mute",
    props.channel.id
  )
}
</script>


<template>
  <div class="
      flex
      w-full
      min-w-0
      items-center
      gap-2
      overflow-hidden
      border-b
      border-white/10
      px-2
      py-1
      transition-colors
      last:border-b-0

      sm:gap-3
      sm:px-4
      sm:py-2

      md:gap-5
      md:px-6
      md:py-3
    " :class="channel.muted ? 'bg-red-500/[0.05]' : ''">

    <!-- CHANNEL INFO -->
    <div class="shrink-0">
      <ChannelInfo :channel="channel" />
    </div>


    <!-- VOLUME -->
    <div class="
        min-w-0
        flex-1
        overflow-hidden
      " :class="channel.muted ? 'opacity-45' : 'opacity-100'">
      <VolumeMeter :volume="channel.volume" :color="meterColor" :muted="channel.muted" @set-volume="setVolume" />
    </div>


    <!-- VALUE -->
    <div class="
        w-8
        shrink-0
        text-right
        text-xs
        font-bold
        tabular-nums

        sm:w-12
        sm:text-base

        md:w-16
        md:text-xl
      " :class="channel.muted
        ? 'text-red-300'
        : 'text-white'
        ">
      {{ Math.round(channel.volume) }}%
    </div>


    <!-- MUTE -->
    <MuteButton :muted="channel.muted" @toggle="toggleMute" />

  </div>
</template>