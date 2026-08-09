<script setup lang="ts">
import { computed } from "vue"

import ChannelRow from "./ChannelRow.vue"

import {
  mixerStore,
} from "../../stores/mixer"

const channels = computed(
  () => mixerStore.channels
)

const isLoading = computed(
  () => mixerStore.loading
)

function setVolume(
  id: number,
  volume: number
) {
  mixerStore.setVolume(
    id,
    volume
  )
}

function toggleMute(
  id: number
) {
  const channel =
    mixerStore.channels.find(
      (item) =>
        item.id === id
    )

  if (!channel) {
    return
  }

  channel.muted =
    !channel.muted
}
</script>

<template>
  <main class="
            min-h-screen
            bg-[#02070b]
            px-4
            py-6
            text-white

            sm:px-6
            lg:px-10
        ">
    <!-- HEADER -->
    <header class="
                mx-auto
                mb-6
                flex
                max-w-[1800px]
                items-center
                justify-between
            ">
      <div>
        <div class="
                        text-sm
                        font-bold
                        uppercase
                        tracking-[0.35em]
                        text-cyan-400
                    ">
          AMEN MIXER
        </div>

        <h1 class="
                        mt-2
                        text-5xl
                        font-black
                        tracking-tight
                        text-white

                        max-md:text-4xl
                    ">
          WINDOWS MIXER
        </h1>
      </div>

      <!-- STATUS -->
      <div class="
                    flex
                    items-center
                    gap-2
                    rounded-full
                    border
                    border-emerald-500/40
                    bg-emerald-500/5
                    px-5
                    py-2
                    text-sm
                    font-bold
                    uppercase
                    tracking-widest
                    text-emerald-400
                ">
        <span class="
                        h-2.5
                        w-2.5
                        rounded-full
                        bg-emerald-400
                        shadow-[0_0_12px_rgba(52,211,153,0.8)]
                    " />

        ACTIVE
      </div>
    </header>

    <!-- MIXER -->
    <section class="
                mx-auto
                max-w-[1800px]
                overflow-hidden
                rounded-3xl
                border
                border-slate-700/70
                bg-[#061019]
                shadow-2xl
            ">
      <!-- LOADING -->
      <div v-if="isLoading" class="
                    flex
                    min-h-[300px]
                    items-center
                    justify-center
                    text-slate-500
                ">
        Loading mixer...
      </div>

      <!-- CHANNELS -->
      <template v-else-if="
        channels.length > 0
      ">
        <ChannelRow v-for="channel in channels" :key="channel.id" :channel="channel" @set-volume="
          setVolume
        " @toggle-mute="
          toggleMute
        " />
      </template>

      <!-- EMPTY -->
      <div v-else class="
                    flex
                    min-h-[300px]
                    items-center
                    justify-center
                    text-slate-500
                ">
        No mixer channels available
      </div>
    </section>
  </main>
</template>