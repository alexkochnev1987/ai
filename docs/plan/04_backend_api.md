# Бэкенд API — Cloud Run (Python/FastAPI)

## Структура проекта

```
backend/
├── main.py                    # точка входа FastAPI
├── Dockerfile
├── requirements.txt
├── routers/
│   ├── ai.py                  # /ai/* — все AI эндпоинты
│   ├── calendar.py            # /calendar/* — синхронизация календарей
│   ├── notifications.py       # /notifications/* — планировщик уведомлений
│   └── analytics.py           # /analytics/* — агрегированная аналитика
├── agents/
│   ├── orchestrator.py        # маршрутизатор агентов
│   ├── planner.py             # агент планирования дня
│   ├── prioritizer.py         # агент приоритизации задач
│   └── insights.py            # агент аналитики и инсайтов
├── providers/
│   ├── base.py                # базовый класс провайдера
│   ├── claude.py              # Anthropic Claude
│   ├── openai.py              # OpenAI GPT
│   └── gemini.py              # Google Gemini
├── services/
│   ├── firestore.py           # работа с Firestore
│   ├── google_calendar.py     # Google Calendar API
│   ├── fcm.py                 # Firebase Cloud Messaging
│   ├── secret_manager.py      # Google Secret Manager
│   └── scheduler.py           # Cloud Tasks / Cloud Scheduler
└── models/
    ├── task.py                # Pydantic модели задач
    ├── habit.py               # Pydantic модели привычек
    └── ai_request.py          # Pydantic модели AI запросов
```

---

## AI Orchestrator

```python
# agents/orchestrator.py

class AIOrchestrator:
    def __init__(self):
        self.providers = {
            "claude":  ClaudeProvider(),
            "openai":  OpenAIProvider(),
            "gemini":  GeminiProvider(),
        }
        self.agents = {
            "planner":      PlannerAgent(),
            "prioritizer":  PrioritizerAgent(),
            "insights":     InsightAgent(),
        }
        # Какой агент → какой провайдер
        self.routing = {
            "planner":     "claude",   # лучшее рассуждение
            "prioritizer": "claude",
            "insights":    "openai",   # дешевле для отчётов
        }

    async def run(self, agent: str, context: dict) -> dict:
        provider = self.providers[self.routing[agent]]
        agent_instance = self.agents[agent]

        messages = agent_instance.build_prompt(context)
        response = await provider.complete(messages)
        return agent_instance.parse_response(response)

    async def run_stream(self, agent: str, context: dict):
        provider = self.providers[self.routing[agent]]
        agent_instance = self.agents[agent]
        messages = agent_instance.build_prompt(context)
        async for chunk in provider.stream(messages):
            yield chunk
```

---

## AI эндпоинты

### POST /ai/plan — план дня

```python
@router.post("/ai/plan")
async def generate_daily_plan(
    request: PlanRequest,
    user: User = Depends(get_current_user)
):
    """
    Принимает:
      - tasks: список задач с приоритетами и оценками времени
      - calendar_events: события из Google/Apple Calendar
      - health_data: сон, активность из HealthKit (опционально)
      - date: дата планирования

    Возвращает:
      - schedule: расписание с временными слотами
      - mit: 3 главные задачи дня
      - warnings: предупреждения (дедлайны, перегрузка)
    """
    tasks = await firestore.get_tasks(user.uid, date=request.date)
    calendar_events = await google_calendar.get_events(user.uid, request.date)

    plan = await orchestrator.run(
        agent="planner",
        context={
            "tasks": tasks,
            "events": calendar_events,
            "health": request.health_data,
            "working_hours": user.profile.workingHours,
            "date": request.date,
        }
    )
    return plan

# Pydantic модель запроса
class PlanRequest(BaseModel):
    date: date
    health_data: HealthData | None = None
```

### POST /ai/prioritize — приоритизация задач

```python
@router.post("/ai/prioritize")
async def prioritize_tasks(
    request: PrioritizeRequest,
    user: User = Depends(get_current_user)
):
    """
    Анализирует backlog, предлагает порядок с объяснением.
    Учитывает: дедлайны, зависимости, оценку времени, приоритеты.
    """
    tasks = await firestore.get_pending_tasks(user.uid)

    result = await orchestrator.run(
        agent="prioritizer",
        context={"tasks": tasks}
    )
    return result  # [{taskId, suggestedOrder, reason}]
```

