# Схема данных

## Firestore — структура коллекций

```
users/{userId}
├── profile                    # профиль пользователя
├── settings                   # настройки приложения
├── tasks/{taskId}             # задачи
├── habits/{habitId}           # привычки
├── habit_logs/{logId}         # отметки выполнения привычек
├── pomodoro_sessions/{id}     # pomodoro сессии
├── mood_logs/{logId}          # записи настроения
└── sync_meta                  # метаданные синхронизации

projects/{projectId}           # совместные проекты (общая коллекция)
├── members/{memberId}
└── tasks/{taskId}
```

---

## Схемы документов Firestore

### UserProfile

```typescript
interface UserProfile {
  uid: string
  email: string
  displayName: string
  avatarUrl: string | null
  timezone: string              // "Europe/Moscow"
  workingHours: {
    start: string               // "09:00"
    end: string                 // "18:00"
    workDays: number[]          // [1,2,3,4,5]
  }
  fcmToken: string | null       // push уведомления
  googleCalendarConnected: boolean
  appleCalendarConnected: boolean
  healthKitConnected: boolean
  createdAt: Timestamp
  updatedAt: Timestamp
}
```

### Task

```typescript
interface Task {
  id: string
  title: string
  description: string | null
  status: 'pending' | 'in_progress' | 'done' | 'archived'
  priority: 'urgent' | 'high' | 'medium' | 'low'
  tags: string[]
  projectId: string | null
  dueDate: Timestamp | null
  estimatedMinutes: number | null
  actualMinutes: number          // суммируется из pomodoro сессий
  subtasks: Array<{
    id: string
    title: string
    completed: boolean
    dueDate: Timestamp | null
  }>
  recurrence: {
    type: 'daily' | 'weekly' | 'monthly' | 'custom'
    daysOfWeek: number[]
    interval: number
    endDate: Timestamp | null
  } | null
  isMIT: boolean                 // главная задача дня
  pomodoroCount: number          // количество помидоров
  syncSource: 'app' | 'google' | 'apple' | null
  externalId: string | null
  createdAt: Timestamp
  updatedAt: Timestamp
  deletedAt: Timestamp | null    // soft delete
}
```

### Habit

```typescript
interface Habit {
  id: string
  title: string
  icon: string                   // emoji или icon name
  color: string                  // hex цвет
  frequency: {
    type: 'daily' | 'weekly' | 'custom'
    daysOfWeek: number[]
    timesPerWeek: number | null
  }
  reminderTime: string | null    // "21:00"
  streakGoal: number             // 30 / 90 / 365
  currentStreak: number
  longestStreak: number
  graceDaysUsed: number          // сброс каждый месяц
  graceDayLimit: number          // default: 1
  category: 'mental' | 'physical' | 'learning' | 'work' | 'custom'
  healthKitId: string | null
  archivedAt: Timestamp | null
  createdAt: Timestamp
  updatedAt: Timestamp
}
```

### HabitLog

```typescript
interface HabitLog {
  id: string
  habitId: string
  date: string                   // "2026-03-07"
  completed: boolean
  isGraceDay: boolean
  note: string | null
  createdAt: Timestamp
}
```

### PomodoroSession

```typescript
interface PomodoroSession {
  id: string
  taskId: string | null
  type: 'work' | 'short_break' | 'long_break'
  durationMinutes: number
  completedAt: Timestamp
  interrupted: boolean
}
```

### MoodLog

```typescript
interface MoodLog {
  id: string
  score: 1 | 2 | 3 | 4 | 5     // 1=плохо, 5=отлично
  note: string | null
  date: string                   // "2026-03-07"
  createdAt: Timestamp
}
```

---

## WatermelonDB — SQLite таблицы (локально)

### tasks

| Поле | Тип | Описание |
|------|-----|----------|
| server_id | string? | ID в Firestore |
| title | string | Название |
| description | string? | Описание |
| status | string | pending/in_progress/done/archived |
| priority | string | urgent/high/medium/low |
| tags_json | string | JSON массив тегов |
| due_date | number? | Unix timestamp |
| estimated_minutes | number? | Оценка времени |
| actual_minutes | number | Реальное время |
| subtasks_json | string | JSON массив подзадач |
| recurrence_json | string? | JSON правила повтора |
| is_mit | boolean | Главная задача дня |
| project_id | string? | ID проекта |
| is_synced | boolean | Синхронизировано с сервером |
| deleted_at | number? | Soft delete |
| created_at | number | Unix timestamp |
| updated_at | number | Unix timestamp |

### habits

| Поле | Тип | Описание |
|------|-----|----------|
| server_id | string? | ID в Firestore |
| title | string | Название |
| icon | string | Иконка |
| color | string | Цвет hex |
| frequency_json | string | JSON настройки частоты |
| reminder_time | string? | "21:00" |
| current_streak | number | Текущий streak |
| longest_streak | number | Рекордный streak |
| streak_goal | number | Цель (30/90/365) |
| grace_days_used | number | Использованные grace days |
| category | string | Категория |
| is_synced | boolean | Синхронизировано |
| created_at | number | Unix timestamp |
| updated_at | number | Unix timestamp |

### habit_logs

| Поле | Тип | Описание |
|------|-----|----------|
| server_id | string? | ID в Firestore |
| habit_id | string | Локальный ID привычки |
| date | string | "2026-03-07" |
| completed | boolean | Выполнено |
| is_grace_day | boolean | Grace day |
| note | string? | Заметка |
| is_synced | boolean | Синхронизировано |
| created_at | number | Unix timestamp |

### sync_queue

| Поле | Тип | Описание |
|------|-----|----------|
| operation | string | create/update/delete |
| collection | string | Название коллекции |
| document_id | string | ID документа |
| payload_json | string | JSON данные |
| created_at | number | Unix timestamp |
| retry_count | number | Количество попыток |

---

## Firestore Security Rules

```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {

    // Пользователь видит только свои данные
    match /users/{userId}/{document=**} {
      allow read, write: if request.auth != null
                         && request.auth.uid == userId;
    }

    // Проекты — только участники
    match /projects/{projectId} {
      allow read: if request.auth != null
                  && request.auth.uid in resource.data.memberIds;
      allow write: if request.auth != null
                   && request.auth.uid == resource.data.ownerId;

      match /tasks/{taskId} {
        allow read, write: if request.auth != null
                           && request.auth.uid in
                              get(/databases/$(database)/documents/projects/$(projectId)).data.memberIds;
      }
    }
  }
}
```
