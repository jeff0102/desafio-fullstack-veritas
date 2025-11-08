import { useEffect, useState } from "react";
import { fetchTasks } from "../lib/api.js";

export function useTasks() {
  const [data, setData] = useState(null);
  const [status, setStatus] = useState("idle"); // idle | loading | ready | error
  const [error, setError] = useState(null);

  useEffect(() => {
    let active = true;
    setStatus("loading");
    fetchTasks()
      .then(tasks => { if (active){ setData(tasks); setStatus("ready"); }})
      .catch(e => { if (active){ setError(String(e)); setStatus("error"); }});
    return () => { active = false; };
  }, []);

  return { data, status, error };
}
