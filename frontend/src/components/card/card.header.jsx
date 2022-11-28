import ConnectionButton from "../button/connection.button";
import DeleteButton from "../button/delete.button";

function CardHeader({
  title,
  host,
  isConnected,
  onClickConnection,
  onClickDelete,
}) {
  return (
    <div className="flex flex-col">
      <div className="flex flex-row justify-between">
        <div className=" text-xl font-bold ">{title}</div>
        <div className="flex flex-row justify-center">
          <div className="flex justify-center items-center px-1">
            <ConnectionButton
              isConnected={isConnected}
              onClick={onClickConnection}
            />
          </div>
          <div className="flex justify-center items-center px-1">
            <DeleteButton onClick={onClickDelete} />
          </div>
        </div>
      </div>
      <div className="text-xs px-2 text-slate-400">host : {host}</div>
      <hr class="my-2 h-px bg-gray-200 border-0 dark:bg-gray-700"></hr>
    </div>
  );
}

export default CardHeader;
