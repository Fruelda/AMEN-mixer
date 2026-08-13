<script setup lang="ts">
import { computed } from "vue"

import ChannelRow from "./ChannelRow.vue"
import MixerHeader from "./MixerHeader.vue"

import { mixerStore } from "../../stores/mixer"

const emit = defineEmits<{
  (event: "open-devices"): void
}>()

const channels = computed(() => mixerStore.channels)
const isLoading = computed(() => mixerStore.loading)
const deviceCount = computed(() => mixerStore.devices.length)

function setVolume(id: number, volume: number) {
  void mixerStore.setVolume(id, volume)
}

function toggleMute(id: number) {
  const channel = mixerStore.channels.find(
    (item) => item.id === id,
  )

  if (!channel) return

  channel.muted = !channel.muted
}

function openDevices() {
  emit("open-devices")
}
</script>

<template>
  <main
    class="relative min-h-screen w-full bg-cover bg-center bg-fixed px-2 py-1 landscape:px-1 landscape:py-0 sm:px-6 sm:py-6 lg:px-10"
    :style="{
      backgroundImage: `
        linear-gradient(
          rgba(2, 8, 13, 0.45),
          rgba(2, 8, 13, 0.65)
        ),
        url('/amen-bg.png')
      `,
    }">
    <div class="relative z-10">

      <MixerHeader :device-count="deviceCount" @open-devices="openDevices" />

      <section
        class="mx-auto w-full max-w-[1800px] overflow-hidden rounded-xl border border-white/15 bg-black/35 backdrop-blur-xl shadow-2xl landscape:rounded-lg">
        <div v-if="isLoading" class="flex min-h-24 items-center justify-center text-slate-300">
          Loading mixer...
        </div>

        <div v-else-if="channels.length > 0" class="flex flex-col">
          <ChannelRow v-for="channel in channels" :key="channel.id" :channel="channel" @set-volume="setVolume"
            @toggle-mute="toggleMute" />
        </div>

        <div v-else class="flex min-h-24 items-center justify-center text-slate-300">
          No mixer channels available
        </div>
      </section>

    </div>
  </main>
</template>