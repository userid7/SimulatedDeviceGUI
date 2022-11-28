import { useState, useEffect } from "react";
import "./App.css";
import { GetAllReader } from "../wailsjs/go/hfreader/App";
import HFList from "./components/hf/hf.list";
import PMGatewayList from "./components/pm/pm.gateway.list";

function App() {
  return (
    <div id="App">
      <div className="min-h-screen px-6 py-8 justify-center">
        <div className="flex flex-col w-max min-w-max">
          <div className="text-white text-2xl font-bold font-mono">
            <h1>Device Simulator ALPHA</h1>
          </div>
          <div className="flex flex-row gap-x-5">
            <HFList />
            <PMGatewayList />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
