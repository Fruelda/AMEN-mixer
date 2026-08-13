<script setup lang="ts">
defineProps<{
    deviceCount: number
}>()

const emit = defineEmits<{
    (event: "back"): void
}>()
</script>

<template>
    <header
        class="mb-3 flex items-center justify-between rounded-xl border border-white/10 bg-black/30 px-3 py-3 backdrop-blur-xl sm:mb-6 sm:px-6 sm:py-5">
        <div class="flex items-center gap-3">
            <button type="button"
                class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border border-white/10 bg-white/5 text-white/70 transition hover:border-cyan-400/30 hover:bg-cyan-400/10 hover:text-cyan-300 active:scale-95 sm:h-11 sm:w-11"
                aria-label="Back to mixer" @click="emit('back')">
                <svg viewBox="0 0 24 24" fill="none" class="h-5 w-5">
                    <path d="M15 18L9 12L15 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                        stroke-linejoin="round" />
                </svg>
            </button>

            <div>
                <div class="text-[8px] font-bold uppercase tracking-[0.3em] text-cyan-400 sm:text-xs">
                    AMEN MIXER
                </div>

                <h1 class="text-lg font-black tracking-tight text-white sm:text-3xl lg:text-4xl">
                    CONNECTED DEVICES
                </h1>
            </div>
        </div>

        <div
            class="flex items-center gap-2 rounded-full border border-cyan-400/25 bg-cyan-400/10 px-3 py-1.5 text-xs font-bold uppercase tracking-wider text-cyan-300 sm:px-5 sm:py-2 sm:text-sm">
            <span class="h-2 w-2 rounded-full bg-cyan-400 shadow-[0_0_12px_rgba(34,211,238,0.8)]" />

            {{ deviceCount }}

            <span class="hidden sm:inline">
                Connected
            </span>
        </div>
    </header>
</template>