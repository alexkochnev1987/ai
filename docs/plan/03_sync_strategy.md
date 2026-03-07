# Стратегия синхронизации

## Принцип: Offline-First

Приложение полностью работает без интернета. Все операции пишутся локально в SQLite (WatermelonDB) и синхронизируются с Firestore при появлении сети.

```
Пользователь создаёт задачу
        │
        ▼
  WatermelonDB       ← немедленно, без задержки
  (is_synced=false)
        │
  Есть сеть?
   ├── Да → SyncService.push() → Firestore
   └── Нет → sync_queue → попытка при reconnect
```

---

## SyncService

```typescript
class SyncService {

  // Запускается при:
  // - старте приложения
  // - переходе в foreground
  // - появлении сети (NetInfo listener)
  // - явном pull-to-refresh пользователем
  async sync(): Promise<void> {
    if (!this.isOnline()) return
    if (this.isSyncing) return      // защита от параллельных вызовов

    this.isSyncing = true
    try {
      await this.pushLocalChanges()   // 1. локальное → сервер
      await this.pullServerChanges()  // 2. сервер → локальное
      await this.updateSyncTimestamp()
    } finally {
      this.isSyncing = false
    }
  }

  // Отправка несинхронизированных записей
  private async pushLocalChanges() {
    const collections = ['tasks', 'habits', 'habit_logs', 'mood_logs', 'pomodoro_sessions']

    for (const collection of collections) {
      const unsynced = await database
        .get(collection)
        .query(Q.where('is_synced', false))
        .fetch()

      for (const record of unsynced) {
        if (record.deletedAt) {
          await firestoreService.delete(collection, record.serverId)
        } else {
          await firestoreService.upsert(collection, record.toFirestore())
        }
        await record.update(r => { r.isSynced = true })
      }
    }

    // Обрабатываем очередь офлайн операций
    await this.processSyncQueue()
  }

  // Получение изменений с сервера (от других устройств)
  private async pullServerChanges() {
    const lastSync = await this.getLastSyncTimestamp()

    const collections = ['tasks', 'habits', 'habit_logs']
    for (const collection of collections) {
      const changes = await firestore
        .collection(`users/${uid}/${collection}`)
        .where('updatedAt', '>', lastSync)
        .get()

      for (const doc of changes.docs) {
        await this.mergeDocument(collection, doc.id, doc.data())
      }
    }
  }
}
```

---

## Стратегия разрешения конфликтов

### Last-Write-Wins по updatedAt

```
Локальный updatedAt  vs  Серверный updatedAt
        │                        │
  Что новее? ──────────────────────
        │
   Локальное новее → пушим на сервер (уже в pushLocalChanges)
   Серверное новее → обновляем локальную запись
```

### Исключения из LWW

| Поле | Стратегия |
|------|-----------|
| `actualMinutes` | Суммирование (merge, не overwrite) |
| `currentStreak` | Берём максимальное значение |
| `subtasks` | Merge по id подзадачи |
| `deletedAt` | Delete всегда побеждает |

```typescript
private async mergeDocument(collection: string, id: string, serverData: any) {
  const local = await database.get(collection)
    .query(Q.where('server_id', id))
    .fetch()

  if (local.length === 0) {
    // Новый документ с сервера — создаём локально
    await database.write(async () => {
      await database.get(collection).create(record => {
        record.fromFirestore(serverData)
        record.isSynced = true
      })
    })
    return
  }

  const localRecord = local[0]
  const serverTime = serverData.updatedAt.toMillis()
  const localTime = localRecord.updatedAt

  if (serverTime > localTime) {
    // Серверная версия новее
    await database.write(async () => {
      await localRecord.update(record => {
        record.fromFirestore(serverData)
        record.isSynced = true
      })
    })
  }
  // Иначе локальная версия новее — она будет запушена в pushLocalChanges
}
```

---

## Реалтайм подписки (онлайн)

Для совместных проектов нужны реалтайм обновления:

```typescript
// Подписка на задачи проекта
const unsubscribe = firestore
  .collection(`projects/${projectId}/tasks`)
  .onSnapshot((snapshot) => {
    snapshot.docChanges().forEach(async (change) => {
      switch (change.type) {
        case 'added':
        case 'modified':
          await syncService.mergeDocument('tasks', change.doc.id, change.doc.data())
          break
        case 'removed':
          await syncService.softDeleteLocal('tasks', change.doc.id)
          break
      }
    })
  })

// Отписываемся при unmount / выходе из проекта
return () => unsubscribe()
```

---

## Очередь офлайн операций (sync_queue)

Если операция произошла офлайн и нет записи в WatermelonDB (напр. удаление):

```typescript
// Добавление в очередь
async addToQueue(operation: 'create' | 'update' | 'delete', collection: string, id: string, payload: object) {
  await database.write(async () => {
    await database.get('sync_queue').create(record => {
      record.operation = operation
      record.collection = collection
      record.documentId = id
      record.payloadJson = JSON.stringify(payload)
      record.retryCount = 0
    })
  })
}

// Обработка очереди при подключении
async processSyncQueue() {
  const queue = await database.get('sync_queue')
    .query(Q.sortBy('created_at', Q.asc))
    .fetch()

  for (const item of queue) {
    try {
      await this.executeQueueItem(item)
      await database.write(async () => { await item.destroyPermanently() })
    } catch (error) {
      if (item.retryCount >= 3) {
        // Логируем ошибку, удаляем из очереди
        await database.write(async () => { await item.destroyPermanently() })
      } else {
        await database.write(async () => {
          await item.update(r => { r.retryCount += 1 })
        })
      }
    }
  }
}
```

---

## Индикатор синхронизации в UI

```typescript
// store/syncStore.ts (Zustand)
interface SyncStore {
  status: 'synced' | 'syncing' | 'pending' | 'offline' | 'error'
  lastSyncedAt: Date | null
  pendingCount: number
}
```

Состояния в хедере:
- `synced` — зелёная галочка, время последней синхронизации
- `syncing` — анимированная иконка
- `pending` — серая иконка, количество несинхронизированных записей
- `offline` — оранжевая иконка облака с крестом
- `error` — красная иконка, кнопка "Повторить"
