<script setup lang="ts">

import type {
    ConnectedDevice
} from "../../types/mixer"


// ============================================================
// PROPS
// ============================================================

defineProps<{
    device: ConnectedDevice
}>()


// ============================================================
// DEVICE LABEL
// ============================================================

function getDeviceLabel(
    clientType: string
) {

    switch (
    clientType
    ) {

        case "desktop":
            return "Desktop"

        case "mobile":
            return "Mobile"

        case "tablet":
            return "Tablet"

        case "hardware":
            return "Hardware"

        default:
            return "Client"

    }

}


// ============================================================
// DEVICE ICON
// ============================================================

function getDeviceIcon(
    clientType: string
) {

    switch (
    clientType
    ) {

        case "desktop":
            return "🖥️"

        case "mobile":
            return "📱"

        case "tablet":
            return "▣"

        case "hardware":
            return "🎛️"

        default:
            return "●"

    }

}

</script>


<template>

    <div class="
      flex
      items-center
      gap-3
      rounded-xl
      bg-black/20
      px-3
      py-3
    ">

        <!-- ==================================================== -->
        <!-- ICON -->
        <!-- ==================================================== -->

        <div class="
        flex
        h-10
        w-10
        shrink-0
        items-center
        justify-center
        rounded-xl
        bg-white/5
        text-lg
      ">
            {{
                getDeviceIcon(
                    device.clientType
                )
            }}
        </div>


        <!-- ==================================================== -->
        <!-- INFO -->
        <!-- ==================================================== -->

        <div class="
        min-w-0
        flex-1
      ">

            <div class="
          truncate
          text-sm
          font-medium
          text-white
        ">
                {{ device.name }}
            </div>


            <div class="
          mt-0.5
          flex
          items-center
          gap-2
          text-xs
          text-white/40
        ">

                <span>
                    {{
                        getDeviceLabel(
                            device.clientType
                        )
                    }}
                </span>

                <span>
                    •
                </span>

                <span class="truncate">
                    {{ device.id }}
                </span>

            </div>

        </div>


        <!-- ==================================================== -->
        <!-- STATUS -->
        <!-- ==================================================== -->

        <div class="
        flex
        shrink-0
        items-center
        gap-2
      ">

            <div class="
          h-2
          w-2
          rounded-full
          bg-emerald-400
        " />

            <span class="
          text-xs
          text-emerald-400
        ">
                Online
            </span>

        </div>

    </div>

</template>