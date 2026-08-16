<script setup lang="ts">
import { computed } from "vue"

import ConnectedDevices from "./ConnectedDevices.vue"
import DevicesHeader from "./DeviceHeader.vue"

import { mixerStore } from "../../stores/mixer.ts"

const emit = defineEmits<{
  (event: "back"): void
}>()

const deviceCount = computed(
  () => mixerStore.devices.length,
)

function goBack() {
  emit("back")
}
</script>

<template>
  <main class="
    relative
    min-h-dvh
    w-full
    overflow-y-auto
    bg-cover
    bg-center
    px-2
    py-2
    pb-20
    sm:px-6
    sm:py-6
    lg:px-10
  ">
    <div class="relative z-10 mx-auto h-auto w-full max-w-[1800px]">
      <DevicesHeader :device-count="deviceCount" @back="goBack" />

      <section
        class="mb-3 rounded-xl border border-white/10 bg-black/20 px-4 py-3 backdrop-blur-md sm:mb-5 sm:px-6 sm:py-4">
        <div class="flex items-center gap-2">
          <span class="h-2 w-2 rounded-full bg-emerald-400 shadow-[0_0_10px_rgba(52,211,153,0.8)]" />

          <span class="text-[10px] font-semibold uppercase tracking-[0.2em] text-emerald-400 sm:text-xs">
            Realtime Network
          </span>
        </div>

        <p class="mt-2 max-w-3xl text-xs leading-relaxed text-white/50 sm:text-sm">
          Devices currently connected to the AMEN Mixer
          realtime network.
        </p>
      </section>

      <ConnectedDevices />

      <div class="mt-4 flex justify-center sm:hidden">
        <button type="button"
          class="rounded-full border border-white/10 bg-black/30 px-5 py-2 text-xs font-semibold text-white/70 backdrop-blur-md active:scale-95"
          @click="goBack">
          ← Back to Mixer
        </button>
      </div>
    </div>
  </main>
</template>