# Фронтенд — Модули и экраны

## Навигационная структура

```
App
├── AuthStack
│   ├── WelcomeScreen
│   ├── SignInScreen
│   └── OnboardingScreen (5 шагов)
│
└── MainTabs
    ├── DashboardTab
    │   └── DashboardScreen
    ├── TasksTab
    │   ├── TaskListScreen
    │   ├── TaskDetailScreen
    │   └── TaskFormScreen (modal)
    ├── HabitsTab
    │   ├── HabitListScreen
    │   ├── HabitDetailScreen (история + тепловая карта)
    │   └── HabitFormScreen (modal)
    ├── CalendarTab
    │   └── CalendarScreen (day/week/month)
    └── ProfileTab
        ├── ProfileScreen
        ├── SettingsScreen
        ├── AnalyticsScreen
        └── IntegrationsScreen
```

---

## M1 — Аутентификация

### Экраны
- **WelcomeScreen** — логотип, кнопки "Войти через Apple" / "Войти через Google"
- **OnboardingScreen** — 5 шагов после первого входа:
  1. Выбор целей (работа / здоровье / учёба / личное)
  2. Рабочие часы (с X до Y, рабочие дни)
  3. Добавить первые 3 привычки
  4. Подключить Google Calendar (опционально)
  5. Включить уведомления

### Логика
```typescript
// services/AuthService.ts
class AuthService {
  async signInWithApple(): Promise<User>
  async signInWithGoogle(): Promise<User>
  async signOut(): Promise<void>
  async deleteAccount(): Promise<void>  // GDPR
}
```

---

## M2 — Dashboard

### Компоненты
```
DashboardScreen
├── HeaderGreeting          — "Доброе утро, Алекс ☀️" + статус синхронизации
├── DailyProgressCard       — задачи выполнено/всего, прогресс-бар
├── MITBlock                — 3 главные задачи дня (выделены)
├── ActivePomodoroCard      — если активна сессия (таймер в реальном времени)
├── HabitStreaksRow         — горизонтальный скролл активных streaks
├── AIBriefCard             — карточка AI плана / подсказки
├── MoodCheckCard           — если настроение не отмечено сегодня
└── FABButton               — плавающая кнопка быстрых действий
```

### FAB меню
- Добавить задачу
- Запустить Pomodoro
- Отметить настроение
- Добавить привычку

---

## M3 — Задачи

### Экраны

**TaskListScreen**
```
├── SearchBar
├── FilterBar               — теги / приоритет / дедлайн / проект
├── SectionedTaskList
│   ├── Просрочено          — красный заголовок
│   ├── Сегодня
│   ├── Завтра
│   └── Позже
└── TaskCard
    ├── Свайп вправо → выполнено (зелёный)
    └── Свайп влево → меню (отложить / удалить)
```

**TaskCard компонент**
```
[приоритет] Название задачи              [дедлайн]
            тег1 тег2                    3/5 подзадач ██░░░
            ⏱ 1ч 30мин оценка           🍅🍅 помидоры
```

**TaskDetailScreen**
- Все поля задачи (редактируемые инлайн)
- Список подзадач с чекбоксами
- Кнопка "Начать работу" → запускает таймер + Pomodoro
- История времени (сессии по дням)
- Комментарии (для совместных проектов)

**TaskFormScreen (modal)**
- Поля: название, описание, дедлайн, приоритет, теги, проект, оценка времени
- Подзадачи (добавить / удалить)
- Повтор (ежедневно / по дням / ежемесячно)

### Функциональность
```typescript
// Drag & drop сортировка
import DraggableFlatList from 'react-native-draggable-flatlist'

// Свайп действия
import { Swipeable } from 'react-native-gesture-handler'

// Трекер времени
class TaskTimerService {
  start(taskId: string): void
  pause(): void
  stop(): Promise<number>    // возвращает потраченные минуты
}
```

---

## M4 — Привычки и Streaks

### HabitListScreen
```
├── DateStrip               — горизонтальный скролл дат (сегодня в центре)
├── HabitCard (для каждой привычки)
│   ├── иконка + цвет + название
│   ├── 🔥 streak число
│   ├── прогресс к цели (30/90 дней)
│   └── чекбокс на сегодня (большой, легко нажать)
└── AddHabitButton
```

### HabitDetailScreen
```
├── Статистика: текущий streak / рекорд / % за месяц
├── Тепловая карта (365 дней, как GitHub)
├── График: streak история по неделям
├── Лучшие дни недели
└── Кнопка архивировать
```

