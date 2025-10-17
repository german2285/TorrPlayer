<template>
  <div class="settings-overlay" @click="onClose">
    <div class="settings-dialog" @click.stop>
      <!-- Header -->
      <div class="settings-header">
        <h2 class="title-large">Настройки</h2>
        <button class="close-btn" @click="onClose">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
            <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
          </svg>
        </button>
      </div>

      <!-- Settings Content -->
      <div class="settings-content">

        <!-- Воспроизведение -->
        <div class="settings-section">
          <h3 class="section-title">Воспроизведение</h3>
          <div class="settings-list">

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Автовоспроизведение</div>
                <div class="setting-description">Начинать воспроизведение автоматически</div>
              </div>
              <label class="switch">
                <input type="checkbox" v-model="settings.autoplay">
                <span class="switch-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Субтитры</div>
                <div class="setting-description">Показывать субтитры по умолчанию</div>
              </div>
              <label class="switch">
                <input type="checkbox" v-model="settings.subtitles">
                <span class="switch-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Громкость</div>
                <div class="setting-description">{{ settings.volume }}%</div>
              </div>
              <input
                type="range"
                class="slider"
                min="0"
                max="100"
                v-model="settings.volume"
              >
            </div>

          </div>
        </div>

        <!-- Качество -->
        <div class="settings-section">
          <h3 class="section-title">Качество</h3>
          <div class="settings-list">

            <div class="setting-item clickable" @click="selectQuality('auto')">
              <div class="setting-info">
                <div class="setting-label">Автоматическое</div>
                <div class="setting-description">Подстраивать под скорость интернета</div>
              </div>
              <div class="radio-indicator" :class="{ active: settings.quality === 'auto' }"></div>
            </div>

            <div class="setting-item clickable" @click="selectQuality('1080p')">
              <div class="setting-info">
                <div class="setting-label">1080p Full HD</div>
                <div class="setting-description">Высокое качество</div>
              </div>
              <div class="radio-indicator" :class="{ active: settings.quality === '1080p' }"></div>
            </div>

            <div class="setting-item clickable" @click="selectQuality('720p')">
              <div class="setting-info">
                <div class="setting-label">720p HD</div>
                <div class="setting-description">Хорошее качество</div>
              </div>
              <div class="radio-indicator" :class="{ active: settings.quality === '720p' }"></div>
            </div>

            <div class="setting-item clickable" @click="selectQuality('480p')">
              <div class="setting-info">
                <div class="setting-label">480p SD</div>
                <div class="setting-description">Экономия трафика</div>
              </div>
              <div class="radio-indicator" :class="{ active: settings.quality === '480p' }"></div>
            </div>

          </div>
        </div>

        <!-- Интерфейс -->
        <div class="settings-section">
          <h3 class="section-title">Интерфейс</h3>
          <div class="settings-list">

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Темная тема</div>
                <div class="setting-description">Включено по умолчанию</div>
              </div>
              <label class="switch">
                <input type="checkbox" v-model="settings.darkTheme" disabled checked>
                <span class="switch-slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Анимации</div>
                <div class="setting-description">Плавные переходы и эффекты</div>
              </div>
              <label class="switch">
                <input type="checkbox" v-model="settings.animations">
                <span class="switch-slider"></span>
              </label>
            </div>

          </div>
        </div>

        <!-- О программе -->
        <div class="settings-section">
          <h3 class="section-title">О программе</h3>
          <div class="settings-list">

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Версия</div>
                <div class="setting-description">TorrPlayer 1.0.0</div>
              </div>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Разработчик</div>
                <div class="setting-description">TorrPlayer Team</div>
              </div>
            </div>

          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Ref } from 'vue'

const emit = defineEmits<{
  (e: 'close'): void
}>()

type QualityOption = 'auto' | '1080p' | '720p' | '480p'

interface SettingsData {
  autoplay: boolean
  subtitles: boolean
  volume: number
  quality: QualityOption
  darkTheme: boolean
  animations: boolean
}

const settings: Ref<SettingsData> = ref({
  autoplay: true,
  subtitles: false,
  volume: 80,
  quality: 'auto',
  darkTheme: true,
  animations: true
})

