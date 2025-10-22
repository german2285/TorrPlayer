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
                <div class="setting-label">Громкость музыки</div>
                <div class="setting-description">{{ bgMusicVolume }}%</div>
              </div>
              <input
                type="range"
                class="slider"
                min="0"
                max="100"
                v-model="bgMusicVolume"
              >
            </div>

          </div>
        </div>

        <!-- Торрент -->
        <div class="settings-section">
          <h3 class="section-title">Торрент</h3>
          <div class="settings-list">

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Размер RAM кеша</div>
                <div class="setting-description">{{ formatCacheSize(cacheSizeMB) }}</div>
              </div>
              <input
                type="range"
                class="slider"
                min="64"
                max="2048"
                step="64"
                v-model="cacheSizeMB"
              >
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Количество соединений</div>
                <div class="setting-description">{{ torrentSettings.connectionsLimit }} соединений</div>
              </div>
              <input
                type="range"
                class="slider"
                min="10"
                max="100"
                step="5"
                v-model="torrentSettings.connectionsLimit"
              >
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Скорость загрузки</div>
                <div class="setting-description">{{ formatSpeed(downloadRateMB) }}</div>
              </div>
              <input
                type="range"
                class="slider"
                min="0"
                max="100"
                step="5"
                v-model="downloadRateMB"
              >
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Скорость отдачи</div>
                <div class="setting-description">{{ formatSpeed(uploadRateMB) }}</div>
              </div>
              <input
                type="range"
                class="slider"
                min="0"
                max="100"
                step="5"
                v-model="uploadRateMB"
              >
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <div class="setting-label">Предзагрузка кеша</div>
                <div class="setting-description">{{ torrentSettings.preloadCache }}%</div>
              </div>
              <input
                type="range"
                class="slider"
                min="0"
                max="100"
                step="5"
                v-model="torrentSettings.preloadCache"
              >
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

          </div>
        </div>

        <!-- Цветовая тема -->
        <div class="settings-section">
          <h3 class="section-title">Цветовая тема</h3>
          <div class="settings-list">

            <div class="setting-item-full">
              <div class="setting-info">
                <div class="setting-label">Основной цвет приложения</div>
                <div class="setting-description">Выберите цвет из предустановленных или создайте свой</div>
              </div>
            </div>

            <!-- Предустановленные цвета -->
            <div class="color-presets">
              <button
                v-for="preset in colorPresets"
                :key="preset.color"
                class="color-preset-btn"
                :class="{ active: settings.themeColor === preset.color }"
                :style="{ backgroundColor: preset.color }"
                @click="selectColorPreset(preset.color)"
                :title="preset.name"
              >
                <svg v-if="settings.themeColor === preset.color" width="24" height="24" viewBox="0 0 24 24" fill="white">
                  <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z"/>
                </svg>
              </button>
            </div>

            <!-- RGB Color Picker -->
            <div class="setting-item-full">
              <div class="setting-info">
                <div class="setting-label">Пользовательский цвет</div>
                <div class="setting-description">Используйте RGB пикер для точной настройки</div>
              </div>
              <div class="color-picker-container">
                <input
                  type="color"
                  class="color-picker"
                  v-model="settings.themeColor"
                  @input="applyThemeColorWrapper(settings.themeColor)"
                >
                <span class="color-value">{{ settings.themeColor.toUpperCase() }}</span>
              </div>
            </div>

          </div>
        </div>

        <!-- Расширенные настройки -->
        <div class="settings-section">
          <h3 class="section-title">Расширенные настройки</h3>
          <div class="settings-list">

            <div class="setting-item clickable" @click="showAdvancedSettings = !showAdvancedSettings">
              <div class="setting-info">
                <div class="setting-label">{{ showAdvancedSettings ? 'Скрыть' : 'Показать' }} расширенные настройки</div>
                <div class="setting-description">Параметры для опытных пользователей</div>
              </div>
              <svg
                class="expand-icon"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="currentColor"
                :style="{ transform: showAdvancedSettings ? 'rotate(180deg)' : 'rotate(0deg)' }"
              >
                <path d="M7 10l5 5 5-5z"/>
              </svg>
            </div>

            <transition name="expand">
              <div v-if="showAdvancedSettings" class="advanced-settings-content">

                <!-- ХРАНИЛИЩЕ -->
                <div class="subsection-divider">ХРАНИЛИЩЕ</div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Использование диска</div>
                    <div class="setting-description">{{ torrentSettings.useDisk ? 'Диск' : 'RAM' }}</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.useDisk">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <div class="setting-item-full">
                  <div class="setting-info">
                    <div class="setting-label">Путь сохранения торрентов</div>
                    <div class="setting-description">Оставьте пустым для автоматического выбора</div>
                  </div>
                  <input
                    type="text"
                    class="text-input"
                    placeholder="По умолчанию"
                    v-model="torrentSettings.torrentsSavePath"
                  >
                </div>

                <!-- РЕТРЕКЕРЫ -->
                <div class="subsection-divider">РЕТРЕКЕРЫ</div>

                <div class="setting-item-with-help">
                  <div class="setting-info">
                    <div class="setting-label-with-icon">
                      <span>Режим ретрекеров</span>
                      <button
                        class="help-icon-button"
                        @click="showRetrackersHelpDialog = true"
                        title="Подробнее о режимах"
                      >
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                          <path d="M11 18h2v-2h-2v2zm1-16C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm0-14c-2.21 0-4 1.79-4 4h2c0-1.1.9-2 2-2s2 .9 2 2c0 2-3 1.75-3 5h2c0-2.25 3-2.5 3-5 0-2.21-1.79-4-4-4z"/>
                        </svg>
                      </button>
                    </div>
                    <div class="setting-description">{{ formatRetrackersMode(torrentSettings.retrackersMode) }}</div>
                  </div>
                  <input
                    type="range"
                    class="slider"
                    min="0"
                    max="3"
                    step="1"
                    v-model="torrentSettings.retrackersMode"
                  >
                </div>

                <!-- СЕТЕВЫЕ НАСТРОЙКИ -->
                <div class="subsection-divider">СЕТЕВЫЕ НАСТРОЙКИ</div>

                <div class="setting-item-full">
                  <div class="setting-info">
                    <div class="setting-label">Порт входящих соединений</div>
                    <div class="setting-description">0 = автоматический выбор порта</div>
                  </div>
                  <input
                    type="number"
                    class="number-input"
                    min="0"
                    max="65535"
                    placeholder="0"
                    v-model.number="torrentSettings.peersListenPort"
                  >
                </div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Отключить DHT</div>
                    <div class="setting-description">Распределенная хеш-таблица</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.disableDHT">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Отключить PEX</div>
                    <div class="setting-description">Обмен пирами</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.disablePEX">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Отключить µTP</div>
                    <div class="setting-description">Микро транспортный протокол</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.disableUTP">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Отключить UPnP</div>
                    <div class="setting-description">Автоматическая проброска портов</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.disableUPNP">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Отключить TCP</div>
                    <div class="setting-description">TCP протокол</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.disableTCP">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Отключить раздачу</div>
                    <div class="setting-description">Полностью отключить upload</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.disableUpload">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Принудительное шифрование</div>
                    <div class="setting-description">Шифровать все соединения</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.forceEncrypt">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Поддержка IPv6</div>
                    <div class="setting-description">Использовать IPv6 адреса</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.enableIPv6">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <!-- ПРОИЗВОДИТЕЛЬНОСТЬ -->
                <div class="subsection-divider">ПРОИЗВОДИТЕЛЬНОСТЬ</div>

                <div class="setting-item setting-warning">
                  <div class="setting-info">
                    <div class="setting-label">Режим отладки</div>
                    <div class="setting-description">Подробное логирование (увеличивает нагрузку)</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.enableDebug">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Опережающее чтение буфера</div>
                    <div class="setting-description">{{ torrentSettings.readerReadAHead }}%</div>
                  </div>
                  <input
                    type="range"
                    class="slider"
                    min="0"
                    max="100"
                    step="5"
                    v-model="torrentSettings.readerReadAHead"
                  >
                </div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Удалять кеш при остановке</div>
                    <div class="setting-description">Освобождать память при закрытии торрента</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.removeCacheOnDrop">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <div class="setting-item">
                  <div class="setting-info">
                    <div class="setting-label">Адаптивный режим</div>
                    <div class="setting-description">Автоматическая настройка производительности</div>
                  </div>
                  <label class="switch">
                    <input type="checkbox" v-model="torrentSettings.responsiveMode">
                    <span class="switch-slider"></span>
                  </label>
                </div>

                <div class="setting-item-full">
                  <div class="setting-info">
                    <div class="setting-label">Таймаут отключения торрента</div>
                    <div class="setting-description">Время ожидания перед отключением (секунды)</div>
                  </div>
                  <input
                    type="number"
                    class="number-input"
                    min="5"
                    max="300"
                    placeholder="30"
                    v-model.number="torrentSettings.torrentDisconnectTimeout"
                  >
                </div>

              </div>
            </transition>

          </div>
        </div>

        <!-- RuTracker -->
        <div class="settings-section">
          <h3 class="section-title">RuTracker</h3>
          <div class="settings-list">

            <!-- Auth Status -->
            <div class="setting-item-full">
              <div class="setting-info">
                <div class="setting-label">Статус авторизации</div>
                <div class="setting-description" :style="{ color: rutrackerAuth.isAuthenticated ? 'var(--md-sys-color-primary)' : 'var(--md-sys-color-error)' }">
                  {{ rutrackerAuth.isAuthenticated ? '✓ Авторизован' : '✗ Не авторизован' }}
                </div>
              </div>
            </div>

            <!-- Auth Form (только если не авторизован) -->
            <div v-if="!rutrackerAuth.isAuthenticated" class="setting-item-full">

              <!-- Auth Mode Toggle -->
              <div class="auth-mode-toggle">
                <button
                  class="auth-mode-btn"
                  :class="{ active: rutrackerAuth.mode === 'login' }"
                  @click="rutrackerAuth.mode = 'login'"
                >
                  Вход
                </button>
                <button
                  class="auth-mode-btn"
                  :class="{ active: rutrackerAuth.mode === 'register' }"
                  @click="rutrackerAuth.mode = 'register'; loadCaptcha()"
                >
                  Регистрация
                </button>
              </div>

              <!-- Login Form -->
              <div v-if="rutrackerAuth.mode === 'login'" class="auth-form">
                <input
                  type="text"
                  class="text-input"
                  placeholder="Имя пользователя"
                  v-model="rutrackerAuth.loginData.username"
                >
                <input
                  type="password"
                  class="text-input"
                  placeholder="Пароль"
                  v-model="rutrackerAuth.loginData.password"
                >
                <button class="auth-submit-btn" @click="handleLogin" :disabled="rutrackerAuth.loading">
                  {{ rutrackerAuth.loading ? 'Вход...' : 'Войти' }}
                </button>
              </div>

              <!-- Registration Form -->
              <div v-if="rutrackerAuth.mode === 'register'" class="auth-form">
                <input
                  type="text"
                  class="text-input"
                  placeholder="Имя пользователя"
                  v-model="rutrackerAuth.registerData.username"
                >
                <input
                  type="password"
                  class="text-input"
                  placeholder="Пароль (макс. 20 символов)"
                  maxlength="20"
                  v-model="rutrackerAuth.registerData.password"
                >
                <input
                  type="email"
                  class="text-input"
                  placeholder="Email"
                  v-model="rutrackerAuth.registerData.email"
                >

                <!-- CAPTCHA -->
                <div v-if="rutrackerAuth.captchaData" class="captcha-container">
                  <img
                    :src="rutrackerAuth.captchaData.imageBase64"
                    alt="CAPTCHA"
                    class="captcha-image"
                  >
                  <button class="captcha-reload-btn" @click="loadCaptcha" title="Обновить CAPTCHA">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/>
                    </svg>
                  </button>
                </div>

                <input
                  type="text"
                  class="text-input"
                  placeholder="Введите код с картинки"
                  v-model="rutrackerAuth.registerData.captchaCode"
                >

                <button class="auth-submit-btn" @click="handleRegister" :disabled="rutrackerAuth.loading || !rutrackerAuth.captchaData">
                  {{ rutrackerAuth.loading ? 'Регистрация...' : 'Зарегистрироваться' }}
                </button>
              </div>

              <!-- Error Message -->
              <div v-if="rutrackerAuth.error" class="auth-error">
                {{ rutrackerAuth.error }}
              </div>
            </div>

            <!-- Logout Button (только если авторизован) -->
            <div v-if="rutrackerAuth.isAuthenticated" class="setting-item-full">
              <button class="auth-logout-btn" @click="handleLogout">
                Выйти
              </button>
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

  <!-- Retrackers Help Dialog -->
  <div v-if="showRetrackersHelpDialog" class="dialog-overlay" @click="showRetrackersHelpDialog = false">
    <div class="dialog-container" @click.stop>
      <div class="dialog-header">
        <h3>Режимы ретрекеров</h3>
        <button class="dialog-close" @click="showRetrackersHelpDialog = false">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
            <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
          </svg>
        </button>
      </div>

      <div class="dialog-content">
        <div class="retrackers-mode-item">
          <div class="mode-number">0</div>
          <div class="mode-info">
            <div class="mode-title">Не добавлять</div>
            <div class="mode-description">
              Не изменять список трекеров торрента. Используйте этот режим для приватных торрентов или когда не нужны дополнительные источники.
            </div>
          </div>
        </div>

        <div class="retrackers-mode-item mode-default">
          <div class="mode-number">1</div>
          <div class="mode-info">
            <div class="mode-title">Добавить (по умолчанию)</div>
            <div class="mode-description">
              Добавить публичные ретрекеры к существующим трекерам. Увеличивает количество источников (пиров и сидов), что ускоряет загрузку. Рекомендуется для большинства торрентов.
            </div>
          </div>
        </div>

        <div class="retrackers-mode-item">
          <div class="mode-number">2</div>
          <div class="mode-info">
            <div class="mode-title">Удалить</div>
            <div class="mode-description">
              Удалить все публичные ретрекеры из списка трекеров. Используйте для ограничения трафика или работы только с оригинальными трекерами торрента.
            </div>
          </div>
        </div>

        <div class="retrackers-mode-item">
          <div class="mode-number">3</div>
          <div class="mode-info">
            <div class="mode-title">Заменить</div>
            <div class="mode-description">
              Заменить все трекеры торрента на публичные ретрекеры. Полная замена списка трекеров. Используйте, если оригинальные трекеры не работают.
            </div>
          </div>
        </div>
      </div>

      <div class="dialog-actions">
        <button class="dialog-button" @click="showRetrackersHelpDialog = false">
          Понятно
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import type { Ref } from 'vue'
import { GetSettings, SetSettings, CheckAuthStatus, LoginToRuTracker, RegisterOnRuTracker, GetRegistrationCaptcha, LogoutFromRuTracker } from '../../wailsjs/go/app/App'
import type { app } from '../../wailsjs/go/models'
import { applyThemeColor, loadSavedTheme } from '../utils/themeUtils'
import type { LoginData, RegistrationData, CaptchaData } from '../types'

