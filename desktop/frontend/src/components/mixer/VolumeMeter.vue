<script setup lang="ts">
import { computed, ref } from "vue"

const props = defineProps<{
    volume: number
    color?: "cyan" | "green" | "purple" | "yellow" | "white"
}>()

const emit = defineEmits<{
    (event: "set-volume", volume: number): void
}>()

const isDragging = ref(false)

const segmentCount = 28

const normalizedVolume = computed(() => {
    return Math.max(
        0,
        Math.min(
            100,
            Math.round(props.volume)
        )
    )
})

const activeSegments = computed(() => {
    return Math.round(
        (normalizedVolume.value / 100) *
        segmentCount
    )
})

const colorClass = computed(() => {
    switch (props.color) {
        case "green":
            return "bg-emerald-400"

        case "purple":
            return "bg-violet-400"

        case "yellow":
            return "bg-yellow-400"

        case "white":
            return "bg-slate-200"

        default:
            return "bg-cyan-400"
    }
})

function setVolumeFromPointer(
    event: PointerEvent
) {
    const element =
        event.currentTarget as HTMLElement

    const rect =
        element.getBoundingClientRect()

    const x =
        event.clientX - rect.left

    const percentage =
        (x / rect.width) * 100

    const volume = Math.round(
        Math.max(
            0,
            Math.min(
                100,
                percentage
            )
        )
    )

    emit(
        "set-volume",
        volume
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

    setVolumeFromPointer(event)
}

function handlePointerMove(
    event: PointerEvent
) {
    if (!isDragging.value) {
        return
    }

    setVolumeFromPointer(event)
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
            cursor-pointer
        " :class="isDragging
            ? 'cursor-grabbing'
            : 'cursor-pointer'
            " @pointerdown="handlePointerDown" @pointermove="handlePointerMove" @pointerup="handlePointerUp"
        @pointercancel="handlePointerCancel">
        <!-- TRACK -->
        <div class="
                flex
                h-12
                w-full
                items-center
                gap-[3px]
                overflow-hidden
                rounded-lg
                border
                border-slate-800
                bg-slate-900
                p-[3px]
                shadow-inner
            ">
            <div v-for="index in segmentCount" :key="index" class="
                    h-full
                    min-w-0
                    flex-1
                    rounded-[3px]
                    transition-colors
                    duration-75
                " :class="index <= activeSegments
                    ? colorClass
                    : 'bg-slate-800'
                    " />
        </div>
    </div>
</template>