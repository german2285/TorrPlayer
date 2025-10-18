<template>
  <div class="app">
    <div class="container">
      <!-- Search Bar -->
      <div class="search-container">
        <div class="search-bar">
          <img :src="searchIcon" class="search-icon" alt="Search" />
          <input
            type="text"
            class="search-input"
            placeholder="Поиск торрентов..."
            v-model="searchQuery"
          >
          <div class="action-icons">
            <RippleEffect>
              <button class="icon-btn" title="Добавить торрент" @click="showAddTorrent = true">
                <img :src="addIcon" class="icon" alt="Add" />
              </button>
            </RippleEffect>
            <RippleEffect>
              <button class="icon-btn" title="Настройки" @click="showSettings = true">
                <img :src="settingsIcon" class="icon" alt="Settings" />
              </button>
            </RippleEffect>
          </div>
        </div>
      </div>

      <!-- Torrents List -->
      <TransitionGroup name="movie-list" tag="div" class="movies-list">
        <MovieCardWithLoading
          v-for="(torrent, index) in filteredTorrents"
          :key="torrent.hash"
          :torrent="torrent"
          :style="{ '--animation-order': index }"
          @click="onTorrentClick(torrent)"
          @play="onPlayClick(torrent)"
          @remove="onRemoveClick(torrent)"
        />
      </TransitionGroup>

      <!-- Empty State -->
      <div v-if="filteredTorrents.length === 0 && !loading" class="empty-state">
        <svg width="64" height="64" fill="currentColor" viewBox="0 0 16 16">
          <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
          <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
        </svg>
        <p>Нет торрентов</p>
        <p class="empty-hint">Нажмите "+" чтобы добавить торрент</p>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Загрузка...</p>
      </div>
    </div>

    <!-- Add Torrent Dialog -->
    <AddTorrentModal
      v-if="showAddTorrent"
      @close="showAddTorrent = false"
      @added="onTorrentAdded"
    />

    <!-- File List Dialog -->
    <FileListModal
      v-if="showFileList && selectedTorrent"
      :torrent="selectedTorrent"
      @close="showFileList = false"
      @play="onFileSelected"
    />

    <!-- Settings Dialog -->
    <Settings v-if="showSettings" @close="showSettings = false" />

    <!-- Background Music -->
    <audio
      ref="bgAudio"
      loop
      preload="auto"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import type { Ref } from 'vue'
import MovieCardWithLoading from './components/MovieCardWithLoading.vue'
import RippleEffect from './components/RippleEffect.vue'
import AddTorrentModal from './components/AddTorrentModal.vue'
import FileListModal from './components/FileListModal.vue'
import Settings from './components/Settings.vue'
import { GetTorrents, RemoveTorrent } from '../wailsjs/go/app/App'
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime'
import type { Torrent } from './types'
import searchIcon from './assets/icons/search.svg'
import addIcon from './assets/icons/add.svg'
import settingsIcon from './assets/icons/settings.svg'
import backgroundMusicUrl from './assets/audio/background-music.mp3'
import { loadSavedTheme } from './utils/themeUtils'

interface TorrentMetadataLoadedEvent {
  hash: string
  title: string
  peers: number
  seeders: number
  fileCount: number
  totalSize: number
  sizeStr: string
  loaded: boolean
}

const torrents: Ref<Torrent[]> = ref([])
const searchQuery: Ref<string> = ref('')
const loading: Ref<boolean> = ref(true)
const showAddTorrent: Ref<boolean> = ref(false)
const showFileList: Ref<boolean> = ref(false)
const showSettings: Ref<boolean> = ref(false)
const selectedTorrent: Ref<Torrent | null> = ref(null)

// Background music
const bgAudio: Ref<HTMLAudioElement | null> = ref(null)
const bgVolume: Ref<number> = ref(0.3) // Default volume 30% (0-1 range for audio element)

// Filter torrents by search query
const filteredTorrents = computed(() => {
  if (!searchQuery.value.trim()) {
    return torrents.value
  }
  const query = searchQuery.value.toLowerCase()
  return torrents.value.filter(t =>
    t.name.toLowerCase().includes(query) ||
    t.title.toLowerCase().includes(query)
  )
})

// Load torrents from backend
const loadTorrents = async () => {
  try {
    loading.value = true
    torrents.value = await GetTorrents()
  } catch (error) {
    console.error('Failed to load torrents:', error)
  } finally {
    loading.value = false
  }
}