const emit = defineEmits<{
  (e: 'close'): void
}>()

interface SettingsData {
  darkTheme: boolean
  themeColor: string
}

interface TorrentSettings {
  // Existing settings
  cacheSize: number        // in bytes
  cacheSizeStr: string
  connectionsLimit: number
  downloadRate: number     // in kb/s, 0 = unlimited
  uploadRate: number       // in kb/s, 0 = unlimited
  preloadCache: number     // in percent 0-100
  retrackersMode: number   // 0-3
  bgMusicVolume: number    // in percent 0-100
  themeColor: string

  // Storage settings
  useDisk: boolean         // Use disk instead of RAM cache
  torrentsSavePath: string // Path to save torrents

  // Network settings
  peersListenPort: number  // Port for incoming connections (0 = auto)

  // Protocol settings
  disableDHT: boolean      // Disable DHT
  disablePEX: boolean      // Disable PEX
  disableUTP: boolean      // Disable µTP
  disableUPNP: boolean     // Disable UPnP
  disableTCP: boolean      // Disable TCP

  // Upload and encryption
  disableUpload: boolean   // Disable upload completely
  forceEncrypt: boolean    // Force encryption
  enableIPv6: boolean      // Enable IPv6 support

  // Advanced settings
  enableDebug: boolean     // Enable debug logging
  readerReadAHead: number  // Reader read ahead percentage 0-100
  removeCacheOnDrop: boolean // Remove cache on drop
  responsiveMode: boolean  // Responsive performance mode
  torrentDisconnectTimeout: number // Timeout in seconds
}

