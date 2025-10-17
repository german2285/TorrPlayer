<template>
  <div class="movie-card" @click="$emit('click', movie)">
    <!-- Banner with title overlay -->
    <div class="movie-banner">
      <img
        v-if="movie.poster"
        :src="movie.poster"
        :alt="movie.title"
        @error="onImageError"
      />
      <div v-else class="banner-placeholder">🎬</div>

      <!-- Gradient overlay -->
      <div class="banner-overlay"></div>

      <!-- Title on banner -->
      <div class="banner-content">
        <div class="movie-title">{{ movie.title }}</div>
        <div class="movie-year">{{ movie.year || 'N/A' }}</div>
      </div>
    </div>

    <!-- Metadata on the right -->
    <div class="movie-metadata">
      <div class="movie-header">
        <div v-if="movie.rating" class="movie-rating">
          <img :src="starIcon" class="rating-icon" alt="Rating" />
          {{ movie.rating }}
        </div>
      </div>

      <div v-if="movie.genres && movie.genres.length > 0" class="movie-genres">
        <span
          v-for="(genre, index) in movie.genres.slice(0, 3)"
          :key="index"
          class="genre-tag"
        >
          {{ genre }}
        </span>
      </div>

      <div v-if="movie.description" class="movie-description">
        {{ movie.description }}
      </div>

      <div class="movie-info">
        <div class="movie-info-left">
          <div v-if="movie.duration" class="info-item">
            <img :src="timelapseIcon" class="info-icon" alt="Duration" />
            {{ movie.duration }}
          </div>
          <div v-if="movie.ageRating" class="info-item">
            {{ movie.ageRating }}
          </div>
        </div>
        <div class="movie-info-right">
          <PlayButton @click="onPlayClick" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Movie } from '../types'
import timelapseIcon from '../assets/icons/timelapse.svg'
import starIcon from '../assets/icons/star.svg'
import PlayButton from './PlayButton.vue'

interface Props {
  movie: Movie
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'click', movie: Movie): void
  (e: 'play', movie: Movie): void
}>()

const onPlayClick = (event?: Event): void => {
  // Останавливаем всплытие чтобы не сработал клик по карточке
  if (event) {
    event.stopPropagation()
  }
  console.log('🎬 MovieCard: Play button clicked, emitting play event for:', props.movie.title)
  emit('play', props.movie)
}

const onImageError = (e: Event): void => {
  const target = e.target as HTMLImageElement
  const parent = target.parentElement
  if (parent) {
    target.style.display = 'none'
    const placeholder = parent.querySelector('.banner-placeholder') as HTMLElement
    if (placeholder) {
      placeholder.style.display = 'flex'
    }
  }
}
</script>

<style scoped>
/* M3 Expressive MovieCard */
.movie-card {
  background: var(--md-sys-color-surface-container-high);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: var(--md-sys-shape-corner-extra-large); /* M3 Expressive: 48px */
  padding: 20px; /* Expressive: увеличен padding */
  display: flex;
  gap: 24px; /* Expressive: увеличен gap */
  cursor: pointer;
  transition: all var(--md-sys-motion-duration-long1) var(--md-sys-motion-easing-emphasized); /* Expressive: 600ms плавная */
  box-shadow: var(--md-sys-elevation-level2); /* Expressive: начинаем с level2 */
  position: relative;
  overflow: hidden;
  will-change: transform, box-shadow, border-radius;
}

.movie-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--md-sys-color-primary);
  opacity: 0;
  transition: opacity var(--md-sys-motion-duration-medium2) var(--md-sys-motion-easing-standard); /* Expressive */
  pointer-events: none;
}

.movie-card:hover::before {
  opacity: 0.08; /* M3 Expressive: увеличена */
}

.movie-card:hover {
  border-color: var(--md-sys-color-primary);
  border-radius: var(--md-sys-shape-corner-large-increased); /* M3 Expressive: морфинг до 28px */
  box-shadow: var(--md-sys-elevation-level4); /* M3 Expressive: драматичная тень */
  transform: translateY(-6px) scale(1.02); /* Expressive: более выраженный подъем */
}

.movie-card:active {
  border-radius: var(--md-sys-shape-corner-large); /* M3 Expressive: еще менее округлая при клике */
  transform: translateY(-3px) scale(0.98);
  transition: all var(--md-sys-motion-duration-short4) var(--md-sys-motion-easing-emphasized-accelerate); /* Expressive: быстрая */
}

.movie-banner {
  width: 320px; /* Expressive: немного увеличен */
  height: 180px;
  border-radius: var(--md-sys-shape-corner-large); /* M3 Expressive: 24px */
  background: linear-gradient(135deg, var(--md-sys-color-primary-container) 0%, var(--md-sys-color-tertiary-container) 100%);
  flex-shrink: 0;
  overflow: hidden;
  position: relative;
  box-shadow: var(--md-sys-elevation-level2);
  transition: all var(--md-sys-motion-duration-long1) var(--md-sys-motion-easing-emphasized); /* Expressive */
  will-change: transform, border-radius;
}