const onClose = (): void => {
  emit('close')
}

const selectQuality = (quality: QualityOption): void => {
  settings.value.quality = quality
}
</script>

<style scoped>
.settings-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

/* M3 Expressive Animations */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.settings-dialog {
  background: var(--md-sys-color-surface-container-high);
  border-radius: var(--md-sys-shape-corner-extra-large); /* M3 Expressive: 48px */
  box-shadow: var(--md-sys-elevation-level5); /* Expressive: драматичная тень */
  width: 90%;
  max-width: 650px; /* Expressive: немного шире */
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  overflow: hidden; /* Предотвращает выход контента за закругленные границы */
  /* M3 Expressive Physics System: Spring-based animations */
  animation:
    slideUp-spatial var(--md-sys-motion-spring-expressive-default-spatial-duration) var(--md-sys-motion-spring-expressive-default-spatial),
    slideUp-effects var(--md-sys-motion-spring-expressive-default-effects-duration) var(--md-sys-motion-spring-expressive-default-effects);
}

@keyframes slideUp {
  from {
    opacity: 0;
    border-radius: var(--md-sys-shape-corner-extra-large-increased); /* M3 Expressive: начинаем с 56px */
    transform: translateY(60px) scale(0.88); /* Expressive: более выраженный */
  }
  to {
    opacity: 1;
    border-radius: var(--md-sys-shape-corner-extra-large);
    transform: translateY(0) scale(1);
  }
}

/* Spatial animation (movement, size, corners) */
@keyframes slideUp-spatial {
  from {
    border-radius: var(--md-sys-shape-corner-extra-large-increased);
    transform: translateY(60px) scale(0.88);
  }
  to {
    border-radius: var(--md-sys-shape-corner-extra-large);
    transform: translateY(0) scale(1);
  }
}

/* Effects animation (opacity) */
@keyframes slideUp-effects {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.settings-header {
  padding: 24px 24px 16px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.title-large {
  margin: 0;
  font-family: var(--md-sys-typescale-title-large-font);
  font-size: var(--md-sys-typescale-title-large-size);
  font-weight: var(--md-sys-typescale-title-large-weight);
  color: var(--md-sys-color-on-surface);
}

.close-btn {
  width: 48px;  /* M3 Expressive: увеличен */
  height: 48px;
  border-radius: var(--md-sys-shape-corner-full);
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--md-sys-color-on-surface-variant);
  /* M3 Expressive Physics: Fast spatial spring for small components */
  transition:
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    border-radius var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    color var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
  position: relative;
  will-change: transform, border-radius;
}

.close-btn::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: var(--md-sys-color-on-surface);
  opacity: 0;
  /* M3 Expressive Physics: Fast effects spring for opacity */
  transition: opacity var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
}

.close-btn:hover {
  border-radius: var(--md-sys-shape-corner-medium); /* M3 Expressive: морфинг из круга в скругленный квадрат */
  transform: scale(1.08) rotate(90deg); /* Expressive: вращение при hover */
  color: var(--md-sys-color-error);
}

.close-btn:hover::before {
  opacity: var(--md-sys-state-hover-opacity); /* Expressive: 0.12 */
}

.close-btn:active {
  transform: scale(0.95) rotate(90deg);
}

.close-btn:active::before {
  opacity: var(--md-sys-state-pressed-opacity); /* Expressive: 0.16 */
}

.close-btn svg {
  position: relative;
  z-index: 1;
}

.settings-content {
  overflow-y: auto;
  padding: 8px 0 24px 0;
  scrollbar-width: none; /* Firefox: скрываем scrollbar */
  -ms-overflow-style: none; /* IE/Edge: скрываем scrollbar */
}

/* WebKit (Chrome, Safari, Edge): скрываем scrollbar */
.settings-content::-webkit-scrollbar {
  display: none;
}

.settings-section {
  padding: 16px 0;
}

.section-title {
  margin: 0;
  padding: 8px 24px;
  font-family: var(--md-sys-typescale-title-medium-font);
  font-size: var(--md-sys-typescale-title-medium-size);
  font-weight: var(--md-sys-typescale-title-medium-weight);
  color: var(--md-sys-color-primary);
}

