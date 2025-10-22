<template>
  <!-- Упрощенный вариант для contained режима -->
  <div v-if="contained" class="shape-only" :class="animationClass" :style="shapeOnlyStyle"></div>

  <!-- Полный вариант с рамкой и контейнером -->
  <div v-else class="frame" :style="frameStyle">
    <div class="loading-indicator" :style="indicatorStyle">
      <div class="container" :style="containerStyle">
        <div class="shape-container">
          <div
            class="shape"
            :class="animationClass"
            :style="shapeStyle"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed, PropType } from "vue";

type AnimationType = 'pulse' | 'smooth' | 'rotate' | 'scale' | 'morph';
type ColorScheme = 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'on-primary';

export default defineComponent({
  name: "LoadingIndicatorAdvanced",

  props: {
    // Тип анимации
    animation: {
      type: String as PropType<AnimationType>,
      default: 'morph',
      validator: (value: string) => ['pulse', 'smooth', 'rotate', 'scale', 'morph'].includes(value)
    },

    // Цветовая схема
    colorScheme: {
      type: String as PropType<ColorScheme>,
      default: 'primary',
      validator: (value: string) => ['primary', 'secondary', 'success', 'warning', 'error', 'on-primary'].includes(value)
    },

    // Размер (extra-small, small, medium, large)
    size: {
      type: String,
      default: 'medium',
      validator: (value: string) => ['extra-small', 'small', 'medium', 'large'].includes(value)
    },

    // Скорость анимации (в секундах)
    speed: {
      type: Number,
      default: 1.2
    },

    // Убрать рамку и тень (для использования внутри кнопки)
    contained: {
      type: Boolean,
      default: false
    }
  },

  setup(props) {
    // Цветовые схемы Material Design 3
    const colorSchemes = {
      primary: {
        container: 'var(--md-sys-color-primary-container, #EADDFF)',
        shape: 'var(--md-sys-color-primary, #4F378A)'
      },
      secondary: {
        container: 'var(--md-sys-color-secondary-container, #E8DEF8)',
        shape: 'var(--md-sys-color-secondary, #625B71)'
      },
      success: {
        container: '#D3F0D4',
        shape: '#1D6F42'
      },
      warning: {
        container: '#FFE7B8',
        shape: '#825500'
      },
      error: {
        container: 'var(--md-sys-color-error-container, #FFDAD6)',
        shape: 'var(--md-sys-color-error, #BA1A1A)'
      },
      'on-primary': {
        container: 'transparent',
        shape: 'var(--md-sys-color-on-primary, #FFFFFF)'
      }
    };

    // Размеры
    const sizes = {
      'extra-small': {
        frame: 24,
        indicator: 24,
        padding: 0,
        shapeSize: 24
      },
      small: {
        frame: 48,
        indicator: 32,
        padding: 8,
        shapeSize: 24
      },
      medium: {
        frame: 68,
        indicator: 48,
        padding: 10,
        shapeSize: 32
      },
      large: {
        frame: 88,
        indicator: 64,
        padding: 12,
        shapeSize: 48
      }
    };

    const currentSize = computed(() => sizes[props.size as keyof typeof sizes]);
    const currentColors = computed(() => colorSchemes[props.colorScheme as keyof typeof colorSchemes]);

    const frameStyle = computed(() => ({
      width: `${currentSize.value.frame}px`,
      height: `${currentSize.value.frame}px`,
      padding: `${currentSize.value.padding}px`
    }));

    const indicatorStyle = computed(() => ({
      width: `${currentSize.value.indicator}px`,
      height: `${currentSize.value.indicator}px`
    }));

    const indicatorClasses = computed(() => ({
      'no-border': props.contained
    }));

    const containerStyle = computed(() => ({
      background: currentColors.value.container
    }));

    const shapeStyle = computed(() => ({
      background: currentColors.value.shape,
      animationDuration: `${props.speed}s`,
      width: `${currentSize.value.shapeSize}px`,
      height: `${currentSize.value.shapeSize}px`
    }));

    // Упрощенный стиль для contained режима
    const shapeOnlyStyle = computed(() => ({
      background: currentColors.value.shape,
      animationDuration: `${props.speed}s`,
      width: `${currentSize.value.shapeSize}px`,
      height: `${currentSize.value.shapeSize}px`
    }));

    const animationClass = computed(() => `animation-${props.animation}`);

    return {
      frameStyle,
      indicatorStyle,
      indicatorClasses,
      containerStyle,
      shapeStyle,
      shapeOnlyStyle,
      animationClass
    };
  }
});
</script>

