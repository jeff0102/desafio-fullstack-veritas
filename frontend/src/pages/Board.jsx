import ConnectivityBadge from "../components/ConnectivityBadge.jsx";
import { useTasks } from "../hooks/useTasks.js";

export default function Board() {
  const { data, status, error } = useTasks();

  return (
    <div style={{ padding: 16, fontFamily: "system-ui, sans-serif" }}>
      <h1>Kanban (skeleton)</h1>
      <ConnectivityBadge />
      {status === "loading" && <p>Loading tasks�</p>}
      {status === "error" && <p style={{ color: "red" }}>{error}</p>}
      {status === "ready" && (
        <pre style={{ background:"#111", color:"#ddd", padding:12, borderRadius:8, overflow:"auto" }}>
{JSON.stringify(data, null, 2)}
        </pre>
      )}
    </div>
  );
}
