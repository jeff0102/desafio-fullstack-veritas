import React from "react";

export default class ErrorBoundary extends React.Component {
  constructor(props){ super(props); this.state = { hasError:false, error:null }; }
  static getDerivedStateFromError(error){ return { hasError:true, error }; }
  componentDidCatch(error, info){ console.error("[ErrorBoundary]", error, info); }

  render(){
    if (this.state.hasError){
      return (
        <div style={{ padding:16, color:"white", background:"#b00020" }}>
          <h2>UI crashed</h2>
          <pre>{String(this.state.error)}</pre>
        </div>
      );
    }
    return this.props.children;
  }
}
