# TorrPlayer Merged

Полнофункциональное десктопное приложение для потокового просмотра торрентов с красивым GUI интерфейсом в стиле стриминговых сервисов.

## ✨ Возможности

### 🎬 Торрент Streaming
- ✅ Загрузка торрентов из magnet-ссылок, .torrent файлов или хеша
- ✅ Потоковое воспроизведение без полной загрузки
- ✅ Кэширование в оперативной памяти (RAM)
- ✅ Статистика в реальном времени (скорость, пиры, кэш)
- ✅ Просмотр файлов внутри торрента
- ✅ Встроенный MPV плеер

### 🎨 Графический интерфейс
- ✅ Красивый темный UI в стиле Amazon Prime Video
- ✅ Vue.js 3 frontend с анимациями
- ✅ Адаптивный дизайн
- ✅ Карточки торрентов с информацией
- ✅ Страница настроек

### ⚙️ Настройки
- ✅ Размер кэша
- ✅ Лимиты скорости (загрузка/отдача)
- ✅ Лимит соединений
- ✅ Режим ретрекеров
- ✅ Предзагрузка кэша

### 💾 База данных
- ✅ Хранение торрентов и настроек
- ✅ История просмотров
- ✅ BoltDB + JSON

## 📋 Требования

### Для сборки:

