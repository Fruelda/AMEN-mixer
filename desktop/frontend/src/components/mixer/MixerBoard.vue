<script setup lang="ts">
import { computed } from "vue"
// import bgImage from "../../assets/images/amen-bg.png"

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
      relative
      min-h-screen
      w-full

      bg-cover
      bg-center
      bg-fixed

      px-2
      py-1

      landscape:px-1
      landscape:py-0

      sm:px-6
      sm:py-6

      lg:px-10
    " :style="{
      backgroundImage: `
        linear-gradient(
          rgba(2,8,13,0.45),
          rgba(2,8,13,0.65)
        ),
        url('/amen-bg.png')
      `
    }">


    <div class="relative z-10">


      <!-- HEADER -->

      <header class="
          mx-auto

          mb-1

          flex
          max-w-[1800px]

          items-center
          justify-between


          rounded-lg

          border
          border-white/10

          bg-black/20

          px-2
          py-1


          backdrop-blur-md


          landscape:mb-0


          sm:mb-6
          sm:px-6
          sm:py-4
        ">


        <div>


          <div class="
              text-[8px]

              font-bold

              uppercase

              tracking-[0.3em]

              text-cyan-400


              landscape:text-[7px]


              sm:text-sm
            ">
            AMEN MIXER
          </div>



          <h1 class="
              mt-0

              text-lg

              font-black

              tracking-tight

              text-white


              landscape:text-base


              sm:text-4xl

              lg:text-5xl
            ">
            WINDOWS MIXER
          </h1>


        </div>





        <!-- STATUS -->


        <div class="
            flex

            items-center

            gap-1


            rounded-full


            border

            border-emerald-400/30


            bg-emerald-400/10


            px-2

            py-0.5


            text-[8px]

            font-bold

            uppercase

            tracking-widest


            text-emerald-400


            backdrop-blur-md


            landscape:px-2

            landscape:py-0


            sm:px-5

            sm:py-2

            sm:text-sm
          ">


          <span class="
              h-1.5

              w-1.5


              rounded-full


              bg-emerald-400


              shadow-[0_0_12px_rgba(52,211,153,0.8)]


              sm:h-2.5

              sm:w-2.5
            " />


          ACTIVE


        </div>


      </header>






      <!-- MIXER -->


      <section class="
          mx-auto

          w-full

          max-w-[1800px]


          overflow-hidden


          rounded-xl


          border

          border-white/15


          bg-black/35


          backdrop-blur-xl


          shadow-2xl


          landscape:rounded-lg
        ">



        <div v-if="isLoading" class="
            flex

            min-h-24

            items-center

            justify-center

            text-slate-300
          ">

          Loading mixer...

        </div>





        <div v-else-if="channels.length > 0" class="
            flex

            flex-col
          ">


          <ChannelRow v-for="channel in channels" :key="channel.id" :channel="channel" @set-volume="setVolume"
            @toggle-mute="toggleMute" />


        </div>





        <div v-else class="
            flex

            min-h-24

            items-center

            justify-center

            text-slate-300
          ">

          No mixer channels available

        </div>



      </section>


    </div>


  </main>


</template>