<style scoped>
/* Frame - Внешний контейнер */
.frame {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 10px;
  position: relative;
}

/* Loading indicator - Основной индикатор */
.loading-indicator {
  box-sizing: border-box;
  background: #ffffff;
  border: 1px solid #000000;
  box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25);
  flex: none;
  position: relative;
}

/* Убрать рамку и тень для contained варианта */
.loading-indicator.no-border {
  background: transparent;
  border: none;
  box-shadow: none;
}

/* Container - Контейнер с фоновым цветом */
.container {
  position: absolute;
  width: 100%;
  height: 100%;
  left: 0;
  top: 0;
  border-radius: inherit;
}

/* Shape container - Контейнер для фигуры */
.shape-container {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Shape - Базовые стили фигуры */
.shape {
  position: absolute;
  will-change: transform, width, height, border-radius;
}

/* Shape-only - Упрощенный вариант для contained режима */
.shape-only {
  position: relative;
  z-index: 1;
  will-change: transform, width, height, border-radius;
}

/* === АНИМАЦИИ === */

/* Анимация: Pulse (пульсация по кадрам из Figma) */
.animation-pulse {
  animation: loading-pulse 2s ease-in-out infinite;
}

@keyframes loading-pulse {
  0% {
    width: 33px;
    height: 34px;
  }

  14% {
    width: 34px;
    height: 33.58px;
  }

  28% {
    width: 34px;
    height: 32px;
    transform: translateY(-1px);
  }

  42% {
    width: 34px;
    height: 34px;
    transform: translateY(0);
  }

  57% {
    width: 34px;
    height: 34px;
  }

  71% {
    width: 28px;
    height: 28px;
  }

  85% {
    width: 30px;
    height: 30px;
  }

  100% {
    width: 33px;
    height: 34px;
  }
}

/* Анимация: Smooth (плавная пульсация) */
.animation-smooth {
  animation: loading-smooth 2s ease-in-out infinite;
}

@keyframes loading-smooth {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }

  50% {
    opacity: 0.8;
    transform: scale(0.85);
  }
}

/* Анимация: Rotate (вращение) */
.animation-rotate {
  animation: loading-rotate 2s linear infinite;
}

@keyframes loading-rotate {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

/* Анимация: Scale (масштабирование) */
.animation-scale {
  animation: loading-scale 2s ease-in-out infinite;
}

@keyframes loading-scale {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }

  25% {
    transform: scale(1.2);
    opacity: 0.9;
  }

  50% {
    transform: scale(0.8);
    opacity: 0.7;
  }

  75% {
    transform: scale(1.1);
    opacity: 0.85;
  }
}

/* Анимация: Morph (морфинг формы - из FileListModal) */
/* Material Design 3 Expressive - морфинг фигуры с вращением */
.animation-morph {
  animation: loading-morph 1.2s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

@keyframes loading-morph {
  0% {
    border-radius: 50%;
    transform: rotate(0deg) scale(1);
  }
  25% {
    border-radius: 20%;
    transform: rotate(90deg) scale(0.85);
  }
  50% {
    border-radius: 10%;
    transform: rotate(180deg) scale(1);
  }
  75% {
    border-radius: 20%;
    transform: rotate(270deg) scale(0.85);
  }
  100% {
    border-radius: 50%;
    transform: rotate(360deg) scale(1);
  }
}
</style>
