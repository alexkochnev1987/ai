# Инфраструктура — GCP + Firebase

## Управление через API и CLI

### Инструменты

| Инструмент | Назначение |
|-----------|-----------|
| `gcloud CLI` | Cloud Run, Cloud Scheduler, Cloud Tasks, IAM |
| `firebase CLI` | Firestore rules, Firebase деплой |
| `terraform` | IaC — вся инфраструктура как код |
| `Secret Manager API` | Хранение API ключей |
| `Firebase Admin SDK` | Работа с Firebase из бэкенда (Python) |

---

## Terraform — вся инфраструктура

```hcl
# infrastructure/main.tf

terraform {
  required_providers {
    google = { source = "hashicorp/google", version = "~> 5.0" }
  }
  backend "gcs" {
    bucket = "your-project-terraform-state"
    prefix = "productivityai"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# ─── Cloud Run ───────────────────────────────────────
resource "google_cloud_run_v2_service" "api" {
  name     = "ai-backend"
  location = var.region

  template {
    scaling { min_instance_count = 1 }  # нет холодного старта

    containers {
      image = "gcr.io/${var.project_id}/ai-backend:latest"

      resources {
        limits = { cpu = "1", memory = "512Mi" }
      }

      env {
        name = "ANTHROPIC_API_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.anthropic_key.secret_id
            version = "latest"
          }
        }
      }
      env {
        name  = "FIREBASE_PROJECT_ID"
        value = var.project_id
      }
    }
  }
}

# Публичный доступ к Cloud Run
resource "google_cloud_run_service_iam_member" "public" {
  service  = google_cloud_run_v2_service.api.name
  location = var.region
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# ─── Secret Manager ──────────────────────────────────
resource "google_secret_manager_secret" "anthropic_key" {
  secret_id = "anthropic-api-key"
  replication { auto {} }
}

resource "google_secret_manager_secret" "openai_key" {
  secret_id = "openai-api-key"
  replication { auto {} }
}

resource "google_secret_manager_secret" "google_calendar_secret" {
  secret_id = "google-calendar-client-secret"
  replication { auto {} }
}

# ─── Cloud Scheduler ─────────────────────────────────
resource "google_cloud_scheduler_job" "daily_brief" {
  name      = "daily-ai-brief"
  schedule  = "0 8 * * *"
  time_zone = "Europe/Moscow"

  http_target {
    uri         = "${google_cloud_run_v2_service.api.uri}/notifications/daily-brief-all"
    http_method = "POST"
    oidc_token {
      service_account_email = google_service_account.scheduler_sa.email
    }
  }
}

resource "google_cloud_scheduler_job" "weekly_report" {
  name      = "weekly-report"
  schedule  = "0 19 * * 0"   # Воскресенье 19:00
  time_zone = "Europe/Moscow"

  http_target {
    uri         = "${google_cloud_run_v2_service.api.uri}/notifications/weekly-report-all"
    http_method = "POST"
    oidc_token {
      service_account_email = google_service_account.scheduler_sa.email
    }
  }
}

resource "google_cloud_scheduler_job" "streak_check" {
  name      = "streak-check"
  schedule  = "0 21 * * *"   # Ежедневно 21:00
  time_zone = "Europe/Moscow"

  http_target {
    uri         = "${google_cloud_run_v2_service.api.uri}/notifications/streak-check-all"
    http_method = "POST"
    oidc_token {
      service_account_email = google_service_account.scheduler_sa.email
    }
  }
}

# ─── Cloud Tasks ─────────────────────────────────────
resource "google_cloud_tasks_queue" "notifications" {
  name     = "notifications"
  location = var.region

  rate_limits {
    max_concurrent_dispatches = 100
    max_dispatches_per_second = 50
  }

  retry_config {
    max_attempts  = 3
    max_backoff   = "300s"
    min_backoff   = "5s"
  }
}

# ─── Firestore ───────────────────────────────────────
resource "google_firestore_database" "default" {
  project     = var.project_id
  name        = "(default)"
  location_id = var.region
  type        = "FIRESTORE_NATIVE"
}

# ─── Service Accounts ────────────────────────────────
resource "google_service_account" "cloud_run_sa" {
  account_id   = "cloud-run-api"
  display_name = "Cloud Run API Service Account"
}

resource "google_service_account" "scheduler_sa" {
  account_id   = "cloud-scheduler"
  display_name = "Cloud Scheduler Service Account"
}

# Права Cloud Run SA
resource "google_project_iam_member" "cloud_run_firestore" {
  project = var.project_id
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.cloud_run_sa.email}"
}

resource "google_project_iam_member" "cloud_run_secrets" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.cloud_run_sa.email}"
}

resource "google_project_iam_member" "cloud_run_firebase_auth" {
  project = var.project_id
  role    = "roles/firebaseauth.admin"
  member  = "serviceAccount:${google_service_account.cloud_run_sa.email}"
}
```

