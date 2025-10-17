<template>
  <div class="movie-card-container">
    <!-- Карточка фильма с Material 3 Expressive дизайном -->
    <div class="movie-card" @click="handleCardClick">
      <!-- Постер (слева) -->
      <div class="movie-poster">
        <img
          v-if="torrent.poster"
          :src="torrent.poster"
          :alt="torrent.title"
          @error="onImageError"
          class="poster-image"
        />
        <div v-else class="poster-placeholder">
          <span class="poster-icon">🎬</span>
        </div>
      </div>

      <!-- Контент (справа) -->
      <div class="movie-content">
        <!-- Название -->
        <h3 class="movie-title">{{ torrent.title || 'Untitled' }}</h3>

        <!-- Категория -->
        <div v-if="torrent.category" class="movie-category">
          <span class="category-chip">{{ formatCategory(torrent.category) }}</span>
        </div>

        <!-- Кнопка Play или Loading Indicator -->
        <div class="play-section">
          <div v-if="isMetadataLoading" class="play-loading">
            <CircularProgressIndicator
              :size="48"
              :thickness="4"
              :indeterminate="true"
              variant="retreat"
            />
            <span class="loading-text">Подключение...</span>
          </div>
          <div v-else class="action-buttons">
            <button
              class="play-button"
              @click.stop="handlePlayClick"
              :disabled="!hasMetadata"
            >
              <svg class="play-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M8 5v14l11-7z"/>
              </svg>
              <span>Воспроизвести</span>
            </button>
            <button
              class="remove-button"
              @click.stop="handleRemoveClick"
              title="Удалить торрент"
            >
              <svg class="remove-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/>
              </svg>
            </button>
          </div>
        </div>

        <!-- Метаданные с Progress Indicators -->
        <div class="metadata-section">
          <!-- Пиры -->
          <div class="metadata-item">
            <svg class="metadata-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/>
            </svg>
            <div class="metadata-value">
              <div v-if="isMetadataLoading" class="metadata-loading">
                <LinearProgressIndicator
                  :indeterminate="true"
                  :thickness="3"
                  animation-type="disjoint"
                />
                <span class="metadata-loading-text">загрузка...</span>
              </div>
              <span v-else class="metadata-text">
                {{ formatPeers(metadata.peers) }}
              </span>
            </div>
          </div>

          <!-- Количество файлов -->
          <div class="metadata-item">
            <svg class="metadata-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M20 6h-8l-2-2H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm-1 12H5c-.55 0-1-.45-1-1V9c0-.55.45-1 1-1h14c.55 0 1 .45 1 1v8c0 .55-.45 1-1 1z"/>
            </svg>
            <div class="metadata-value">
              <div v-if="isMetadataLoading" class="metadata-loading">
                <LinearProgressIndicator
                  :indeterminate="true"
                  :thickness="3"
                  animation-type="disjoint"
                />
                <span class="metadata-loading-text">загрузка...</span>
              </div>
              <span v-else class="metadata-text">
                {{ formatFileCount(metadata.fileCount) }}
              </span>
            </div>
          </div>

          <!-- Размер -->
          <div class="metadata-item">
            <svg class="metadata-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM9 17H7v-7h2v7zm4 0h-2V7h2v10zm4 0h-2v-4h2v4z"/>
            </svg>
            <div class="metadata-value">
              <div v-if="isMetadataLoading" class="metadata-loading">
                <LinearProgressIndicator
                  :indeterminate="true"
                  :thickness="3"
                  animation-type="disjoint"
                />
                <span class="metadata-loading-text">загрузка...</span>
              </div>
              <span v-else class="metadata-text">
                {{ formatSize(metadata.totalSize) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import CircularProgressIndicator from './CircularProgressIndicator.vue'
import LinearProgressIndicator from './LinearProgressIndicator.vue'
import type { Torrent } from '../types'

interface TorrentMetadata {
  peers: number
  fileCount: number
  totalSize: number
}

interface Props {
  torrent: Torrent
  initialLoading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  initialLoading: true
})

const emit = defineEmits<{
  (e: 'play', torrent: Torrent): void
  (e: 'click', torrent: Torrent): void
  (e: 'remove', torrent: Torrent): void
}>()

// Состояние загрузки метаданных - берем из пропса torrent
const isMetadataLoading = computed(() => props.torrent.loadingMeta ?? true)

// Метаданные берем напрямую из торрента
const metadata = computed(() => ({
  peers: props.torrent.peers ?? 0,
  fileCount: props.torrent.fileCount ?? 0,
  totalSize: props.torrent.size ?? 0
}))

const hasMetadata = computed(() => metadata.value.fileCount > 0)

// Форматирование данных
const formatCategory = (category: string): string => {
  const categories: Record<string, string> = {
    'movie': 'Фильм',
    'series': 'Сериал',
    'tv': 'ТВ-шоу',
    'documentary': 'Документальный',
    'anime': 'Аниме'
  }
  return categories[category] || category
}

const formatPeers = (count: number): string => {
  if (count === 0) return 'Нет пиров'
  if (count === 1) return '1 пир'
  if (count < 5) return `${count} пира`
  return `${count} пиров`
}

const formatFileCount = (count: number): string => {
  if (count === 0) return 'Нет файлов'
  if (count === 1) return '1 файл'
  if (count < 5) return `${count} файла`
  return `${count} файлов`
}

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

const onImageError = (e: Event): void => {
  const target = e.target as HTMLImageElement
  target.style.display = 'none'
}

const handleCardClick = (): void => {
  emit('click', props.torrent)
}

const handlePlayClick = (): void => {
  if (hasMetadata.value) {
    emit('play', props.torrent)
  }
}

const handleRemoveClick = (): void => {
  emit('remove', props.torrent)
}
</script>

<style scoped>
/* Material 3 Expressive - Карточка фильма */
.movie-card-container {
  width: 100%;
}

.movie-card {
  display: flex;
  gap: 24px;
  padding: 20px;
  background: var(--md-sys-color-surface-container, #F3EDF7);
  border: 1px solid var(--md-sys-color-outline-variant, #C4C7C5);
  border-radius: 28px; /* M3 Expressive: Extra Large */
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.2, 0, 0, 1); /* Emphasized easing */
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  position: relative;
  overflow: hidden;
}

.movie-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--md-sys-color-primary, #6750A4);
  opacity: 0;
  transition: opacity 0.3s ease;
  pointer-events: none;
}

.movie-card:hover {
  transform: translateY(-4px) scale(1.01);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  border-color: var(--md-sys-color-primary, #6750A4);
  border-radius: 24px; /* M3 Expressive: морфинг скругления */
}

.movie-card:hover::before {
  opacity: 0.05;
}

.movie-card:active {
  transform: translateY(-2px) scale(0.99);
  transition: all 0.15s cubic-bezier(0.4, 0, 1, 1);
}

/* Постер */
.movie-poster {
  width: 160px;
  height: 240px;
  border-radius: 16px; /* M3 Expressive */
  overflow: hidden;
  flex-shrink: 0;
  background: linear-gradient(135deg, #6750A4 0%, #7F67BE 100%);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.4s cubic-bezier(0.2, 0, 0, 1);
}

.movie-card:hover .movie-poster {
  transform: scale(1.05) rotate(1deg);
  border-radius: 12px;
}

.poster-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.poster-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #6750A4 0%, #7F67BE 100%);
}

.poster-icon {
  font-size: 64px;
  filter: brightness(1.2);
}

/* Контент */
.movie-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.movie-title {
  font-family: 'Roboto Flex', 'Roboto', sans-serif;
  font-size: 24px;
  font-weight: 600;
  line-height: 1.3;
  color: var(--md-sys-color-on-surface, #1C1B1F);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.movie-category {
  display: flex;
  gap: 8px;
}

.category-chip {
  padding: 6px 16px;
  background: var(--md-sys-color-secondary-container, #E8DEF8);
  color: var(--md-sys-color-on-secondary-container, #1D192B);
  border-radius: 16px; /* M3 Expressive */
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.2, 0, 0, 1);
}

.category-chip:hover {
  background: var(--md-sys-color-primary-container, #EADDFF);
  color: var(--md-sys-color-on-primary-container, #21005D);
  transform: scale(1.05);
}

/* Секция воспроизведения */
.play-section {
  margin-top: auto;
}

.play-loading {
  display: flex;
  align-items: center;
  gap: 16px;
}

.loading-text {
  font-size: 14px;
  color: var(--md-sys-color-on-surface-variant, #49454F);
  font-weight: 500;
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 12px;
}

.play-button {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 32px;
  background: var(--md-sys-color-primary, #6750A4);
  color: var(--md-sys-color-on-primary, #FFFFFF);
  border: none;
  border-radius: 24px; /* M3 Expressive: Large */
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.2, 0, 0, 1);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.play-button:hover:not(:disabled) {
  background: var(--md-sys-color-primary-container, #EADDFF);
  color: var(--md-sys-color-on-primary-container, #21005D);
  transform: scale(1.08) translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.25);
  border-radius: 28px;
}

.play-button:active:not(:disabled) {
  transform: scale(1.02);
  transition: all 0.15s ease;
}

.play-button:disabled {
  opacity: 0.38;
  cursor: not-allowed;
}

.play-icon {
  width: 24px;
  height: 24px;
}

.remove-button {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 48px;
  min-height: 48px;
  padding: 12px;
  background: transparent;
  color: var(--md-sys-color-error, #BA1A1A);
  border: 1px solid var(--md-sys-color-error, #BA1A1A);
  border-radius: var(--md-sys-shape-corner-full, 24px);
  cursor: pointer;
  transition: all var(--md-sys-motion-duration-medium1, 0.3s) var(--md-sys-motion-easing-emphasized, cubic-bezier(0.2, 0, 0, 1));
}

.remove-button:hover {
  background: var(--md-sys-color-error-container, #FFDAD6);
  color: var(--md-sys-color-on-error-container, #410002);
  transform: scale(1.08);
  border-color: transparent;
}

.remove-button:active {
  transform: scale(0.95);
  transition: all 0.15s ease;
}

.remove-icon {
  width: 24px;
  height: 24px;
}

/* Метаданные */
.metadata-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.metadata-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.metadata-icon {
  width: 20px;
  height: 20px;
  color: var(--md-sys-color-primary, #6750A4);
  flex-shrink: 0;
}

.metadata-value {
  flex: 1;
  min-width: 0;
}

.metadata-loading {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.metadata-loading-text {
  font-size: 12px;
  color: var(--md-sys-color-on-surface-variant, #49454F);
  font-style: italic;
}

.metadata-text {
  font-size: 14px;
  font-weight: 500;
  color: var(--md-sys-color-on-surface, #1C1B1F);
}

/* Адаптивность */
@media (max-width: 768px) {
  .movie-card {
    flex-direction: column;
    gap: 16px;
  }

  .movie-poster {
    width: 100%;
    height: 200px;
  }
}
</style>
