<template>
  <span class="ripple-container" @mousedown="createRipple">
    <slot></slot>
  </span>
</template>

<script setup lang="ts">
const createRipple = (event: MouseEvent): void => {
  const button = event.currentTarget as HTMLElement
  if (!button) return

  const ripple = document.createElement('span')
  const diameter = Math.max(button.clientWidth, button.clientHeight)
  const radius = diameter / 2

  const rect = button.getBoundingClientRect()
  ripple.style.width = ripple.style.height = `${diameter}px`
  ripple.style.left = `${event.clientX - rect.left - radius}px`
  ripple.style.top = `${event.clientY - rect.top - radius}px`
  ripple.classList.add('ripple')

  const oldRipple = button.querySelector('.ripple')
  if (oldRipple) {
    oldRipple.remove()
  }

  button.appendChild(ripple)
}
</script>

<style scoped>
.ripple-container {
  position: relative;
  overflow: hidden;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

/* M3 Expressive Ripple */
.ripple-container :deep(.ripple) {
  position: absolute;
  border-radius: 50%;
  background: var(--md-sys-color-on-surface);
  opacity: 0.35; /* Expressive: более видимый */
  transform: scale(0);
  animation: ripple-animation var(--md-sys-motion-duration-long1) var(--md-sys-motion-easing-emphasized-decelerate); /* Expressive: 600ms, плавная */
  pointer-events: none;
}

@keyframes ripple-animation {
  0% {
    transform: scale(0);
    opacity: 0.35;
  }
  50% {
    opacity: 0.2; /* Expressive: пик в середине */
  }
  100% {
    transform: scale(5); /* Expressive: больший размер */
    opacity: 0;
  }
}
</style>