const settings: Ref<SettingsData> = ref({
  darkTheme: true,
  themeColor: '#ffffff' // Synced with JSON config default
})

// Предустановленные цветовые схемы M3
const colorPresets = [
  { name: 'Фиолетовый (по умолчанию)', color: '#6750A4' },
  { name: 'Синий', color: '#0061A4' },
  { name: 'Зеленый', color: '#006D3B' },
  { name: 'Красный', color: '#BA1A1A' },
  { name: 'Оранжевый', color: '#A05000' },
  { name: 'Розовый', color: '#A44A7A' },
  { name: 'Бирюзовый', color: '#006874' },
  { name: 'Желтый', color: '#755B00' }
]

const torrentSettings: Ref<TorrentSettings> = ref({
  // Existing settings
  cacheSize: 67108864, // 64 MB по умолчанию
  cacheSizeStr: '64 MB',
  connectionsLimit: 25,
  downloadRate: 0,
  uploadRate: 0,
  preloadCache: 50,
  retrackersMode: 1,
  bgMusicVolume: 30,
  themeColor: '#ffffff',

  // Storage settings
  useDisk: false,
  torrentsSavePath: '',

  // Network settings
  peersListenPort: 0,

  // Protocol settings
  disableDHT: false,
  disablePEX: false,
  disableUTP: false,
  disableUPNP: false,
  disableTCP: false,

  // Upload and encryption
  disableUpload: false,
  forceEncrypt: false,
  enableIPv6: false,

  // Advanced settings
  enableDebug: false,
  readerReadAHead: 95,
  removeCacheOnDrop: false,
  responsiveMode: false,
  torrentDisconnectTimeout: 30
})

