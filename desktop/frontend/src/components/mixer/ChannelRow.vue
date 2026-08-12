<script setup lang="ts">

import AppIcon from "./AppIcon.vue"
import VolumeMeter from "./VolumeMeter.vue"
import MuteButton from "./MuteButton.vue"

import type {
  MixerChannel
} from "../../types/mixer"


// ============================================================
// PROPS
// ============================================================

const props =
  defineProps<{
    channel: MixerChannel
  }>()


// ============================================================
// EMIT
// ============================================================

const emit =
  defineEmits<{

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
// VOLUME COLOR
// ============================================================

function getVolumeColor() {

  switch (
  props.channel.app
  ) {

    case "browser":
      return "cyan"

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

}

</script>


<template>

  <div class="
      flex
      w-full
      items-center

      gap-2

      border-b
      border-white/10

      px-2
      py-1.5

      transition-colors

      last:border-b-0

      sm:gap-3
      sm:px-4
      sm:py-2

      md:gap-5
      md:px-6
      md:py-3
    " :class="channel.muted
        ? 'bg-red-500/[0.05]'
        : ''
      ">

    <!-- ==================================================== -->
    <!-- APP INFO -->
    <!-- ==================================================== -->

    <div class="
        flex
        w-[145px]
        shrink-0
        items-center

        gap-2

        sm:w-[190px]
        sm:gap-3

        md:w-[230px]
      ">

      <!-- ICON -->

      <div class="
          flex
          h-9
          w-9
          shrink-0
          items-center
          justify-center

          sm:h-10
          sm:w-10

          md:h-11
          md:w-11
        ">

        <AppIcon :channel="channel" />

      </div>


      <!-- NAME -->

      <div class="
          min-w-0
          flex-1
        ">

        <div class="
            truncate

            text-sm
            font-bold

            sm:text-base

            md:text-xl
          " :class="channel.muted
              ? 'text-red-300'
              : 'text-white'
            ">

          {{ channel.name }}

        </div>


        <!-- STATUS -->

        <div class="
            text-[7px]

            font-bold
            uppercase

            tracking-[0.12em]

            sm:text-[8px]

            md:text-[9px]
          " :class="!channel.connected
              ? 'text-slate-500'
              : channel.muted
                ? 'text-red-400'
                : 'text-emerald-400'
            ">

          {{
            !channel.connected
              ? "Offline"
              : channel.muted
                ? "Muted"
                : "Connected"
          }}

        </div>

      </div>

    </div>


    <!-- ==================================================== -->
    <!-- VOLUME METER -->
    <!-- ==================================================== -->

    <div class="
        min-w-0
        flex-1

        transition-opacity
      " :class="channel.muted
          ? 'opacity-45'
          : 'opacity-100'
        ">

      <VolumeMeter :volume="channel.volume" :color="getVolumeColor()" @set-volume="
        volume =>
          emit(
            'set-volume',
            channel.id,
            volume
          )
      " />

    </div>


    <!-- ==================================================== -->
    <!-- VALUE -->
    <!-- ==================================================== -->

    <div class="
        w-10
        shrink-0

        text-right

        text-sm
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


    <!-- ==================================================== -->
    <!-- MUTE BUTTON -->
    <!-- ==================================================== -->

    <MuteButton :muted="channel.muted" @toggle="
      emit(
        'toggle-mute',
        channel.id
      )
      " />

  </div>

</template>