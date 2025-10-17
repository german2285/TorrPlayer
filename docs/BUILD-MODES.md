# Build Modes

## Production Build (build.sh)

**Использование:**
```bash
./build.sh
```

**Характеристики:**
- ✅ Оптимизирован для производства
- ✅ Размер: ~17 MB
- ✅ Без консоли
- ✅ Без DevTools
- ✅ Быстрая загрузка

**Результат:**
- `build/bin/torrplayer-merged.exe` (17 MB)
- `build/bin/libmpv-2.dll` (113 MB)

---

## Development Build (build-dev.sh)

**Использование:**
```bash
./build-dev.sh
```

**Характеристики:**
- 🔧 Режим разработки/отладки
- 🔧 Размер: ~32 MB (без оптимизаций)
- ✅ **КОНСОЛЬ ВИДНА** - все логи в консольном окне
- ✅ **DevTools включены** - F12 для открытия
- ✅ **Подробное логирование**

**Что увидите в консоли:**
```
TorrPlayer starting...
Initializing BitTorrent client...
BitTorrent client initialized successfully
Loading torrents from database...
Loading torrent from DB: <hash> Title: <name>
Torrent loaded successfully: <hash>
Loaded X torrents from database
GetTorrents called, found X torrents
Checking DB for torrent <hash>
Found in DB: Title=..., Size=...
Torrent: <hash>, Name: <name>, Status: <status>, Size: <size>
Returning X torrents
```

**Результат:**
- `build/bin/torrplayer-merged.exe` (32 MB - DEBUG)
- `build/bin/libmpv-2.dll` (113 MB)

---

## Когда использовать какой режим?

### Production (build.sh)
- ✅ Для конечных пользователей
- ✅ Для релиза
- ✅ Когда всё работает
- ✅ Меньший размер файла

### Development (build-dev.sh)
- 🔧 Для отладки
- 🔧 Когда что-то не работает
- 🔧 Чтобы увидеть все логи
- 🔧 Для разработки новых функций
- 🔧 Чтобы понять что происходит

---

## Отладка проблем

**Если что-то не работает:**

1. Соберите DEV версию:
   ```bash
   ./build-dev.sh
   ```

2. Запустите на Windows:
   ```cmd
   cd build\bin
   torrplayer-merged.exe
   ```

3. Смотрите логи в консоли:
   - При старте: загрузка из БД
   - При добавлении: сохранение в БД
   - При показе: поиск данных
   - Все ошибки и предупреждения

4. Нажмите F12 для DevTools:
   - Console: JavaScript логи
   - Network: запросы к backend
   - Elements: структура DOM

---

## Различия в поведении

| Функция | Production | Development |
|---------|-----------|-------------|
| Размер .exe | 17 MB | 32 MB |
| Консоль | Скрыта | Видна |
| DevTools | Нет | Да (F12) |
| Логирование | Минимальное | Подробное |
| Оптимизация | Да | Нет |
| Скорость | Быстрее | Медленнее |
| Отладка | Сложнее | Легче |

---

## Файлы

```
torrplayer-merged/
├── build.sh          ← Production build
├── build-dev.sh      ← Development build (NEW!)
├── build/
│   └── bin/
│       ├── torrplayer-merged.exe
│       └── libmpv-2.dll
```

---

## Примечания

- **Оба режима** создают одинаковые файлы в `build/bin/`
- **DEV версия** перезаписывает production версию и наоборот
- **libmpv-2.dll** одинаковая для обоих режимов
- **Для релиза** всегда используйте `build.sh`
- **Для отладки** всегда используйте `build-dev.sh`
