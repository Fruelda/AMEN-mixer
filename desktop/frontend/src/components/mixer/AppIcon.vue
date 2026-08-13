<script setup lang="ts">
import type { MixerChannel } from "../../types/mixer"

defineProps<{
  channel: MixerChannel
}>()
</script>

<template>
  <div class="app-icon" :class="{ 'app-icon--valeton': channel.app === 'valeton' }" aria-hidden="true">
    <!-- Browser / Chrome -->
    <svg v-if="channel.app === 'browser' || channel.app === 'chrome'" viewBox="0 0 48 48" class="h-full w-full">
      <circle cx="24" cy="24" r="21" fill="#fff" />
      <path d="M24 24 6.2 24A18 18 0 0 1 36.7 11H24Z" fill="#e94235" />
      <path d="m24 24 8.9 15.5A18 18 0 0 1 6.2 24H24Z" fill="#35a853" />
      <path d="m24 24 12.7-13A18 18 0 0 1 32.9 39.5Z" fill="#fabb05" />
      <circle cx="24" cy="24" r="8.5" fill="#4285f4" stroke="#e8f0fe" stroke-width="2" />
    </svg>

    <!-- Spotify -->
    <svg v-else-if="channel.app === 'spotify'" viewBox="0 0 48 48" class="h-full w-full">
      <circle cx="24" cy="24" r="22" fill="#1ed760" />
      <path d="M13 18.2c7.8-2.2 17.8-1.3 24 2.1" fill="none" stroke="#08150d" stroke-width="3.6"
        stroke-linecap="round" />
      <path d="M14.5 24.5c6.9-1.8 15.4-1.1 20.8 1.8" fill="none" stroke="#08150d" stroke-width="3.2"
        stroke-linecap="round" />
      <path d="M16.1 30.4c5.6-1.3 11.9-.7 16.6 1.6" fill="none" stroke="#08150d" stroke-width="2.8"
        stroke-linecap="round" />
    </svg>

    <!-- Discord -->
    <svg v-else-if="channel.app === 'discord'" viewBox="0 0 48 48" class="h-full w-full">
      <rect x="3" y="3" width="42" height="42" rx="13" fill="#5865f2" />
      <path
        d="M16 15.5c5.2-2.4 10.8-2.4 16 0 2.6 3.6 4 7.6 4.4 12-2.6 3.4-5.4 5.6-8.4 6.8l-2-2.8c1.1-.4 2.1-.9 3.1-1.5-5.1 2.3-10.3 2.3-15.3 0 1 .6 2.1 1.1 3.2 1.5l-2 2.8c-3-1.2-5.8-3.4-8.4-6.8.4-4.4 1.8-8.4 4.4-12 1.7-.8 3.3-1.3 5-1.7Z"
        fill="#fff" />
      <circle cx="19" cy="25" r="2.2" fill="#5865f2" />
      <circle cx="29" cy="25" r="2.2" fill="#5865f2" />
    </svg>

    <!-- Valeton -->
    <svg v-else-if="channel.app === 'valeton'" viewBox="0 0 48 48" class="h-full w-full">
      <g transform="rotate(-35 24 24)">
        <path d="M19 27c-5.5-1-9 1-9 5.1 0 4 4.2 7.4 9.3 6.2 4.6-1 7-5.5 6.3-10.2l-1.1-7.2-4 .6Z" fill="#c9781a" />
        <path d="M23.5 22 34.8 6.5l2.8 2.1-11 15.7Z" fill="#df9a2c" />
        <path d="m34.4 7.4 3.8-2.2 2.5 1.9-2.8 3.7Z" fill="#b96e13" />
        <circle cx="17" cy="32" r="3.2" fill="#f0b243" />
        <path d="M20.5 28.2 35.6 7.9" stroke="#fff0be" stroke-width=".8" />
      </g>
    </svg>

    <!-- Master / Default -->
    <svg v-else viewBox="0 0 48 48" class="h-full w-full">
      <circle cx="24" cy="24" r="22" fill="#4a535c" />
      <path d="M12 21h7l8-6v18l-8-6h-7Z" fill="#f1f3f5" />
      <path d="M31 19c2.2 2.4 2.2 7.6 0 10M35 16c4.5 4.6 4.5 11.4 0 16" fill="none" stroke="#f1f3f5" stroke-width="2.2"
        stroke-linecap="round" />
    </svg>
  </div>
</template>

<style scoped>
.app-icon {
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  flex: none;
  filter: drop-shadow(0 2px 4px rgb(0 0 0 / 50%));
}

.app-icon--valeton {
  width: 46px;
  height: 46px;
}
</style>