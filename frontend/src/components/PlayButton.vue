<template>
  <button class="play-button" @click="onClick">
    <span class="button-background"></span>
    <span class="button-ripple"></span>
    <img :src="playIcon" class="button-icon" alt="Play" />
    <span class="button-label">Play</span>
  </button>
</template>

<script setup lang="ts">
import playIcon from '../assets/icons/play.svg'

const emit = defineEmits<{
  (e: 'click'): void
}>()

const onClick = (): void => {
  console.log('🎬 PlayButton: Click event triggered!')
  emit('click')
}
</script>

<style scoped>
/* M3 Expressive Filled Button */
/* Based on Material Design 3 Expressive specifications */

.play-button {
  /* Container */
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  min-height: 36px;
  min-width: auto;
  padding: 8px 20px;
  border: none;
  border-radius: var(--md-sys-shape-corner-full); /* M3 Expressive: полностью округлый */
  outline: none;
  cursor: pointer;
  user-select: none;
  -webkit-tap-highlight-color: transparent;
  overflow: hidden;

  /* Typography */
  font-family: var(--md-sys-typescale-label-large-font);
  font-size: 14px;
  font-weight: var(--md-sys-typescale-label-large-weight);
  line-height: var(--md-sys-typescale-label-large-line-height);
  text-transform: none;

  /* Colors */
  color: var(--md-sys-color-on-primary);

  /* M3 Expressive Physics: Fast spatial spring for transform and size */
  transition:
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    box-shadow var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);

  box-shadow: var(--md-sys-elevation-level1);
  will-change: transform, box-shadow;
}

/* Background layer */
.button-background {
  position: absolute;
  inset: 0;
  background: var(--md-sys-color-primary);
  border-radius: inherit;
  z-index: 0;
  /* M3 Expressive Physics: Fast effects spring for background color */
  transition: background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
}

/* Ripple/State layer */
.button-ripple {
  position: absolute;
  inset: 0;
  background: var(--md-sys-color-on-primary);
  opacity: 0;
  border-radius: inherit;
  z-index: 1;
  pointer-events: none;
  /* M3 Expressive Physics: Fast effects spring for opacity */
  transition: opacity var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
}

/* Content (icon and label) */
.button-icon,
.button-label {
  position: relative;
  z-index: 2;
}

.button-icon {
  width: 20px; /* Smaller icon for card */
  height: 20px;
  filter: brightness(0) invert(1); /* Make icon white to match on-primary color */
  /* M3 Expressive Physics: Fast spatial spring for icon scale */
  transition:
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    filter var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
}

.button-label {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Hover state */
.play-button:hover {
  transform: scale(1.05); /* M3 Expressive: subtle scale up */
  box-shadow: var(--md-sys-elevation-level2);
}

.play-button:hover .button-ripple {
  opacity: var(--md-sys-state-hover-opacity); /* M3 Expressive: 0.12 */
}

.play-button:hover .button-icon {
  transform: scale(1.1) rotate(-5deg); /* M3 Expressive: playful rotation */
}

/* Focus state */
.play-button:focus-visible {
  outline: 3px solid var(--md-sys-color-primary);
  outline-offset: 2px;
  transform: scale(1.05);
}

.play-button:focus-visible .button-ripple {
  opacity: var(--md-sys-state-focus-opacity); /* M3 Expressive: 0.16 */
}

/* Active/Pressed state */
.play-button:active {
  transform: scale(0.95); /* M3 Expressive: press down effect */
  box-shadow: var(--md-sys-elevation-level0);
}

.play-button:active .button-ripple {
  opacity: var(--md-sys-state-pressed-opacity); /* M3 Expressive: 0.16 */
}

.play-button:active .button-icon {
  transform: scale(0.9); /* M3 Expressive: icon shrinks on press */
}

/* Disabled state */
.play-button:disabled {
  cursor: not-allowed;
  pointer-events: none;
  color: var(--md-sys-color-on-surface);
  opacity: 0.38;
}

.play-button:disabled .button-background {
  background: var(--md-sys-color-on-surface);
  opacity: 0.12;
}
</style>
