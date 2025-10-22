<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal">
      <h2>Файлы в торренте: {{ torrent.name }}</h2>

      <div v-if="loading" class="loading">Загрузка списка файлов...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else class="files-list">
        <template v-for="(file, index) in files" :key="file.index">
          <div
            class="file-item"
            @click="playFile(file.index)"
          >
            <div class="file-icon">📄</div>
            <div class="file-info">
              <div class="file-name">{{ file.path }}</div>
              <div class="file-size">{{ file.sizeStr }}</div>
            </div>
            <!-- Кнопка Play с Loading Indicator внутри -->
            <button class="btn-play-file" :aria-label="playingFileIndex === file.index ? 'Loading...' : 'Play file'">
              <!-- Loading Indicator (Contained) при воспроизведении -->
              <LoadingIndicatorAdvanced
                v-if="playingFileIndex === file.index"
                animation="morph"
                size="extra-small"
                color-scheme="on-primary"
                :contained="true"
                :speed="1.2"
                aria-hidden="true"
              />
              <!-- Иконка Play когда не воспроизводится -->
              <img v-else :src="playIcon" class="play-icon" alt="Play" />
            </button>
          </div>
          <!-- Material Design Divider -->
          <div v-if="index < files.length - 1" class="divider"></div>
        </template>
      </div>

      <div class="modal-actions">
        <button @click="$emit('close')" class="btn-close">Закрыть</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Torrent, TorrentFile } from '../types'
import { GetTorrentFiles, PlayTorrentFile } from '../../wailsjs/go/app/App'
import playIcon from '../assets/icons/play_arrow.svg'
import LoadingIndicatorAdvanced from './LoadingIndicatorAdvanced.vue'

const props = defineProps<{
  torrent: Torrent
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'play', fileIndex: number): void
}>()

const files = ref<TorrentFile[]>([])
const loading = ref(true)
const error = ref('')
const playingFileIndex = ref<number | null>(null)

const loadFiles = async () => {
  try {
    loading.value = true
    error.value = ''
    files.value = await GetTorrentFiles(props.torrent.hash)
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err)
  } finally {
    loading.value = false
  }
}

const playFile = async (fileIndex: number) => {
  try {
    playingFileIndex.value = fileIndex
    await PlayTorrentFile(props.torrent.hash, fileIndex)
    emit('play', fileIndex)
    // Небольшая задержка перед закрытием для показа завершения
    await new Promise(resolve => setTimeout(resolve, 300))
    emit('close')
  } catch (err) {
    playingFileIndex.value = null
    alert('Ошибка воспроизведения: ' + (err instanceof Error ? err.message : String(err)))
  }
}

onMounted(loadFiles)
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: var(--md-sys-color-surface);
  border-radius: 24px;
  padding: 32px;
  width: 90%;
  max-width: 800px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

h2 {
  margin: 0 0 24px;
  color: var(--md-sys-color-on-surface);
  font-size: 20px;
}

.loading,
.error {
  padding: 24px;
  text-align: center;
  border-radius: 12px;
}

.loading {
  background: var(--md-sys-color-primary-container);
  color: var(--md-sys-color-on-primary-container);
}

.error {
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
}

.files-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  margin-bottom: 24px;
  padding-right: 16px; /* Отступ для scrollbar */
}

.file-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 0;
  cursor: pointer;
  transition: all 0.2s ease;
}

/* Material Design State Layer */
.file-item::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--md-sys-color-on-surface);
  opacity: 0;
  transition: opacity 0.2s ease;
  pointer-events: none;
  border-radius: 4px;
}

.file-item:hover::before {
  opacity: 0.08;
}

/* Material Design Divider (Inset) */
.divider {
  height: 1px;
  background: var(--md-sys-color-outline-variant, rgba(0, 0, 0, 0.12));
  margin: 0 16px;
  flex-shrink: 0;
  opacity: 1;
}

.file-icon {
  font-size: 32px;
  flex-shrink: 0;
  position: relative;
  z-index: 1;
}

.file-info {
  flex: 1;
  min-width: 0;
  position: relative;
  z-index: 1;
}

.file-name {
  font-weight: 600;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  font-size: 14px;
  opacity: 0.7;
}

/* === M3 Expressive Icon Button Styles === */

/* M3 Filled Icon Button (Small, Round) - Play file */
.btn-play-file {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  padding: 0;
  background: var(--md-sys-color-primary);
  border: none;
  border-radius: 50%; /* Round shape */
  cursor: pointer;
  overflow: hidden;
  box-shadow: var(--md-sys-elevation-level1);
  transition: box-shadow var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
              transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
  flex-shrink: 0;
  z-index: 1;
}

/* M3 State Layer */
.btn-play-file::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--md-sys-color-on-primary);
  opacity: 0;
  border-radius: 50%;
  transition: opacity 0.15s ease;
  pointer-events: none;
}

.btn-play-file:hover {
  box-shadow: var(--md-sys-elevation-level2);
}

.btn-play-file:hover::before {
  opacity: var(--md-sys-state-hover-opacity); /* 0.08 */
}

.btn-play-file:active {
  box-shadow: var(--md-sys-elevation-level0);
  transform: scale(0.95); /* M3 Expressive press effect */
}

.btn-play-file:active::before {
  opacity: var(--md-sys-state-pressed-opacity); /* 0.12 */
}

/* Play icon */
.play-icon {
  width: 24px;
  height: 24px;
  position: relative;
  z-index: 1;
  filter: brightness(0) invert(1); /* Make icon white for on-primary color */
  pointer-events: none;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
}

/* M3 Text Button (Close) */
.btn-close {
  position: relative;
  padding: 10px 12px;
  background: transparent;
  color: var(--md-sys-color-primary);
  border: none;
  border-radius: var(--md-sys-shape-corner-full);
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.1px;
  cursor: pointer;
  overflow: hidden;
  transition: box-shadow var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
}

.btn-close::before {
  content: '';
  position: absolute;
  inset: 0;
  background: currentColor;
  opacity: 0;
  transition: opacity 0.15s ease;
  pointer-events: none;
}

.btn-close:hover::before {
  opacity: var(--md-sys-state-hover-opacity);
}

.btn-close:active::before {
  opacity: var(--md-sys-state-pressed-opacity);
}
</style>
