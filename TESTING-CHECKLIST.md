# ✅ Чеклист для тестирования

## 🎯 Ваш тест-кейс

### Сценарий:
1. ✅ Открыл приложение
2. ✅ Добавил торрент
3. ✅ Торрент появился
4. ✅ Закрыл приложение
5. ✅ Подождал 10 сек
6. ✅ Открыл приложение
7. ❓ **Что должно быть:** Торренты сразу показываются

---

## 🔧 Как проверить (DEV режим)

### 1. Соберите DEV версию:
```bash
./build-dev.sh
```

### 2. Скопируйте на Windows:
```
build/bin/torrplayer-merged.exe (32 MB)
build/bin/libmpv-2.dll (113 MB)
```

### 3. Первый запуск (чистая БД):
```
Запустите torrplayer-merged.exe

В КОНСОЛИ должно быть:
✅ "TorrPlayer starting..."
✅ "Loading torrents from database..."
✅ "Loaded 0 torrents from database"  ← БД пустая - ОК!
✅ "GetTorrents called, found 0 torrents"
✅ "Returning 0 torrents"

В ПРИЛОЖЕНИИ:
✅ Пустой список - ОК!
```

### 4. Добавьте торрент:
```
Нажмите "+" → Добавьте magnet/файл

В КОНСОЛИ должно быть:
✅ "Adding torrent: ..."
✅ "Torrent added to client: <hash>"
✅ "Torrent saved to database: <hash>"
✅ (В фоне) "Metadata received for: <hash> - <name>"
✅ "GetTorrents called, found 1 torrents"
✅ "Torrent: <hash>, Name: <name>, Status: loading/ready, Size: ..."

В ПРИЛОЖЕНИИ:
✅ Торрент появился СРАЗУ!
```

### 5. Закройте приложение:
```
Закройте окно приложения
(Консоль тоже закроется)
```

### 6. Подождите 10 секунд:
```
(Просто ждём)
```

### 7. Откройте снова:
```
Запустите torrplayer-merged.exe

В КОНСОЛИ должно быть:
✅ "TorrPlayer starting..."
✅ "Loading torrents from database..."
✅ "Loading torrent from DB: <hash> Title: <name>"  ← Загружает!
✅ "Torrent loaded successfully: <hash>"
✅ "Loaded 1 torrents from database"  ← Загрузился!
✅ "GetTorrents called, found 1 torrents"
✅ "Torrent: <hash>, Name: <name>, Status: ..., Size: ..."

В ПРИЛОЖЕНИИ:
✅ Торрент СРАЗУ показывается!
```

---

## 🐛 Что делать если НЕ работает

### Если при перезапуске "Loaded 0 torrents":

**Проверьте:**
```
1. В папке приложения есть файл TorrServer.db?
   → Нет: БД не создалась
   → Есть: Размер > 0?

2. В логах при добавлении было:
   "Torrent saved to database: <hash>"?
   → Нет: Не сохранилось в БД
   → Да: Проблема с чтением

3. При запуске есть:
   "Loading torrents from database..."?
   → Нет: Функция не вызывается
   → Да: Идём дальше

4. После этого есть:
   "Loading torrent from DB: <hash>"?
   → Нет: Торрента нет в БД
   → Да: Загружается
```

### Если при перезапуске "Loaded 1 torrents" но не показывается:

**Проверьте:**
```
1. "GetTorrents called, found 1 torrents"?
   → Нет: Frontend не запрашивает
   → Да: Backend отвечает

2. "Checking DB for torrent <hash>"?
   → Есть: Проверяет БД
   → "Found in DB: Title=..., Size=..."?
      → Да: Нашёл данные

3. "Torrent: <hash>, Name: <name>, Status: <status>, Size: <size>"?
   → Name пустое? Status = loading?
   → Проблема с данными

4. "Returning 1 torrents"?
   → Да: Backend вернул данные
   → Проблема во Frontend
```

---

## 📊 Ожидаемые логи (полный цикл)

### Первый запуск:
```
[INFO] TorrPlayer starting...
[INFO] Initializing BitTorrent client...
[INFO] BitTorrent client initialized successfully
[INFO] Loading torrents from database...
[TLog] Loading torrents from DB...
[TLog] Loaded 0 torrents from DB
[INFO] Loaded 0 torrents from database
[INFO] GetTorrents called, found 0 torrents
[INFO] Returning 0 torrents
```

### Добавление торрента:
```
[INFO] Adding torrent: magnet:?xt=...
[TLog] New torrent <hash>
[INFO] Torrent added to client: <hash>
[TLog] save to db: <hash>
[INFO] Torrent saved to database: <hash>
[INFO] GetTorrents called, found 1 torrents
[INFO] Checking DB for torrent <hash> (name empty: false, size zero: true)
[TLog] Loading torrents from DB...
[INFO] Found in DB: Title=<name>, Size=0
[INFO] Torrent: <hash>, Name: <name>, Status: loading, Size: 0
[INFO] Returning 1 torrents

(Через несколько секунд в фоне):
[INFO] Metadata received for: <hash> - <name>
[TLog] save to db: <hash>
```

### Второй запуск:
```
[INFO] TorrPlayer starting...
[INFO] Initializing BitTorrent client...
[INFO] BitTorrent client initialized successfully
[INFO] Loading torrents from database...
[TLog] Loading torrents from DB...
[TLog] Loading torrent from DB: <hash> Title: <name>
[TLog] New torrent <hash>
[TLog] Torrent loaded successfully: <hash>
[TLog] Loaded 1 torrents from DB
[INFO] Loaded 1 torrents from database
[INFO] GetTorrents called, found 1 torrents
[INFO] Torrent: <hash>, Name: <name>, Status: ready, Size: <size>
[INFO] Returning 1 torrents
```

---

## 🎉 Если всё работает

**Вы должны видеть:**
1. ✅ Торренты загружаются при старте
2. ✅ Торренты показываются мгновенно
3. ✅ Торренты сохраняются между сессиями
4. ✅ Никакой "бесконечной загрузки"

**Можете собрать production версию:**
```bash
./build.sh
```

---

## 📁 Файлы для отладки

- `build-dev.sh` - Скрипт сборки DEV режима
- `BUILD-MODES.md` - Описание режимов сборки
- `DEBUG-GUIDE.md` - Подробное руководство по отладке
- `TESTING-CHECKLIST.md` - Этот файл

---

## 💡 Полезные советы

1. **Всегда начинайте с чистой БД** для теста
2. **Сохраняйте логи** если что-то не работает
3. **Используйте F12** для JavaScript ошибок
4. **Проверяйте наличие TorrServer.db** в папке приложения
5. **Dev версия медленнее** - это нормально
