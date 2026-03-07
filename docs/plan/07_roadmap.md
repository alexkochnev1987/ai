# Roadmap — Этапы реализации

## Общий план

```
Фаза 1 — MVP        (Спринты 1–3,  ~6 недель)
Фаза 2 — Core       (Спринты 4–5,  ~4 недели)
Фаза 3 — AI         (Спринты 6–7,  ~4 недели)
Фаза 4 — Social     (Спринты 8–9,  ~4 недели)
─────────────────────────────────────────────
Итого:  ~18 недель до полного продукта
App Store (MVP):    после Спринта 3
```

---

## Фаза 1 — MVP

### Спринт 1 — Фундамент (2 нед)

**Цель:** Запустить приложение с Firebase и локальным хранилищем

**Бэкенд / Инфраструктура:**
- [ ] Создать Firebase проект (dev + prod)
- [ ] Настроить Firestore + Security Rules
- [ ] Terraform: Cloud Run, Secret Manager, Firestore
- [ ] GitHub Actions: CI/CD пайплайн
- [ ] Firebase Auth: Apple Sign-In + Google Sign-In

**Фронтенд:**
- [ ] Expo проект + React Navigation (tabs + stack)
- [ ] Firebase SDK + WatermelonDB установка и конфигурация
- [ ] WatermelonDB схема: tasks, habits, habit_logs, sync_queue
- [ ] SyncService: push несинхронизированных записей
- [ ] NetInfo listener: триггер синхронизации при появлении сети
- [ ] Zustand stores: useTaskStore, useSyncStore, useUIStore
- [ ] AuthService: вход / выход
- [ ] WelcomeScreen + OnboardingScreen (5 шагов)

**Готовность:** Пользователь может войти, данные хранятся локально и синхронизируются

---

### Спринт 2 — Задачи (2 нед)

**Цель:** Полноценный CRUD задач с офлайн поддержкой

**Фронтенд:**
- [ ] TaskListScreen: список с секциями (просрочено/сегодня/завтра/позже)
- [ ] TaskCard: свайп вправо (выполнить), свайп влево (меню)
- [ ] TaskFormScreen (modal): все поля создания/редактирования
- [ ] TaskDetailScreen: подзадачи, редактирование инлайн
- [ ] Drag & drop сортировка
- [ ] Фильтры: теги, приоритет, дедлайн
- [ ] Поиск по задачам
- [ ] Повторяющиеся задачи: автосоздание следующей копии
- [ ] Трекер времени: кнопка "Начать работу" + фоновый таймер

**Синхронизация:**
- [ ] Pull изменений с сервера (другие устройства)
- [ ] Реалтайм подписка Firestore (online)
- [ ] Конфликты: last-write-wins по updatedAt

**Готовность:** Полноценное управление задачами, офлайн + онлайн

---

### Спринт 3 — Привычки + Dashboard + Уведомления (2 нед)

**Цель:** Привычки со streak, Dashboard, push уведомления → App Store MVP

**Фронтенд:**
- [ ] HabitListScreen: список привычек, DateStrip, чекбоксы
- [ ] HabitDetailScreen: тепловая карта (365 дней), статистика
- [ ] HabitFormScreen (modal): создание привычки
- [ ] Streak алгоритм: grace day логика
- [ ] DashboardScreen: DailyProgressCard, MITBlock, HabitStreaksRow
- [ ] FAB быстрые действия
- [ ] MoodCheckCard (заглушка — данные без AI)
- [ ] ProfileScreen + SettingsScreen

**Бэкенд (Cloud Run — минимум):**
- [ ] FastAPI проект + Dockerfile
- [ ] Деплой на Cloud Run
- [ ] POST /notifications/schedule-streak-check
- [ ] POST /notifications/daily-brief (без AI — простой текст)
- [ ] Firebase FCM: отправка push уведомлений

**Уведомления:**
- [ ] Streak под угрозой (21:00)
- [ ] Дедлайн задачи (за 2 часа и за 1 день)
- [ ] Напоминание о настроении (настраиваемое время)
- [ ] Cloud Scheduler: ежедневные задачи

**Готовность:** MVP — публикация в App Store TestFlight

---

## Фаза 2 — Core

### Спринт 4 — Pomodoro + Трекер времени (2 нед)

