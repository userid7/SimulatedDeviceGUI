import { useState, useEffect } from "react";
import "./App.css";
import { GetReader, CreateReader } from "../wailsjs/go/main/App";
import HFCard from "./components/hf/hf.card";
import CreateHFCard from "./components/hf/create.hf.card";

function App() {
  const [resultText, setResultText] = useState(
    "Please enter your name below 👇"
  );
  const [name, setName] = useState("");
  const [readers, setReaders] = useState([]);
  const updateName = (e) => setName(e.target.value);
  const updateResultText = (result) => setResultText(result);

  // function greet() {
  //   Greet(name).then(updateResultText);
  // }

  useEffect(() => {
    const interval = setInterval(async () => {
      console.log("Interval Pooling");
      console.log(readers);
      var newReaders = await GetReader();
      console.log(newReaders);
      setReaders(newReaders || []);
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div id="App">
      <div className="min-h-screen px-6 py-8 justify-center">
        <div className="flex flex-col w-max min-w-max">
          <div className="text-white text-2xl font-bold font-mono">
            <h1>Komatsu HF Reader Simulator</h1>
          </div>
          {readers.map((reader, i) => {
            return <HFCard reader={reader} />;
          })}
          <CreateHFCard />
        </div>
      </div>
    </div>
  );
}

export default App;