### GET /ai/insights — еженедельные инсайты

```python
@router.get("/ai/insights")
async def get_insights(
    period: str = "week",  # week / month
    user: User = Depends(get_current_user)
):
    """
    Анализирует данные за период.
    Возвращает инсайты, паттерны, рекомендации.
    """
    tasks = await firestore.get_completed_tasks(user.uid, period)
    habits = await firestore.get_habit_logs(user.uid, period)
    moods = await firestore.get_mood_logs(user.uid, period)
    pomodoros = await firestore.get_pomodoro_sessions(user.uid, period)

    insights = await orchestrator.run(
        agent="insights",
        context={
            "tasks": tasks,
            "habits": habits,
            "moods": moods,
            "pomodoros": pomodoros,
            "period": period,
        }
    )
    return insights
```

### POST /ai/stream — стриминг (для AI чата v2)

```python
@router.post("/ai/stream")
async def stream_response(
    request: ChatRequest,
    user: User = Depends(get_current_user)
):
    """Server-Sent Events стриминг AI ответа"""
    async def generator():
        async for chunk in orchestrator.run_stream("chat", {"messages": request.messages}):
            yield f"data: {json.dumps({'chunk': chunk})}\n\n"
        yield "data: [DONE]\n\n"

    return StreamingResponse(generator(), media_type="text/event-stream")
```

---

## Уведомления — Cloud Tasks

```python
# routers/notifications.py

@router.post("/notifications/schedule-streak-check")
async def schedule_streak_check(user: User = Depends(get_current_user)):
    """
    Ежедневно в 21:00 проверяет незавершённые привычки
    и отправляет push тем, у кого streak под угрозой
    """
    habits = await firestore.get_active_habits(user.uid)
    today_logs = await firestore.get_today_habit_logs(user.uid)

    for habit in habits:
        if not any(log.habitId == habit.id for log in today_logs):
            await fcm.send_push(
                token=user.fcmToken,
                title=f"Streak под угрозой 🔥",
                body=f"Не забудь: {habit.title}. Осталось 3 часа!",
                data={"habitId": habit.id, "type": "streak_reminder"}
            )

@router.post("/notifications/daily-brief")
async def send_daily_brief(user: User = Depends(get_current_user)):
    """
    8:00 утра — отправить AI план дня как push уведомление
    """
    plan = await generate_daily_plan(PlanRequest(date=date.today()), user)
    summary = plan.get("summary", "Твой план на сегодня готов")

    await fcm.send_push(
        token=user.fcmToken,
        title="Доброе утро! ☀️",
        body=summary,
        data={"type": "daily_brief", "planId": plan["id"]}
    )
```

---

## Календарь — Google Calendar API

```python
# routers/calendar.py

@router.post("/calendar/google/sync")
async def sync_google_calendar(
    request: CalendarSyncRequest,
    user: User = Depends(get_current_user)
):
    """
    Двусторонняя синхронизация:
    1. Импорт событий Google → Firestore
    2. Экспорт задач с датой → Google Calendar
    """
    credentials = await get_google_credentials(user.uid)

    # Тянем события за период
    events = await google_calendar.get_events(
        credentials,
        start=request.start,
        end=request.end
    )

    # Сохраняем в Firestore
    for event in events:
        await firestore.upsert_calendar_event(user.uid, event)

    # Пушим задачи с дедлайном в Google Calendar
    tasks_with_deadline = await firestore.get_tasks_with_deadline(user.uid)
    for task in tasks_with_deadline:
        if not task.externalId:
            event_id = await google_calendar.create_event(credentials, task)
            await firestore.update_task(user.uid, task.id, {"externalId": event_id})

    return {"synced": len(events), "exported": len(tasks_with_deadline)}
```

---

## Аутентификация запросов

```python
# Все эндпоинты защищены Firebase ID Token
async def get_current_user(authorization: str = Header(...)) -> User:
    token = authorization.replace("Bearer ", "")
    decoded = auth.verify_id_token(token)
    user_data = await firestore.get_user(decoded["uid"])
    return User(**user_data)
```

---

## Dockerfile

```dockerfile
FROM python:3.11-slim

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
```
