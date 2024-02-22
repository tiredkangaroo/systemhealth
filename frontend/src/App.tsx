import { useEffect, useState } from "react";
import "./App.css";

interface SystemHealth {
  CPUTemp: number;
  BatteryTemp: number;
  BatteryStatus: string;
  BatteryCapacity: number;
  CPUUtilization: number;
  MemoryUsage: string;
  StorageUsage: string;
}
function App() {
  const [health, setHealth] = useState<SystemHealth | null | undefined>(
    undefined,
  );
  useEffect(() => {
    async function fetchData() {
      let res;
      try {
        res = await fetch("/api/get");
      } catch (_) {
        return setHealth(null);
      }
      const data = (await res.json()) as SystemHealth;
      setHealth(data);
    }
    fetchData();
  });
  if (health === undefined) {
    return <h1> loading </h1>;
  }
  if (health === null) {
    return <h1> failed to recieve system health info </h1>;
  }
  return (
    <>
      <div className="cpu component">
        <ul>
          <li> CPU Temp: {health.CPUTemp} </li>
          <li> CPU Utilization: {health.CPUUtilization}% </li>
        </ul>
      </div>
      <div className="battery component">
        <ul>
          <li> Battery Temp: {health.BatteryTemp} </li>
          <li> Battery Status: {health.BatteryStatus} </li>
          <li> Battery Capacity: {health.BatteryCapacity} </li>
        </ul>
      </div>
      <div className="memory component">
        <ul>
          <li> Memory Utilization: {health.MemoryUsage} </li>
        </ul>
      </div>
      <div className="storage component">
        <ul>
          <li> Storage Utilization: {health.StorageUsage} </li>
        </ul>
      </div>
    </>
  );
}

export default App;
