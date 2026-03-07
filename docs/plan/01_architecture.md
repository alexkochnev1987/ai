# Архитектура системы

## Общая схема

```
┌─────────────────────────────────────┐
│       React Native (Expo)           │
│  ┌──────────┐  ┌──────────────────┐ │
│  │ Zustand  │  │  React Query     │ │
│  │ UI state │  │  Server cache    │ │
│  └──────────┘  └──────────────────┘ │
│  ┌───────────────────────────────┐  │
│  │   WatermelonDB (SQLite)       │  │
│  │   Offline-first local store   │  │
│  └───────────────┬───────────────┘  │
└──────────────────┼──────────────────┘
                   │ sync
    ┌──────────────┼──────────────────┐
    │              │                  │
┌───▼────┐  ┌──────▼──────┐  ┌──────▼──────┐
│Firebase│  │ Cloud Run   │  │  Google /   │
│Auth    │  │ FastAPI     │  │  Apple      │
│Firestore│  │ AI Layer    │  │  Calendar   │
│FCM     │  │ Python      │  │  HealthKit  │
└────────┘  └─────────────┘  └─────────────┘
```

---

## Технический стек

| Слой | Технология | Назначение |
|------|-----------|-----------|
| Mobile | React Native + Expo | UI, навигация |
| Navigation | React Navigation v6 | Стек, табы, модалки |
| UI State | Zustand | Модалки, активная сессия, фильтры |
| Server Cache | React Query | Кэш AI ответов, аналитики |
| Local DB | WatermelonDB (SQLite) | Офлайн хранение всех данных |
| Auth | Firebase Auth | Apple/Google Sign-In |
| Cloud DB | Firestore | Облачная копия данных |
| Push | Firebase FCM | Уведомления |
| AI API | Cloud Run (Python) | FastAPI + агенты |
| AI Models | Claude API (primary) | Планирование, приоритеты |
| Календарь | Google Calendar API | Двусторонняя синхронизация |
| Календарь iOS | EventKit | Apple Calendar |
| Здоровье | HealthKit | Сон, шаги, активность |
| Графики | Victory Native | Аналитика |
| Инфраструктура | Terraform + gcloud | IaC, деплой |
| Секреты | Google Secret Manager | API ключи |

---

## AI слой — Provider Pattern

```
src/ai/
├── core/
│   ├── AIProvider.interface.ts   # контракт провайдера
│   ├── AIMessage.types.ts        # типы сообщений
│   └── AIOrchestrator.ts         # маршрутизация агентов
├── providers/
│   ├── ClaudeProvider.ts         # Anthropic Claude
│   ├── OpenAIProvider.ts         # GPT-4o
│   ├── GeminiProvider.ts         # Google Gemini
│   └── LocalProvider.ts          # заглушка для офлайн
├── agents/
│   ├── PlannerAgent.ts           # план на день
│   ├── PrioritizerAgent.ts       # приоритизация задач
│   ├── InsightAgent.ts           # аналитика + подсказки
│   └── CalendarAgent.ts          # работа с календарём
└── tools/
    ├── TaskTools.ts              # инструменты задач
    ├── CalendarTools.ts          # инструменты календаря
    └── HealthTools.ts            # инструменты здоровья
```

### Интерфейс провайдера

```typescript
interface AIProvider {
  name: string
  complete(messages: AIMessage[]): Promise<AIResponse>
  stream(messages: AIMessage[]): AsyncGenerator<string>
  supportsFunctions: boolean
}
```

### Маршрутизация агентов

| Агент | Провайдер | Причина |
|-------|-----------|---------|
| PlannerAgent | Claude | Лучшее пошаговое рассуждение |
| PrioritizerAgent | Claude | Анализ зависимостей |
| InsightAgent | GPT-4o | Дешевле для отчётов |
| CalendarAgent | Gemini | Нативная интеграция Google |

---

## Бэкенд — Cloud Run (Python)

### Почему Cloud Run, не Firebase Functions

- Нативный Python → LangChain, Pydantic AI, любые AI библиотеки
- WebSocket / SSE для стриминга AI ответов
- Нет холодного старта при `min-instances: 1`
- Запускается как Docker контейнер — полный контроль

### Почему GCP, не AWS

- Firebase уже на Google инфраструктуре
- Нет cross-cloud latency
- Нативная интеграция Firestore ↔ Cloud Run через Service Account
- Google Calendar API в той же экосистеме

---

## Структура проекта

```
productivityai/
├── mobile/                  # React Native приложение
│   ├── src/
│   │   ├── screens/
│   │   ├── components/
│   │   ├── ai/              # AI слой (providers, agents)
│   │   ├── database/        # WatermelonDB схема и модели
│   │   ├── sync/            # SyncService
│   │   ├── services/        # Firebase, Calendar, Health
│   │   └── store/           # Zustand stores
│   └── package.json
├── backend/                 # FastAPI на Cloud Run
│   ├── routers/
│   ├── agents/
│   ├── providers/
│   ├── services/
│   ├── models/
│   ├── Dockerfile
│   └── main.py
├── infrastructure/          # Terraform
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
└── docs/
    └── plan/
```
