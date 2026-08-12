<script setup lang="ts">

import {
    computed
} from "vue"

import ConnectedDevices from "./ConnectedDevices.vue"

import {
    mixerStore
} from "../../stores/mixer"


// ============================================================
// EMIT
// ============================================================

const emit =
    defineEmits<{
        (
            e: "back"
        ): void
    }>()


// ============================================================
// STATE
// ============================================================

const deviceCount =
    computed(
        () =>
            mixerStore.devices.length
    )


// ============================================================
// BACK
// ============================================================

function goBack() {

    emit(
        "back"
    )

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
      py-2

      sm:px-6
      sm:py-6

      lg:px-10
    " :style="{
        backgroundImage: `
        linear-gradient(
          rgba(2,8,13,0.55),
          rgba(2,8,13,0.75)
        ),
        url('/amen-bg.png')
      `
    }">

        <div class="
        relative
        z-10

        mx-auto

        w-full
        max-w-[1800px]
      ">

            <!-- ================================================= -->
            <!-- HEADER -->
            <!-- ================================================= -->

            <header class="
          mb-3

          flex

          items-center
          justify-between

          rounded-xl

          border
          border-white/10

          bg-black/30

          px-3
          py-3

          backdrop-blur-xl

          sm:mb-6
          sm:px-6
          sm:py-5
        ">

                <!-- ================================================= -->
                <!-- LEFT -->
                <!-- ================================================= -->

                <div class="
            flex
            items-center

            gap-3
          ">

                    <!-- =============================================== -->
                    <!-- BACK BUTTON -->
                    <!-- =============================================== -->

                    <button type="button" class="
              flex

              h-9
              w-9

              shrink-0

              items-center
              justify-center

              rounded-xl

              border
              border-white/10

              bg-white/5

              text-white/70

              transition

              hover:border-cyan-400/30
              hover:bg-cyan-400/10
              hover:text-cyan-300

              active:scale-95

              sm:h-11
              sm:w-11
            " aria-label="Back to mixer" @click="goBack">

                        <svg viewBox="0 0 24 24" fill="none" class="
                h-5
                w-5
              ">

                            <path d="M15 18L9 12L15 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                                stroke-linejoin="round" />

                        </svg>

                    </button>


                    <!-- =============================================== -->
                    <!-- TITLE -->
                    <!-- =============================================== -->

                    <div>

                        <div class="
                text-[8px]

                font-bold

                uppercase

                tracking-[0.3em]

                text-cyan-400

                sm:text-xs
              ">
                            AMEN MIXER
                        </div>


                        <h1 class="
                text-lg

                font-black

                tracking-tight

                text-white

                sm:text-3xl

                lg:text-4xl
              ">
                            CONNECTED DEVICES
                        </h1>

                    </div>

                </div>


                <!-- ================================================= -->
                <!-- TOTAL -->
                <!-- ================================================= -->

                <div class="
            flex
            items-center

            gap-2

            rounded-full

            border
            border-cyan-400/25

            bg-cyan-400/10

            px-3
            py-1.5

            text-xs

            font-bold

            uppercase

            tracking-wider

            text-cyan-300

            sm:px-5
            sm:py-2
            sm:text-sm
          ">

                    <span class="
              h-2
              w-2

              rounded-full

              bg-cyan-400

              shadow-[0_0_12px_rgba(34,211,238,0.8)]
            " />


                    {{ deviceCount }}


                    <span class="
              hidden

              sm:inline
            ">
                        Connected
                    </span>

                </div>

            </header>


            <!-- ================================================= -->
            <!-- PAGE DESCRIPTION -->
            <!-- ================================================= -->

            <div class="
          mb-3

          rounded-xl

          border
          border-white/10

          bg-black/20

          px-4
          py-3

          backdrop-blur-md

          sm:mb-5
          sm:px-6
          sm:py-4
        ">

                <div class="
            flex
            items-center

            gap-2
          ">

                    <span class="
              h-2
              w-2

              rounded-full

              bg-emerald-400

              shadow-[0_0_10px_rgba(52,211,153,0.8)]
            " />


                    <span class="
              text-[10px]

              font-semibold

              uppercase

              tracking-[0.2em]

              text-emerald-400

              sm:text-xs
            ">
                        Realtime Network
                    </span>

                </div>


                <p class="
            mt-2

            max-w-3xl

            text-xs

            leading-relaxed

            text-white/50

            sm:text-sm
          ">
                    Devices currently connected to the AMEN Mixer realtime network.
                </p>

            </div>


            <!-- ================================================= -->
            <!-- DEVICE LIST -->
            <!-- ================================================= -->

            <ConnectedDevices />


            <!-- ================================================= -->
            <!-- BACK MOBILE -->
            <!-- ================================================= -->

            <div class="
          mt-4

          flex
          justify-center

          sm:hidden
        ">

                <button type="button" class="
            rounded-full

            border
            border-white/10

            bg-black/30

            px-5
            py-2

            text-xs

            font-semibold

            text-white/70

            backdrop-blur-md

            active:scale-95
          " @click="goBack">
                    ← Back to Mixer
                </button>

            </div>

        </div>

    </main>

</template>