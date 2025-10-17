<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal">
      <h2>Добавить торрент</h2>

      <div class="input-group">
        <label>Magnet-ссылка, путь к .torrent или хеш:</label>
        <textarea v-model="input" placeholder="magnet:?xt=urn:btih:...&#10;или C:\path\to\file.torrent&#10;или ABC123..." rows="4"></textarea>
      </div>

      <!-- Drag & Drop Zone -->
      <div
        class="drop-zone"
        :class="{ 'drop-zone-active': isDragging }"
        @dragover.prevent="handleDragOver"
        @dragleave.prevent="handleDragLeave"
        @drop.prevent="handleDrop"
      >
        <div class="drop-zone-content">
          <svg class="drop-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M7 10L12 15L17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M12 15V3" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M20 21H4C3.44772 21 3 20.5523 3 20V10C3 9.44772 3.44772 9 4 9H8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M16 9H20C20.5523 9 21 9.44772 21 10V20C21 20.5523 20.5523 21 20 21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <p class="drop-text">Перетащите .torrent файл сюда</p>
          <p class="drop-subtext">или нажмите для выбора файла</p>
          <input
            type="file"
            ref="fileInput"
            accept=".torrent"
            @change="handleFileSelect"
            style="display: none"
          />
          <button @click="openFilePicker" class="btn-browse" type="button">Выбрать файл</button>
        </div>
      </div>

      <div v-if="error" class="error">{{ error }}</div>
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <div class="loading-text">
          <p>{{ loadingMessage }}</p>
          <p class="loading-hint">Метаданные будут загружены в фоне</p>
        </div>
      </div>

      <div class="modal-actions">
        <button @click="$emit('close')" class="btn-cancel">Отмена</button>
        <button @click="addTorrent" class="btn-add" :disabled="!input.trim() || loading">Добавить</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { AddTorrent } from '../../wailsjs/go/app/App'

const input = ref('')
const loading = ref(false)
const error = ref('')
const loadingMessage = ref('Добавление торрента...')
const isDragging = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'added'): void
}>()

const addTorrent = async () => {
  if (!input.value.trim()) return

  try {
    loading.value = true
    error.value = ''
    loadingMessage.value = 'Добавление торрента...'

    await AddTorrent(input.value.trim())
    emit('added')
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err)
  } finally {
    loading.value = false
  }
}

const handleDragOver = (e: DragEvent) => {
  isDragging.value = true
}

const handleDragLeave = (e: DragEvent) => {
  isDragging.value = false
}

const handleDrop = async (e: DragEvent) => {
  isDragging.value = false

  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    await processFile(files[0])
  }
}

const handleFileSelect = async (e: Event) => {
  const target = e.target as HTMLInputElement
  const files = target.files
  if (files && files.length > 0) {
    await processFile(files[0])
  }
}

const openFilePicker = () => {
  fileInput.value?.click()
}

const processFile = async (file: File) => {
  if (!file.name.endsWith('.torrent')) {
    error.value = 'Пожалуйста, выберите файл .torrent'
    return
  }

  try {
    loading.value = true
    error.value = ''
    loadingMessage.value = 'Чтение файла...'

    // Read file as base64
    const reader = new FileReader()
    reader.onload = async (e) => {
      const result = e.target?.result as string
      // Convert to data URL format that backend expects
      const base64Data = result.split(',')[1]
      const dataURL = `data:application/x-bittorrent;base64,${base64Data}`

      loadingMessage.value = 'Добавление торрента...'

      try {
        await AddTorrent(dataURL)
        emit('added')
      } catch (err) {
        error.value = err instanceof Error ? err.message : String(err)
      } finally {
        loading.value = false
      }
    }
    reader.onerror = () => {
      error.value = 'Ошибка чтения файла'
      loading.value = false
    }
    reader.readAsDataURL(file)
  } catch (err) {
    error.value = err instanceof Error ? err.message : String(err)
    loading.value = false
  }
}
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
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal {
  background: var(--md-sys-color-surface);
  border-radius: 24px;
  padding: 32px;
  width: 90%;
  max-width: 600px;
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

h2 {
  margin: 0 0 24px;
  color: var(--md-sys-color-on-surface);
}

.input-group {
  margin-bottom: 24px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: var(--md-sys-color-on-surface-variant);
}

textarea {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid var(--md-sys-color-outline);
  border-radius: 12px;
  background: var(--md-sys-color-surface-container-high);
  color: var(--md-sys-color-on-surface);
  font-family: monospace;
  font-size: 14px;
  resize: vertical;
  transition: border-color 0.2s ease;
}

textarea:focus {
  outline: none;
  border-color: var(--md-sys-color-primary);
}

.error {
  padding: 12px;
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
  border-radius: 8px;
  margin-bottom: 16px;
}

.loading {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: var(--md-sys-color-primary-container);
  color: var(--md-sys-color-on-primary-container);
  border-radius: 12px;
  margin-bottom: 16px;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--md-sys-color-primary);
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  flex-shrink: 0;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  flex: 1;
  text-align: left;
}

.loading-text p {
  margin: 0;
  font-weight: 600;
}

.loading-hint {
  font-size: 13px;
  opacity: 0.8;
  margin-top: 4px !important;
  font-weight: 400 !important;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

button {
  padding: 12px 24px;
  border: none;
  border-radius: 24px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-cancel {
  background: var(--md-sys-color-surface-container-highest);
  color: var(--md-sys-color-on-surface);
}

.btn-cancel:hover {
  transform: scale(1.05);
}

.btn-add {
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
}

.btn-add:hover:not(:disabled) {
  transform: scale(1.05);
}

.btn-add:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Drag & Drop Zone */
.drop-zone {
  border: 2px dashed var(--md-sys-color-outline);
  border-radius: 16px;
  padding: 32px;
  text-align: center;
  background: var(--md-sys-color-surface-container-low);
  transition: all 0.3s ease;
  margin-bottom: 24px;
  cursor: pointer;
}

.drop-zone:hover {
  border-color: var(--md-sys-color-primary);
  background: var(--md-sys-color-surface-container);
}

.drop-zone-active {
  border-color: var(--md-sys-color-primary);
  background: var(--md-sys-color-primary-container);
  transform: scale(1.02);
}

.drop-zone-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.drop-icon {
  width: 48px;
  height: 48px;
  color: var(--md-sys-color-primary);
  opacity: 0.7;
  transition: all 0.3s ease;
}

.drop-zone:hover .drop-icon,
.drop-zone-active .drop-icon {
  opacity: 1;
  transform: translateY(-4px);
}

.drop-text {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.drop-subtext {
  margin: 0;
  font-size: 14px;
  color: var(--md-sys-color-on-surface-variant);
  opacity: 0.8;
}

.btn-browse {
  margin-top: 8px;
  padding: 10px 20px;
  background: var(--md-sys-color-primary-container);
  color: var(--md-sys-color-on-primary-container);
  border: none;
  border-radius: 20px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-browse:hover {
  transform: scale(1.05);
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
}
</style>
