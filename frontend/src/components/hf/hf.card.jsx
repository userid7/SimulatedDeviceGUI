import { useState, useEffect, useRef } from "react";
import { _ } from "lodash";
import Switch from "@mui/material/Switch";
import {
  DeleteReader,
  SetReaderCardPresent,
  SetReaderConnection,
  SetReaderEpc,
} from "../../../wailsjs/go/main/App";

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
    }, 3000);

    return () => clearTimeout(delayDebounceFn);
  }, [epc]);

  const handleSwitchChange = (e) => {
    setIsCardPresent(e.target.checked);
    SetReaderCardPresent(props.reader.Id, e.target.checked);
  };

  const handleConnectionChange = (e) => {
    setIsConnected(e.target.value);
    SetReaderConnection(props.reader.Id, e.target.value);
  };

  // const usePrevious = (value) => {
  //   const ref = useRef();
  //   useEffect(() => {
  //     ref.current = value;
  //   });
  //   return ref.current;
  // };

  // const myPreviousReader = usePrevious(props.reader);

  // useEffect(() => {
  //   console.log("useEffect");
  //   console.log(myPreviousReader);
  //   console.log(props.reader);
  //   if (myPreviousReader) {
  //     if (!_.isEqual(myPreviousReader.UidBuffer, props.reader.UidBuffer)) {
  //       console.log("Epc Change detected");
  //       setEpc(props.reader.UidBuffer);
  //     }
  //     if (
  //       !_.isEqual(myPreviousReader.IsCardPresent, props.reader.IsCardPresent)
  //     ) {
  //       console.log("IsCardPresenet Change detected");
  //       setIsCardPresent(props.reader.IsCardPresent);
  //     }
  //   }
  // }, [props.reader]);

  return (
    <div className="block p-4 rounded-lg shadow-lg bg-white border-gray-200 hover:bg-gray-100 my-2">
      <div className="flex flex-row justify-between">
        <div className=" text-xl font-bold ">
          Reader{" "}
          {props.reader.Line +
            "-" +
            props.reader.Post +
            "-" +
            props.reader.Code}
        </div>
        <button onClick={deleteReaderHandler}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
            fill="currentColor"
            className="w-5 h-5"
          >
            <path
              fillRule="evenodd"
              d="M8.75 1A2.75 2.75 0 006 3.75v.443c-.795.077-1.584.176-2.365.298a.75.75 0 10.23 1.482l.149-.022.841 10.518A2.75 2.75 0 007.596 19h4.807a2.75 2.75 0 002.742-2.53l.841-10.52.149.023a.75.75 0 00.23-1.482A41.03 41.03 0 0014 4.193V3.75A2.75 2.75 0 0011.25 1h-2.5zM10 4c.84 0 1.673.025 2.5.075V3.75c0-.69-.56-1.25-1.25-1.25h-2.5c-.69 0-1.25.56-1.25 1.25v.325C8.327 4.025 9.16 4 10 4zM8.58 7.72a.75.75 0 00-1.5.06l.3 7.5a.75.75 0 101.5-.06l-.3-7.5zm4.34.06a.75.75 0 10-1.5-.06l-.3 7.5a.75.75 0 101.5.06l.3-7.5z"
              clipRule="evenodd"
            />
          </svg>
        </button>
      </div>
      <div id="hf" className="flex flex-row justify-between py-2">
        <div className="flex flex-col justify-center pr-2">
          {/* EPC */}
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
        {/* TODO : TRY TO USE FLOWBITE */}
        <div className="flex flex-row justify-center">
          <Switch checked={isCardPresent} onChange={handleSwitchChange} />
          <div className="flex justify-center px-2">
            <button>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                className="w-5 h-5"
              >
                <path d="M5.25 3A2.25 2.25 0 003 5.25v9.5A2.25 2.25 0 005.25 17h9.5A2.25 2.25 0 0017 14.75v-9.5A2.25 2.25 0 0014.75 3h-9.5z" />
              </svg>
            </button>
          </div>
          <div className="flex justify-center px-2">
            <button>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                className="w-5 h-5"
              >
                <path
                  fillRule="evenodd"
                  d="M7.84 1.804A1 1 0 018.82 1h2.36a1 1 0 01.98.804l.331 1.652a6.993 6.993 0 011.929 1.115l1.598-.54a1 1 0 011.186.447l1.18 2.044a1 1 0 01-.205 1.251l-1.267 1.113a7.047 7.047 0 010 2.228l1.267 1.113a1 1 0 01.206 1.25l-1.18 2.045a1 1 0 01-1.187.447l-1.598-.54a6.993 6.993 0 01-1.929 1.115l-.33 1.652a1 1 0 01-.98.804H8.82a1 1 0 01-.98-.804l-.331-1.652a6.993 6.993 0 01-1.929-1.115l-1.598.54a1 1 0 01-1.186-.447l-1.18-2.044a1 1 0 01.205-1.251l1.267-1.114a7.05 7.05 0 010-2.227L1.821 7.773a1 1 0 01-.206-1.25l1.18-2.045a1 1 0 011.187-.447l1.598.54A6.993 6.993 0 017.51 3.456l.33-1.652zM10 13a3 3 0 100-6 3 3 0 000 6z"
                  clipRule="evenodd"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default HFCard;