// Handle torrent click - show file list
const onTorrentClick = (torrent: Torrent): void => {
  console.log('Torrent clicked:', torrent)
  selectedTorrent.value = torrent
  showFileList.value = true
}

// Handle play button - show file list
const onPlayClick = (torrent: Torrent): void => {
  console.log('Play button clicked for torrent:', torrent)
  selectedTorrent.value = torrent
  showFileList.value = true
}

// Handle file selected from file list
const onFileSelected = async (fileIndex: number): Promise<void> => {
  showFileList.value = false
  // FileListModal will handle the playback
}

// Handle torrent removal
const onRemoveClick = async (torrent: Torrent): Promise<void> => {
  if (!confirm(`Удалить торрент "${torrent.name}"?`)) {
    return
  }

  try {
    await RemoveTorrent(torrent.hash)
    await loadTorrents()
  } catch (error) {
    console.error('Failed to remove torrent:', error)
    alert('Ошибка удаления торрента: ' + error)
  }
}

// Handle new torrent added
const onTorrentAdded = async (): Promise<void> => {
  showAddTorrent.value = false
  await loadTorrents()
}

// Handle metadata loaded event from backend
const onMetadataLoaded = (event: TorrentMetadataLoadedEvent) => {
  console.log('Metadata loaded event:', event)

  // Find and update the torrent in the list
  const torrentIndex = torrents.value.findIndex(t => t.hash === event.hash)
  if (torrentIndex !== -1) {
    // Update torrent with new metadata
    torrents.value[torrentIndex] = {
      ...torrents.value[torrentIndex],
      fileCount: event.fileCount,
      size: event.totalSize,
      sizeStr: event.sizeStr,
      peers: event.peers,
      seeders: event.seeders,
      loadingMeta: false // Metadata loaded, stop showing loading indicators
    }
    console.log('Updated torrent with metadata:', torrents.value[torrentIndex])
  } else {
    // Torrent not in list yet, reload the list
    console.log('Torrent not found in list, reloading...')
    loadTorrents()
  }
}

// Auto-refresh torrents every 5 seconds
let refreshInterval: number | null = null

// Watch volume changes
watch(bgVolume, (newVolume) => {
  if (!bgAudio.value) return

  bgAudio.value.volume = newVolume

  if (newVolume === 0) {
    bgAudio.value.pause()
  } else if (bgAudio.value.paused) {
    bgAudio.value.play().catch(err => console.error('Failed to play background music:', err))
  }
})

onMounted(async () => {
  // Load saved theme colors BEFORE rendering UI
  await loadSavedTheme()

  await loadTorrents()

  // Listen to metadata loaded events
  EventsOn('torrent:metadataLoaded', onMetadataLoaded)

  // Auto-refresh every 5 seconds (for peers/seeders updates)
  refreshInterval = window.setInterval(async () => {
    try {
      torrents.value = await GetTorrents()
    } catch (error) {
      console.error('Failed to refresh torrents:', error)
    }
  }, 5000)

  // Load background music volume from localStorage (0-100 range)
  const savedBgVolume = localStorage.getItem('bgMusicVolume')
  if (savedBgVolume) {
    bgVolume.value = parseInt(savedBgVolume, 10) / 100 // Convert 0-100 to 0-1
  }

  // Listen to background music volume changes from Settings
  const handleVolumeChange = (event: CustomEvent) => {
    bgVolume.value = event.detail / 100 // Convert 0-100 to 0-1
  }
  window.addEventListener('bgMusicVolumeChanged', handleVolumeChange as EventListener)

  // Initialize background music
  if (bgAudio.value) {
    bgAudio.value.src = backgroundMusicUrl
    bgAudio.value.volume = bgVolume.value
    if (bgVolume.value > 0) {
      bgAudio.value.play().catch(err => console.error('Failed to play background music:', err))
    }
  }
})

// Cleanup on unmount
onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }

  // Unsubscribe from events
  EventsOff('torrent:metadataLoaded')
})
</script>

<style scoped>
.app {
  font-family: var(--md-sys-typescale-body-large-font);
  background: var(--md-sys-color-background);
  color: var(--md-sys-color-on-background);
  padding: 24px;
  min-height: 100vh;
  display: flex;
  justify-content: center;
}

.container {
  width: 100%;
  max-width: 1200px;
}

.search-container {
  position: relative;
  margin-bottom: 24px;
}

.search-bar {
  background: var(--md-sys-color-surface-container-high);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: var(--md-sys-shape-corner-extra-large);
  padding: 12px 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: var(--md-sys-elevation-level1);
  transition: all var(--md-sys-motion-duration-medium2) var(--md-sys-motion-easing-emphasized);
  position: relative;
  z-index: 101;
  will-change: transform, border-radius, box-shadow;
}

