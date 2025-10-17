<template>
  <div class="torrent-card" @click="$emit('click', torrent)">
    <!-- Banner -->
    <div class="torrent-banner">
      <div class="banner-placeholder">🎬</div>
      <div class="banner-overlay"></div>
      <div class="banner-content">
        <div class="torrent-title">{{ torrent.name }}</div>
        <div class="torrent-size">{{ torrent.sizeStr }}</div>
      </div>
    </div>

    <!-- Metadata -->
    <div class="torrent-metadata">
      <div class="torrent-stats">
        <div class="stat-item">
          <span class="stat-label">↓</span>
          <span class="stat-value">{{ torrent.downSpeedStr }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">↑</span>
          <span class="stat-value">{{ torrent.upSpeedStr }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Пиры:</span>
          <span class="stat-value">{{ torrent.peers }} ({{ torrent.seeders }} сидов)</span>
        </div>
      </div>

      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: torrent.progress + '%' }"></div>
      </div>

      <div class="torrent-actions">
        <button class="btn-play" @click.stop="$emit('play', torrent)">
          ▶ Воспроизвести
        </button>
        <button class="btn-remove" @click.stop="$emit('remove', torrent)">
          🗑️ Удалить
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Torrent } from '../types'

defineProps<{
  torrent: Torrent
}>()

defineEmits<{
  (e: 'click', torrent: Torrent): void
  (e: 'play', torrent: Torrent): void
  (e: 'remove', torrent: Torrent): void
}>()
</script>

<style scoped>
.torrent-card {
  display: flex;
  background: var(--md-sys-color-surface-container);
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  height: 200px;
}

.torrent-card:hover {
  transform: scale(1.02);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}

.torrent-banner {
  position: relative;
  width: 300px;
  flex-shrink: 0;
}

.banner-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 64px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.banner-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0,0,0,0.8), transparent);
}

.banner-content {
  position: absolute;
  bottom: 16px;
  left: 16px;
  right: 16px;
  color: white;
  z-index: 1;
}

.torrent-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.torrent-size {
  font-size: 14px;
  opacity: 0.9;
}

.torrent-metadata {
  flex: 1;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.torrent-stats {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  gap: 8px;
  align-items: center;
  font-size: 14px;
}

.stat-label {
  font-weight: 600;
  opacity: 0.7;
}

.stat-value {
  color: var(--md-sys-color-primary);
}

.progress-bar {
  height: 8px;
  background: var(--md-sys-color-surface-container-highest);
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--md-sys-color-primary);
  transition: width 0.3s ease;
}

.torrent-actions {
  display: flex;
  gap: 12px;
  margin-top: auto;
}

.btn-play,
.btn-remove {
  padding: 12px 24px;
  border: none;
  border-radius: 24px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-play {
  flex: 1;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
}

.btn-play:hover {
  transform: scale(1.05);
}

.btn-remove {
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
}

.btn-remove:hover {
  transform: scale(1.05);
}
</style>
