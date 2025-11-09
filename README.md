# desafio-fullstack-veritas

Mini Kanban board (Todo / Doing / Done) with **Frontend: React (Vite)** and **Backend: Go (chi)**.  

## Technical decisions

This section explains *why* choices were made and the **trade-offs** accepted.

### Backend (Go + chi)
- **Why Go + chi**
  - Small, explicit router; easy to reason about handlers/middleware.
  - Go stdlib covers most needs; fast startup.
- **Layering: `core` / `store` / `http`**
  - Separates domain from transport/persistence; swapping the store is trivial.
  - *Trade-off:* Slight overhead for a small codebase, but enables growth.
- **JSON persistence (file) instead of DB**
  - Instant local persistence, zero infra; fine for single-user demo.
  - *Trade-off:* No concurrency guarantees or queries; would migrate to **PostgreSQL** if multi-user.
- **Ordering model: `status` + `order`**
  - Server reindexes on move → stable ordering per column; simpler frontend.
  - *Trade-off:* No “true” concurrent edit handling.
- **Endpoint `PUT /tasks/{id}/reorder`**
  - Explicit intent (move + index); keeps update vs reorder semantics clean.
  - *Trade-off:* One extra endpoint vs overloading `PUT /tasks/{id}`.

---

### Frontend (React + Vite + dnd-kit)
- **Why plain JS (no TypeScript)**
  - Smaller surface and faster setup.
  - *Trade-off:* Fewer compile-time checks.
- **Why dnd-kit (not react-beautiful-dnd)**
  - Maintained, flexible sensors, `DragOverlay`, solid collision strategies.
  - *Trade-off:* Slightly lower-level API → a few more lines, more control.
- **UX choices**
  - Column-priority targeting + `DragOverlay`; edit modal (title + description).
  - *Trade-off:* Minimal accessibility (keyboard DnD omitted to stay in scope).

---

## 1) How to run

### A) With Docker (recommended)

**Prerequisite:** Docker Desktop running.

From the repository root:

```bash
# 0) Clean previous state (optional)
docker compose down -v

# 1) Fresh build (ensures updated CMD/entrypoints are used)
docker compose build --no-cache

# 2) Hydrate frontend deps into the named volume
docker compose run --rm frontend pnpm install

# 3) Approve pnpm postinstall builds once
docker compose run --rm frontend pnpm approve-builds

# 4) Bring services up
docker compose up -d
```

- API → `http://localhost:8080`  
- App → `http://localhost:5173`

> Note: Compose uses a **named volume** for `/app/node_modules` so the bind-mount `./frontend:/app` doesn’t hide dependencies.

### B) Local (manual, no Docker)

**Prerequisites:** Go 1.21+, Node 18+.

#### Install pnpm

Windows fallback 
```powershell
# Install pnpm globally using npm.cmd
npm.cmd install -g pnpm
```

**Backend (terminal A)**
```bash
go mod download -C backend
go run -C backend ./cmd/api   # http://localhost:8080
```

**Frontend (terminal B)**
```bash
cd frontend
cp .env.example .env          # ensure VITE_API_BASE_URL=http://localhost:8080
pnpm install
pnpm approve-builds           # if prompted (esbuild)
pnpm dev                      # http://localhost:5173
```

---

## 2) API (minimal summary)

Base URL: `http://localhost:8080`

- `GET    /tasks`               → list (`?status=todo|doing|done`; ordered by `order`)
- `POST   /tasks`               → create `{ title, description? }`
- `GET    /tasks/{id}`          → by id
- `PUT    /tasks/{id}`          → partial update `{ title?, description?, status? }`
- `DELETE /tasks/{id}`          → delete
- `PUT    /tasks/{id}/reorder`  → move + reindex `{ status, index }`

---

## 3) Known limitations

- No auth (local demo).
- No pagination/search.
- No automated tests.
- JSON file is single-process oriented.

---

## 4) Docs

- `docs/user-flow.png`

## 5) Next steps

- Migrate store to **PostgreSQL**; add migrations and repository interfaces.
- Add TypeScript, basic unit tests (store/http) and CI.
- Improve accessibility and keyboard interactions; audit focus management.
- Add filtering/search and pagination if the dataset grows.
