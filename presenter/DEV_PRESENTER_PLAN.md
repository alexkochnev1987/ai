# План реализации в ветке `dev-presenter`

> Основа: продуктовый план из `presenter/IMPLEMENTATION_PLAN.md`.
> Цель этой версии — превратить общий roadmap в инженерный пошаговый execution-план для разработки.

## 1) Цель ветки

Собрать рабочий MVP AI Presenter:
- веб-приложение (создание презентации, статус, запуск/остановка),
- Chrome Extension (выполнение команд на живом сайте),
- серверная оркестрация (анализ, план шагов, очередь команд, TTS).

## 2) Definition of Done (MVP)

MVP считается готовым, если выполнены все условия:
1. Пользователь может зарегистрироваться/войти и создать презентацию по URL.
2. Система строит analysis + план шагов и показывает их в UI.
3. План можно отредактировать перед запуском.
4. При запуске extension выполняет шаги на реальной вкладке: `navigate/click/scroll/fill_input`.
5. Озвучка и субтитры синхронизированы по шагам.
6. Доступны `play/pause/stop/replay`.
7. Все критические ошибки фиксируются в логах с привязкой к `sessionId`.

## 3) Этапы реализации

## Этап A — Инфраструктурный каркас (2–3 дня)

**Задачи:**
- Подготовить монорепо-структуру:
  - `presenter/apps/web` (Next.js App Router)
  - `presenter/packages/extension` (MV3)
- Подключить PostgreSQL через `docker-compose`.
- Настроить ORM (Prisma/Drizzle) + первую миграцию.
- Поднять базовый Auth (signup/login/session).

**Артефакты:**
- Рабочий запуск `web + db` локально.
- Таблицы пользователей и базовой сущности презентации.

**Критерий приёмки:**
- Новый пользователь создаётся, логинится и видит пустой dashboard.

## Этап B — Протокол Web ↔ Extension (2–4 дня)

**Задачи:**
- Спроектировать `PlaybackSession` и `CommandQueue`.
- Реализовать API:
  - `POST /api/sessions` — создать сессию,
  - `GET /api/sessions/:id/next-command` — polling,
  - `POST /api/sessions/:id/command-result` — результат команды.
- В extension реализовать polling loop и обработку `noop/ping`.

**Артефакты:**
- Сквозной канал «сервер отправил команду → extension вернул результат».

**Критерий приёмки:**
- Для тестовой сессии видно последовательность команд и статусы выполнения в БД.

## Этап C — Сбор контекста сайта (3–4 дня)

**Задачи:**
- `extract_dom` команда:
  - URL, title, интерактивные элементы, формы, заголовки.
- Базовая нормализация селекторов: `data-testid → id → aria-label → text`.
- Команда `screenshot` (capture visible tab) и сохранение метаданных.
- Сохранение snapshot в JSONB.

**Артефакты:**
- История снимков состояния страницы для сессии.

**Критерий приёмки:**
- Для заданного URL можно получить как минимум 1 DOM snapshot + 1 screenshot metadata.

## Этап D — AI анализ и генерация плана (3–5 дней)

**Задачи:**
- Ввести abstraction для AI-провайдеров (OpenAI/Claude).
- Реализовать pipeline:
  - собрать snapshots,
  - сгенерировать `analysis`,
  - сгенерировать `presentation_steps`.
- Сделать UI статуса (SSE/polling) и экран редактирования шагов.

**Артефакты:**
- Сохранённые analysis/plan в БД.
- Редактор шагов в веб-интерфейсе.

**Критерий приёмки:**
- Пользователь видит автогенерированный план из N шагов и может сохранить правки.

## Этап E — TTS и субтитры (2–3 дня)

**Задачи:**
- Генерация аудио на шаг (сервером).
- Хранение аудио URL (локально в dev / object storage).
- Команда `play_audio` + `show_subtitle` для extension.
- Минимальная синхронизация длительностей шага.

**Артефакты:**
- Для каждого шага есть текст, аудио, длительность.

**Критерий приёмки:**
- При проигрывании шага слышно озвучку и отображается соответствующий subtitle overlay.

## Этап F — Оркестратор воспроизведения (3–4 дня)

**Задачи:**
- Серверный state machine: `idle/running/paused/stopped/completed/failed`.
- Команды исполнения шага: `navigate`, `click`, `scroll`, `fill_input`, `wait`.
- UI-плеер: `Play/Pause/Stop/Replay`.
- Обработка таймаутов и ретраев для нестабильных действий.

**Артефакты:**
- End-to-end воспроизведение плана на live-сайте.

**Критерий приёмки:**
- Полный прогон 5+ шагов с возможностью `pause` и `replay`.

## Этап G — Надёжность и hardening (2–3 дня)

**Задачи:**
- Логирование по `sessionId/presentationId/stepId`.
- Стандартизированный error model + user-friendly сообщения.
- Fallback-стратегии селекторов и skip-step при fatal error.
- Базовые e2e smoke-тесты и чек-лист QA.

**Критерий приёмки:**
- В случае падения шага система не «зависает» бесконечно и завершает сессию предсказуемым статусом.

## 4) Технический бэклог (приоритет P0/P1)

### P0 (обязательно для MVP)
- Схема БД: presentations, analyses, plans, steps, sessions, commands, command_results.
- Командный протокол polling.
- Генерация и редактирование плана.
- Исполнение шагов + TTS + subtitles.
- Управление playback state.

### P1 (после MVP)
- Улучшение устойчивости селекторов через эвристики и scoring.
- Поддержка нескольких голосов/языков и auto-duration fit.
- История версий плана.
- Метрики качества исполнения (step success rate).

## 5) Риски и меры

