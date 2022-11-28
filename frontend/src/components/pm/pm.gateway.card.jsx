import { useState, useEffect, useRef } from "react";
import { _ } from "lodash";
import {
  DeletePMGateway,
  SetPMGatewayConnection,
} from "../../../wailsjs/go/pm/App";
import CreatePMCard from "./create.pm.card";
import PMCard from "./pm.card";
import CardHeader from "../card/card.header";

function PMGatewayCard(props) {
  const [isConnected, setIsConnected] = useState(props.pmGateway.IsConnected);

  const [pmGateway, setPMGateway] = useState(props.pmGateway);

  const deletePMGatewayHandler = (e) => {
    console.log("Delete PM Gateway");
    e.preventDefault();
    DeletePMGateway(props.pmGateway.Id);
  };

  const handleConnectionChange = (e) => {
    console.log("Connection changed");
    console.log(!isConnected);
    setIsConnected(!isConnected);
    SetPMGatewayConnection(props.pmGateway.Id, !isConnected);
  };

  useEffect(() => {
    console.log("props.pmGateway");
    console.log(props.pmGateway);
    setPMGateway(props.pmGateway);
  }, [JSON.stringify(props.pmGateway)]);

  return (
    <div className="block p-4 rounded-lg shadow-lg bg-white border-gray-200 hover:bg-gray-100 my-2">
      {/* <div className="flex flex-row justify-between">
        <div className=" text-xl font-bold ">
          PM Gateway {props.pmGateway.Line + "-" + props.pmGateway.Code}
        </div>
        <div className="flex flex-row justify-center">
          <div className="flex justify-center items-center px-1">
            <button>
              <svg
                viewBox="0 0 20 20"
                className={`w-5 h-5  fill-current ${
                  isConnected ? "text-green-600" : "text-red-600"
                }`}
              >
                <path d="M5.25 3A2.25 2.25 0 003 5.25v9.5A2.25 2.25 0 005.25 17h9.5A2.25 2.25 0 0017 14.75v-9.5A2.25 2.25 0 0014.75 3h-9.5z" />
              </svg>
            </button>
          </div>
          <div className="flex justify-center items-center px-1">
            <button onClick={deletePMGatewayHandler}>
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
        </div>
      </div>
      <div className="text-xs px-2 text-slate-400">
        target : {props.pmGateway.TargetUrl}
      </div>

      <hr class="my-2 h-px bg-gray-200 border-0 dark:bg-gray-700"></hr> */}

      <CardHeader
        title={
          "PM Gateway " + props.pmGateway.Line + "-" + props.pmGateway.Code
        }
        host={props.pmGateway.TargetUrl}
        isConnected={isConnected}
        onClickDelete={deletePMGatewayHandler}
        onClickConnection={handleConnectionChange}
      />

      <div className="px-2">
        {pmGateway.PMs.map((PM, i) => {
          return <PMCard pm={PM} />;
        })}
      </div>

      <hr class="my-2 h-px bg-gray-200 border-0 dark:bg-gray-700"></hr>

      {props.pmGateway.PMs.length < 8 && (
        <CreatePMCard pmGatewayId={pmGateway.Id} />
      )}
    </div>
  );
}

export default PMGatewayCard;
