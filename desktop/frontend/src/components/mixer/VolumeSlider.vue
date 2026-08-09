<template>
  <label class="volume-slider" :aria-label="`Volume ${modelValue}%`">
    <span class="volume-slider__track">
      <span class="volume-slider__fill" :style="fillStyle" />
      <span class="volume-slider__segments" />
    </span>

    <input
      type="range"
      min="0"
      max="100"
      step="1"
      :value="modelValue"
      @input="onInput"
    >
  </label>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  modelValue: number
  color: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const palette: Record<string, string> = {
  blue: '#168ee9',
  green: '#27cc67',
  purple: '#8955de',
  yellow: '#ffc400',
  white: '#e5e9ec',
}

const fillStyle = computed(() => ({
  width: `${Math.max(0, Math.min(100, props.modelValue))}%`,
  backgroundColor: palette[props.color] ?? palette.blue,
  boxShadow: `0 0 9px ${(palette[props.color] ?? palette.blue)}44`,
}))

function onInput(event: Event) {
  emit('update:modelValue', Number((event.target as HTMLInputElement).value))
}
</script>

<style scoped>
.volume-slider {
  position: relative;
  display: block;
  width: 100%;
  height: 25px;
  cursor: pointer;
  user-select: none;
}

.volume-slider__track {
  position: absolute;
  inset: 3px 0;
  overflow: hidden;
  border-radius: 2px;
  background: #252b30;
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, .8), 0 1px 0 rgba(255, 255, 255, .025);
}

.volume-slider__fill {
  position: absolute;
  inset: 0 auto 0 0;
  transition: width 90ms linear;
}

.volume-slider__segments {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image: repeating-linear-gradient(
    90deg,
    transparent 0 4px,
    rgba(5, 10, 14, .72) 4px 6px
  );
}

input {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}
</style>