**Фронтенд:**
- [ ] PomodoroScreen: круговой таймер, выбор задачи, тип сессии
- [ ] Фоновый таймер (AppState: foreground/background)
- [ ] Live Activity (iOS 16+): таймер на экране блокировки
- [ ] Звуковые уведомления при окончании сессии
- [ ] Счётчик помидоров на TaskCard (🍅🍅🍅)
- [ ] Логирование времени в задачу (actualMinutes)
- [ ] Настройки: длительность work/break
- [ ] Статистика Pomodoro в Analytics

---

### Спринт 5 — Календарь + Apple Health (2 нед)

**Фронтенд:**
- [ ] CalendarScreen: day/week/month view
- [ ] Цветовое кодирование по источнику
- [ ] Создание события из приложения
- [ ] IntegrationsScreen: подключение Calendar и Health

**Бэкенд:**
- [ ] Google OAuth: хранение + refresh tokens в Firestore
- [ ] GET /calendar/google/events
- [ ] POST /calendar/google/sync (двусторонняя)
- [ ] Apple Calendar: EventKit на фронтенде (нет бэкенда)

**Apple Health:**
- [ ] HealthKit: запрос разрешений (шаги, сон, активность)
- [ ] Чтение данных: react-native-health
- [ ] Запись: Pomodoro → Mindful Minutes

---

## Фаза 3 — AI

### Спринт 6 — AI Orchestrator + Plan (2 нед)

**Бэкенд:**
- [ ] Provider Pattern: ClaudeProvider, OpenAIProvider, GeminiProvider
- [ ] AIOrchestrator: маршрутизация, fallback при ошибке провайдера
- [ ] PlannerAgent: промпт + парсинг ответа
- [ ] POST /ai/plan: план дня со слотами
- [ ] POST /ai/prioritize: приоритизация backlog
- [ ] Server-Sent Events: /ai/stream

**Фронтенд:**
- [ ] AIBriefCard на Dashboard
- [ ] Экран "Принять план": слоты дня, редактирование
- [ ] React Query: кэш плана (5 мин), инвалидация при изменении задач
- [ ] Обработка офлайн: показываем вчерашний план из кэша

---

### Спринт 7 — Аналитика + AI Insights (2 нед)

**Бэкенд:**
- [ ] InsightAgent: анализ данных за период
- [ ] GET /ai/insights: инсайты + рекомендации
- [ ] POST /notifications/weekly-report: еженедельный отчёт

**Фронтенд:**
- [ ] AnalyticsScreen: все графики (Victory Native)
- [ ] ProductivityHeatmap: матрица день×час
- [ ] MoodTrendChart: линейный график
- [ ] CategoryPieChart: время по категориям
- [ ] Трекер настроения: MoodCheckModal + история
- [ ] AIInsightsCard: проактивные подсказки на Dashboard
- [ ] WeeklyReportScreen: детальный еженедельный отчёт

---

## Фаза 4 — Social

### Спринт 8–9 — Совместные проекты (4 нед)

**Firestore:**
- [ ] projects коллекция + Security Rules для участников
- [ ] Реалтайм подписки на изменения проекта

**Фронтенд:**
- [ ] ProjectListScreen: список проектов пользователя
- [ ] ProjectScreen: Kanban доска (drag & drop между колонками)
- [ ] Назначение задач участникам
- [ ] ActivityFeed: история изменений
- [ ] Инвайт по ссылке (Dynamic Links)
- [ ] Уведомления: FCM при изменении задачи другим участником

---

## Приоритеты и риски

### Критические риски

| Риск | Вероятность | Митигация |
|------|-------------|-----------|
| Apple HealthKit отклонит приложение | Средняя | Сделать HealthKit полностью опциональным |
| Высокие расходы на AI API | Высокая | Кэширование, rate limiting, дешёвый провайдер для инсайтов |
| Синхронизация конфликты данных | Средняя | Тесты edge-cases, grace period 5 мин для конфликтов |
| App Store Review > 2 недель | Средняя | TestFlight для бета, подготовить документацию заранее |

### Технический долг — разобрать после MVP

- E2E тесты (Detox)
- Performance: виртуализация списков при > 500 задачах
- Шифрование локальной SQLite базы
- GDPR: экспорт и удаление данных

---

## Definition of Done для каждого спринта

- [ ] Юнит-тесты на бизнес-логику (streak, sync conflicts, recurrence)
- [ ] Ручное тестирование офлайн сценария (режим самолёта)
- [ ] Firestore Security Rules протестированы
- [ ] Код прошёл review
- [ ] Деплой на Staging окружение
- [ ] Обновлены схемы в docs/plan/
