<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal">
      <h2>Файлы в торренте: {{ torrent.name }}</h2>

      <div v-if="loading" class="loading">Загрузка списка файлов...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else class="files-list">
        <div
          v-for="file in files"
          :key="file.index"
          class="file-item"
          @click="playFile(file.index)"
        >
          <div class="file-icon">📄</div>
          <div class="file-info">
            <div class="file-name">{{ file.path }}</div>
            <div class="file-size">{{ file.sizeStr }}</div>
          </div>
          <button class="btn-play-file">▶</button>
        </div>
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
    await PlayTorrentFile(props.torrent.hash, fileIndex)
    emit('play', fileIndex)
  } catch (err) {
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
  gap: 8px;
  margin-bottom: 24px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: var(--md-sys-color-surface-container-high);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.file-item:hover {
  background: var(--md-sys-color-surface-container-highest);
  transform: translateX(4px);
}

.file-icon {
  font-size: 32px;
  flex-shrink: 0;
}

.file-info {
  flex: 1;
  min-width: 0;
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

.btn-play-file {
  padding: 8px 16px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  border: none;
  border-radius: 20px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.btn-play-file:hover {
  transform: scale(1.1);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
}

.btn-close {
  padding: 12px 24px;
  background: var(--md-sys-color-surface-container-highest);
  color: var(--md-sys-color-on-surface);
  border: none;
  border-radius: 24px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.btn-close:hover {
  transform: scale(1.05);
}
</style>