// Для отображения в UI (конвертация из bytes в MB/GB)
const cacheSizeMB = ref(64)
const downloadRateMB = ref(0)
const uploadRateMB = ref(0)
const bgMusicVolume = ref(30) // Proxy для UI

// Collapsible advanced settings state
const showAdvancedSettings = ref(false)

// Retrackers help dialog state
const showRetrackersHelpDialog = ref(false)

// RuTracker authentication state
const rutrackerAuth = ref({
  isAuthenticated: false,
  mode: 'login' as 'login' | 'register',
  loading: false,
  error: '',
  loginData: {
    username: '',
    password: ''
  } as LoginData,
  registerData: {
    username: '',
    password: '',
    email: '',
    captchaCode: '',
    captchaSid: '',
    codeField: ''
  } as RegistrationData,
  captchaData: null as CaptchaData | null
})

// Watch bgMusicVolume changes and dispatch event
watch(bgMusicVolume, (newVolume) => {
  torrentSettings.value.bgMusicVolume = newVolume
  // Dispatch custom event for App.vue to listen
  window.dispatchEvent(new CustomEvent('bgMusicVolumeChanged', { detail: newVolume }))
})

onMounted(async () => {
  try {
    const btSettings = await GetSettings()
    if (btSettings) {
      // Existing settings
      torrentSettings.value.cacheSize = btSettings.cacheSize
      torrentSettings.value.cacheSizeStr = btSettings.cacheSizeStr
      torrentSettings.value.connectionsLimit = btSettings.connectionsLimit
      torrentSettings.value.downloadRate = btSettings.downloadRate
      torrentSettings.value.uploadRate = btSettings.uploadRate
      torrentSettings.value.preloadCache = btSettings.preloadCache
      torrentSettings.value.retrackersMode = btSettings.retrackersMode
      torrentSettings.value.bgMusicVolume = btSettings.bgMusicVolume

      // Storage settings
      torrentSettings.value.useDisk = btSettings.useDisk || false
      torrentSettings.value.torrentsSavePath = btSettings.torrentsSavePath || ''

      // Network settings
      torrentSettings.value.peersListenPort = btSettings.peersListenPort || 0

      // Protocol settings
      torrentSettings.value.disableDHT = btSettings.disableDHT || false
      torrentSettings.value.disablePEX = btSettings.disablePEX || false
      torrentSettings.value.disableUTP = btSettings.disableUTP || false
      torrentSettings.value.disableUPNP = btSettings.disableUPNP || false
      torrentSettings.value.disableTCP = btSettings.disableTCP || false

      // Upload and encryption
      torrentSettings.value.disableUpload = btSettings.disableUpload || false
      torrentSettings.value.forceEncrypt = btSettings.forceEncrypt || false
      torrentSettings.value.enableIPv6 = btSettings.enableIPv6 || false

      // Advanced settings
      torrentSettings.value.enableDebug = btSettings.enableDebug || false
      torrentSettings.value.readerReadAHead = btSettings.readerReadAHead || 95
      torrentSettings.value.removeCacheOnDrop = btSettings.removeCacheOnDrop || false
      torrentSettings.value.responsiveMode = btSettings.responsiveMode || false
      torrentSettings.value.torrentDisconnectTimeout = btSettings.torrentDisconnectTimeout || 30

      // Конвертируем для UI
      cacheSizeMB.value = Math.round(btSettings.cacheSize / (1024 * 1024))
      downloadRateMB.value = Math.round(btSettings.downloadRate / 1024)
      uploadRateMB.value = Math.round(btSettings.uploadRate / 1024)
      bgMusicVolume.value = btSettings.bgMusicVolume

      // Загрузить цвет темы из backend
      if (btSettings.themeColor) {
        settings.value.themeColor = btSettings.themeColor
      }
    }
  } catch (error) {
    console.error('Failed to load torrent settings:', error)
  }

  // Check RuTracker auth status
  try {
    const isAuth = await CheckAuthStatus()
    rutrackerAuth.value.isAuthenticated = isAuth
  } catch (error) {
    console.error('Failed to check RuTracker auth status:', error)
    rutrackerAuth.value.isAuthenticated = false
  }
})

