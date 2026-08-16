<script setup lang="ts">
defineProps<{
    deviceCount: number
}>()

const emit = defineEmits<{
    (event: "open-devices"): void
}>()
</script>

<template>
    <header
        class="relative mx-auto mb-1 flex max-w-[1800px] items-center justify-between rounded-lg border border-white/10 bg-black/20 px-2 py-1 backdrop-blur-md landscape:mb-0 sm:mb-6 sm:px-6 sm:py-4">

        <!-- LOGO CENTER -->
        <div class="absolute left-1/2 -translate-x-1/2">
            <span class="
                text-4xl
                font-black
                tracking-[0.35em]
                text-cyan-400
                drop-shadow-[0_0_15px_rgba(34,211,238,0.9)]
                animate-breathing-glow
                ">
                АМИНЬ
            </span>
        </div>


        <!-- DEVICE STATUS RIGHT -->
        <div class="ml-auto flex items-center gap-1.5 sm:gap-3">

            <button type="button"
                class="group flex items-center gap-1.5 rounded-full border border-cyan-400/25 bg-cyan-400/10 px-2 py-0.5 text-[8px] font-bold uppercase tracking-wider text-cyan-300 backdrop-blur-md transition hover:border-cyan-400/50 hover:bg-cyan-400/20 active:scale-95 sm:gap-2 sm:px-4 sm:py-2 sm:text-xs"
                @click="emit('open-devices')">
                <svg viewBox="0 0 24 24" fill="none" class="h-3 w-3 sm:h-4 sm:w-4">
                    <rect x="3" y="4" width="18" height="12" rx="2" stroke="currentColor" stroke-width="1.8" />
                    <path d="M8 20H16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
                    <path d="M12 16V20" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
                </svg>

                <span>
                    {{ deviceCount }}
                </span>

                <span class="hidden sm:inline">
                    Devices
                </span>
            </button>


            <template>
                <header
                    class="relative mx-auto mb-1 flex max-w-[1800px] items-center rounded-lg border border-white/10 bg-black/20 px-2 py-1 backdrop-blur-md landscape:mb-0 sm:mb-6 sm:px-6 sm:py-4">

                    <!-- DEVICE LEFT -->
                    <button type="button"
                        class="flex items-center gap-1.5 rounded-full border border-cyan-400/25 bg-cyan-400/10 px-2 py-0.5 text-[8px] font-bold uppercase tracking-wider text-cyan-300 backdrop-blur-md transition hover:border-cyan-400/50 hover:bg-cyan-400/20 active:scale-95 sm:gap-2 sm:px-4 sm:py-2 sm:text-xs"
                        @click="emit('open-devices')">
                        <svg viewBox="0 0 24 24" fill="none" class="h-3 w-3 sm:h-4 sm:w-4">
                            <rect x="3" y="4" width="18" height="12" rx="2" stroke="currentColor" stroke-width="1.8" />
                            <path d="M8 20H16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
                            <path d="M12 16V20" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
                        </svg>

                        <span>
                            {{ deviceCount }} DEVICES
                        </span>
                    </button>


                    <!-- LOGO CENTER -->
                    <div class="absolute left-1/2 -translate-x-1/2">
                        <span class="
                    text-4xl
                    font-black
                    tracking-[0.35em]
                    text-cyan-400
                    drop-shadow-[0_0_15px_rgba(34,211,238,0.9)]
                    animate-breathing-glow
                ">
                            АМИНЬ
                        </span>
                    </div>

                </header>
            </template>

        </div>

    </header>
</template>