- **Ломкие селекторы на динамических UI** → fallback-цепочка + повторный `extract_dom`.
- **Service Worker MV3 выгружается** → хранение session state в `chrome.storage.session` + reconnect.
- **Ограничения iframes/canvas/CAPTCHA/2FA** → явные ограничения MVP в UX.
- **Долгие операции AI/TTS** → job-модель + SSE/polling прогресс.

## 6) План работ на ближайший спринт (Sprint 1)

1. Каркас `apps/web` + `packages/extension`.
2. Миграции БД и auth.
3. API сессии и command polling.
4. Extension loop с `noop/ping`.
5. Smoke e2e: «создать сессию → получить команду → вернуть результат».

**Ожидаемый результат Sprint 1:**
- Готов технический «скелет» платформы и подтверждён рабочий транспорт между сервером и extension.


## 7) Детальный runtime flow: после вставки ссылки до озвучки и исполнения

Ниже — конкретный сценарий, который отвечает на вопрос «что происходит после ввода URL».

### 7.1 Шаги пользователя в UI

1. Пользователь вводит:
   - `targetUrl`
   - опционально: описание цели, язык, длительность, креды для входа.
2. Нажимает «Сгенерировать план».
3. UI делает `POST /api/presentations` и получает `presentationId` + `jobId`.
4. UI подписывается на прогресс через `GET /api/jobs/:jobId/events` (SSE) или polling.

### 7.2 Что делает сервер после `POST /api/presentations`

Сервер создаёт запись `presentation` со статусом `queued` и ставит задачу в очередь `analysis_job`.

**Пайплайн job-а:**
1. `collect_context`:
   - создаёт `session` для extension,
   - отправляет команды: `navigate(targetUrl)`, `extract_dom`, `screenshot`.
2. `build_analysis`:
   - из snapshot-ов формирует сжатый контекст страницы,
   - вызывает AI провайдера и сохраняет `presentation_analysis`.
3. `build_plan`:
   - второй AI вызов (или тот же с отдельным prompt) генерирует `presentation_steps[]`.
4. `persist_plan`:
   - нормализует шаги, duration, actions,
   - сохраняет `presentation_plan` + `presentation_steps`.
5. Меняет статус на `plan_ready`.

### 7.3 Как именно генерируется план (AI)

Рекомендуемый двухпроходный подход:

1. **Анализ (pass 1):**
   - вход: DOM snapshot + title + URL + цель пользователя,
   - выход: `analysis` (главные блоки интерфейса, потенциальный user journey, риски).
2. **Генерация шагов (pass 2):**
   - вход: `analysis` + ограничение по времени + язык,
   - выход: массив шагов формата:
   - `title`, `narration_text`, `actions[]`, `estimated_duration_seconds`, `page_url`.

**Валидация перед сохранением:**
- шагов не 0,
- есть `narration_text`,
- у action корректный тип и payload,
- суммарная длительность в допустимом диапазоне.

### 7.4 Что происходит при нажатии Play

1. UI: `POST /api/playback/start`.
2. Сервер создаёт `playback_session` со статусом `running`.
3. Оркестратор проходит шаги по порядку и кладёт в queue команды:
   - `show_subtitle`
   - `play_audio`
   - `navigate/click/scroll/fill_input/wait`
4. Extension по polling забирает следующую команду и исполняет её на целевой вкладке.
5. По завершении каждой команды extension отправляет `command_result`.
6. Сервер обновляет прогресс шага и двигается к следующей команде.

### 7.5 Как озвучивается текст (TTS)

Озвучка делается на **сервере**, а не в extension:

1. После генерации/редактирования плана сервер берёт `narration_text` каждого шага.
2. Вызывает TTS provider с параметрами `voice`, `language`, `speaking_rate`.
3. Сохраняет аудиофайл и метаданные (`audio_url`, `duration_ms`) в `audio_assets`.
4. При старте playback команда `play_audio` передаёт extension ссылку на готовый файл.
5. Extension проигрывает `<audio src="audio_url">` и параллельно показывает subtitle overlay.

### 7.6 Синхронизация озвучки и действий

Базовая стратегия MVP:

- `show_subtitle` и `play_audio` отправляются перед интерактивными действиями шага.
- Если действие короткое, сервер ждёт завершения аудио (`duration_ms`) и затем переходит дальше.
- Если действие длинное, используется `wait_after_actions_ms` + контрольный таймаут.

Улучшение после MVP:
- timeline per step (миллисекундная шкала),
- фоновые и foreground actions,
- ускорение/замедление голоса под целевую длительность шага.

### 7.7 Управление состоянием (play/pause/stop/replay)

- `pause`: сервер перестаёт выдавать новые команды; extension завершает текущую и ждёт.
- `stop`: сервер переводит сессию в `stopped`, очередь закрывается.
- `replay`: создаётся новая `playback_session`, шаги проигрываются с начала (или с выбранного шага).

### 7.8 Минимальные API контракты для этого flow

- `POST /api/presentations` → создать презентацию и analysis job.
- `GET /api/jobs/:jobId/events` → прогресс генерации.
- `GET /api/presentations/:id/plan` + `PATCH /api/presentations/:id/plan` → просмотр/редактирование.
- `POST /api/playback/start|pause|stop|replay` → управление проигрыванием.
- `GET /api/sessions/:id/next-command` + `POST /api/sessions/:id/command-result` → канал с extension.

### 7.9 Наблюдаемость (что логируем обязательно)

- `presentationId`, `jobId`, `sessionId`, `stepId`, `commandId`.
- AI latency, TTS latency, command execution latency.
- Причины ошибок: selector_not_found, navigation_timeout, audio_playback_error, provider_error.
