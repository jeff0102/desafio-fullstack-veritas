import React from "react";
import { createRoot } from "react-dom/client";
import App from "./App.jsx";
import ErrorBoundary from "./components/ErrorBoundary.jsx";

console.log("[main.jsx] mounting React");
const el = document.getElementById("root");
if (!el) throw new Error("Root element #root not found in index.html");

createRoot(el).render(
  <ErrorBoundary>
    <App />
  </ErrorBoundary>
);
