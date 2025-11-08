// UTF-8
import { useMemo, useState } from "react";
import {
  DndContext,
  pointerWithin,           // B) prioritize the column under the pointer
  DragOverlay,
  MouseSensor,
  TouchSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { useTasks } from "../hooks/useTasks";
import Column from "../components/Column";
import TaskCard from "../components/TaskCard";
import EditTaskModal from "../components/EditTaskModal";

export default function Board() {
  const { grouped, status, error, createTask, updateTask, deleteTask, reorderTask } = useTasks();
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");

  // A) local UI state for editing
  const [editing, setEditing] = useState(null); // task or null

  // B) DnD overlay + cleaner column targeting
  const [activeTaskId, setActiveTaskId] = useState(null);

  const sensors = useSensors(
    useSensor(MouseSensor, { activationConstraint: { distance: 6 } }),
    useSensor(TouchSensor, { pressDelay: 120, activationConstraint: { distance: 6 } })
  );

  const columns = useMemo(() => ["todo", "doing", "done"], []);

  async function handleCreate(e) {
    e.preventDefault();
    if (!title.trim()) return;
    await createTask({ title: title.trim(), description: desc.trim() || undefined });
    setTitle(""); setDesc("");
  }

  function findTaskById(id) {
    for (const col of columns) {
      const idx = grouped[col].findIndex((t) => t.id === id);
      if (idx >= 0) return { col, idx, task: grouped[col][idx] };
    }
    return null;
  }

  // B) DnD handlers
  function handleDragStart(event) {
    setActiveTaskId(event.active?.id ?? null);
  }
  function handleDragCancel() {
    setActiveTaskId(null);
  }
  function handleDragEnd(event) {
    const { active, over } = event;
    setActiveTaskId(null);
    if (!over) return;

    const activeId = active.id;
    const originCol = active.data?.current?.columnId;

    // Resolve destination column:
    // - If dropping over a task → use that task's column
    // - If dropping over a column droppable → use that column (append)
    let destCol = over.data?.current?.columnId;
    if (!destCol) {
      const maybe = String(over.id || "");
      if (maybe.startsWith("column:")) destCol = maybe.split(":")[1];
    }
    if (!destCol) return;

    // Compute destination index
    let destIndex = 0;
    if (over.data?.current?.type === "task") {
      const meta = findTaskById(over.id);
      destIndex = meta ? meta.idx : 0;
    } else {
      // Dropped on column surface → append at end
      destIndex = grouped[destCol].length;
    }

    // If same column and same position, ignore
    if (originCol === destCol) {
      const meta = findTaskById(activeId);
      if (meta && meta.idx === destIndex) return;
    }

    reorderTask(activeId, destCol, destIndex).catch((err) => {
      console.error("reorder failed:", err);
      alert("Reorder failed");
    });
  }

  // A) edit helpers
  function openEdit(task) { setEditing(task); }
  async function saveEdit(patch) {
    if (!editing) return;
    await updateTask(editing.id, patch);
    setEditing(null);
  }

  return (
    <div style={page()}>
      <h1 style={{ marginBottom: 8 }}>Kanban (Drag & Drop)</h1>

      <form onSubmit={handleCreate} style={form()}>
        <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="New task title…" style={input()} />
        <input value={desc}  onChange={(e) => setDesc(e.target.value)}  placeholder="Description (optional)…" style={input()} />
        <button type="submit" style={btnPrimary()}>Add</button>
      </form>

      {status === "loading" && <p>Loading tasks…</p>}
      {status === "error" && <p style={{ color: "salmon" }}>{error}</p>}

      <DndContext
        sensors={sensors}
        collisionDetection={pointerWithin}   // B) more intuitive column targeting
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        onDragCancel={handleDragCancel}
      >
        <div style={board()}>
          <Column id="todo"  tasks={grouped.todo}  onEdit={openEdit} onDelete={deleteTask}/>
          <Column id="doing" tasks={grouped.doing} onEdit={openEdit} onDelete={deleteTask}/>
          <Column id="done"  tasks={grouped.done}  onEdit={openEdit} onDelete={deleteTask}/>
        </div>

        {/* B) Drag overlay to reduce “snap” feeling while dragging */}
        <DragOverlay>
          {activeTaskId ? (
            (() => {
              const meta = findTaskById(activeTaskId);
              return meta ? (
                <div style={{ width: 280 }}>
                  <TaskCard task={meta.task} onEdit={() => {}} onDelete={() => {}} columnId={meta.col} index={meta.idx} />
                </div>
              ) : null;
            })()
          ) : null}
        </DragOverlay>
      </DndContext>

      {/* A) Edit modal */}
      {editing && (
        <EditTaskModal
          task={editing}
          onSave={saveEdit}
          onClose={() => setEditing(null)}
        />
      )}
    </div>
  );
}

// styles
function page(){ return { padding:16, fontFamily:"system-ui, sans-serif", color:"#e5e7eb", background:"#0a0f1a", minHeight:"100vh" }; }
function board(){ return { display:"grid", gridTemplateColumns:"repeat(3, 1fr)", gap:12, alignItems:"start" }; }
function form(){ return { display:"flex", gap:8, margin:"8px 0 16px 0" }; }
function input(){ return { flex:1, padding:"8px 10px", borderRadius:6, border:"1px solid #233052", background:"#0b1220", color:"#e5e7eb" }; }
function btnPrimary(){ return { padding:"8px 12px", borderRadius:6, border:"1px solid #233052", background:"#1f2937", color:"#e5e7eb", cursor:"pointer" }; }
