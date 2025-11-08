// UTF-8
import { useEffect, useState } from "react";
import Modal from "./Modal";

export default function EditTaskModal({ task, onSave, onClose }) {
  const [title, setTitle] = useState(task?.title || "");
  const [desc, setDesc] = useState(task?.description || "");

  useEffect(() => {
    setTitle(task?.title || "");
    setDesc(task?.description || "");
  }, [task]);

  function submit(e) {
    e.preventDefault();
    const nextTitle = title.trim();
    if (!nextTitle) return;
    onSave({ title: nextTitle, description: desc.trim() || "" });
  }

  return (
    <Modal title="Edit task" onClose={onClose}>
      <form onSubmit={submit} style={{ display: "grid", gap: 10 }}>
        <label style={label()}>
          <span>Title</span>
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Task title"
            style={input()}
            maxLength={140}
            autoFocus
          />
        </label>
        <label style={label()}>
          <span>Description</span>
          <textarea
            value={desc}
            onChange={(e) => setDesc(e.target.value)}
            placeholder="Optional description"
            style={textarea()}
            maxLength={1000}
            rows={4}
          />
        </label>
        <div style={{ display: "flex", gap: 8, justifyContent: "flex-end" }}>
          <button type="button" onClick={onClose} style={btn()}>Cancel</button>
          <button type="submit" style={btnPrimary()}>Save</button>
        </div>
      </form>
    </Modal>
  );
}

function label(){ return { display:"grid", gap:6, fontSize:13, color:"#cbd5e1" }; }
function input(){ return { padding:"8px 10px", borderRadius:6, border:"1px solid #233052", background:"#0b1220", color:"#e5e7eb" }; }
function textarea(){ return { ...input(), resize:"vertical" }; }
function btn(){ return { padding:"8px 12px", borderRadius:6, border:"1px solid #233052", background:"#0b1220", color:"#e5e7eb", cursor:"pointer" }; }
function btnPrimary(){ return { ...btn(), background:"#1f2937" }; }