```hcl
# infrastructure/variables.tf
variable "project_id" { default = "productivityai-prod" }
variable "region"     { default = "europe-west1" }
```

---

## Secret Manager — хранение ключей

```bash
# Добавить секрет
echo -n "sk-ant-xxx" | gcloud secrets create anthropic-api-key \
  --data-file=-

echo -n "sk-openai-xxx" | gcloud secrets create openai-api-key \
  --data-file=-

# Обновить версию
echo -n "sk-ant-new-xxx" | gcloud secrets versions add anthropic-api-key \
  --data-file=-

# Прочитать (для проверки)
gcloud secrets versions access latest --secret="anthropic-api-key"
```

Использование в Python:
```python
# services/secret_manager.py
from google.cloud import secretmanager

def get_secret(secret_id: str) -> str:
    client = secretmanager.SecretManagerServiceClient()
    name = f"projects/{PROJECT_ID}/secrets/{secret_id}/versions/latest"
    response = client.access_secret_version(name=name)
    return response.payload.data.decode("UTF-8")

# При старте приложения
ANTHROPIC_API_KEY = get_secret("anthropic-api-key")
OPENAI_API_KEY    = get_secret("openai-api-key")
```

---

## CI/CD — GitHub Actions

```yaml
# .github/workflows/deploy.yml
name: Deploy to Cloud Run

on:
  push:
    branches: [main]
    paths: ['backend/**']

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: google-github-actions/auth@v2
        with:
          credentials_json: ${{ secrets.GCP_SA_KEY }}

      - uses: google-github-actions/setup-gcloud@v2

      - name: Build and push Docker image
        run: |
          gcloud builds submit backend/ \
            --tag gcr.io/${{ vars.PROJECT_ID }}/ai-backend:${{ github.sha }}

      - name: Deploy to Cloud Run
        run: |
          gcloud run deploy ai-backend \
            --image gcr.io/${{ vars.PROJECT_ID }}/ai-backend:${{ github.sha }} \
            --region europe-west1 \
            --platform managed
```

---

## Firebase Rules деплой

```bash
# firebase.json
{
  "firestore": {
    "rules": "firestore.rules",
    "indexes": "firestore.indexes.json"
  }
}

# Деплой правил
firebase deploy --only firestore:rules --project productivityai-prod
```

---

## Окружения

| Окружение | Firebase проект | Cloud Run URL |
|-----------|----------------|---------------|
| Development | productivityai-dev | localhost:8080 |
| Staging | productivityai-staging | staging-api.xxx.run.app |
| Production | productivityai-prod | api.xxx.run.app |

```bash
# Переключение окружения
gcloud config set project productivityai-dev
firebase use dev
```

---

## Мониторинг

- **Cloud Run** — встроенные метрики (latency, errors, instances) в GCP Console
- **Firebase** — Firebase Console (DAU, crashes, Firestore usage)
- **Ошибки** — Cloud Error Reporting → уведомления в Slack
- **Логи** — Cloud Logging, запросы к AI логируются с userId (без контента)
