<script setup lang="ts">
import { computed } from "vue"

import DeviceCard from "./DeviceCard.vue"

import {
    mixerStore
} from "../../stores/mixer"


// ============================================================
// DEVICES
// ============================================================

const devices =
    computed(
        () =>
            mixerStore.devices
    )


// ============================================================
// DEVICE COUNT
// ============================================================

const deviceCount =
    computed(
        () =>
            devices.value.length
    )

</script>


<template>

    <div class="
      w-full
      rounded-2xl
      border
      border-white/10
      bg-white/5
      p-4
    ">

        <!-- ==================================================== -->
        <!-- HEADER -->
        <!-- ==================================================== -->

        <div class="
        mb-4
        flex
        items-center
        justify-between
      ">

            <div>

                <div class="
            text-sm
            font-semibold
            text-white
          ">
                    Connected Devices
                </div>

                <div class="
            mt-1
            text-xs
            text-white/40
          ">
                    {{ deviceCount }} device connected
                </div>

            </div>


            <div class="
          flex
          h-8
          min-w-8
          items-center
          justify-center
          rounded-full
          bg-white/10
          px-2
          text-xs
          font-semibold
          text-white
        ">
                {{ deviceCount }}
            </div>

        </div>


        <!-- ==================================================== -->
        <!-- EMPTY STATE -->
        <!-- ==================================================== -->

        <div v-if="devices.length === 0" class="
        rounded-xl
        border
        border-dashed
        border-white/10
        px-4
        py-5
        text-center
        text-xs
        text-white/40
      ">
            No connected devices
        </div>


        <!-- ==================================================== -->
        <!-- DEVICE LIST -->
        <!-- ==================================================== -->

        <div v-else class="space-y-2">

            <DeviceCard v-for="device in devices" :key="device.id" :device="device" />

        </div>

    </div>

</template>