# 🔧 Руководство по отладке TorrPlayer

## 🚀 Быстрый старт отладки

### 1. Соберите DEV версию:
```bash
cd /root/torrplayer-merged
./build-dev.sh
```

### 2. Скопируйте файлы на Windows:
```
build/bin/torrplayer-merged.exe
build/bin/libmpv-2.dll
```

### 3. Запустите на Windows:
- Двойной клик на `torrplayer-merged.exe`
- **Откроется консольное окно с логами!**

---

## 📋 Что смотреть в логах

### При запуске приложения:
```
✅ "TorrPlayer starting..."
✅ "Initializing BitTorrent client..."
✅ "BitTorrent client initialized successfully"
✅ "Loading torrents from database..."
✅ "Loaded X torrents from database"
```

**Если видите "Loaded 0 torrents" - БД пустая (норма для первого запуска)**

### При добавлении торрента:
```
✅ "Adding torrent: <magnet/file>"
✅ "Torrent added to client: <hash>"
✅ "Torrent saved to database: <hash>"
✅ "GetTorrents called, found X torrents"
```

### При показе торрентов:
```
✅ "GetTorrents called, found X torrents"
✅ "Checking DB for torrent <hash>"
✅ "Found in DB: Title=..., Size=..."
✅ "Torrent: <hash>, Name: <name>, Status: <status>, Size: <size>"
✅ "Returning X torrents"
```

---

## 🐛 Типичные проблемы и что смотреть

### Проблема: "Бесконечная загрузка при старте"

**Смотрите в логах:**
```
Loading torrents from database...
Loaded X torrents from database  ← Сколько загружено?
```

**Если 0:**
- БД пустая или не читается
- Проверьте файл БД в папке приложения

**Если > 0, но не показываются:**
```
GetTorrents called, found X torrents
Checking DB for torrent <hash>
Found in DB: Title=..., Size=...  ← Находит в БД?
```

### Проблема: "Торрент добавился, но не показывается"

**Смотрите:**
```
✅ Torrent saved to database: <hash>  ← Сохранился?
✅ GetTorrents called, found X torrents  ← Нашёл?
✅ Checking DB for torrent <hash>  ← Проверил БД?
✅ Found in DB: Title=..., Size=...  ← Нашёл в БД?
```

### Проблема: "После перезапуска торренты пропали"

**При старте смотрите:**
```
Loading torrents from database...
Loading torrent from DB: <hash> Title: <name>  ← Каждый торрент
Torrent loaded successfully: <hash>
Loaded X torrents from database  ← Итого
```

---

## 🛠️ DevTools (F12)

После запуска нажмите **F12** чтобы открыть DevTools:

### Console Tab:
- JavaScript ошибки
- Вызовы API
- Frontend логи

### Network Tab:
- Запросы к backend
- Время ответа
- Ошибки API

### Application Tab:
- LocalStorage
- SessionStorage
- Cookies

---

## 📝 Полезные команды

### Очистить БД (для теста):
```bash
# В папке приложения на Windows
del TorrServer.db
```

### Проверить размер БД:
```bash
dir TorrServer.db
```

### Сохранить логи в файл:
```cmd
torrplayer-merged.exe > logs.txt 2>&1
```

---

## 🔍 Что делать если проблема не решается

1. **Запустите DEV версию**
2. **Сохраните все логи из консоли**
3. **Откройте DevTools (F12) → Console**
4. **Скопируйте ошибки JavaScript**
5. **Опишите шаги для воспроизведения**

---

## 📦 Переключение между режимами

### Production → Development:
```bash
./build-dev.sh
```
- Файл станет больше (32 MB)
- Появится консоль
- Включатся DevTools

### Development → Production:
```bash
./build.sh
```
- Файл станет меньше (17 MB)
- Консоль скроется
- DevTools выключатся

---

## ⚠️ Важно

- **DEV версия работает медленнее** - это нормально
- **Консоль нельзя закрыть** - закроется и приложение
- **Логи пишутся в консоль в реальном времени**
- **F12 работает только в DEV режиме**
- **Для релиза используйте build.sh**

---

## 🎯 Текущая архитектура загрузки

```
1. Старт приложения
   ↓
2. BTServer.Connect() - инициализация BT клиента
   ↓
3. LoadTorrentsFromDB() - загрузка ВСЕХ торрентов из БД в память
   ↓
4. Frontend: GetTorrents() - запрос списка
   ↓
5. Backend: ListTorrent() объединяет память + БД
   ↓
6. GetTorrents() проверяет данные, если нет - берёт из БД
   ↓
7. Возврат списка во Frontend
   ↓
8. Автообновление каждые 5 секунд
```

**Все этапы логируются в DEV режиме!**
