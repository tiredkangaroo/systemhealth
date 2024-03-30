import { useEffect, useState } from "react";
import "./App.css";

interface ServiceHealth {
  name: string;
  status: string | null;
}
interface SystemHealth {
  cpu_temp: number;
  battery_temp: number;
  battery_status: string;
  battery_capacity: number;
  cpu_utilization: number;
  memory_usage: string;
  storage_usage: string;
  services: Array<ServiceHealth>;
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
      console.log("Set Health Data to");
      console.log(data);
    }
    fetchData();
  }, []);
  if (health === undefined) {
    return <h1> loading </h1>;
  }
  if (health === null) {
    return <h1> failed to recieve system health info </h1>;
  }
  return (
    <>
      <div className="cpu component">
        <h1> cpu </h1>
        <ul>
          <li> CPU Temp: {health.cpu_temp}° F </li>
          <li> CPU Utilization: {health.cpu_utilization.toFixed(1)}% </li>
        </ul>
      </div>
      <div className="battery component">
        <h1> battery </h1>
        <ul>
          <li> Battery Temp: {health.battery_temp}° F </li>
          <li> Battery Status: {health.battery_status} </li>
          <li> Battery Percentage: {health.battery_capacity}% </li>
        </ul>
      </div>
      <div className="memory component">
        <h1> memory </h1>
        <ul>
          <li> Memory Utilization: {health.memory_usage} </li>
        </ul>
      </div>
      <div className="storage component">
        <h1> storage </h1>
        <ul>
          <li> Storage Utilization: {health.storage_usage} </li>
        </ul>
      </div>
      <div className="service component">
        <h1> service status </h1>
        <ul>
          {health.services.map((service) => {
            return (
              <li className={`service-${service.status?.replace(" ", "")}`}>
                {service.name}
              </li>
            );
          })}
        </ul>
      </div>
    </>
  );
}

export default App;
