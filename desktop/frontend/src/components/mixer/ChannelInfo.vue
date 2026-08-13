<script setup lang="ts">
import { computed } from "vue"

import AppIcon from "./AppIcon.vue"

import type {
    MixerChannel
} from "../../types/mixer"


const props = defineProps<{
    channel: MixerChannel
}>()


// ============================================================
// STATUS
// ============================================================

const statusText = computed(() => {
    if (!props.channel.connected) {
        return "Offline"
    }

    if (props.channel.muted) {
        return "Muted"
    }

    return "Connected"
})


const statusClass = computed(() => {
    if (!props.channel.connected) {
        return "text-slate-500"
    }

    if (props.channel.muted) {
        return "text-red-400"
    }

    return "text-emerald-400"
})


const nameClass = computed(() => {
    return props.channel.muted
        ? "text-red-300"
        : "text-white"
})
</script>


<template>
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


        <!-- INFO -->

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
        " :class="nameClass">
                {{ channel.name }}
            </div>


            <div class="
          text-[7px]

          font-bold
          uppercase

          tracking-[0.12em]

          sm:text-[8px]

          md:text-[9px]
        " :class="statusClass">
                {{ statusText }}
            </div>

        </div>

    </div>
</template>