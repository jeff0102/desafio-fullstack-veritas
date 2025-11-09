# desafio-fullstack-veritas

Mini Kanban board (Todo / Doing / Done) with **Frontend: React (Vite)** and **Backend: Go (chi)**.  
Supports create, edit (title + description), drag & drop across/within columns, and delete.  
Persistence via a simple JSON store behind a REST API.

> This is a **skeleton README** (no diagrams yet). Diagrams will be added later under `/docs`.

---

## 1) Project layout

```
/backend
  cmd/api/main.go
  internal/
    http/        # router + handlers
    core/        # domain (Task)
    store/       # memory + JSON persistence
  tasks.json     # runtime data (ignored via .gitignore)

/frontend
  package.json
  src/
    lib/api.js
    hooks/useTasks.js
    components/{Column.jsx,TaskCard.jsx,Modal.jsx,EditTaskModal.jsx}
    pages/Board.jsx
    App.jsx
  .env.example   # VITE_API_BASE_URL=http://localhost:8080

/docs
  # user-flow.png (required) - to be added later
  # data-flow.png (optional) - to be added later

README.md
```

---

## 2) How to run

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

**Notes**
- The compose file uses a **named volume** `frontend_node_modules` mounted at `/app/node_modules` to persist dependencies inside the container.
- The frontend container auto-installs on startup **only if** `node_modules/.bin/vite` is missing; step (2) guarantees a clean, stable dev experience.

---
## 3) API (summary)

Base URL: `http://localhost:8080`

- `GET    /tasks`               → list (supports `?status=todo|doing|done`; sorted by `order`)
- `POST   /tasks`               → create `{ title, description? }`
- `GET    /tasks/{id}`          → fetch by id
- `PUT    /tasks/{id}`          → update partial `{ title?, description?, status? }`
- `DELETE /tasks/{id}`          → delete
- `PUT    /tasks/{id}/reorder`  → move + reindex `{ status, index }`

Notes:
- CORS enabled for `http://localhost:5173` (Vite).
- `order` is stable within each status; the server reindexes when items move.

---

## 4) Frontend UX (high-level)

- Three fixed columns: **Todo**, **Doing**, **Done**.
- Add tasks (title + optional description).
- Edit via **modal** (title + description).
- Drag & drop (dnd-kit) with **column-priority** targeting and **DragOverlay**.
- Basic feedback: loading / error messages.

---

## 5) Technical decisions (short)

- **Go + chi** router; layered backend:
  - `core` (domain `Task` with `status`, `order`, timestamps)
  - `store` (in-memory + JSON persisted)
  - `http` (router + handlers)
- **JSON persistence** (bonus) for quick local state without a DB.
- **React + Vite**, **pnpm** for dependency management.
- **dnd-kit** for modern drag & drop; **pointerWithin** collision and **DragOverlay**.
- **Optimistic UI** on reorder; server-side reindex to guarantee final order.

---

## 6) Known limitations

- No auth (local demo).
- No pagination/search.
- No automated tests yet.
- Single-process writing model (JSON is not a multi-user DB).

---

## 7) Future work

- Minimal API/store tests + CI.
- Accessibility improvements (non-drag interactions).
- Search/filter UI; inline edits; keyboard shortcuts.
- Migrate persistence to SQLite/Postgres if multi-user or larger datasets.

---

## 8) Docs & deliverables

- **Required**: `docs/user-flow.png` (to be added).
- **Optional**: `docs/data-flow.png` (to be added).

---

## 9) Quick tips

- Reset data: stop API, delete `backend/tasks.json`, restart API.
- If frontend can’t reach API, ensure `frontend/.env` contains:
  ```
  VITE_API_BASE_URL=http://localhost:8080
  ```
- If `pnpm` isn’t recognized, run:
  ```
  corepack enable
  corepack prepare pnpm@latest --activate
  ```
