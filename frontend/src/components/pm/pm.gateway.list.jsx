import { useState, useEffect } from "react";
import { GetAllActivePMGateway } from "../../../wailsjs/go/pm/App";
import PMGatewayCard from "./pm.gateway.card";
import CreatePMCard from "./create.pm.gateway.card";

function PMGatewayList() {
  const [pmGateways, setPMGateways] = useState([]);

  useEffect(() => {
    const interval = setInterval(async () => {
      console.log("Interval Pooling");
      var newPMGateways = await GetAllActivePMGateway();
      console.log(newPMGateways);
      setPMGateways(newPMGateways || []);
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="flex flex-col w-96">
      {pmGateways.map((pmGateway, i) => {
        return <PMGatewayCard pmGateway={pmGateway} />;
      })}
      <CreatePMCard />
    </div>
  );
}

export default PMGatewayList;
