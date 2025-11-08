// UTF-8
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";

export default function TaskCard({ task, onEdit, onDelete, columnId, index }) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } =
    useSortable({ id: task.id, data: { columnId, index, type: "task" } });

  const style = {
    padding: 12,
    borderRadius: 8,
    background: isDragging ? "#1f2937" : "#111827",
    color: "#e5e7eb",
    border: "1px solid #374151",
    boxShadow: isDragging ? "0 6px 12px rgba(0,0,0,.4)" : "none",
    transform: CSS.Transform.toString(transform),
    transition,
    cursor: "grab",
  };

  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
      <div style={{ display: "flex", justifyContent: "space-between", gap: 8 }}>
        <strong>{task.title}</strong>
        <div style={{ display: "flex", gap: 6 }}>
          <button onClick={() => onEdit(task)} title="Edit" style={btn()}>✏️</button>
          <button onClick={() => onDelete(task.id)} title="Delete" style={btn()}>🗑️</button>
        </div>
      </div>
      {task.description ? (
        <p style={{ marginTop: 6, fontSize: 13, color: "#cbd5e1" }}>{task.description}</p>
      ) : null}
    </div>
  );
}

function btn() {
  return {
    border: "1px solid #374151",
    background: "#0b1220",
    color: "#e5e7eb",
    borderRadius: 6,
    padding: "2px 6px",
    cursor: "pointer",
  };
}
