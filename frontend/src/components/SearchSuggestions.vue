<template>
  <Transition name="suggestions-fade">
    <div v-if="show && suggestions.length > 0" class="suggestions-dropdown">
      <div class="suggestions-list">
        <div
          v-for="(suggestion, index) in suggestions"
          :key="index"
          class="suggestion-item"
          :class="{ 'selected': selectedIndex === index }"
          @click="onSelect(suggestion)"
          @mouseenter="selectedIndex = index"
        >
          <div class="suggestion-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
              <path d="M15.5 14h-.79l-.28-.27A6.471 6.471 0 0 0 16 9.5 6.5 6.5 0 1 0 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
            </svg>
          </div>
          <div class="suggestion-content">
            <div class="suggestion-text">{{ suggestion.text }}</div>
            <div v-if="suggestion.subtitle" class="suggestion-subtitle">{{ suggestion.subtitle }}</div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Ref } from 'vue'

interface Suggestion {
  text: string
  subtitle?: string
}

interface Props {
  show: boolean
  query: string
}

const props = withDefaults(defineProps<Props>(), {
  show: false,
  query: ''
})

const emit = defineEmits<{
  (e: 'select', text: string): void
}>()

const selectedIndex: Ref<number> = ref(-1)
const suggestions: Ref<Suggestion[]> = ref([])

// Генерируем подсказки на основе запроса
watch(() => props.query, (newQuery: string) => {
  if (newQuery.trim().length > 0) {
    suggestions.value = generateSuggestions(newQuery)
  } else {
    suggestions.value = []
  }
})

const generateSuggestions = (query: string): Suggestion[] => {
  // Здесь будет логика получения подсказок от API
  // Пока просто генерируем примеры
  const mockSuggestions: Suggestion[] = [
    { text: query, subtitle: 'Искать фильмы' },
    { text: `${query} 2024`, subtitle: 'Новинки' },
    { text: `${query} hd`, subtitle: 'Высокое качество' },
    { text: `${query} сериал`, subtitle: 'Сериалы' },
    { text: `${query} боевик`, subtitle: 'Жанр: Боевик' }
  ]

  return mockSuggestions.slice(0, 5)
}

const onSelect = (suggestion: Suggestion): void => {
  emit('select', suggestion.text)
}
</script>

<style scoped>
.suggestions-dropdown {
  position: absolute;
  top: calc(100% - 1px);
  left: 0;
  right: 0;
  background: var(--md-sys-color-surface-container-high);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-top: none;
  border-radius: 0 0 var(--md-sys-shape-corner-extra-large) var(--md-sys-shape-corner-extra-large);
  box-shadow: var(--md-sys-elevation-level3);
  overflow: hidden;
  z-index: 100;
  max-height: 400px;
  overflow-y: auto;
}

.suggestions-list {
  padding: 8px 0;
}

.suggestion-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 24px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
}

.suggestion-item::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--md-sys-color-on-surface);
  opacity: 0;
  transition: opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.suggestion-item:hover::before,
.suggestion-item.selected::before {
  opacity: var(--md-sys-state-hover-opacity);
}

.suggestion-item:active::before {
  opacity: var(--md-sys-state-pressed-opacity);
}

.suggestion-icon {
  width: 20px;
  height: 20px;
  color: var(--md-sys-color-on-surface-variant);
  flex-shrink: 0;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  z-index: 1;
}

.suggestion-item:hover .suggestion-icon,
.suggestion-item.selected .suggestion-icon {
  color: var(--md-sys-color-primary);
  transform: scale(1.1);
}

.suggestion-content {
  flex: 1;
  min-width: 0;
  position: relative;
  z-index: 1;
}

.suggestion-text {
  font-family: var(--md-sys-typescale-body-large-font);
  font-size: var(--md-sys-typescale-body-large-size);
  color: var(--md-sys-color-on-surface);
  font-weight: 400;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.suggestion-item:hover .suggestion-text,
.suggestion-item.selected .suggestion-text {
  color: var(--md-sys-color-primary);
}

.suggestion-subtitle {
  font-family: var(--md-sys-typescale-body-small-font);
  font-size: 12px;
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Scrollbar styling */
.suggestions-dropdown::-webkit-scrollbar {
  width: 8px;
}

.suggestions-dropdown::-webkit-scrollbar-track {
  background: transparent;
}

.suggestions-dropdown::-webkit-scrollbar-thumb {
  background: var(--md-sys-color-outline-variant);
  border-radius: 4px;
}

.suggestions-dropdown::-webkit-scrollbar-thumb:hover {
  background: var(--md-sys-color-outline);
}

/* Animations - M3 Expressive */
.suggestions-fade-enter-active {
  animation: suggestions-enter var(--md-sys-motion-duration-medium3) var(--md-sys-motion-easing-emphasized-decelerate); /* Expressive: 500ms */
}

.suggestions-fade-leave-active {
  animation: suggestions-leave var(--md-sys-motion-duration-short4) var(--md-sys-motion-easing-emphasized-accelerate); /* Expressive: 250ms */
}

@keyframes suggestions-enter {
  from {
    opacity: 0;
    border-radius: var(--md-sys-shape-corner-medium); /* M3 Expressive: морфинг */
    transform: translateY(-12px) scale(0.96); /* Expressive */
  }
  to {
    opacity: 1;
    border-radius: 0 0 var(--md-sys-shape-corner-extra-large) var(--md-sys-shape-corner-extra-large);
    transform: translateY(0) scale(1);
  }
}

@keyframes suggestions-leave {
  from {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
  to {
    opacity: 0;
    transform: translateY(-8px) scale(0.98); /* Expressive */
  }
}
</style>
