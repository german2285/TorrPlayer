# 🔍 Анализ работы RuTracker.org - Результаты

## 📊 Краткое резюме

**RuTracker НЕ использует JavaScript API для загрузки данных о торрентах.**

Вся информация возвращается в виде **готового HTML** при первом запросе к `tracker.php`.
Данные встроены непосредственно в HTML таблицу и парсятся браузером на стороне клиента.

---

## 🎯 Основной запрос

### URL
```
POST https://rutracker.net/forum/tracker.php?nm=%D0%92%D0%BE%20%D0%B2%D1%81%D0%B5%20%D1%82%D1%8F%D0%B6%D0%BA%D0%B8%D0%B5
```

### Параметры запроса

**Query параметры (в URL):**
- `nm` - поисковый запрос (URL-encoded)

**POST данные (в теле):**
```
max=1&nm=%C2%EE+%E2%F1%E5+%F2%FF%E6%EA%E8%E5
```
- `max=1` - флаг максимального количества результатов или режима отображения
- `nm` - поисковый запрос в кодировке Windows-1251

### Важные заголовки

```http
Content-Type: application/x-www-form-urlencoded
Accept-Language: ru,en;q=0.9
Referer: https://rutracker.net/forum/index.php
User-Agent: Mozilla/5.0 ...
```

---

## 📦 Структура ответа

**Тип:** `text/html; charset=Windows-1251`
**Размер:** ~247 KB

### HTML структура данных

Торренты находятся в таблице с классом `.hl-tr`:

```html
<tr id="trs-tr-6754594" class="tCenter hl-tr" data-topic_id="6754594">
    <td class="row1 t-ico">...</td>
    <td class="row1 t-ico" title="проверено">...</td>
    <td class="row1 f-name-col">
        <div class="f-name">
            <a class="gen f ts-text" href="tracker.php?f=266">Категория</a>
        </div>
    </td>
    <td class="row4 med tLeft t-title-col tt">
        <div class="wbr t-title">
            <a data-topic_id="6754594" class="med tLink tt-text ts-text hl-tags bold" 
               href="viewtopic.php?t=6754594">Название торрента</a>
        </div>
    </td>
    <td class="row1 u-name-col">
        <div class="wbr u-name">
            <a class="med ts-text" href="tracker.php?pid=50671503">Автор</a>
        </div>
    </td>
    <td class="row4 small nowrap tor-size" data-ts_text="201885959599">
        <a class="small tr-dl dl-stub" href="dl.php?t=6754594">188.02 GB ↓</a>
    </td>
    <td class="row4 nowrap" data-ts_text="4">
        <b class="seedmed">4</b>  <!-- Сиды -->
    </td>
    <td class="row4 leechmed bold" title="Личи">3</td>  <!-- Личи -->
    <td class="row4 small number-format">107</td>  <!-- Скачивания -->
    <td class="row4 small nowrap" data-ts_text="1759956075">
        <p>8-Окт-25</p>  <!-- Дата -->
    </td>
</tr>
```

---

## 🔧 Как получить данные

### Вариант 1: Парсинг HTML (необходим)

Так как API нет, единственный способ - парсить HTML:

```python
import requests
from bs4 import BeautifulSoup
from urllib.parse import quote

# Поисковый запрос
search_query = "Во все тяжкие"

# Кодируем в Windows-1251 для POST данных
post_data_encoded = f"max=1&nm={quote(search_query.encode('windows-1251'), safe='')}"

# Параметры запроса
url = f"https://rutracker.net/forum/tracker.php?nm={quote(search_query)}"
headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    'Content-Type': 'application/x-www-form-urlencoded',
    'Referer': 'https://rutracker.net/forum/index.php'
}

# Делаем POST запрос
response = requests.post(url, data=post_data_encoded, headers=headers)

# Парсим HTML
soup = BeautifulSoup(response.content, 'html.parser', from_encoding='windows-1251')

# Находим все торренты
torrents = soup.find_all('tr', class_='hl-tr')

for torrent in torrents:
    topic_id = torrent.get('data-topic_id')
    
    # Категория
    category = torrent.find('td', class_='f-name-col').find('a').text.strip()
    
    # Название
    title = torrent.find('a', class_='tLink').text.strip()
    
    # Автор
    author = torrent.find('td', class_='u-name-col').find('a').text.strip()
    
    # Размер
    size = torrent.find('td', class_='tor-size').find('a').text.strip()
    
    # Сиды/Личи
    seeds = torrent.find('b', class_='seedmed').text.strip()
    leeches = torrent.find('td', class_='leechmed').text.strip()
    
    print(f"{topic_id}: {title}")
    print(f"  Категория: {category}")
    print(f"  Размер: {size}, S/L: {seeds}/{leeches}")
```

### Вариант 2: Использование существующих библиотек

Проверьте наличие готовых Python библиотек для работы с RuTracker:
- `rutracker-api` (если существует)
- `torrent-parser` + кастомный парсер

---

## ⚠️ Важные особенности

1. **Кодировка Windows-1251** - необходимо корректно кодировать/декодировать текст
2. **Нет JSON API** - все данные только в HTML
3. **POST параметры дублируются** - и в URL и в теле запроса
4. **Cloudflare защита** - возможна блокировка при частых запросах
5. **Cookies могут быть важны** - для авторизации и доступа к некоторым разделам
6. **Rate limiting** - необходимо делать паузы между запросами

---

## 🎯 Структура данных торрента

Каждая строка таблицы содержит:

| Поле | CSS селектор | Описание |
|------|-------------|----------|
| Topic ID | `tr[data-topic_id]` | Уникальный идентификатор |
| Категория | `.f-name-col > .f-name > a` | Раздел форума |
| Название | `.t-title-col .tLink` | Полное название раздачи |
| Автор | `.u-name-col a` | Имя автора раздачи |
| Размер | `.tor-size a` | Размер файлов |
| Сиды | `.seedmed` | Количество сидов |
| Личи | `.leechmed[title="Личи"]` | Количество личеров |
| Скачивания | `.number-format` | Всего скачиваний |
| Дата | `.small.nowrap p` | Дата последнего обновления |
| Ссылка | `.tLink[href]` | Ссылка на торрент |

---

## 💡 Рекомендации

1. **Используйте парсинг HTML** - другого способа нет
2. **Добавьте обработку ошибок** - для Cloudflare защиты
3. **Кэшируйте результаты** - чтобы не перегружать сервер
4. **Используйте сессии** - для сохранения cookies
5. **Добавьте задержки** - между запросами (2-5 секунд)
6. **Проверяйте robots.txt** - соблюдайте правила сайта

---

## 📝 Дополнительные эндпоинты

- `viewtopic.php?t={topic_id}` - страница торрента
- `dl.php?t={topic_id}` - скачивание .torrent файла
- `tracker.php?f={forum_id}` - фильтрация по разделу
- `tracker.php?pid={user_id}` - раздачи пользователя

---

## ⚖️ Юридические аспекты

⚠️ **Внимание:** RuTracker может быть заблокирован в некоторых странах.
Использование VPN может быть необходимо для доступа к сайту.
Убедитесь, что вы соблюдаете законы вашей юрисдикции при использовании торрентов.
