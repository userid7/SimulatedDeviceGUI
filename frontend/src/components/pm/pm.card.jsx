import { useState, useEffect, useRef } from "react";
import { _ } from "lodash";
import Switch from "@mui/material/Switch";
import {
  DeletePM,
  SetPM,
  SetPMIsOk,
  SetPMIsRandom,
  SetPMKw,
} from "../../../wailsjs/go/pm/App";

function PMCard(props) {
  const [pm, setPm] = useState(props.pm);

  const deletePMHandler = (e) => {
    console.log("Delete PM");
    e.preventDefault();
    DeletePM(props.pm.Id);
  };

  useEffect(() => {
    console.log("props.pm");
    console.log(props.pm);
    setPm(props.pm);
  }, [JSON.stringify(props.pm)]);

  const updateKWValue = (e) => {
    var newPm = pm;
    var kw = parseFloat(e.target.value);
    newPm.Kw = kw;
    setPm(newPm);
  };

  const handleIsOkChange = (e) => {
    console.log("IsOk Change");
    var newPm = pm;
    var isOk = e.target.checked;
    newPm.IsOk = isOk;
    setPm(newPm);
    SetPMIsOk(newPm.Id, newPm.IsOk);
  };

  const handleIsRandomChange = (e) => {
    console.log("IsRandom Change");

    var newPm = pm;
    var isRandom = e.target.checked;
    console.log(isRandom);
    newPm.IsRandom = isRandom;
    setPm(newPm);
    SetPMIsRandom(newPm.Id, newPm.IsRandom);
  };

  useEffect(() => {
    const delayDebounceFn = setTimeout(() => {
      console.log(pm.Kw);
      SetPMKw(pm.Id, pm.Kw);
    }, 1000);

    return () => clearTimeout(delayDebounceFn);
  }, [pm.Kw]);

  return (
    <div className="my-2">
      <div className="flex flex-row justify-between items-center">
        <div className=" text-md font-bold ">
          PM {pm.Post + "-" + pm.Code}
          {/* {"TW" + "-" + "JIG1"} */}
        </div>
        <div className="flex flex-row">
          <div className="flex justify-center items-center px-1">
            <Switch checked={pm.IsOk} onChange={handleIsOkChange} />
          </div>
          <div className="flex justify-center items-center px-1">
            <button onClick={deletePMHandler}>
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
      <div className="flex flex-col justify-between items-start">
        <table class="w-full table-auto">
          <thead className="text-sm">
            <tr>
              <th>Type</th>
              <th>Value</th>
              <th></th>
              <th>Random</th>
            </tr>
          </thead>
          <tbody className="text-center align-center text-sm">
            {/* <tr>
              <td>V</td>
              <td>500</td>
              <td>
                <div className="">
                  <input
                    id="default-range"
                    type="range"
                    class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
                  ></input>
                </div>
              </td>
              <td>
                <div className="h-full justify-center items-center align-center">
                  <label class="inline-flex relative items-center cursor-pointer">
                    <input
                      type="checkbox"
                      value=""
                      class="sr-only peer"
                    ></input>
                    <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
                  </label>
                </div>
              </td>
            </tr>
            <tr>
              <td>I</td>
              <td>500</td>
              <td>
                <input
                  id="default-range"
                  type="range"
                  class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
                ></input>
              </td>
              <td>
                <label class="inline-flex relative items-center cursor-pointer">
                  <input type="checkbox" value="" class="sr-only peer"></input>
                  <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
                </label>
              </td>
            </tr> */}
            <tr>
              <td>kw</td>
              <td>{pm.Kw}</td>
              <td>
                <input
                  id="default-range"
                  type="range"
                  class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
                  onChange={updateKWValue}
                  defaultValue={pm.Kw}
                ></input>
              </td>
              <td>
                <label class="inline-flex relative items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={pm.IsRandom}
                    onChange={handleIsRandomChange}
                    class="sr-only peer"
                  ></input>
                  <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
                </label>
              </td>
            </tr>
            <tr>
              <td>kwh</td>
              <td>1000</td>
            </tr>
          </tbody>
        </table>
        {/* <div className="flex flex-row justify-center items-center px-4 gap-x-5">
          <label class="block text-sm font-medium text-gray-900 dark:text-white">
            kw
          </label>
          <input
            type="text"
            id="kwh"
            class="w-12 bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            placeholder="0"
            required
          ></input>
          <input
            id="default-range"
            type="range"
            class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
          ></input>
          <label class="inline-flex relative items-center cursor-pointer">
            <input type="checkbox" value="" class="sr-only peer"></input>
            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
          </label>
        </div>
        <div className="flex flex-row justify-center items-center px-4 gap-x-5">
          <label class="block text-sm font-medium text-gray-900 dark:text-white">
            kwh
          </label>
          <input
            type="text"
            id="kwh"
            class="w-12 bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            placeholder="0"
            required
          ></input>
        </div> */}
      </div>
    </div>
  );
}

export default PMCard;
