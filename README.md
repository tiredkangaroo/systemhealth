# System Health
specifically for laptop machines running the linux kernel <5.

#### This app tracks
- CPU temperature*
- CPU utilization
- Battery Temperature*
- Battery Status
- Battery Capacity**
- Memory Utilization
- Storage Utilization

#### Running the Application
Requires <a href="https://go.dev">go</a> 1.22.

Build: `go build`

The command generates a bytecode file that can be run in the terminal.

The application runs on port 8106.

#### Routes
`/` -> The frontend page for system health.

`/api/get` -> The API (the frontend relies on this).

#### API Response Structure
```
{
  "cpu_temp": float,
  "battery_temp": float,
  "battery_status": "Full" | "Discharging" | "Charging",
  "battery_capacity": (0-100),
  "cpu_utilization": float (0-100),
  "memory_usage": floatGB/floatGB,
  "storage_usage": floatGB/floatGB,
  "services": [
    {
      "name": service_name,
      "Status": service_status
    },
  ]
}
```
<img width="1510" alt="Screenshot 2024-04-24 at 2 55 33 PM" src="https://github.com/tiredkangaroo/systemhealth/assets/81335306/bad7ae2e-d07c-480a-a381-08110013b8de">

<b>*Side note</b>: The CPU and Battery temperature are being fetched from their respective files, however on the Macbook Pro 2017 running Ubuntu Server, both files seem to always return the same values.

<b>**Side note 2</b>: The Battery Percentage listed is your Battery Capacity. It is the capacity your battery is holding now. If the battery has wear and tear, the status may be `Full` but the percentage will not be 100. Batteries naturally are unable to keep as much battery as they kept when they were new.
