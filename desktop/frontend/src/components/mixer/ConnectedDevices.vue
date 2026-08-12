<script setup lang="ts">

import {
    computed
} from "vue"

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
            w-full
            rounded-2xl
            border
            border-white/10
            bg-white/5
            p-4
        ">

        <!-- ================================================= -->
        <!-- HEADER -->
        <!-- ================================================= -->

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


        <!-- ================================================= -->
        <!-- EMPTY STATE -->
        <!-- ================================================= -->

        <div v-if="
            devices.length === 0
        " class="
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


        <!-- ================================================= -->
        <!-- DEVICE LIST -->
        <!-- ================================================= -->

        <div v-else class="
                space-y-2
            ">

            <div v-for="
device in devices
                " :key="device.id
                    " class="
                    flex
                    items-center
                    gap-3
                    rounded-xl
                    bg-black/20
                    px-3
                    py-3
                ">

                <!-- ICON -->

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


                <!-- INFO -->

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

                        <span class="
                                truncate
                            ">
                            {{ device.id }}
                        </span>

                    </div>

                </div>


                <!-- STATUS -->

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

        </div>

    </div>

</template>