.settings-list {
  display: flex;
  flex-direction: column;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  gap: 16px;
  transition: background 0.2s cubic-bezier(0.2, 0, 0, 1);
  position: relative;
}

.setting-item.clickable {
  cursor: pointer;
}

.setting-item.clickable::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--md-sys-color-on-surface);
  opacity: 0;
  transition: opacity 0.2s cubic-bezier(0.2, 0, 0, 1);
}

.setting-item.clickable:hover::before {
  opacity: var(--md-sys-state-hover-opacity);
}

.setting-item.clickable:active::before {
  opacity: var(--md-sys-state-pressed-opacity);
}

.setting-info {
  flex: 1;
  min-width: 0;
  position: relative;
  z-index: 1;
}

.setting-label {
  font-family: var(--md-sys-typescale-body-large-font);
  font-size: var(--md-sys-typescale-body-large-size);
  font-weight: 400;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 4px;
}

.setting-description {
  font-family: var(--md-sys-typescale-body-medium-font);
  font-size: var(--md-sys-typescale-body-medium-size);
  color: var(--md-sys-color-on-surface-variant);
}

/* Custom Switch - M3 Expressive */
.switch {
  position: relative;
  display: inline-block;
  width: 60px;  /* Expressive: увеличен */
  height: 36px; /* Expressive: увеличен */
  flex-shrink: 0;
  z-index: 1;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.switch-slider {
  position: absolute;
  cursor: pointer;
  inset: 0;
  background: var(--md-sys-color-surface-container-highest);
  border: 2px solid var(--md-sys-color-outline);
  /* M3 Expressive Physics: Fast effects spring for colors only */
  transition:
    background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    border-color var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
  border-radius: var(--md-sys-shape-corner-full); /* M3 Expressive: всегда полностью округлый */
  will-change: background, border-color;
}

.switch-slider:before {
  position: absolute;
  content: "";
  height: 20px; /* Expressive: увеличен thumb */
  width: 20px;
  left: 6px;
  top: 6px;
  background: var(--md-sys-color-outline);
  /* M3 Expressive Physics: Fast spatial spring for thumb movement */
  transition:
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    width var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    height var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
  border-radius: 50%;
  box-shadow: var(--md-sys-elevation-level1);
}

.switch-slider:hover {
  border-color: var(--md-sys-color-primary);
}

.switch input:checked + .switch-slider {
  background: var(--md-sys-color-primary);
  border-color: var(--md-sys-color-primary);
}

.switch input:checked + .switch-slider:before {
  transform: translateX(24px); /* Expressive: увеличенное смещение */
  background: var(--md-sys-color-on-primary);
  width: 22px; /* Expressive: thumb немного увеличивается */
  height: 22px;
  top: 5px;
  box-shadow: var(--md-sys-elevation-level2);
}

.switch input:disabled + .switch-slider {
  opacity: 0.38;
  cursor: not-allowed;
}

/* Custom Slider - M3 Expressive (Medium Size) */
/* Основано на Android tokens.xml: m3_comp_slider_medium_* */
.slider {
  width: 220px; /* Expressive: немного шире */
  height: 8px; /* M3 Expressive: увеличен track (адаптировано от 40dp Android) */
  border-radius: var(--md-sys-shape-corner-medium); /* M3 Expressive: 16px */
  background: var(--md-sys-color-surface-container-highest);
  outline: none;
  -webkit-appearance: none;
  appearance: none;
  position: relative;
  z-index: 1;
  /* M3 Expressive Physics: Fast spatial spring for track height */
  transition: height var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
}

.slider:hover {
  height: 10px; /* M3 Expressive: track увеличивается при hover */
}

/* WebKit (Chrome, Safari, Edge) - Вертикальная линия вместо круга */
.slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 4px; /* M3 Expressive: узкий thumb как вертикальная линия (как в Android) */
  height: 28px; /* M3 Expressive: высокий thumb (адаптировано от 44dp Android) */
  border-radius: var(--md-sys-shape-corner-extra-small); /* M3 Expressive: 6px */
  background: var(--md-sys-color-primary);
  cursor: grab;
  box-shadow: var(--md-sys-elevation-level2);
  /* M3 Expressive Physics: Fast spatial spring for thumb size */
  transition:
    width var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    height var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    box-shadow var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
}

