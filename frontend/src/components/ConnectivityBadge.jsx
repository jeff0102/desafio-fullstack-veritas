import { useEffect, useState } from "react";
import { pingTasks } from "../lib/api.js";

export default function ConnectivityBadge() {
  const [state, setState] = useState("checking"); // checking | ok | fail

  useEffect(() => {
    let mounted = true;
    const check = async () => {
      const r = await pingTasks();
      if (mounted) setState(r === "ok" ? "ok" : "fail");
    };
    check();
    const id = setInterval(check, 3000);
    return () => { mounted = false; clearInterval(id); };
  }, []);

  const style = { padding:"4px 8px", borderRadius:6, fontSize:12, display:"inline-block", border:"1px solid", marginBottom:12 };

  if (state === "checking") return <span style={style}>Checking API�</span>;
  if (state === "ok") return <span style={{...style, borderColor:"green"}}>API: OK</span>;
  return <span style={{...style, borderColor:"red"}}>API: FAIL</span>;
}
