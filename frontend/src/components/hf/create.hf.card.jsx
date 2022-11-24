import { useState } from "react";
import Switch from "@mui/material/Switch";
import { CreateReader } from "../../../wailsjs/go/hfreader/App";

function CreateHFCard() {
  const [line, setLine] = useState("");
  const [post, setPost] = useState("");
  const [code, setCode] = useState("");
  const [targetUrl, setTargetUrl] = useState("");

  const updateLine = (e) => setLine(e.target.value);
  const updatePost = (e) => setPost(e.target.value);
  const updateCode = (e) => setCode(e.target.value);
  const updateTargetUrl = (e) => setTargetUrl(e.target.value);

  const submitNewDevice = (e) => {
    console.log("Submit New Device");
    e.preventDefault();
    CreateReader(line, post, code, targetUrl);
  };

  return (
    <div className="block p-4 rounded-lg shadow-lg bg-white border-gray-200 hover:bg-gray-100 my-2">
      <div id="hf" className="flex flex-col py-2">
        <div className="flex flex-col justify-center py-2 text-xl font-bold ">
          Create New Reader
        </div>
        <div className="flex flex-col justify-center py-2">
          <input
            id="line"
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
            onChange={updateLine}
            autoComplete="off"
            placeholder="Line"
            name="input"
            type="text"
          />
        </div>
        <div className="flex flex-col justify-center py-2">
          <input
            id="post"
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
            onChange={updatePost}
            autoComplete="off"
            placeholder="Post"
            name="input"
            type="text"
          />
        </div>
        <div className="flex flex-col justify-center py-2">
          <input
            id="code"
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
            onChange={updateCode}
            autoComplete="off"
            placeholder="Code"
            name="input"
            type="text"
          />
        </div>
        <div className="flex flex-col justify-center py-2">
          <input
            id="targetUrl"
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
            onChange={updateTargetUrl}
            autoComplete="off"
            placeholder="targetUrl"
            name="input"
            type="text"
          />
        </div>
        <div className="py-3">
          <button
            type="button"
            class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 mr-2 mb-2 w-full dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800"
            onClick={submitNewDevice}
          >
            Add
          </button>
        </div>
      </div>
    </div>
  );
}

export default CreateHFCard;