const onClose = async (): Promise<void> => {
  // Сохраняем настройки перед закрытием
  await saveSettings()
  emit('close')
}

const saveSettings = async (): Promise<void> => {
  try {
    // Конвертируем из UI в bytes/kb
    torrentSettings.value.cacheSize = cacheSizeMB.value * 1024 * 1024
    torrentSettings.value.downloadRate = downloadRateMB.value * 1024
    torrentSettings.value.uploadRate = uploadRateMB.value * 1024
    torrentSettings.value.bgMusicVolume = bgMusicVolume.value

    // Создаем объект с всеми настройками включая цвет темы
    const allSettings = {
      ...torrentSettings.value,
      themeColor: settings.value.themeColor
    }

    await SetSettings(allSettings as any)
  } catch (error) {
    console.error('Failed to save settings:', error)
  }
}

const formatCacheSize = (mb: number): string => {
  if (mb >= 1024) {
    return `${(mb / 1024).toFixed(1)} GB`
  }
  return `${mb} MB`
}

const formatSpeed = (mb: number): string => {
  if (mb === 0) return 'Безлимит'
  if (mb >= 1024) {
    return `${(mb / 1024).toFixed(1)} GB/s`
  }
  return `${mb} MB/s`
}

const formatRetrackersMode = (mode: number): string => {
  const modes = ['Не добавлять', 'Добавить', 'Удалить', 'Заменить']
  return modes[mode] || 'Неизвестно'
}

// Wrapper для applyThemeColor с обновлением settings
const applyThemeColorWrapper = (color: string): void => {
  settings.value.themeColor = color
  applyThemeColor(color)
}

// Выбрать предустановленный цвет
const selectColorPreset = (color: string): void => {
  applyThemeColorWrapper(color)
}

// RuTracker Methods

// Load CAPTCHA for registration
const loadCaptcha = async (): Promise<void> => {
  try {
    rutrackerAuth.value.error = ''
    const captcha = await GetRegistrationCaptcha()
    if (captcha) {
      rutrackerAuth.value.captchaData = captcha
      rutrackerAuth.value.registerData.captchaSid = captcha.sid
      rutrackerAuth.value.registerData.codeField = captcha.codeField
      rutrackerAuth.value.registerData.captchaCode = '' // Clear previous code
    }
  } catch (error) {
    console.error('Failed to load CAPTCHA:', error)
    rutrackerAuth.value.error = 'Не удалось загрузить CAPTCHA. Попробуйте еще раз.'
  }
}

// Handle login
const handleLogin = async (): Promise<void> => {
  try {
    rutrackerAuth.value.loading = true
    rutrackerAuth.value.error = ''

    if (!rutrackerAuth.value.loginData.username || !rutrackerAuth.value.loginData.password) {
      rutrackerAuth.value.error = 'Заполните все поля'
      return
    }

    await LoginToRuTracker(rutrackerAuth.value.loginData)

    // Check auth status after login
    const isAuth = await CheckAuthStatus()
    rutrackerAuth.value.isAuthenticated = isAuth

    if (isAuth) {
      // Clear form
      rutrackerAuth.value.loginData.username = ''
      rutrackerAuth.value.loginData.password = ''
    }
  } catch (error: any) {
    console.error('Login failed:', error)
    rutrackerAuth.value.error = error?.message || 'Ошибка входа. Проверьте логин и пароль.'
  } finally {
    rutrackerAuth.value.loading = false
  }
}

// Handle registration
const handleRegister = async (): Promise<void> => {
  try {
    rutrackerAuth.value.loading = true
    rutrackerAuth.value.error = ''

    const data = rutrackerAuth.value.registerData

    if (!data.username || !data.password || !data.email || !data.captchaCode) {
      rutrackerAuth.value.error = 'Заполните все поля'
      return
    }

    if (data.password.length > 20) {
      rutrackerAuth.value.error = 'Пароль должен быть не более 20 символов'
      return
    }

    await RegisterOnRuTracker(data)

    // Check auth status after registration
    const isAuth = await CheckAuthStatus()
    rutrackerAuth.value.isAuthenticated = isAuth

    if (isAuth) {
      // Clear form
      rutrackerAuth.value.registerData.username = ''
      rutrackerAuth.value.registerData.password = ''
      rutrackerAuth.value.registerData.email = ''
      rutrackerAuth.value.registerData.captchaCode = ''
      rutrackerAuth.value.captchaData = null
    }
  } catch (error: any) {
    console.error('Registration failed:', error)
    rutrackerAuth.value.error = error?.message || 'Ошибка регистрации. Попробуйте еще раз.'
    // Reload CAPTCHA on error
    await loadCaptcha()
  } finally {
    rutrackerAuth.value.loading = false
  }
}

