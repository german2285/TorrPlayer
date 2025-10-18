// === Утилиты для работы с цветами и темой ===

// Конвертация HEX в HSL
export const hexToHSL = (hex: string): { h: number; s: number; l: number } => {
  // Убрать #
  hex = hex.replace('#', '')

  // Конвертировать в RGB
  const r = parseInt(hex.substring(0, 2), 16) / 255
  const g = parseInt(hex.substring(2, 4), 16) / 255
  const b = parseInt(hex.substring(4, 6), 16) / 255

  const max = Math.max(r, g, b)
  const min = Math.min(r, g, b)
  let h = 0
  let s = 0
  const l = (max + min) / 2

  if (max !== min) {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)

    switch (max) {
      case r: h = ((g - b) / d + (g < b ? 6 : 0)) / 6; break
      case g: h = ((b - r) / d + 2) / 6; break
      case b: h = ((r - g) / d + 4) / 6; break
    }
  }

  return { h: h * 360, s: s * 100, l: l * 100 }
}

// Конвертация HSL в HEX
export const hslToHex = (h: number, s: number, l: number): string => {
  s /= 100
  l /= 100

  const c = (1 - Math.abs(2 * l - 1)) * s
  const x = c * (1 - Math.abs((h / 60) % 2 - 1))
  const m = l - c / 2
  let r = 0, g = 0, b = 0

  if (h >= 0 && h < 60) { r = c; g = x; b = 0 }
  else if (h >= 60 && h < 120) { r = x; g = c; b = 0 }
  else if (h >= 120 && h < 180) { r = 0; g = c; b = x }
  else if (h >= 180 && h < 240) { r = 0; g = x; b = c }
  else if (h >= 240 && h < 300) { r = x; g = 0; b = c }
  else if (h >= 300 && h < 360) { r = c; g = 0; b = x }

  const toHex = (n: number) => {
    const hex = Math.round((n + m) * 255).toString(16)
    return hex.length === 1 ? '0' + hex : hex
  }

  return `#${toHex(r)}${toHex(g)}${toHex(b)}`
}

// Генерация полной M3 палитры из одного цвета (DARK THEME)
export const generateM3Palette = (sourceColor: string) => {
  const base = hexToHSL(sourceColor)

  // Primary - яркий акцентный цвет для темной темы
  const primary = sourceColor
  const onPrimary = '#000000' // Темный текст на ярком primary
  const primaryContainer = hslToHex(base.h, Math.max(base.s - 10, 40), 30) // Темный контейнер
  const onPrimaryContainer = hslToHex(base.h, Math.max(base.s, 80), 90) // Светлый текст

  // Secondary - сдвиг оттенка на +30-40 градусов
  const secHue = (base.h + 35) % 360
  const secSat = Math.max(base.s - 20, 30)
  const secondary = hslToHex(secHue, secSat, 60) // Средняя яркость
  const onSecondary = '#000000'
  const secondaryContainer = hslToHex(secHue, secSat, 30) // Темный контейнер
  const onSecondaryContainer = hslToHex(secHue, secSat + 20, 90) // Светлый текст

  // Tertiary - сдвиг оттенка на +60-90 градусов
  const terHue = (base.h + 75) % 360
  const terSat = Math.max(base.s - 15, 35)
  const tertiary = hslToHex(terHue, terSat, 65)
  const onTertiary = '#000000'
  const tertiaryContainer = hslToHex(terHue, terSat, 30) // Темный контейнер
  const onTertiaryContainer = hslToHex(terHue, terSat + 20, 90) // Светлый текст

  // Surface - ТЕМНЫЕ нейтральные тона (Dark Mode)
  const surfaceHue = base.h
  const surface = hslToHex(surfaceHue, 8, 8) // Очень темный surface
  const onSurface = hslToHex(surfaceHue, 8, 90) // Светлый текст
  const surfaceVariant = hslToHex(surfaceHue, 10, 25) // Темный variant
  const onSurfaceVariant = hslToHex(surfaceHue, 10, 80) // Светлый текст на variant

  // Surface containers - градация темных оттенков
  const surfaceContainerLowest = hslToHex(surfaceHue, 8, 6) // Самый темный
  const surfaceContainerLow = hslToHex(surfaceHue, 8, 11)
  const surfaceContainer = hslToHex(surfaceHue, 8, 13)
  const surfaceContainerHigh = hslToHex(surfaceHue, 8, 15)
  const surfaceContainerHighest = hslToHex(surfaceHue, 8, 17) // Самый светлый из темных

  // Outline - границы
  const outline = hslToHex(surfaceHue, 8, 40)
  const outlineVariant = hslToHex(surfaceHue, 10, 30)

  // Background
  const background = surface
  const onBackground = onSurface

  return {
    primary, onPrimary, primaryContainer, onPrimaryContainer,
    secondary, onSecondary, secondaryContainer, onSecondaryContainer,
    tertiary, onTertiary, tertiaryContainer, onTertiaryContainer,
    surface, onSurface, surfaceVariant, onSurfaceVariant,
    surfaceContainer, surfaceContainerHigh, surfaceContainerHighest,
    surfaceContainerLow, surfaceContainerLowest,
    outline, outlineVariant,
    background, onBackground
  }
}