.movie-card:hover .movie-banner {
  border-radius: var(--md-sys-shape-corner-medium); /* M3 Expressive: морфинг до 16px */
  transform: scale(1.08) rotate(1deg); /* Expressive: добавлен rotate */
  box-shadow: var(--md-sys-elevation-level3);
}

.movie-banner img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.movie-card:hover .movie-banner img {
  transform: scale(1.1);
}

.banner-placeholder {
  width: 100%;
  height: 100%;
  display: none;
  align-items: center;
  justify-content: center;
  font-size: 56px;
  position: absolute;
  top: 0;
  left: 0;
  color: var(--md-sys-color-on-primary-container);
}

.banner-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 100%;
  background: linear-gradient(
    to top,
    rgba(0, 0, 0, 0.85) 0%,
    rgba(0, 0, 0, 0.5) 40%,
    rgba(0, 0, 0, 0.2) 70%,
    rgba(0, 0, 0, 0) 100%
  );
  pointer-events: none;
}

.banner-content {
  position: absolute;
  bottom: 16px;
  left: 16px;
  right: 16px;
  z-index: 1;
}

.movie-title {
  font-family: var(--md-sys-typescale-title-large-font);
  font-size: var(--md-sys-typescale-title-large-size);
  font-weight: var(--md-sys-typescale-title-large-weight);
  line-height: var(--md-sys-typescale-title-large-line-height);
  color: #ffffff;
  margin-bottom: 4px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.9);
}

.movie-year {
  font-family: var(--md-sys-typescale-label-large-font);
  font-size: var(--md-sys-typescale-label-large-size);
  color: rgba(255, 255, 255, 0.9);
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.8);
}

.movie-metadata {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
}

.movie-header {
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  gap: 12px;
}

.movie-rating {
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  padding: 8px 16px; /* Expressive: увеличен padding */
  border-radius: var(--md-sys-shape-corner-large); /* M3 Expressive: 24px */
  font-family: var(--md-sys-typescale-label-large-font);
  font-size: var(--md-sys-typescale-label-large-size);
  font-weight: var(--md-sys-typescale-label-large-weight);
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: var(--md-sys-elevation-level2);
  transition: all var(--md-sys-motion-duration-medium2) var(--md-sys-motion-easing-emphasized); /* Expressive */
  will-change: transform, border-radius;
}

.movie-card:hover .movie-rating {
  border-radius: var(--md-sys-shape-corner-full); /* M3 Expressive: морфинг в pill */
  transform: scale(1.1) translateY(-2px); /* Expressive: более выраженный */
  box-shadow: var(--md-sys-elevation-level3);
}

.movie-genres {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.genre-tag {
  background: var(--md-sys-color-surface-container-highest);
  color: var(--md-sys-color-on-surface-variant);
  padding: 8px 18px; /* Expressive: увеличен padding */
  border-radius: var(--md-sys-shape-corner-medium); /* M3 Expressive: 16px */
  font-family: var(--md-sys-typescale-label-large-font);
  font-size: 13px; /* Expressive: немного крупнее */
  font-weight: 500;
  border: 1px solid var(--md-sys-color-outline-variant);
  transition: all var(--md-sys-motion-duration-medium1) var(--md-sys-motion-easing-emphasized); /* Expressive */
  will-change: transform, border-radius;
}

.genre-tag:hover {
  background: var(--md-sys-color-primary-container);
  color: var(--md-sys-color-on-primary-container);
  border-color: var(--md-sys-color-primary);
  border-radius: var(--md-sys-shape-corner-large); /* M3 Expressive: морфинг до 24px */
  transform: translateY(-3px) scale(1.08) rotate(2deg); /* Expressive: добавлен rotate */
  box-shadow: var(--md-sys-elevation-level2);
}

.movie-description {
  font-family: var(--md-sys-typescale-body-medium-font);
  font-size: var(--md-sys-typescale-body-medium-size);
  line-height: var(--md-sys-typescale-body-medium-line-height);
  color: var(--md-sys-color-on-surface-variant);
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.movie-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  font-family: var(--md-sys-typescale-body-medium-font);
  font-size: var(--md-sys-typescale-body-medium-size);
  color: var(--md-sys-color-on-surface-variant);
  margin-top: auto;
}

.movie-info-left {
  display: flex;
  gap: 20px;
  align-items: center;
}

.movie-info-right {
  margin-left: auto;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.info-icon {
  width: 18px;
  height: 18px;
  opacity: 0.7;
  color: var(--md-sys-color-primary);
}

.rating-icon {
  width: 18px;
  height: 18px;
  color: var(--md-sys-color-on-primary);
}
</style>