// Handle logout
const handleLogout = async (): Promise<void> => {
  try {
    await LogoutFromRuTracker()
    rutrackerAuth.value.isAuthenticated = false
    rutrackerAuth.value.error = ''
  } catch (error: any) {
    console.error('Logout failed:', error)
    rutrackerAuth.value.error = error?.message || 'Ошибка выхода'
  }
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
  will-change: transform, width, height, background; /* Оптимизация анимации */
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
  height: 16px; /* M3 Expressive Medium: 16dp track height (официальная спецификация) */
  border-radius: var(--md-sys-shape-corner-medium); /* M3 Expressive: 16px */
  background: var(--md-sys-color-surface-container-highest);
  outline: none;
  -webkit-appearance: none;
  appearance: none;
  position: relative;
  z-index: 1;
  /* M3 Expressive Physics: Fast spatial spring for track height */
  transition: height var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
  will-change: height; /* Оптимизация анимации */
}

.slider:hover {
  height: 18px; /* M3 Expressive: track увеличивается при hover */
}

/* WebKit (Chrome, Safari, Edge) - Вертикальная линия вместо круга */
.slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 4px; /* M3 Expressive Medium: 4dp узкий thumb (официальная спецификация) */
  height: 44px; /* M3 Expressive Medium: 44dp высокий thumb (официальная спецификация Android) */
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
  will-change: width, height, box-shadow; /* Оптимизация анимации */
}

.slider::-webkit-slider-thumb:hover {
  height: 48px; /* M3 Expressive Medium: увеличивается при hover */
  width: 5px;
  box-shadow: var(--md-sys-elevation-level3);
  background: var(--md-sys-color-primary);
}

.slider::-webkit-slider-thumb:active {
  cursor: grabbing;
  height: 52px; /* M3 Expressive Medium: максимальная высота при active */
  width: 6px;
  box-shadow: var(--md-sys-elevation-level4);
  background: var(--md-sys-color-primary);
}

/* Firefox - Вертикальная линия */
.slider::-moz-range-thumb {
  width: 4px; /* M3 Expressive Medium: 4dp */
  height: 44px; /* M3 Expressive Medium: 44dp (официальная спецификация) */
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
  will-change: width, height, box-shadow; /* Оптимизация анимации */
}

.slider::-moz-range-thumb:hover {
  height: 48px; /* M3 Expressive Medium: увеличивается при hover */
  width: 5px;
  box-shadow: var(--md-sys-elevation-level3);
}

.slider::-moz-range-thumb:active {
  cursor: grabbing;
  height: 52px; /* M3 Expressive Medium: максимальная высота при active */
  width: 6px;
  box-shadow: var(--md-sys-elevation-level4);
}

/* Track styles для Firefox */
.slider::-moz-range-track {
  height: 16px; /* M3 Expressive Medium: 16dp (соответствие основному track) */
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

/* Full Width Setting Item */
.setting-item-full {
  padding: 16px 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* Color Presets Grid */
.color-presets {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(60px, 1fr));
  gap: 12px;
  padding: 0 24px 16px 24px;
}

.color-preset-btn {
  width: 60px;
  height: 60px;
  border-radius: var(--md-sys-shape-corner-large); /* M3 Expressive: 20px */
  border: 3px solid transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  /* M3 Expressive Physics: Fast spatial spring */
  transition:
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    border-color var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    box-shadow var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
  box-shadow: var(--md-sys-elevation-level2);
}

.color-preset-btn:hover {
  transform: scale(1.1);
  box-shadow: var(--md-sys-elevation-level3);
}

.color-preset-btn:active {
  transform: scale(0.95);
}

.color-preset-btn.active {
  border-color: white;
  box-shadow: var(--md-sys-elevation-level4), 0 0 0 2px var(--md-sys-color-primary);
  transform: scale(1.05);
}

/* Color Picker Container */
.color-picker-container {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 0 24px 8px 24px;
}

.color-picker {
  width: 80px;
  height: 80px;
  border-radius: var(--md-sys-shape-corner-extra-large); /* M3 Expressive: 48px */
  border: 3px solid var(--md-sys-color-outline);
  cursor: pointer;
  /* M3 Expressive Physics: Fast spatial spring */
  transition:
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial),
    border-color var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    box-shadow var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
  box-shadow: var(--md-sys-elevation-level2);
}

.color-picker:hover {
  transform: scale(1.08);
  border-color: var(--md-sys-color-primary);
  box-shadow: var(--md-sys-elevation-level3);
}

.color-picker:active {
  transform: scale(0.98);
}

.color-value {
  font-family: var(--md-sys-typescale-body-large-font);
  font-size: var(--md-sys-typescale-body-large-size);
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
  letter-spacing: 0.5px;
  background: var(--md-sys-color-surface-container-high);
  padding: 12px 20px;
  border-radius: var(--md-sys-shape-corner-medium);
  box-shadow: var(--md-sys-elevation-level1);
}

/* Text Input - Material Design 3 Outlined TextField */
.text-input {
  width: 100%;
  padding: 16px;
  font-family: var(--md-sys-typescale-body-large-font);
  font-size: var(--md-sys-typescale-body-large-size);
  color: var(--md-sys-color-on-surface);
  background: transparent;
  border: 2px solid var(--md-sys-color-outline);
  border-radius: var(--md-sys-shape-corner-medium);
  outline: none;
  transition:
    border-color var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    box-shadow var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
}

.text-input:hover {
  border-color: var(--md-sys-color-on-surface);
}

.text-input:focus {
  border-color: var(--md-sys-color-primary);
  box-shadow: 0 0 0 1px var(--md-sys-color-primary);
}

.text-input::placeholder {
  color: var(--md-sys-color-on-surface-variant);
  opacity: 0.6;
}

