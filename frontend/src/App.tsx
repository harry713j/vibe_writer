import { useEffect, useState } from "react";

function App() {
  const [status, setStatus] = useState<string>("");

  useEffect(() => {
    const getHealth = async () => {
      const res = await fetch(import.meta.env.VITE_API_URL + "/api/v1/health");

      const data = await res.json();
      setStatus(data.message);
    };

    getHealth();
  }, []);

  console.log();
  return (
    <div>
      <h1>Vibewriter - {status}</h1>
    </div>
  );
}

export default App;
