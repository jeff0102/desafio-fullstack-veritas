// UTF-8
import { useEffect, useMemo, useState } from "react";
import { apiList, apiCreate, apiUpdate, apiDelete, apiReorder } from "../lib/api";

export function useTasks() {
  const [items, setItems] = useState([]);
  const [status, setStatus] = useState("idle");
  const [error, setError] = useState("");

  async function refresh() {
    try {
      setStatus("loading");
      const data = await apiList();
      setItems(data);
      setStatus("ready");
      setError("");
    } catch (err) {
      setStatus("error");
      setError(String(err?.message || err));
    }
  }

  useEffect(() => { refresh(); }, []);

  const grouped = useMemo(() => {
    const by = { todo: [], doing: [], done: [] };
    for (const t of items) if (by[t.status]) by[t.status].push(t);
    for (const k of Object.keys(by)) {
      by[k].sort((a, b) => a.order - b.order || (new Date(b.updatedAt) - new Date(a.updatedAt)));
    }
    return by;
  }, [items]);

  async function createTask({ title, description }) {
    const created = await apiCreate({ title, description });
    setItems(prev => [...prev, created]);
    return created;
  }

  async function updateTask(id, patch) {
    const updated = await apiUpdate(id, patch);
    setItems(prev => prev.map(t => (t.id === id ? updated : t)));
    return updated;
  }

  async function deleteTask(id) {
    await apiDelete(id);
    setItems(prev => prev.filter(t => t.id !== id));
  }

  async function reorderTask(id, destStatus, destIndex) {
    const prev = items;

    const optimistic = (() => {
      const moving = prev.find(t => t.id === id);
      if (!moving) return prev;
      const without = prev.filter(t => t.id !== id);
      const dest = without
        .filter(t => t.status === destStatus)
        .sort((a, b) => a.order - b.order || (new Date(b.updatedAt) - new Date(a.updatedAt)));
      const i = Math.max(0, Math.min(destIndex, dest.length));
      dest.splice(i, 0, { ...moving, status: destStatus });

      const reindexedDest = dest.map((t, idx) => ({ ...t, order: idx + 1 }));
      const others = without.filter(t => t.status !== destStatus);
      return [...others, ...reindexedDest];
    })();

    setItems(optimistic);
    try {
      const saved = await apiReorder(id, { status: destStatus, index: destIndex });
      setItems(cur => cur.map(t => (t.id === id ? saved : t)));
      const fresh = await apiList();
      setItems(fresh);
      return saved;
    } catch (err) {
      setItems(prev);
      throw err;
    }
  }

  return { data: items, grouped, status, error, refresh, createTask, updateTask, deleteTask, reorderTask };
}
