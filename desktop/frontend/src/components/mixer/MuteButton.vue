<script setup lang="ts">
import { computed } from "vue"

const props = defineProps<{
    muted: boolean
}>()

const emit = defineEmits<{
    (event: "toggle"): void
}>()

const buttonClass = computed(() => {
    return props.muted
        ? [
            "border-red-400/60",
            "bg-red-500/80",
            "text-white",
            "shadow-[0_0_18px_rgba(239,68,68,0.45)]"
        ]
        : [
            "border-white/15",
            "bg-white/5",
            "text-white/50",
            "hover:border-white/25",
            "hover:bg-white/10",
            "hover:text-white"
        ]
})

const label = computed(() => {
    return props.muted
        ? "MUTED"
        : "MUTE"
})

const title = computed(() => {
    return props.muted
        ? "Unmute"
        : "Mute"
})

function toggle() {
    emit("toggle")
}
</script>

<template>
    <button type="button" :aria-pressed="muted" :title="title" class="
      flex
      h-8
      w-[78px]
      shrink-0
      items-center
      justify-center
      gap-1.5
      rounded-lg
      border

      text-[9px]
      font-black
      uppercase

      transition-colors
      duration-150

      active:scale-95

      sm:h-10
      sm:w-[90px]
      sm:text-xs

      md:h-12
      md:w-[100px]
      md:text-sm
    " :class="buttonClass" @click="toggle">
        <svg viewBox="0 0 24 24" fill="none" class="
        h-3.5
        w-3.5
        shrink-0

        sm:h-4
        sm:w-4
      ">
            <!-- SPEAKER -->
            <path d="M5 9V15H9L14 19V5L9 9H5Z" fill="currentColor" />

            <!-- SOUND -->
            <template v-if="!muted">
                <path d="M17 9C18.3 10.5 18.3 13.5 17 15" stroke="currentColor" stroke-width="1.8"
                    stroke-linecap="round" />

                <path d="M19.5 6.5C22 9.5 22 14.5 19.5 17.5" stroke="currentColor" stroke-width="1.8"
                    stroke-linecap="round" />
            </template>

            <!-- MUTED -->
            <template v-else>
                <path d="M17 9L21 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" />

                <path d="M21 9L17 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
            </template>
        </svg>

        <span>
            {{ label }}
        </span>
    </button>
</template>