.search-bar:hover {
  transform: scale(1.01);
  box-shadow: var(--md-sys-elevation-level2);
}

.search-bar:focus-within {
  border-color: var(--md-sys-color-primary);
  box-shadow: var(--md-sys-elevation-level3);
  transform: scale(1.02);
}

.search-icon {
  width: 24px;
  height: 24px;
  opacity: 0.6;
  flex-shrink: 0;
  color: var(--md-sys-color-on-surface-variant);
  transition: all var(--md-sys-motion-duration-medium2) var(--md-sys-motion-easing-emphasized-decelerate);
}

.search-bar:focus-within .search-icon {
  opacity: 1;
  color: var(--md-sys-color-primary);
  transform: scale(1.15);
}

.icon {
  width: 24px;
  height: 24px;
  opacity: 0.8;
  color: var(--md-sys-color-on-surface-variant);
  transition: all var(--md-sys-motion-duration-medium2) var(--md-sys-motion-easing-emphasized);
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  font-family: var(--md-sys-typescale-body-large-font);
  font-size: var(--md-sys-typescale-body-large-size);
  color: var(--md-sys-color-on-surface);
  background: transparent;
}

.search-input::placeholder {
  color: var(--md-sys-color-on-surface-variant);
  opacity: 0.6;
}

.action-icons {
  display: flex;
  gap: 8px;
}

.icon-btn {
  min-width: 48px;
  min-height: 48px;
  padding: 12px;
  border-radius: var(--md-sys-shape-corner-full);
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--md-sys-motion-duration-medium1) var(--md-sys-motion-easing-standard);
  color: var(--md-sys-color-on-surface-variant);
  flex-shrink: 0;
  position: relative;
  will-change: transform, border-radius;
}

.icon-btn::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: var(--md-sys-color-on-surface);
  opacity: 0;
  transition: opacity var(--md-sys-motion-duration-short3) var(--md-sys-motion-easing-standard);
}

.icon-btn:hover {
  color: var(--md-sys-color-primary);
  transform: scale(1.08);
  border-radius: var(--md-sys-shape-corner-large);
}

.icon-btn:hover::before {
  opacity: 0.12;
}

.icon-btn:active {
  border-radius: var(--md-sys-shape-corner-medium);
  transform: scale(0.95);
}

.icon-btn:hover .icon {
  opacity: 1;
  transform: scale(1.2) rotate(8deg);
}

.movies-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.movie-list-enter-active {
  animation: movie-enter 0.6s ease-out;
  animation-delay: calc(var(--animation-order) * 0.08s);
}

.movie-list-leave-active {
  animation: movie-leave 0.4s ease-in;
}

.movie-list-move {
  transition: transform 0.6s ease;
}

@keyframes movie-enter {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes movie-leave {
  from {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
  to {
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 96px 32px;
  color: var(--md-sys-color-on-surface-variant);
  text-align: center;
  background: var(--md-sys-color-surface-container);
  border-radius: var(--md-sys-shape-corner-extra-large);
  animation: emptyStateEnter 0.8s ease-out;
}

@keyframes emptyStateEnter {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.empty-state svg {
  margin-bottom: 32px;
  opacity: 0.5;
  width: 80px;
  height: 80px;
  animation: iconFloat 3s ease-in-out infinite;
}

@keyframes iconFloat {
  0%, 100% {
    transform: translateY(0) scale(1);
  }
  50% {
    transform: translateY(-8px) scale(1.05);
  }
}

.empty-state p {
  font-family: var(--md-sys-typescale-title-large-font);
  font-size: var(--md-sys-typescale-title-large-size);
  font-weight: var(--md-sys-typescale-title-large-weight);
}

.empty-hint {
  font-family: var(--md-sys-typescale-body-medium-font);
  font-size: var(--md-sys-typescale-body-medium-size);
  opacity: 0.7;
  margin-top: 8px;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: var(--md-sys-color-on-surface-variant);
  background: var(--md-sys-color-surface-container);
  border-radius: var(--md-sys-shape-corner-extra-large);
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid var(--md-sys-color-surface-container-highest);
  border-top-color: var(--md-sys-color-primary);
  border-right-color: var(--md-sys-color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 24px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-state p {
  font-family: var(--md-sys-typescale-body-large-font);
  font-size: var(--md-sys-typescale-body-large-size);
}
</style>
