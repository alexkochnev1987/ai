# AI Web Presenter — Implementation Plan

## Product

A universal AI tool: the user provides a URL (and options), and the system analyzes the application, builds a scenario, and runs a **live** presentation on the real website — auto-navigation, voice-over, subtitles; playback can be paused and replayed.

### Flow

```
Login → URL (+description, credentials, language, duration) → analysis → plan → review/edit → Play → actions on site + TTS + subtitles
```

## Architecture (without separate FastAPI)

| Layer | Responsibility |
|------|------|
| **Next.js (App Router)** | UI, auth, REST API, AI/TTS calls, DB; progress in UI via **SSE or polling** |
| **PostgreSQL** | Users, presentations, analyses, plans, steps, sessions, command queue |
| **Chrome Extension (MV3)** | Execution on target tab: DOM, clicks, scroll, navigation, subtitles, audio |

**Why extension:** bypass Same-Origin Policy constraints, control any open website without using “screen video streaming” as the primary mode.

### Communication without mandatory WebSocket

- **Browser (web app) ↔ Next.js:** HTTP + for long tasks, **SSE** (`/api/.../events`) or polling by `jobId` / `sessionId`.
- **Extension ↔ Next.js:** `GET` next command (short **polling** or long polling) + `POST` result (`command_completed`, `dom_snapshot`, errors).

AI providers (Claude/OpenAI) are called **server-side only** (HTTPS request/response; streaming provider responses is still HTTP). Extension **must not store** API keys.

### Diagram

```
┌─────────────────┐     HTTPS (poll + POST)      ┌──────────────────────┐
│ Chrome Extension │ ◄──────────────────────────► │ Next.js              │
│ (content + SW)   │                               │ API + Auth + DB      │
└────────┬────────┘                               │ AI + TTS             │
         │ executes on tab                        └──────────┬───────────┘
         ▼                                                    │
   Live target website                              ┌──────────▼───────────┐
                                                    │ PostgreSQL           │
                                                    └──────────────────────┘
┌─────────────────┐     HTTPS + SSE/poll         ┌──────────────────────┐
│ Web App (Next)  │ ◄──────────────────────────► │ (same Next.js app)   │
└─────────────────┘                               └──────────────────────┘
```

## Extension modules (MV3)

| Module | Responsibility |
|--------|--------|
| Service Worker | Polling API, message routing to content script, `chrome.alarms` + keep-alive (avoid SW termination) |
| Content script | DOM extraction, action execution, Shadow DOM subtitles overlay, `<audio>` playback by URL |
| Popup | Status, link to web app, hint: “open a tab with the target website” |

**Selectors (priority):** `data-testid` → `id` → `aria-label` → text heuristics → CSS path.

**Service Worker:** keep playback state in `chrome.storage.session`; when SW restarts, reconnect using saved `sessionId`.

## Data and API (reference)

### Entities

- `User`, `Presentation`, `PresentationAnalysis`, `PresentationPlan`, `PresentationStep`, `AudioAsset`, `PlaybackSession`
- Extension command queue: type (`navigate`, `click`, `scroll`, `fill_input`, `extract_dom`, `screenshot`, `show_subtitle`, `play_audio`), payload, status

### Plan step (type)

```typescript
interface PresentationStep {
  step_number: number
  page_url: string
  title: string
  narration_text: string
  estimated_duration_seconds: number
  actions: StepAction[] // navigate, click, scroll, fill_input, highlight
  wait_after_actions_ms: number
}
```

### Credentials for target-site login

Store on server (encrypted, e.g. AES-256-GCM). Pass to extension **one-time only** in `fill_input` command payload; extension **does not persist** them.

## Suggested repository structure

```
presenter/
├── apps/
│   └── web/                    # Next.js (App Router)
│       ├── app/
│       ├── lib/                # db, ai, tts, encryption, orchestration
│       └── ...
├── packages/
│   └── extension/              # MV3
│       ├── src/background/
│       ├── src/content/
│       ├── src/popup/
│       └── manifest.json
├── docker-compose.yml          # PostgreSQL (local)
├── IMPLEMENTATION_PLAN.md
└── ...
```

## Implementation phases

### Phase 0 — Foundation

- [ ] Repository scaffold: `apps/web` + `packages/extension` (or single-directory alternative if preferred)
- [ ] Next.js: routing, DB (Prisma/Drizzle), migrations, signup/login
- [ ] `docker-compose.yml` — Postgres
- [ ] Extension: MV3 manifest, SW + content + popup stubs

### Phase 1 — Extension protocol without WS

- [ ] Playback session and command models
- [ ] API: create session, poll “next command”, POST result
- [ ] Extension: polling loop + no-op command execution (end-to-end pipe test)

### Phase 2 — Live-site context collection

- [ ] DOM snapshot (navigation, interactive elements, headings, forms, truncated text)
- [ ] Screenshot: `chrome.tabs.captureVisibleTab` by server command
- [ ] Save snapshots to DB (JSONB)

### Phase 3 — AI: analysis and plan

- [ ] Provider abstraction (Claude/OpenAI), user key on server
- [ ] Use case: crawl/collect pages via extension commands → AI → `analysis` + draft plan
- [ ] UI: create presentation; details page with status (SSE/poll)
- [ ] Plan editing before launch

### Phase 4 — TTS

- [ ] Generate per-step audio (e.g., Google Cloud TTS), store files (S3-compatible / local in dev)
- [ ] Duration scaling to match selected total time; SSML `prosody` when needed

### Phase 5 — Playback

- [ ] Extension execution of `navigate`, `click`, `scroll`, `fill_input`, subtitles, `play_audio` by URL
- [ ] Server-side step orchestrator and command queue
- [ ] Player UI: play / pause / stop / replay

### Phase 6 — Reliability and settings

- [ ] Fallback selectors, timeouts, skip-step on fatal errors
- [ ] Settings: AI provider, API key
- [ ] Extension install docs (load unpacked)

## MVP verification

1. Start Postgres (`docker compose up`) and Next.js (`pnpm dev` / `npm run dev`)
2. Install extension (Chrome → Extensions → Load unpacked)
3. Sign up, create presentation for a public URL
4. Validate: analysis → plan → edits → playback with voice and subtitles on a live tab
5. Validate pause/stop/replay
6. (Optional) simple login scenario with injected credentials

## Risks

- Selector fragility after layout updates; SPA and delayed content require explicit waits and repeated `extract_dom`
- iframe / canvas / closed Shadow DOM have limited automation capabilities
- CAPTCHA, 2FA are out of MVP automation scope
- Website legal constraints (ToS) — provide user warning

## Hosting note

On **serverless** (e.g., classic Vercel), long-running background tasks and frequent polling per session should use a **job queue** and/or separate long-running worker as load grows. For MVP, a single Node process (Docker/VPS) or careful split “request → DB job → worker” is sufficient.