/* Number Input - Material Design 3 */
.number-input {
  width: 100%;
  padding: 16px;
  font-family: var(--md-sys-typescale-body-large-font);
  font-size: var(--md-sys-typescale-body-large-size);
  color: var(--md-sys-color-on-surface);
  background: transparent;
  border: 2px solid var(--md-sys-color-outline);
  border-radius: var(--md-sys-shape-corner-medium);
  outline: none;
  transition:
    border-color var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    box-shadow var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
}

.number-input:hover {
  border-color: var(--md-sys-color-on-surface);
}

.number-input:focus {
  border-color: var(--md-sys-color-primary);
  box-shadow: 0 0 0 1px var(--md-sys-color-primary);
}

.number-input::placeholder {
  color: var(--md-sys-color-on-surface-variant);
  opacity: 0.6;
}

/* Remove spinner arrows for number inputs */
.number-input::-webkit-outer-spin-button,
.number-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.number-input[type=number] {
  -moz-appearance: textfield;
}

/* Expand Icon - M3 Expressive */
.expand-icon {
  color: var(--md-sys-color-on-surface-variant);
  transition: transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
  flex-shrink: 0;
  z-index: 1;
}

/* Warning Setting Item */
.setting-warning {
  border-left: 3px solid var(--md-sys-color-error);
  background: linear-gradient(
    90deg,
    var(--md-sys-color-error-container) 0%,
    transparent 100%
  );
  opacity: 0.95;
}

.setting-warning .setting-label {
  color: var(--md-sys-color-error);
  font-weight: 500;
}

/* Advanced Settings Content */
.advanced-settings-content {
  display: flex;
  flex-direction: column;
  width: 100%;
}

/* Expand Transition - M3 Expressive */
.expand-enter-active,
.expand-leave-active {
  transition:
    opacity var(--md-sys-motion-spring-expressive-default-effects-duration) var(--md-sys-motion-spring-expressive-default-effects),
    max-height var(--md-sys-motion-spring-expressive-default-spatial-duration) var(--md-sys-motion-spring-expressive-default-spatial);
  overflow: hidden;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
}

.expand-enter-to,
.expand-leave-from {
  opacity: 1;
  max-height: 1000px; /* Large enough to contain all advanced settings */
}

/* Subsection Dividers */
.subsection-divider {
  font-family: var(--md-sys-typescale-label-small-font);
  font-size: 11px;
  font-weight: 500;
  color: var(--md-sys-color-primary);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin: 20px 24px 12px 24px;
  padding-bottom: 8px;
  border-bottom: 1px dashed var(--md-sys-color-outline-variant);
  opacity: 0.8;
}

/* Setting Item with Help Button */
.setting-item-with-help {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  gap: 16px;
  position: relative;
}

.setting-label-with-icon {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* Help Icon Button - Material Design */
.help-icon-button {
  width: 20px;
  height: 20px;
  padding: 0;
  margin: 0;
  background: transparent;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  color: var(--md-sys-color-on-surface-variant);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition:
    background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    color var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
}

.help-icon-button:hover {
  background: var(--md-sys-color-surface-container-highest);
  color: var(--md-sys-color-primary);
  transform: scale(1.2);
}

.help-icon-button:active {
  transform: scale(1.1);
}

.help-icon-button svg {
  width: 16px;
  height: 16px;
}

/* Material Design Dialog Overlay */
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  animation: fadeIn 0.2s ease;
}

/* Material Design Dialog Container */
.dialog-container {
  background: var(--md-sys-color-surface-container-high);
  border-radius: var(--md-sys-shape-corner-extra-large);
  box-shadow: var(--md-sys-elevation-level5);
  width: 90%;
  max-width: 550px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  animation:
    dialogSlideUp-spatial var(--md-sys-motion-spring-expressive-default-spatial-duration) var(--md-sys-motion-spring-expressive-default-spatial),
    dialogSlideUp-effects var(--md-sys-motion-spring-expressive-default-effects-duration) var(--md-sys-motion-spring-expressive-default-effects);
}

@keyframes dialogSlideUp-spatial {
  from {
    transform: translateY(40px) scale(0.92);
  }
  to {
    transform: translateY(0) scale(1);
  }
}

@keyframes dialogSlideUp-effects {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* Dialog Header */
.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 24px 16px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.dialog-header h3 {
  margin: 0;
  font-family: var(--md-sys-typescale-headline-small-font);
  font-size: var(--md-sys-typescale-headline-small-size);
  font-weight: var(--md-sys-typescale-headline-small-weight);
  color: var(--md-sys-color-on-surface);
}

.dialog-close {
  width: 40px;
  height: 40px;
  border-radius: var(--md-sys-shape-corner-full);
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--md-sys-color-on-surface-variant);
  transition:
    background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    color var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
}

.dialog-close:hover {
  background: var(--md-sys-color-surface-container-highest);
  color: var(--md-sys-color-error);
  transform: rotate(90deg);
}

.dialog-close:active {
  transform: rotate(90deg) scale(0.95);
}

/* Dialog Content */
.dialog-content {
  padding: 16px 24px;
  overflow-y: auto;
  flex: 1;
  scrollbar-width: thin;
  scrollbar-color: var(--md-sys-color-outline-variant) transparent;
}

.dialog-content::-webkit-scrollbar {
  width: 8px;
}

.dialog-content::-webkit-scrollbar-track {
  background: transparent;
}

.dialog-content::-webkit-scrollbar-thumb {
  background: var(--md-sys-color-outline-variant);
  border-radius: 4px;
}

