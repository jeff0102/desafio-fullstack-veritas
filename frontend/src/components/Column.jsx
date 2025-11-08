// UTF-8
import { useDroppable } from "@dnd-kit/core";
import { SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable";
import TaskCard from "./TaskCard";

const TITLES = { todo: "Todo", doing: "Doing", done: "Done" };

export default function Column({ id, tasks, onEdit, onDelete }) {
  const { setNodeRef, isOver } = useDroppable({
    id: `column:${id}`,
    data: { columnId: id, type: "column" },
  });

  return (
    <div style={col(isOver)} ref={setNodeRef}>
      <div style={colHeader()}>{TITLES[id] || id}</div>
      <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
        <SortableContext items={tasks.map((t) => t.id)} strategy={verticalListSortingStrategy}>
          {tasks.map((t, idx) => (
            <TaskCard
              key={t.id}
              task={t}
              index={idx}
              columnId={id}
              onEdit={onEdit}
              onDelete={onDelete}
            />
          ))}
        </SortableContext>
      </div>
    </div>
  );
}

function col(isOver) {
  return {
    flex: 1,
    minHeight: 300,
    padding: 12,
    borderRadius: 10,
    background: isOver ? "#0b1220" : "#0b0f19",
    border: "1px solid #233052",
    outline: isOver ? "2px dashed #3b82f6" : "none",
  };
}
function colHeader() {
  return {
    fontWeight: 700,
    paddingBottom: 8,
    marginBottom: 8,
    borderBottom: "1px solid #233052",
    color: "#e5e7eb",
  };
}
