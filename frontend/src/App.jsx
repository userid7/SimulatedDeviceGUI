import { useState, useEffect } from "react";
import "./App.css";
import { GetAllReader } from "../wailsjs/go/hfreader/App";
import HFCard from "./components/hf/hf.card";
import CreateHFCard from "./components/hf/create.hf.card";

function App() {
  const [readers, setReaders] = useState([]);

  useEffect(() => {
    const interval = setInterval(async () => {
      console.log("Interval Pooling");
      // console.log(readers);
      var newReaders = await GetAllReader();
      console.log(newReaders);
      setReaders(newReaders || []);
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  // useEffect(() => {
  //   console.log("readers");
  //   console.log(readers);
  // }, [readers]);

  return (
    <div id="App">
      <div className="min-h-screen px-6 py-8 justify-center">
        <div className="flex flex-col w-max min-w-max">
          <div className="text-white text-2xl font-bold font-mono">
            <h1>Device Simulator</h1>
          </div>
          <div className="w-96">
            {readers.map((reader, i) => {
              return <HFCard reader={reader} />;
            })}
            <CreateHFCard />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
