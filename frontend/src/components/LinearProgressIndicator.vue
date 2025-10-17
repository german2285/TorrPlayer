<template>
  <div class="linear-progress-container" :style="containerStyle">
    <div class="linear-progress-track" :style="trackStyle">
      <div
        v-if="!indeterminate"
        class="linear-progress-indicator"
        :style="indicatorStyle"
      />
      <div
        v-else
        class="linear-progress-indicator indeterminate"
        :class="{ 'disjoint': animationType === 'disjoint' }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  progress?: number // 0-100
  indeterminate?: boolean
  thickness?: number
  indicatorColor?: string
  trackColor?: string
  animationType?: 'contiguous' | 'disjoint' // Material 3 animation types
  rounded?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  progress: 0,
  indeterminate: true,
  thickness: 4,
  indicatorColor: 'var(--md-sys-color-primary, #6750A4)',
  trackColor: 'var(--md-sys-color-secondary-container, #E8DEF8)',
  animationType: 'disjoint',
  rounded: true
})

const containerStyle = computed(() => ({
  height: `${props.thickness}px`
}))

const trackStyle = computed(() => ({
  backgroundColor: props.trackColor,
  borderRadius: props.rounded ? '50px' : '0'
}))

const indicatorStyle = computed(() => ({
  width: `${props.progress}%`,
  backgroundColor: props.indicatorColor,
  borderRadius: props.rounded ? '50px' : '0',
  transition: 'width 0.3s cubic-bezier(0.4, 0, 0.2, 1)'
}))
</script>

<style scoped>
.linear-progress-container {
  width: 100%;
  position: relative;
  overflow: hidden;
}

.linear-progress-track {
  width: 100%;
  height: 100%;
  position: relative;
  opacity: 0.24;
}

.linear-progress-indicator {
  height: 100%;
  position: absolute;
  left: 0;
  top: 0;
  opacity: 1;
}

/* Material 3 Expressive - Indeterminate Animation (Disjoint) */
.linear-progress-indicator.indeterminate.disjoint {
  width: 0%;
  animation: linear-disjoint 2s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

@keyframes linear-disjoint {
  0% {
    left: -35%;
    width: 0%;
  }
  30% {
    left: -15%;
    width: 50%;
  }
  60% {
    left: 50%;
    width: 50%;
  }
  100% {
    left: 115%;
    width: 0%;
  }
}

/* Material 3 Expressive - Contiguous Animation */
.linear-progress-indicator.indeterminate:not(.disjoint) {
  width: 100%;
  animation: linear-contiguous 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes linear-contiguous {
  0% {
    left: -100%;
  }
  100% {
    left: 100%;
  }
}

/* Material 3 Expressive - Enhanced motion */
@media (prefers-reduced-motion: no-preference) {
  .linear-progress-indicator:not(.indeterminate) {
    transition: width 0.3s cubic-bezier(0.2, 0, 0, 1);
  }
}

/* Accessibility */
@media (prefers-reduced-motion: reduce) {
  .linear-progress-indicator.indeterminate {
    animation: linear-reduced 3s linear infinite;
  }

  @keyframes linear-reduced {
    0% { left: -50%; }
    100% { left: 100%; }
  }
}
</style>
