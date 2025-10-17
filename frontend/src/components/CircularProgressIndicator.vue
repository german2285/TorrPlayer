<template>
  <div class="circular-progress-container" :style="containerStyle">
    <svg
      :width="size"
      :height="size"
      class="circular-progress"
      :class="{ 'indeterminate': indeterminate }"
    >
      <!-- Track (background circle) -->
      <circle
        v-if="showTrack"
        class="progress-track"
        :cx="center"
        :cy="center"
        :r="radius"
        :stroke-width="thickness"
        :stroke="trackColor"
        fill="none"
      />

      <!-- Indicator (progress circle) -->
      <circle
        class="progress-indicator"
        :cx="center"
        :cy="center"
        :r="radius"
        :stroke-width="thickness"
        :stroke="indicatorColor"
        fill="none"
        :stroke-dasharray="circumference"
        :stroke-dashoffset="dashOffset"
        :style="indicatorStyle"
      />
    </svg>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  size?: number
  thickness?: number
  progress?: number // 0-100
  indeterminate?: boolean
  indicatorColor?: string
  trackColor?: string
  showTrack?: boolean
  variant?: 'default' | 'retreat' // Material 3 Expressive variants
}

const props = withDefaults(defineProps<Props>(), {
  size: 40,
  thickness: 4,
  progress: 0,
  indeterminate: true,
  indicatorColor: 'var(--md-sys-color-primary, #6750A4)',
  trackColor: 'var(--md-sys-color-secondary-container, #E8DEF8)',
  showTrack: true,
  variant: 'retreat'
})

// Calculations
const center = computed(() => props.size / 2)
const radius = computed(() => (props.size - props.thickness) / 2)
const circumference = computed(() => 2 * Math.PI * radius.value)

const dashOffset = computed(() => {
  if (props.indeterminate) return 0
  const offset = circumference.value - (props.progress / 100) * circumference.value
  return offset
})

const containerStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`
}))

const indicatorStyle = computed(() => ({
  transition: props.indeterminate ? 'none' : 'stroke-dashoffset 0.3s cubic-bezier(0.4, 0, 0.2, 1)'
}))
</script>

<style scoped>
.circular-progress-container {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.circular-progress {
  transform: rotate(-90deg);
}

.progress-track {
  opacity: 0.24;
}

.progress-indicator {
  stroke-linecap: round;
  transform-origin: center;
}

/* Material 3 Expressive - Indeterminate Animation (Retreat) */
.circular-progress.indeterminate .progress-indicator {
  animation: circular-retreat 2s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

@keyframes circular-retreat {
  0% {
    stroke-dasharray: 1, 200;
    stroke-dashoffset: 0;
    transform: rotate(0deg);
  }
  50% {
    stroke-dasharray: 100, 200;
    stroke-dashoffset: -50;
    transform: rotate(450deg);
  }
  100% {
    stroke-dasharray: 100, 200;
    stroke-dashoffset: -150;
    transform: rotate(1080deg);
  }
}

/* Material 3 Expressive - Enhanced motion with emphasized easing */
@media (prefers-reduced-motion: no-preference) {
  .progress-indicator {
    transition: stroke-dashoffset 0.3s cubic-bezier(0.2, 0, 0, 1);
  }
}

/* Accessibility */
@media (prefers-reduced-motion: reduce) {
  .circular-progress.indeterminate .progress-indicator {
    animation: circular-retreat-reduced 3s linear infinite;
  }

  @keyframes circular-retreat-reduced {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
}
</style>