### Streak алгоритм
```typescript
function calculateStreak(logs: HabitLog[], habit: Habit): number {
  let streak = 0
  const today = format(new Date(), 'yyyy-MM-dd')
  let checkDate = today

  while (true) {
    const log = logs.find(l => l.date === checkDate && l.completed)
    const isGraceDay = logs.find(l => l.date === checkDate && l.isGraceDay)

    if (log || isGraceDay) {
      streak++
    } else if (checkDate === today) {
      // Сегодня ещё не отмечено — не ломаем streak
    } else {
      break
    }
    checkDate = format(subDays(parseISO(checkDate), 1), 'yyyy-MM-dd')
  }
  return streak
}
```

---

## M5 — Календарь

### CalendarScreen
```
├── ViewToggle              — день / неделя / месяц
├── WeekView (default)
│   ├── TimeGrid            — временная сетка с событиями
│   ├── TaskBlocks          — блоки задач (из MIT и с дедлайном)
│   └── EventBlocks         — события из Google/Apple Calendar
├── MonthView
│   ├── МесячнаяSетка
│   └── Точки под датами (количество задач/событий)
└── DayView
    └── Подробное расписание дня
```

Цветовое кодирование:
- Синий — Google Calendar
- Серый — Apple Calendar
- Фиолетовый — задачи приложения
- Красный — просрочено

---

## M6 — Pomodoro

### PomodoroScreen (modal или отдельный таб)
```
├── TaskSelector            — выбрать задачу (опционально)
├── CircularTimer           — большой круговой таймер
│   └── анимированный прогресс + время
├── SessionIndicators       — 🍅🍅🍅○ (4 сессии)
├── Controls                — старт / пауза / стоп / пропустить
├── SessionType             — Работа / Короткий перерыв / Длинный перерыв
└── TodayStats              — сессий сегодня: 6, минут фокуса: 150
```

Live Activity (iOS 16+):
- Таймер на экране блокировки
- В Dynamic Island

---

## M7 — AI планировщик (интеграция в Dashboard)

### AIBriefCard
```
┌─────────────────────────────────┐
│ 🤖 AI план на сегодня           │
│                                 │
│ 9:00  Презентация (1.5ч)        │
│ 11:00 Встреча с командой        │
│ 12:30 Email рассылка (1ч)       │
│ ...                             │
│                                 │
│ ⚠️ Дедлайн по "Отчёту" завтра   │
│                                 │
│ [Принять план]  [Изменить]      │
└─────────────────────────────────┘
```

### AIInsightsCard (проактивные подсказки)
```
💡 "Ты наиболее продуктивен по вторникам с 10 до 12"
⚠️ "3 задачи с дедлайном завтра не запланированы"
🔥 "Streak по английскому под угрозой — 3 часа"
📊 "Продуктивность на 20% ниже прошлой недели"
```

---

## M8 — Аналитика

### AnalyticsScreen
```
├── PeriodSelector          — 7 дней / 30 дней / 90 дней
├── SummaryCards
│   ├── Выполнено задач: 47 (+12% vs прошлый период)
│   ├── Pomodoro сессий: 38
│   ├── Лучший streak: 21 день
│   └── Среднее настроение: 3.8 / 5
├── ProductivityHeatmap     — матрица день×час
├── CategoryPieChart        — время по категориям
├── HabitsCompletionChart   — % выполнения привычек по дням
├── MoodTrendChart          — линейный график настроения
└── WeeklyReportCard        — последний еженедельный отчёт от AI
```

---

## M9 — Настроение

### MoodCheckModal (появляется по уведомлению)
```
Как ты себя чувствуешь сегодня?

😔  😕  😐  🙂  😊
 1   2   3   4   5

[Добавить заметку...]

[Сохранить]
```

---

## M10 — Apple Health

Нет отдельного экрана. Данные используются:
- В AI плане (сон → нагрузка дня)
- В аналитике (корреляция настроения и сна)
- В уведомлениях о перерывах (стоячие минуты)

Подключение в **IntegrationsScreen** → "Подключить Apple Health" → HealthKit permission request

---

## M11 — Совместные проекты

### ProjectScreen
```
├── ProjectHeader           — название, прогресс, участники
├── KanbanBoard
│   ├── Backlog колонка
│   ├── В работе колонка
│   ├── На проверке колонка
│   └── Готово колонка
├── TaskCard (drag & drop между колонками)
│   └── аватар ответственного
└── ActivityFeed            — кто что сделал
```

---

## Zustand Stores

```typescript
// Отдельный store на каждый домен
useTaskStore       — список задач, фильтры, активная задача
useHabitStore      — привычки, streaks, логи
usePomodoroStore   — активная сессия, настройки
useSyncStore       — статус синхронизации, pendingCount
useAIStore         — план дня, инсайты, isLoading
useUIStore         — модалки, тема, активный таб
```
