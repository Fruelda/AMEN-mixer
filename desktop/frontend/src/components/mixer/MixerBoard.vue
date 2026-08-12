<script setup lang="ts">

import {
  computed
} from "vue"

import ChannelRow from "./ChannelRow.vue"

import {
  mixerStore
} from "../../stores/mixer"


// ============================================================
// EMIT
// ============================================================

const emit =
  defineEmits<{
    (
      e: "open-devices"
    ): void
  }>()


// ============================================================
// STATE
// ============================================================

const channels =
  computed(
    () =>
      mixerStore.channels
  )


const isLoading =
  computed(
    () =>
      mixerStore.loading
  )


const deviceCount =
  computed(
    () =>
      mixerStore.devices.length
  )


// ============================================================
// OPEN DEVICES
// ============================================================

function openDevices() {

  emit(
    "open-devices"
  )

}


// ============================================================
// VOLUME
// ============================================================

function setVolume(
  id: number,
  volume: number
) {

  mixerStore.setVolume(
    id,
    volume
  )

}


// ============================================================
// MUTE
// ============================================================

function toggleMute(
  id: number
) {

  const channel =
    mixerStore.channels.find(
      item =>
        item.id === id
    )


  if (
    !channel
  ) {
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

    <div class="
        relative
        z-10
      ">

      <!-- ================================================= -->
      <!-- HEADER -->
      <!-- ================================================= -->

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

        <!-- ================================================= -->
        <!-- TITLE -->
        <!-- ================================================= -->

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


        <!-- ================================================= -->
        <!-- HEADER ACTIONS -->
        <!-- ================================================= -->

        <div class="
            flex
            items-center

            gap-1.5

            sm:gap-3
          ">

          <!-- =============================================== -->
          <!-- DEVICES BUTTON -->
          <!-- =============================================== -->

          <button type="button" class="
              group

              flex
              items-center

              gap-1.5

              rounded-full

              border
              border-cyan-400/25

              bg-cyan-400/10

              px-2
              py-0.5

              text-[8px]

              font-bold

              uppercase

              tracking-wider

              text-cyan-300

              backdrop-blur-md

              transition

              hover:border-cyan-400/50
              hover:bg-cyan-400/20

              active:scale-95

              sm:gap-2
              sm:px-4
              sm:py-2
              sm:text-xs
            " @click="openDevices">

            <!-- ICON -->

            <svg viewBox="0 0 24 24" fill="none" class="
                h-3
                w-3

                sm:h-4
                sm:w-4
              ">

              <rect x="3" y="4" width="18" height="12" rx="2" stroke="currentColor" stroke-width="1.8" />

              <path d="M8 20H16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />

              <path d="M12 16V20" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />

            </svg>


            <!-- COUNT -->

            <span>
              {{ deviceCount }}
            </span>


            <!-- LABEL -->

            <span class="
                hidden

                sm:inline
              ">
              Devices
            </span>

          </button>


          <!-- =============================================== -->
          <!-- ACTIVE STATUS -->
          <!-- =============================================== -->

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

        </div>

      </header>


      <!-- ================================================= -->
      <!-- MIXER -->
      <!-- ================================================= -->

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

        <!-- ================================================= -->
        <!-- LOADING -->
        <!-- ================================================= -->

        <div v-if="isLoading" class="
            flex

            min-h-24

            items-center
            justify-center

            text-slate-300
          ">

          Loading mixer...

        </div>


        <!-- ================================================= -->
        <!-- CHANNELS -->
        <!-- ================================================= -->

        <div v-else-if="channels.length > 0" class="
            flex
            flex-col
          ">

          <ChannelRow v-for="channel in channels" :key="channel.id" :channel="channel" @set-volume="setVolume"
            @toggle-mute="toggleMute" />

        </div>


        <!-- ================================================= -->
        <!-- EMPTY -->
        <!-- ================================================= -->

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