// Simple fetch client for backend connectivity checks.
const BASE_URL = (import.meta.env.VITE_API_BASE_URL || "").replace(/\/$/, "");

export async function pingTasks() {
  try {
    const res = await fetch(`${BASE_URL}/tasks`, { method: "GET" });
    return res.ok ? "ok" : "fail";
  } catch {
    return "fail";
  }
}

export async function fetchTasks() {
  const res = await fetch(`${BASE_URL}/tasks`, { method: "GET" });
  if (!res.ok) throw new Error(`GET /tasks failed: ${res.status}`);
  return res.json();
}
