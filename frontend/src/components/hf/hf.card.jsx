import { useState, useEffect, useRef } from "react";
import { _ } from "lodash";
import Switch from "@mui/material/Switch";
import {
  DeleteReader,
  SetReaderCardPresent,
  SetReaderConnection,
  SetReaderEpc,
} from "../../../wailsjs/go/hfreader/App";
import CardHeader from "../card/card.header";

function HFCard(props) {
  const [epc, setEpc] = useState(props.reader.UidBuffer);
  const [isCardPresent, setIsCardPresent] = useState(
    props.reader.IsCardPresent
  );
  const [isConnected, setIsConnected] = useState(props.reader.IsConnected);

  const deleteReaderHandler = (e) => {
    console.log("Delete Reader");
    e.preventDefault();
    DeleteReader(props.reader.Id);
  };

  const updateEpc = (e) => setEpc(e.target.value);
  useEffect(() => {
    const delayDebounceFn = setTimeout(() => {
      console.log(epc);
      SetReaderEpc(props.reader.Id, epc);
    }, 1000);

    return () => clearTimeout(delayDebounceFn);
  }, [epc]);

  const handleSwitchChange = (e) => {
    setIsCardPresent(e.target.checked);
    SetReaderCardPresent(props.reader.Id, e.target.checked);
  };

  const handleConnectionChange = (e) => {
    console.log("Connection changed");
    console.log(!isConnected);
    setIsConnected(!isConnected);
    SetReaderConnection(props.reader.Id, !isConnected);
  };

  useEffect(() => {
    console.log("props.reader");
    console.log(props.reader);
    setEpc(props.reader.UidBuffer);
    setIsConnected(props.reader.IsConnected);
    setIsCardPresent(props.reader.IsCardPresent);
  }, [JSON.stringify(props.reader)]);

  return (
    <div className="block p-4 rounded-lg shadow-lg bg-white border-gray-200 hover:bg-gray-100 my-2">
      {/* <div className="flex flex-row justify-between">
        <div className=" text-xl font-bold ">
          Reader{" "}
          {props.reader.Line +
            "-" +
            props.reader.Post +
            "-" +
            props.reader.Code}
        </div>
        <div className="flex flex-row justify-center">
          <div className="flex justify-center items-center px-1">
            <ConnectionButton
              isConnected={isConnected}
              onClick={handleConnectionChange}
            />
          </div>
          <div className="flex justify-center items-center px-1">
            <DeleteButton onClick={deleteReaderHandler} />
          </div>
        </div>
      </div>
      <div className="text-xs px-2 text-slate-400">
        target : {props.reader.TargetUrl}
      </div>
      <hr class="my-2 h-px bg-gray-200 border-0 dark:bg-gray-700"></hr> */}

      <CardHeader
        title={
          "Reader " +
          props.reader.Line +
          "-" +
          props.reader.Post +
          "-" +
          props.reader.Code
        }
        host={props.reader.TargetUrl}
        isConnected={isConnected}
        onClickConnection={handleConnectionChange}
        onClickDelete={deleteReaderHandler}
      />
      <div id="hf" className="flex flex-row justify-between py-1 px-1">
        <div className="flex flex-row justify-center items-center gap-2 w-60">
          <input
            id="name"
            className="form-control
            block
            w-full
            px-4
            py-2
            text-md
            font-normal
            text-gray-700
            bg-white bg-clip-padding
            border border-solid border-gray-300
            rounded
            transition
            ease-in-out
            m-0
            focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none
          "
            onChange={updateEpc}
            autoComplete="off"
            placeholder="EPC"
            name="input"
            type="text"
            disabled={isCardPresent}
            value={epc}
          />
        </div>
        <div className="flex flex-row justify-center">
          <div className="flex justify-center items-center px-1">
            <Switch checked={isCardPresent} onChange={handleSwitchChange} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default HFCard;
