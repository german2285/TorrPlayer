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
            placeholder="Поиск торрентов... (Начните с '!' для поиска по RuTracker)"
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
      <div v-if="filteredTorrents.length === 0 && !loading && !rutrackerLoading" class="empty-state">
        <svg width="64" height="64" fill="currentColor" viewBox="0 0 16 16">
          <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
          <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
        </svg>
        <p v-if="isRutrackerSearch">Ничего не найдено на RuTracker</p>
        <p v-else>Нет торрентов</p>
        <p class="empty-hint" v-if="!isRutrackerSearch">Нажмите "+" чтобы добавить торрент</p>
        <p class="empty-hint" v-else>Попробуйте изменить запрос</p>
      </div>

      <!-- Loading State -->
      <div v-if="loading || rutrackerLoading" class="loading-state">
        <div class="spinner"></div>
        <p v-if="rutrackerLoading">Поиск по RuTracker...</p>
        <p v-else>Загрузка...</p>
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
import { GetTorrents, RemoveTorrent, SearchRuTracker, GetRutrackerMagnetLink, AddTorrent } from '../wailsjs/go/app/App'
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime'
import type { Torrent, RutrackerTorrent } from './types'
import searchIcon from './assets/icons/search.svg'
import addIcon from './assets/icons/add.svg'
import settingsIcon from './assets/icons/settings.svg'
import backgroundMusicUrl from './assets/audio/background-music.mp3'
import { loadSavedTheme } from './utils/themeUtils'

const torrents: Ref<Torrent[]> = ref([])
const searchQuery: Ref<string> = ref('')
const loading: Ref<boolean> = ref(true)
const showAddTorrent: Ref<boolean> = ref(false)
const showFileList: Ref<boolean> = ref(false)
const showSettings: Ref<boolean> = ref(false)
const selectedTorrent: Ref<Torrent | null> = ref(null)

// RuTracker search
const rutrackerResults: Ref<RutrackerTorrent[]> = ref([])
const isRutrackerSearch: Ref<boolean> = ref(false)
const rutrackerLoading: Ref<boolean> = ref(false)

// Background music
const bgAudio: Ref<HTMLAudioElement | null> = ref(null)
const bgVolume: Ref<number> = ref(0.3) // Default volume 30% (0-1 range for audio element)

// Filter torrents by search query (or show RuTracker results)
const filteredTorrents = computed(() => {
  // If RuTracker search is active, convert RutrackerTorrent to Torrent format
  if (isRutrackerSearch.value && rutrackerResults.value.length > 0) {
    return rutrackerResults.value.map(rt => ({
      hash: rt.topicId,
      name: rt.title,
      title: rt.title,
      size: 0,
      sizeStr: rt.size,
      status: `Seeds: ${rt.seeds} | Leeches: ${rt.leeches}`,
      progress: 0,
      downSpeed: 0,
      upSpeed: 0,
      downSpeedStr: '',
      upSpeedStr: '',
      peers: rt.leeches,
      seeders: rt.seeds,
      category: rt.category,
      poster: '',
      timestamp: 0,
      // Mark as RuTracker result for special handling
      _isRutrackerResult: true
    } as Torrent & { _isRutrackerResult?: boolean }))
  }

  // Local search
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

// Handle torrent click - show file list or add RuTracker torrent
const onTorrentClick = async (torrent: Torrent & { _isRutrackerResult?: boolean }): Promise<void> => {
  console.log('Torrent clicked:', torrent)

  // If it's a RuTracker result, add it as torrent
  if (torrent._isRutrackerResult) {
    await handleRutrackerClick(torrent)
    return
  }

  // Otherwise, show file list for local torrent
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

// Search RuTracker
const searchRuTracker = async (query: string): Promise<void> => {
  try {
    rutrackerLoading.value = true
    const results = await SearchRuTracker(query)
    rutrackerResults.value = results || []
    isRutrackerSearch.value = true
  } catch (error) {
    console.error('RuTracker search failed:', error)
    alert('Ошибка поиска по RuTracker: ' + error)
    rutrackerResults.value = []
    isRutrackerSearch.value = false
  } finally {
    rutrackerLoading.value = false
  }
}

// Handle RuTracker torrent click
const handleRutrackerClick = async (torrent: Torrent & { _isRutrackerResult?: boolean }): Promise<void> => {
  if (!torrent._isRutrackerResult) {
    return
  }

  try {
    loading.value = true

    // Get magnet link
    const magnetLink = await GetRutrackerMagnetLink(torrent.hash)

    if (!magnetLink) {
      throw new Error('Не удалось получить magnet-ссылку')
    }

    // Add torrent
    await AddTorrent(magnetLink)

    // Reload torrents list
    await loadTorrents()

    // Clear search to show newly added torrent
    searchQuery.value = ''
    rutrackerResults.value = []
    isRutrackerSearch.value = false

    alert('Торрент успешно добавлен!')
  } catch (error) {
    console.error('Failed to add RuTracker torrent:', error)
    alert('Ошибка добавления торрента: ' + error)
  } finally {
    loading.value = false
  }
}

// Watch search query for RuTracker search
watch(searchQuery, async (newQuery) => {
  if (newQuery.startsWith('!')) {
    const query = newQuery.substring(1).trim()
    if (query.length >= 2) {
      await searchRuTracker(query)
    } else {
      rutrackerResults.value = []
      isRutrackerSearch.value = false
    }
  } else {
    rutrackerResults.value = []
    isRutrackerSearch.value = false
  }
})

// Handle video playback starting - pause background music and cleanup
const onVideoPlaybackStarting = () => {
  console.log('Video playback starting - cleaning up resources')

  // Pause and cleanup background music
  if (bgAudio.value) {
    bgAudio.value.pause()
    bgAudio.value.src = '' // Release audio resource
  }

  console.log('Resources cleaned up, ready for window reload after playback')
}

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

  // Listen to video playback starting event (page will reload after playback)
  EventsOn('video:playbackStarting', onVideoPlaybackStarting)

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
  // Unsubscribe from events
  EventsOff('video:playbackStarting')
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
