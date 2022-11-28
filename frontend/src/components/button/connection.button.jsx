function ConnectionButton({ onClick, isConnected }) {
  return (
    <button>
      <svg
        viewBox="0 0 20 20"
        className={`w-5 h-5  fill-current ${
          isConnected ? "text-green-600" : "text-red-600"
        }`}
        onClick={onClick}
      >
        <path d="M5.25 3A2.25 2.25 0 003 5.25v9.5A2.25 2.25 0 005.25 17h9.5A2.25 2.25 0 0017 14.75v-9.5A2.25 2.25 0 0014.75 3h-9.5z" />
      </svg>
    </button>
  );
}

export default ConnectionButton;
