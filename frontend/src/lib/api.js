// UTF-8
const BASE = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

export async function apiList({ status } = {}) {
  const url = new URL(`${BASE}/tasks`);
  if (status) url.searchParams.set("status", status);
  const res = await fetch(url.toString());
  if (!res.ok) throw new Error(`list_failed: ${res.status}`);
  return res.json();
}

export async function apiCreate({ title, description }) {
  const res = await fetch(`${BASE}/tasks`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title, description }),
  });
  if (!res.ok) throw new Error(`create_failed: ${res.status}`);
  return res.json();
}

export async function apiUpdate(id, patch) {
  const res = await fetch(`${BASE}/tasks/${encodeURIComponent(id)}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(patch),
  });
  if (!res.ok) throw new Error(`update_failed: ${res.status}`);
  return res.json();
}

export async function apiDelete(id) {
  const res = await fetch(`${BASE}/tasks/${encodeURIComponent(id)}`, { method: "DELETE" });
  if (!res.ok) throw new Error(`delete_failed: ${res.status}`);
}

export async function apiReorder(id, { status, index }) {
  const res = await fetch(`${BASE}/tasks/${encodeURIComponent(id)}/reorder`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ status, index }),
  });
  if (!res.ok) throw new Error(`reorder_failed: ${res.status}`);
  return res.json();
}
