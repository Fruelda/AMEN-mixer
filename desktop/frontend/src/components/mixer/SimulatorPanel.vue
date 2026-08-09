<script setup lang="ts">

import {
    mixerStore,
} from "../../stores/mixer"


function sendEncoder(
    channel: number,
    value: number
) {

    mixerStore.handleCommand({
        type: "ENC",
        channel,
        value,
    })

}


function sendButton(
    channel: number
) {

    mixerStore.handleCommand({
        type: "BTN",
        channel,
        value: 1,
    })

}

</script>


<template>

    <div class="
            mt-4
            rounded-2xl
            border
            border-slate-700
            bg-[#08131e]
            p-4
            shadow-xl
        ">

        <!-- HEADER -->

        <div class="
                mb-4
                flex
                items-center
                justify-between
            ">

            <div>

                <div class="
                        text-xs
                        font-black
                        uppercase
                        tracking-[0.2em]
                        text-cyan-400
                    ">
                    Hardware Simulator
                </div>

                <div class="
                        mt-1
                        text-xs
                        text-slate-500
                    ">
                    Test EC11 tanpa ESP32
                </div>

            </div>


            <div class="
                    rounded-full
                    border
                    border-yellow-500/30
                    bg-yellow-500/10
                    px-3
                    py-1
                    text-[10px]
                    font-bold
                    text-yellow-400
                ">
                TEST MODE
            </div>

        </div>


        <!-- EC1 -->

        <div class="
                mb-3
                rounded-xl
                border
                border-slate-700/70
                bg-slate-900/40
                p-3
            ">

            <div class="
                    mb-3
                    flex
                    items-center
                    justify-between
                ">

                <span class="
                        text-sm
                        font-bold
                        text-white
                    ">
                    EC1 — Master
                </span>

                <span class="
                        text-xs
                        font-mono
                        text-cyan-400
                    ">
                    CH 1
                </span>

            </div>


            <div class="
                    grid
                    grid-cols-3
                    gap-2
                ">

                <!-- MINUS -->

                <button type="button" class="
                        rounded-xl
                        border
                        border-slate-700
                        bg-slate-800
                        py-3
                        text-xl
                        font-black
                        text-white
                        transition

                        hover:border-cyan-400
                        hover:bg-slate-700

                        active:scale-95
                        active:bg-cyan-500
                    " @click="sendEncoder(1, -1)">
                    −
                </button>


                <!-- BUTTON -->

                <button type="button" class="
                        rounded-xl
                        border
                        border-cyan-500/40
                        bg-cyan-500/10
                        py-3
                        text-xs
                        font-black
                        uppercase
                        tracking-widest
                        text-cyan-400
                        transition

                        hover:bg-cyan-500/20

                        active:scale-95
                    " @click="sendButton(1)">
                    SW
                </button>


                <!-- PLUS -->

                <button type="button" class="
                        rounded-xl
                        border
                        border-slate-700
                        bg-slate-800
                        py-3
                        text-xl
                        font-black
                        text-white
                        transition

                        hover:border-cyan-400
                        hover:bg-slate-700

                        active:scale-95
                        active:bg-cyan-500
                    " @click="sendEncoder(1, 1)">
                    +
                </button>

            </div>

        </div>


        <!-- EC4 -->

        <div class="
                rounded-xl
                border
                border-slate-700/70
                bg-slate-900/40
                p-3
            ">

            <div class="
                    mb-3
                    flex
                    items-center
                    justify-between
                ">

                <span class="
                        text-sm
                        font-bold
                        text-white
                    ">
                    EC4 — Valeton
                </span>

                <span class="
                        text-xs
                        font-mono
                        text-orange-400
                    ">
                    CH 4
                </span>

            </div>


            <div class="
                    grid
                    grid-cols-3
                    gap-2
                ">

                <!-- MINUS -->

                <button type="button" class="
                        rounded-xl
                        border
                        border-slate-700
                        bg-slate-800
                        py-3
                        text-xl
                        font-black
                        text-white
                        transition

                        hover:border-orange-400
                        hover:bg-slate-700

                        active:scale-95
                        active:bg-orange-500
                    " @click="sendEncoder(4, -1)">
                    −
                </button>


                <!-- BUTTON -->

                <button type="button" class="
                        rounded-xl
                        border
                        border-orange-500/40
                        bg-orange-500/10
                        py-3
                        text-xs
                        font-black
                        uppercase
                        tracking-widest
                        text-orange-400
                        transition

                        hover:bg-orange-500/20

                        active:scale-95
                    " @click="sendButton(4)">
                    SW
                </button>


                <!-- PLUS -->

                <button type="button" class="
                        rounded-xl
                        border
                        border-slate-700
                        bg-slate-800
                        py-3
                        text-xl
                        font-black
                        text-white
                        transition

                        hover:border-orange-400
                        hover:bg-slate-700

                        active:scale-95
                        active:bg-orange-500
                    " @click="sendEncoder(4, 1)">
                    +
                </button>

            </div>

        </div>

    </div>

</template>