<script setup lang="ts">
import ChannelRow from './ChannelRow.vue'
import { mixerStore } from '../../stores/mixer'

function setVolume(id: number, value: number) {
  const channel = mixerStore.channels.find(
    (item) => item.id === id
  )

  if (!channel) return

  channel.volume = Math.max(
    0,
    Math.min(100, Number(value))
  )
}

function toggleMute(id: number) {
  const channel = mixerStore.channels.find(
    (item) => item.id === id
  )

  if (!channel) return

  channel.muted = !channel.muted
}
</script>

<template>
  <section class="
      w-full
      max-w-[980px]
      rounded-[22px]
      border
      border-slate-700
      bg-[#07121b]
      p-7
      shadow-2xl
    ">
    <!-- HEADER -->
    <div class="
        mb-6
        flex
        items-center
        justify-between
        gap-4
      ">
      <h1 class="
          text-[30px]
          font-black
          uppercase
          tracking-wide
          text-white
        ">
        <span class="text-sky-400">
          Mode:
        </span>

        Windows Mixer
      </h1>

      <div class="
          flex
          items-center
          gap-2
          rounded-full
          border
          border-slate-700
          bg-slate-900/60
          px-3
          py-1.5
          text-xs
          font-semibold
          uppercase
          tracking-wider
          text-slate-400
        ">
        <span class="
            h-2
            w-2
            rounded-full
            bg-emerald-400
          " />

        Active
      </div>
    </div>

    <!-- MIXER BODY -->
    <div class="
        overflow-hidden
        rounded-[18px]
        border
        border-slate-700
        bg-[#08131d]
        px-5
        py-4
      ">
      <!-- EMPTY STATE -->
      <div v-if="mixerStore.channels.length === 0" class="
          flex
          min-h-[220px]
          items-center
          justify-center
          text-center
          text-sm
          text-slate-500
        ">
        No mixer channels available
      </div>

      <!-- CHANNEL LIST -->
      <div v-else class="
          divide-y
          divide-slate-800/80
        ">
        <ChannelRow v-for="channel in mixerStore.channels" :key="channel.id" :channel="channel" @volume-change="
          setVolume(channel.id, $event)
          " @mute-toggle="
            toggleMute(channel.id)
            " />
      </div>
    </div>
  </section>
</template>