// Применить цветовую тему
export const applyThemeColor = (color: string): void => {
  // Генерировать полную M3 палитру
  const palette = generateM3Palette(color)

  // Применить все цвета M3
  const root = document.documentElement

  // Primary
  root.style.setProperty('--md-sys-color-primary', palette.primary)
  root.style.setProperty('--md-sys-color-on-primary', palette.onPrimary)
  root.style.setProperty('--md-sys-color-primary-container', palette.primaryContainer)
  root.style.setProperty('--md-sys-color-on-primary-container', palette.onPrimaryContainer)

  // Secondary
  root.style.setProperty('--md-sys-color-secondary', palette.secondary)
  root.style.setProperty('--md-sys-color-on-secondary', palette.onSecondary)
  root.style.setProperty('--md-sys-color-secondary-container', palette.secondaryContainer)
  root.style.setProperty('--md-sys-color-on-secondary-container', palette.onSecondaryContainer)

  // Tertiary
  root.style.setProperty('--md-sys-color-tertiary', palette.tertiary)
  root.style.setProperty('--md-sys-color-on-tertiary', palette.onTertiary)
  root.style.setProperty('--md-sys-color-tertiary-container', palette.tertiaryContainer)
  root.style.setProperty('--md-sys-color-on-tertiary-container', palette.onTertiaryContainer)

  // Surface
  root.style.setProperty('--md-sys-color-surface', palette.surface)
  root.style.setProperty('--md-sys-color-on-surface', palette.onSurface)
  root.style.setProperty('--md-sys-color-surface-variant', palette.surfaceVariant)
  root.style.setProperty('--md-sys-color-on-surface-variant', palette.onSurfaceVariant)
  root.style.setProperty('--md-sys-color-surface-container', palette.surfaceContainer)
  root.style.setProperty('--md-sys-color-surface-container-high', palette.surfaceContainerHigh)
  root.style.setProperty('--md-sys-color-surface-container-highest', palette.surfaceContainerHighest)
  root.style.setProperty('--md-sys-color-surface-container-low', palette.surfaceContainerLow)
  root.style.setProperty('--md-sys-color-surface-container-lowest', palette.surfaceContainerLowest)

  // Outline
  root.style.setProperty('--md-sys-color-outline', palette.outline)
  root.style.setProperty('--md-sys-color-outline-variant', palette.outlineVariant)

  // Background
  root.style.setProperty('--md-sys-color-background', palette.background)
  root.style.setProperty('--md-sys-color-on-background', palette.onBackground)
}

// Загрузить и применить сохраненную тему из backend
// Вызывается из App.vue при монтировании
export const loadSavedTheme = async (): Promise<void> => {
  try {
    // Динамический импорт для избежания циркулярной зависимости
    const { GetSettings } = await import('../../wailsjs/go/app/App')
    const settings = await GetSettings()
    if (settings && settings.themeColor) {
      applyThemeColor(settings.themeColor)
    }
  } catch (error) {
    console.error('Failed to load theme from backend:', error)
    // Fallback к дефолтному цвету
    applyThemeColor('#6750A4')
  }
}