- **Go 1.23+** - [Скачать](https://golang.org/dl/)
- **Node.js 16+** - [Скачать](https://nodejs.org/)
- **Wails CLI** - [Установка](https://wails.io/docs/gettingstarted/installation)
- **MinGW-w64** (для компиляции CGO на Windows) - [Скачать](https://www.mingw-w64.org/)

### Для запуска:

- **Windows 10/11 (amd64)**
- **libmpv-2.dll** - [Скачать](https://sourceforge.net/projects/mpv-player-windows/files/libmpv/)
  - Положите `libmpv-2.dll` в папку с `torrplayer-merged.exe`

## 🔧 Установка зависимостей

### 1. Установите Go
```bash
# Скачайте с https://golang.org/dl/
# Убедитесь что Go в PATH
go version
```

### 2. Установите Node.js
```bash
# Скачайте с https://nodejs.org/
node --version
npm --version
```

### 3. Установите Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 4. Установите MinGW-w64 (для Windows)
```bash
# Скачайте с https://www.mingw-w64.org/
# Добавьте в PATH: C:\mingw64\bin
```

Для кросс-компиляции на Linux:
```bash
# Ubuntu/Debian
sudo apt-get install mingw-w64

# Arch Linux
sudo pacman -S mingw-w64-gcc

# Fedora
sudo dnf install mingw64-gcc
```

## 🛠️ Сборка

### На Windows:

```batch
build.bat
```

### На Linux (кросс-компиляция для Windows):

```bash
bash build.sh
```

Готовое приложение будет в `build/bin/torrplayer-merged.exe`

**ВАЖНО:** После сборки скопируйте `libmpv-2.dll` в папку `build/bin/`

## 🚀 Использование

### Запуск приложения:

1. Запустите `torrplayer-merged.exe`
2. Убедитесь что `libmpv-2.dll` находится в той же папке

### Добавление торрента:

1. Нажмите кнопку "Add Torrent" в интерфейсе
2. Введите:
   - Magnet-ссылку: `magnet:?xt=urn:btih:HASH...`
   - Путь к .torrent файлу: `C:\Downloads\movie.torrent`
   - Хеш торрента (40 символов)
3. Дождитесь загрузки информации о торренте

### Просмотр файлов:

1. Кликните на карточку торрента
2. Выберите файл из списка
3. Нажмите "Play"
4. Приложение начнёт кэширование и запустит MPV плеер

### Настройки:

1. Откройте страницу настроек (иконка шестерёнки)
2. Измените параметры:
   - **Cache Size** - размер кэша в RAM (64 MB - 2 GB)
   - **Connections Limit** - макс. соединений (10-100)
   - **Download Rate** - лимит скорости загрузки (KB/s, 0 = без лимита)
   - **Upload Rate** - лимит скорости отдачи (KB/s, 0 = без лимита)
   - **Preload Cache** - предзагрузка кэша (0-100%)
   - **Retrackers Mode**:
     - 0 = не добавлять
     - 1 = добавить ретрекеры (по умолчанию)
     - 2 = удалить ретрекеры
     - 3 = заменить ретрекеры
3. Нажмите "Save Settings"

## 📁 Структура проекта

```
torrplayer-merged/
├── main.go              # Точка входа Wails приложения
├── app.go               # Backend логика (торренты, настройки)
├── player.go            # MPV интеграция (Windows)
├── player_stub.go       # MPV заглушка (не-Windows)
├── go.mod               # Go зависимости
├── wails.json           # Конфигурация Wails
├── build.bat            # Скрипт сборки для Windows
├── build.sh             # Скрипт сборки для Linux
├── mpv/
│   └── client.h         # Заголовочный файл libmpv
├── server/              # Модули из TorrServer
│   ├── torr/            # BitTorrent клиент
│   ├── settings/        # Настройки и БД
│   ├── log/             # Логирование
│   ├── utils/           # Утилиты
│   ├── mimetype/        # MIME типы
│   └── ffprobe/         # Анализ медиа
├── frontend/            # Vue.js интерфейс
│   ├── dist/            # Собранный frontend
│   ├── index.html       # HTML шаблон
│   └── package.json     # NPM зависимости
└── build/
    └── windows/
        └── icon.ico     # Иконка приложения
```

## 🎮 Управление MPV плеером

Во время воспроизведения:
- `Пробел` - пауза/воспроизведение
- `←` / `→` - перемотка на 5 секунд
- `↑` / `↓` - громкость
- `F` - полноэкранный режим
- `M` - отключить звук
- `S` - скриншот
- `Q` или закрыть окно - выход

## 🔍 Backend API (Go ↔ Vue.js)

Доступные методы:

### Торренты:
- `AddTorrent(input string) (*Torrent, error)` - добавить торрент
- `GetTorrents() []Torrent` - список торрентов
- `GetTorrentFiles(hash string) ([]TorrentFile, error)` - файлы торрента
- `GetTorrentStats(hash string) (*TorrentStats, error)` - статистика торрента
- `PlayTorrentFile(hash string, fileIndex int) error` - воспроизвести файл
- `RemoveTorrent(hash string) error` - удалить торрент

### Настройки:
- `GetSettings() *Settings` - получить настройки
- `SetSettings(s *Settings) error` - сохранить настройки

## 🐛 Отладка

### Проблемы с libmpv-2.dll:

Если MPV не запускается:

1. Убедитесь что `libmpv-2.dll` в той же папке с `.exe`
2. Проверьте версию DLL (нужна 64-битная)
3. Установите [Visual C++ Redistributable](https://aka.ms/vs/17/release/vc_redist.x64.exe)
4. Проверьте логи в DevTools (F12 в приложении)

### Проблемы со сборкой:

**Error: `gcc not found`**
- Установите MinGW-w64 и добавьте в PATH

**Error: `wails: command not found`**
- Установите Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

**Error: `mpv/client.h: No such file`**
- Убедитесь что файл `mpv/client.h` существует в проекте

**Error: `npm not found`**
- Установите Node.js и добавьте в PATH

## 📊 Технические детали

### Архитектура:

- **Frontend:** Vue.js 3 + Vite
- **Backend:** Go + Wails v2
- **BitTorrent:** anacrolix/torrent (форк от tsynik)
- **Кэш:** RAM с автоматической очисткой
- **Streaming:** Локальный HTTP сервер
- **MPV:** Интеграция через CGO и libmpv-2.dll
- **База данных:** BoltDB + JSON файлы

### Отличия от TorrServer:

**Убрано:**
- ❌ Веб-интерфейс (заменён на десктопный GUI)
- ❌ REST API (заменён на прямые вызовы функций)
- ❌ DLNA сервер
- ❌ Telegram бот

**Оставлено:**
- ✅ Торрент-клиент
- ✅ Кэш в RAM
- ✅ Стриминг
- ✅ Статистика
- ✅ База данных
- ✅ Настройки

**Добавлено:**
- ✅ Красивый GUI на Vue.js
- ✅ Wails десктопное приложение
- ✅ Улучшенная MPV интеграция

## 📝 Лицензия

Основано на:
- [TorrServer](https://github.com/YouROK/TorrServer) - GPL-3.0
- [Wails](https://wails.io/) - MIT

## 👥 Авторы

- **TorrServer:** [YouROK](https://github.com/YouROK)
- **TorrPlayer Merged:** Объединённая версия с GUI

## 🔗 Ссылки

- [TorrServer (оригинал)](https://github.com/YouROK/TorrServer)
- [Wails Framework](https://wails.io/)
- [MPV Player](https://mpv.io/)
- [libmpv для Windows](https://sourceforge.net/projects/mpv-player-windows/files/libmpv/)
- [anacrolix/torrent](https://github.com/anacrolix/torrent)
- [Vue.js 3](https://vuejs.org/)

## 🎯 Roadmap

- [ ] Поиск торрентов
- [ ] Интеграция с кинопоисками (TMDB, IMDB)
- [ ] Автоматическая загрузка постеров
- [ ] Субтитры (автопоиск и загрузка)
- [ ] Плейлисты
- [ ] Трей-иконка
- [ ] Автозапуск с системой
- [ ] Горячие клавиши
- [ ] Темы оформления
- [ ] Экспорт/импорт настроек

## 💬 Поддержка

Если у вас возникли проблемы или вопросы, создайте Issue в репозитории.

---

**Приятного просмотра! 🎬🍿**
