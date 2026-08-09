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
      flex
      items-center

      gap-1


      border-b
      border-white/10


      px-2
      py-[3px]


      landscape:py-0


      last:border-b-0


      sm:gap-3
      sm:px-3
      sm:py-2


      md:gap-6
      md:px-8
      md:py-5
    ">


    <!-- NAME -->


    <div class="
        w-16
        shrink-0


        landscape:w-14


        sm:w-24


        md:w-40
      ">


      <div class="
          truncate


          text-xs

          font-bold


          text-white


          landscape:text-[11px]


          sm:text-base


          md:text-2xl
        ">

        {{ channel.name }}

      </div>



      <div class="
          text-[6px]

          font-bold

          uppercase

          tracking-[0.12em]


          landscape:text-[5px]


          sm:text-[9px]


          md:text-xs
        " :class="channel.connected
            ? 'text-emerald-400'
            : 'text-slate-500'
          ">

        {{
          channel.connected
            ? 'Connected'
            : 'Offline'
        }}

      </div>


    </div>





    <!-- METER -->


    <div class="
        min-w-0
        flex-1
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





    <!-- VALUE -->


    <div class="
        w-9

        shrink-0


        text-right


        text-sm

        font-bold


        tabular-nums


        text-white


        landscape:text-xs


        sm:w-12

        sm:text-xl


        md:w-24

        md:text-3xl
      ">

      {{ Math.round(channel.volume) }}%

    </div>


  </div>


</template>