/* Retrackers Mode Item */
.retrackers-mode-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  margin-bottom: 12px;
  background: var(--md-sys-color-surface-container);
  border-radius: var(--md-sys-shape-corner-large);
  border: 2px solid transparent;
  transition:
    border-color var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects);
}

.retrackers-mode-item.mode-default {
  border-color: var(--md-sys-color-primary);
  background: var(--md-sys-color-primary-container);
}

.retrackers-mode-item:last-child {
  margin-bottom: 0;
}

.mode-number {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--md-sys-typescale-title-large-font);
  font-size: var(--md-sys-typescale-title-large-size);
  font-weight: 700;
  flex-shrink: 0;
}

.mode-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mode-title {
  font-family: var(--md-sys-typescale-title-small-font);
  font-size: var(--md-sys-typescale-title-small-size);
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.mode-description {
  font-family: var(--md-sys-typescale-body-small-font);
  font-size: var(--md-sys-typescale-body-small-size);
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.5;
}

/* Dialog Actions */
.dialog-actions {
  padding: 16px 24px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
  display: flex;
  justify-content: flex-end;
}

.dialog-button {
  padding: 10px 24px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  border: none;
  border-radius: var(--md-sys-shape-corner-full);
  font-family: var(--md-sys-typescale-label-large-font);
  font-size: var(--md-sys-typescale-label-large-size);
  font-weight: 500;
  cursor: pointer;
  transition:
    background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    box-shadow var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
}

.dialog-button:hover {
  box-shadow: var(--md-sys-elevation-level2);
  transform: scale(1.02);
}

.dialog-button:active {
  transform: scale(0.98);
}

/* RuTracker Auth Styles */

/* Auth Mode Toggle */
.auth-mode-toggle {
  display: flex;
  gap: 8px;
  padding: 4px;
  background: var(--md-sys-color-surface-container-highest);
  border-radius: var(--md-sys-shape-corner-full);
  margin-bottom: 16px;
}

.auth-mode-btn {
  flex: 1;
  padding: 12px 24px;
  background: transparent;
  color: var(--md-sys-color-on-surface-variant);
  border: none;
  border-radius: var(--md-sys-shape-corner-full);
  font-family: var(--md-sys-typescale-label-large-font);
  font-size: var(--md-sys-typescale-label-large-size);
  font-weight: 500;
  cursor: pointer;
  transition:
    background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    color var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
}

.auth-mode-btn:hover {
  background: var(--md-sys-color-surface-container);
}

.auth-mode-btn.active {
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  box-shadow: var(--md-sys-elevation-level2);
}

.auth-mode-btn:active {
  transform: scale(0.98);
}

/* Auth Form */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Auth Submit Button */
.auth-submit-btn {
  width: 100%;
  padding: 16px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  border: none;
  border-radius: var(--md-sys-shape-corner-medium);
  font-family: var(--md-sys-typescale-label-large-font);
  font-size: var(--md-sys-typescale-label-large-size);
  font-weight: 500;
  cursor: pointer;
  transition:
    background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    box-shadow var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
  box-shadow: var(--md-sys-elevation-level1);
}

.auth-submit-btn:hover:not(:disabled) {
  box-shadow: var(--md-sys-elevation-level2);
  transform: scale(1.02);
}

.auth-submit-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.auth-submit-btn:disabled {
  opacity: 0.38;
  cursor: not-allowed;
}

/* CAPTCHA Container */
.captcha-container {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--md-sys-color-surface-container);
  border-radius: var(--md-sys-shape-corner-medium);
  border: 2px solid var(--md-sys-color-outline-variant);
}

.captcha-image {
  flex: 1;
  max-width: 100%;
  height: auto;
  border-radius: var(--md-sys-shape-corner-small);
  box-shadow: var(--md-sys-elevation-level1);
}

.captcha-reload-btn {
  width: 40px;
  height: 40px;
  padding: 0;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  border: none;
  border-radius: var(--md-sys-shape-corner-full);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition:
    background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
  box-shadow: var(--md-sys-elevation-level2);
}

.captcha-reload-btn:hover {
  background: var(--md-sys-color-primary-container);
  transform: scale(1.1) rotate(180deg);
}

.captcha-reload-btn:active {
  transform: scale(0.95) rotate(180deg);
}

/* Auth Error */
.auth-error {
  padding: 12px 16px;
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
  border-radius: var(--md-sys-shape-corner-medium);
  border-left: 3px solid var(--md-sys-color-error);
  font-family: var(--md-sys-typescale-body-medium-font);
  font-size: var(--md-sys-typescale-body-medium-size);
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Auth Logout Button */
.auth-logout-btn {
  width: 100%;
  padding: 16px;
  background: var(--md-sys-color-error);
  color: var(--md-sys-color-on-error);
  border: none;
  border-radius: var(--md-sys-shape-corner-medium);
  font-family: var(--md-sys-typescale-label-large-font);
  font-size: var(--md-sys-typescale-label-large-size);
  font-weight: 500;
  cursor: pointer;
  transition:
    background var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    box-shadow var(--md-sys-motion-spring-expressive-fast-effects-duration) var(--md-sys-motion-spring-expressive-fast-effects),
    transform var(--md-sys-motion-spring-expressive-fast-spatial-duration) var(--md-sys-motion-spring-expressive-fast-spatial);
  box-shadow: var(--md-sys-elevation-level1);
}

.auth-logout-btn:hover {
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-on-error-container);
  box-shadow: var(--md-sys-elevation-level2);
  transform: scale(1.02);
}

.auth-logout-btn:active {
  transform: scale(0.98);
}
</style>
