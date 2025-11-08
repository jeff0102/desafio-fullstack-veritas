// UTF-8
export default function Modal({ title, children, onClose }) {
  // Simple, dependency-free modal (ESC to close)
  function handleKey(e) {
    if (e.key === "Escape") onClose?.();
  }
  return (
    <div
      onKeyDown={handleKey}
      tabIndex={-1}
      style={backdrop()}
      role="dialog"
      aria-modal="true"
    >
      <div style={sheet()}>
        <div style={header()}>
          <strong>{title}</strong>
          <button onClick={onClose} style={xbtn()} aria-label="Close">×</button>
        </div>
        <div style={{ padding: 12 }}>{children}</div>
      </div>
    </div>
  );
}

function backdrop() {
  return {
    position: "fixed", inset: 0, background: "rgba(0,0,0,.5)",
    display: "grid", placeItems: "center", zIndex: 50
  };
}
function sheet() {
  return {
    width: "min(520px, 96vw)", background: "#0b1220", color: "#e5e7eb",
    border: "1px solid #233052", borderRadius: 10, boxShadow: "0 12px 30px rgba(0,0,0,.4)"
  };
}
function header() {
  return {
    display: "flex", justifyContent: "space-between", alignItems: "center",
    padding: "10px 12px", borderBottom: "1px solid #233052"
  };
}
function xbtn() {
  return {
    border: "1px solid #233052", background: "transparent", color: "#e5e7eb",
    borderRadius: 6, padding: "2px 8px", cursor: "pointer", fontSize: 18, lineHeight: 1
  };
}
