# Трекер времени

Приложение для отслеживания времени, затраченного на проекты.

## Возможности

- Создание проектов
- Отслеживание времени работы над проектами
- Добавление описаний к выполненным задачам
- Уведомления о необходимости сделать перерыв
- Интеграция с системным треем
- Просмотр статистики по всем проектам
- Настраиваемые параметры через флаги командной строки
- Логирование действий приложения
- Визуальная индикация активных проектов
- Удобная навигация с кнопкой "Назад"
- Интеллектуальный выбор проектов при остановке отслеживания

## Установка

```bash
go get github.com/MWT-proger/time-tracking
```

## Сборка

### Требования
- Go 1.21 или выше

### Простая сборка
```bash
go build -o time-tracker ./cmd/time-tracker
```

### Сборка с информацией о версии
```bash
go build -ldflags "-X github.com/MWT-proger/time-tracking/internal/app.BuildDate=$(date -u +%Y-%m-%d) -X github.com/MWT-proger/time-tracking/internal/app.GitCommit=$(git rev-parse --short HEAD)" -o time-tracker ./cmd/time-tracker
```

## Использование

```bash
time-tracking [флаги]
```

### Флаги командной строки

| Флаг | Описание | Значение по умолчанию |
|------|----------|------------------------|
| `-data` | Путь к файлу данных | `~/учет_времени.json` |
| `-log-dir` | Директория для логов | `~/.time-tracker/logs` |
| `-log-level` | Уровень логирования (debug, info, warn, error, fatal) | `info` |
| `-notify-time` | Время для уведомления в секундах | `1500` (25 минут) |
| `-help`, `-h` | Показать справку и выйти | - |

### Примеры

```bash
# Показать справку
time-tracking -help

# Запуск с указанием пути к файлу данных
time-tracking -data /path/to/data.json

# Запуск с указанием уровня логирования
time-tracking -log-level debug

# Запуск с указанием времени для уведомления (30 минут)
time-tracking -notify-time 1800
```

### Команды

- **Создать проект** - создание нового проекта
- **Начать отслеживание** - запуск таймера для выбранного проекта
  - Проекты отображаются с индикацией активности (▶)
  - Активные проекты отображаются в начале списка
  - Доступна кнопка "Назад" для отмены выбора
- **Остановить отслеживание** - остановка таймера с добавлением описания выполненной работы
  - Отображаются только активные проекты
  - Доступна кнопка "Назад" для отмены выбора
- **Сводка** - просмотр статистики по всем проектам
- **Выход** - завершение работы приложения

## История разработки

### v0.1.0 (Базовая версия)
- Базовая функциональность отслеживания времени
- Интерфейс командной строки
- Интеграция с системным треем
- Создание и управление проектами
- Отслеживание времени работы
- Сохранение данных в JSON-файл

### v0.2.0 (Улучшение пользовательского опыта)
- Добавлены уведомления о необходимости сделать перерыв
- Улучшен интерфейс системного трея
- Добавлена возможность просмотра статистики по проектам
- Улучшена обработка ошибок

### v0.3.0 (Конфигурация и логирование)
- Добавлена поддержка флагов командной строки
- Добавлена система логирования с разными уровнями
- Добавлена возможность настройки пути к файлу данных
- Добавлена возможность настройки времени для уведомлений
- Улучшена обработка ошибок при работе с файлами

### v0.4.0 (Улучшение надежности)
- Добавлена встроенная иконка для системного трея
- Улучшена обработка ошибок при запуске системного трея
- Добавлена автоматическое создание директорий для данных и логов
- Улучшена документация

### v0.5.0 (Улучшение интерфейса)
- Добавлена кнопка "Назад" в список проектов
- Добавлена визуальная индикация активных проектов (▶)
- Добавлена сортировка проектов (активные отображаются первыми)
- При остановке отслеживания показываются только активные проекты
- Улучшена обработка ошибок при выборе проектов

### v0.5.1 (Улучшение валидации)
- Добавлена проверка на пустое имя при создании проекта
- Добавлена проверка на уникальность имени проекта
- Улучшена обработка ошибок при создании проекта

### v0.6.0 (Управление спринтами проектов)
- Добавлена поддержка спринтов для проектов
- Добавлено меню управления спринтами
- Возможность создания новых спринтов
- Возможность выбора активного спринта
- Просмотр статистики по спринтам
- Учет времени в разрезе спринтов

### v0.7.0 (Реорганизация интерфейса)
- Переработана структура меню для более логичной работы с проектами
- Добавлено меню управления проектом
- Улучшена навигация между проектами и их функциями
- Добавлена возможность выбора спринта при начале отслеживания
- Улучшена статистика по проектам и спринтам

### v0.8.0 (Реорганизация кодовой базы)
- Разделение большого файла app.go на логические модули
- Создание отдельных обработчиков для различных функций
- Улучшение структуры проекта
- Повышение читаемости и поддерживаемости кода
- Устранение циклических зависимостей между пакетами

## Планы на будущее

### Улучшение отображения времени
- Отображение времени в формате "часы минуты" (например, 5 ч 34 мин)
- Показ текущего времени работы в системном трее
- Отображение ожидаемого времени выполнения задачи и сравнение с фактическим

### Улучшение управления проектами
- Добавление этапов проектов с отдельным учетом времени
- Возможность архивирования проектов
- Отображение только активных проектов в статистике
- Показ списка последних действий по проекту
- Отображение текущей задачи в системном трее

### Финансовый учет
- Добавление стоимости часа работы
- Расчет стоимости затраченного времени
- Учет полученных платежей и расчет остатка
- Формирование финансовых отчетов

### Улучшение интерфейса
- Добавление графического интерфейса
- В списке "Остановить отслеживание" показывать только активные проекты
- Добавление графиков и диаграмм для визуализации затраченного времени

### Экспорт и интеграция
- Экспорт данных в различные форматы (CSV, Excel)
- Интеграция с календарем
- Добавление тегов для задач
- Поддержка нескольких пользователей
- Синхронизация данных между устройствами

## Вклад в проект

Мы приветствуем вклад в развитие проекта! Если у вас есть идеи или вы хотите помочь с реализацией запланированных функций, пожалуйста:

1. Создайте форк репозитория
2. Создайте ветку для вашей функции (`git checkout -b feature/amazing-feature`)
3. Зафиксируйте ваши изменения (`git commit -m 'Add some amazing feature'`)
4. Отправьте изменения в ваш форк (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

Также вы можете помочь, сообщая о проблемах или предлагая идеи в разделе Issues.

## Лицензия

MIT