.slider::-webkit-slider-thumb:hover {
  height: 32px; /* M3 Expressive: увеличивается при hover */
  width: 5px;
  box-shadow: var(--md-sys-elevation-level3);
  background: var(--md-sys-color-primary);
}

.slider::-webkit-slider-thumb:active {
  cursor: grabbing;
  height: 36px; /* M3 Expressive: максимальная высота при active */
  width: 6px;
  box-shadow: var(--md-sys-elevation-level4);
  background: var(--md-sys-color-primary);
}

/* Firefox - Вертикальная линия */
.slider::-moz-range-thumb {
  width: 4px; /* M3 Expressive: вертикальная линия */
  height: 28px;
  border-radius: var(--md-sys-shape-corner-extra-small); /* 6px */
  background: var(--md-sys-color-primary);
  cursor: grab;
  border: none;
  box-shadow: var(--md-sys-elevation-level2);
  /* M3 Expressive Physics: Fast spatial spring for thumb size */
  transition:
    width var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    height var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    box-shadow var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
}

.slider::-moz-range-thumb:hover {
  height: 32px;
  width: 5px;
  box-shadow: var(--md-sys-elevation-level3);
}

.slider::-moz-range-thumb:active {
  cursor: grabbing;
  height: 36px;
  width: 6px;
  box-shadow: var(--md-sys-elevation-level4);
}

/* Track styles для Firefox */
.slider::-moz-range-track {
  height: 8px;
  border-radius: var(--md-sys-shape-corner-medium);
  background: var(--md-sys-color-surface-container-highest);
  border: none;
}

/* Radio Indicator - M3 Expressive */
/* Основано на Android: m3_comp_radio_button_* */
.radio-indicator {
  width: 24px;  /* M3 Expressive: увеличен с 20px до 24px (+20%) */
  height: 24px;
  border-radius: 50%;
  border: 2px solid var(--md-sys-color-on-surface-variant);
  display: flex;
  align-items: center;
  justify-content: center;
  /* M3 Expressive Physics: Fast spatial spring for small component */
  transition:
    border-width var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    border-color var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
  position: relative;
  z-index: 1;
  cursor: pointer;
}

/* M3 Expressive: State layer (hover/press) */
.radio-indicator::before {
  content: '';
  position: absolute;
  inset: -8px;
  border-radius: 50%;
  background: var(--md-sys-color-primary);
  opacity: 0;
  /* M3 Expressive Physics: Fast effects spring for opacity */
  transition: opacity var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
}

.setting-item.clickable:hover .radio-indicator::before {
  opacity: var(--md-sys-state-hover-opacity); /* M3 Expressive: 0.12 */
}

.setting-item.clickable:active .radio-indicator::before {
  opacity: var(--md-sys-state-pressed-opacity); /* M3 Expressive: 0.16 */
}

/* Inner dot */
.radio-indicator::after {
  content: '';
  width: 12px;  /* M3 Expressive: увеличен с 10px */
  height: 12px;
  border-radius: 50%;
  background: var(--md-sys-color-primary);
  transform: scale(0);
  /* M3 Expressive Physics: Fast spatial spring for scale animation */
  transition: transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
}

.radio-indicator.active {
  border-color: var(--md-sys-color-primary);
  border-width: 2.5px; /* M3 Expressive: немного толще при active */
}

/* M3 Expressive: Pulsating animation при активации */
.radio-indicator.active::after {
  transform: scale(1);
  /* M3 Expressive Physics: Default spatial spring for pulse effect */
  animation: radioPulse var(--md-sys-motion-spring-expressive-default-spatial-duration) var(--md-sys-motion-spring-expressive-default-spatial) 1;
}

@keyframes radioPulse {
  0% {
    transform: scale(0);
  }
  50% {
    transform: scale(1.3);
  }
  100% {
    transform: scale(1);
  }
}
</style>
