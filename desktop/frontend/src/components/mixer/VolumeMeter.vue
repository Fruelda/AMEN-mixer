<script setup lang="ts">

import {
    computed,
    ref
} from "vue"


type MeterColor =

    | "cyan"

    | "green"

    | "purple"

    | "yellow"

    | "white"



const props = defineProps<{

    volume: number

    color?: MeterColor

    muted?: boolean

}>()



const emit = defineEmits<{

    (
        event: "set-volume",
        volume: number
    ): void

}>()



const isDragging = ref(false)



const SEGMENT_COUNT = 28




// ============================================================
// ACTIVE SEGMENTS
// ============================================================

const activeSegments = computed(() => {

    const volume = Math.max(

        0,

        Math.min(

            100,

            Math.round(props.volume)

        )

    )


    return Math.round(

        (volume / 100) *

        SEGMENT_COUNT

    )

})




// ============================================================
// NORMAL COLOR
// ============================================================

const colorClass = computed(() => {


    const colors: Record<MeterColor, string> = {


        cyan:

            "bg-cyan-400",


        green:

            "bg-emerald-400",


        purple:

            "bg-violet-400",


        yellow:

            "bg-yellow-400",


        white:

            "bg-slate-200"


    }


    return colors[
        props.color ?? "cyan"
    ]

})





// ============================================================
// SEGMENT COLOR
// ============================================================

function segmentClass(
    index: number
) {


    if (
        props.muted
    ) {

        return (

            index <= activeSegments.value

                ? "bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.8)]"

                : "bg-red-500/20"

        )

    }



    return (

        index <= activeSegments.value

            ? colorClass.value

            : "bg-white/10"

    )

}





// ============================================================
// POINTER VOLUME
// ============================================================

function getVolumeFromPointer(

    event: PointerEvent

) {


    const element =

        event.currentTarget as HTMLElement



    const rect =

        element.getBoundingClientRect()



    if (
        rect.width <= 0
    ) {

        return 0

    }



    const percentage =

        (

            (event.clientX - rect.left)

            /

            rect.width

        )

        *

        100



    return Math.round(

        Math.max(

            0,

            Math.min(

                100,

                percentage

            )

        )

    )

}




function updateVolume(

    event: PointerEvent

) {


    emit(

        "set-volume",

        getVolumeFromPointer(event)

    )

}





function handlePointerDown(

    event: PointerEvent

) {


    isDragging.value = true



    const element =

        event.currentTarget as HTMLElement



    element.setPointerCapture(

        event.pointerId

    )



    updateVolume(event)

}




function handlePointerMove(

    event: PointerEvent

) {


    if (
        !isDragging.value
    ) {

        return

    }



    updateVolume(event)

}





function handlePointerUp(

    event: PointerEvent

) {


    isDragging.value = false



    const element =

        event.currentTarget as HTMLElement



    if (

        element.hasPointerCapture(

            event.pointerId

        )

    ) {


        element.releasePointerCapture(

            event.pointerId

        )

    }

}





function handlePointerCancel() {


    isDragging.value = false

}


</script>


<template>


    <div class="
            relative
            w-full
            select-none
            touch-none
        " :class="isDragging
                ? 'cursor-grabbing'
                : 'cursor-pointer'
            " @pointerdown="handlePointerDown" @pointermove="handlePointerMove" @pointerup="handlePointerUp"
        @pointercancel="handlePointerCancel">


        <div class="
                flex
                h-5
                w-full
                items-center
                gap-[1px]
                overflow-hidden
                rounded
                border
                border-white/10
                bg-black/30
                p-[2px]
                sm:h-7
                md:h-12
                md:gap-[3px]
                md:rounded-lg
                md:p-[3px]
            ">


            <div v-for="index in SEGMENT_COUNT" :key="index" class="
                    h-full
                    min-w-0
                    flex-1
                    rounded-sm
                    transition-all
                    duration-150
                    md:rounded-[3px]
                " :class="segmentClass(index)
                    " />



        </div>


    </div>


</template>