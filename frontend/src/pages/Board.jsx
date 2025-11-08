// UTF-8
import { useMemo, useState } from "react";
import { DndContext, closestCorners, MouseSensor, TouchSensor, useSensor, useSensors } from "@dnd-kit/core";
import { useTasks } from "../hooks/useTasks";
import Column from "../components/Column";

export default function Board() {
  const { grouped, status, error, createTask, updateTask, deleteTask, reorderTask } = useTasks();
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");

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

  function handleDragEnd(event) {
    const { active, over } = event;
    if (!over) return;

    const activeId = active.id;
    const originCol = active.data?.current?.columnId;
    let destCol = over.data?.current?.columnId;

    if (!destCol) {
      const maybeCol = String(over.id || "");
      if (maybeCol.startsWith("column:")) destCol = maybeCol.split(":")[1];
    }
    if (!destCol) return;

    let destIndex = 0;
    if (over.data?.current?.type === "task") {
      const meta = findTaskById(over.id);
      destIndex = meta ? meta.idx : 0;
    } else {
      destIndex = grouped[destCol].length;
    }

    if (originCol === destCol) {
      const meta = findTaskById(activeId);
      if (meta && meta.idx === destIndex) return;
    }

    reorderTask(activeId, destCol, destIndex).catch((err) => {
      console.error("reorder failed:", err);
      alert("Reorder failed");
    });
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

      <DndContext sensors={sensors} collisionDetection={closestCorners} onDragEnd={handleDragEnd}>
        <div style={board()}>
          <Column id="todo"  tasks={grouped.todo}  onEdit={(t)=>updateTask(t.id,{title: prompt("Edit title", t.title) || t.title})} onDelete={deleteTask}/>
          <Column id="doing" tasks={grouped.doing} onEdit={(t)=>updateTask(t.id,{title: prompt("Edit title", t.title) || t.title})} onDelete={deleteTask}/>
          <Column id="done"  tasks={grouped.done}  onEdit={(t)=>updateTask(t.id,{title: prompt("Edit title", t.title) || t.title})} onDelete={deleteTask}/>
        </div>
      </DndContext>
    </div>
  );
}

// styles
function page(){ return { padding:16, fontFamily:"system-ui, sans-serif", color:"#e5e7eb", background:"#0a0f1a", minHeight:"100vh" }; }
function board(){ return { display:"grid", gridTemplateColumns:"repeat(3, 1fr)", gap:12, alignItems:"start" }; }
function form(){ return { display:"flex", gap:8, margin:"8px 0 16px 0" }; }
function input(){ return { flex:1, padding:"8px 10px", borderRadius:6, border:"1px solid #233052", background:"#0b1220", color:"#e5e7eb" }; }
function btnPrimary(){ return { padding:"8px 12px", borderRadius:6, border:"1px solid #233052", background:"#1f2937", color:"#e5e7eb", cursor:"pointer" }; }
