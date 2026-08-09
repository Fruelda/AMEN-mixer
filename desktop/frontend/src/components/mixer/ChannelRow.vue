<script setup lang="ts">
import VolumeMeter from "./VolumeMeter.vue"

interface MixerChannel {
  id: number
  name: string
  app: string
  volume: number
  muted: boolean
  connected: boolean
  selected: boolean
}

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

function getVolumeColor() {
  switch (props.channel.app) {
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
            grid
            grid-cols-[160px_minmax(300px,1fr)_100px]
            items-center
            gap-6
            border-b
            border-slate-800/80
            px-4
            py-5
            last:border-b-0

            max-md:grid-cols-[120px_minmax(220px,1fr)_80px]
            max-md:gap-4
        ">
    <!-- CHANNEL -->
    <div class="
                min-w-0
            ">
      <!-- NAME -->
      <div class="
                    truncate
                    text-2xl
                    font-bold
                    text-white

                    max-md:text-xl
                ">
        {{ channel.name }}
      </div>

      <!-- STATUS -->
      <div class="
                    mt-1
                    text-xs
                    font-bold
                    uppercase
                    tracking-[0.2em]
                " :class="channel.connected
                    ? 'text-emerald-400'
                    : 'text-slate-500'
                  ">
        {{
          channel.connected
            ? "Connected"
            : "Offline"
        }}
      </div>
    </div>

    <!-- VOLUME -->
    <VolumeMeter :volume="channel.volume" :color="getVolumeColor()" @set-volume="
      (volume) =>
        emit(
          'set-volume',
          channel.id,
          volume
        )
    " />

    <!-- VALUE -->
    <div class="
                text-right
                text-3xl
                font-bold
                tabular-nums
                text-white

                max-md:text-2xl
            ">
      {{ Math.round(channel.volume) }}%
    </div>
  </div>
</template>