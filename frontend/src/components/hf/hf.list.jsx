import { useState, useEffect } from "react";
import { GetAllReader } from "../../../wailsjs/go/hfreader/App";
import HFCard from "./hf.card";
import CreateHFCard from "./create.hf.card";

function HFList() {
  const [readers, setReaders] = useState([]);

  useEffect(() => {
    const interval = setInterval(async () => {
      console.log("Interval Pooling");
      var newReaders = await GetAllReader();
      console.log(newReaders);
      setReaders(newReaders || []);
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="flex flex-col w-96">
      {readers.map((reader, i) => {
        return <HFCard reader={reader} />;
      })}
      <CreateHFCard />
    </div>
  );
}

export default HFList;
