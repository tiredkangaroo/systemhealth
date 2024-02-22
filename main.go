package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type SystemHealth struct {
	CPUTemp         float64 `json:"cpu_temp"`
	BatteryTemp     float64 `json:"battery_temp"`
	BatteryStatus   string  `json:"battery_status"`
	BatteryCapacity float64 `json:"battery_capacity"`
	CPUUtilization  float64 `json:"cpu_utilization"`
	MemoryUsage     string  `json:"memory_usage"`
	StorageUsage    string  `json:"storage_usage"`
}

func getCPUTemp() (float64, error) {
	tempFile := "/sys/class/thermal/thermal_zone0/temp"
	data, err := os.ReadFile(tempFile)
	if err != nil {
		return 0, err
	}
	tempStr := strings.TrimSpace(string(data))
	tempCelsius, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, err
	}
	tempFahrenheit := (tempCelsius * 9.0 / 5.0) + 32.0
	return tempFahrenheit, nil
}

func getBatteryTemp() (float64, error) {
	tempFile := "/sys/class/power_supply/BAT0/temp"
	data, err := os.ReadFile(tempFile)
	if err != nil {
		return 0, err
	}
	tempStr := strings.TrimSpace(string(data))
	tempCelsius, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, err
	}
	tempFahrenheit := (tempCelsius * 9.0 / 5.0) + 32.0
	return tempFahrenheit, nil
}

func getBatteryStatus() (string, error) {
	statusFile := "/sys/class/power_supply/BAT0/status"
	data, err := os.ReadFile(statusFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func getBatteryCapacity() (float64, error) {
	capacityFile := "/sys/class/power_supply/BAT0/capacity"
	data, err := os.ReadFile(capacityFile)
	if err != nil {
		return 0, err
	}
	capacityStr := strings.TrimSpace(string(data))
	capacity, err := strconv.ParseFloat(capacityStr, 64)
	if err != nil {
		return 0, err
	}
	return capacity, nil
}

func getCPUUtilization() (float64, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(data), "\n")
	cpuLine := lines[0] // First line contains overall CPU statistics
	fields := strings.Fields(cpuLine)
	totalTime := 0.0
	for _, field := range fields[1:] {
		time, err := strconv.ParseFloat(field, 64)
		if err != nil {
			return 0, err
		}
		totalTime += time
	}
	idleTime, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return 0, err
	}
	utilization := 100 * (1 - idleTime/totalTime)
	return utilization, nil
}

func getMemoryUsage() (string, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	var totalMem, freeMem int64

	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			totalMem, err = strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				return "", err
			}
		} else if strings.HasPrefix(line, "MemFree:") {
			fields := strings.Fields(line)
			freeMem, err = strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				return "", err
			}
		}
	}

	usedMem := totalMem - freeMem
	usedMemGB := float64(usedMem) / 1024 / 1024
	totalMemGB := float64(totalMem) / 1024 / 1024
	return fmt.Sprintf("%.2fGB/%.2fGB", usedMemGB, totalMemGB), nil
}

func getStorageUsage() (string, error) {
	cmd := exec.Command("df", "-BG") // Run df command to get disk space usage in GB
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/dev/") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				usedSpaceGB, err := strconv.ParseFloat(strings.TrimSuffix(fields[2], "G"), 64)
				if err != nil {
					return "", err
				}
				totalSpaceGB, err := strconv.ParseFloat(strings.TrimSuffix(fields[1], "G"), 64)
				if err != nil {
					return "", err
				}
				return fmt.Sprintf("%.2fGB/%.2fGB", usedSpaceGB, totalSpaceGB), nil
			}
		}
	}
	return "", fmt.Errorf("unable to get storage usage")
}

func get(w http.ResponseWriter, req *http.Request) {
	cpuTemp, e := getCPUTemp()
	batteryTemp, er := getBatteryTemp()
	batteryStatus, err := getBatteryStatus()
	batteryCapacity, erro := getBatteryCapacity()
	cpuUtilization, erorr := getCPUUtilization()
	memoryUsage, erorrd := getMemoryUsage()
	storageUsage, errordz := getStorageUsage()
	if e != nil {
		fmt.Println(e.Error())
		cpuTemp = -1
	}
	if er != nil {
		fmt.Println(er.Error())
		batteryTemp = -1
	}
	if err != nil {
		fmt.Println(err.Error())
		batteryStatus = "Failed to access battery status."
	}
	if erro != nil {
		fmt.Println(erro.Error())
		batteryCapacity = -1
	}
	if erorr != nil {
		fmt.Println(erorr.Error())
		cpuUtilization = -1
	}
	if erorrd != nil {
		fmt.Println(erorrd.Error())
		memoryUsage = "Failed to access memory usage."
	}
	if errordz != nil {
		fmt.Println(errordz.Error())
		storageUsage = "Failed to access storage usage."
	}
	systemHealth := SystemHealth{
		CPUTemp:         cpuTemp,
		BatteryTemp:     batteryTemp,
		BatteryStatus:   batteryStatus,
		BatteryCapacity: batteryCapacity,
		CPUUtilization:  cpuUtilization,
		MemoryUsage:     memoryUsage,
		StorageUsage:    storageUsage,
	}

	jsonData, err := json.Marshal(systemHealth)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
func main() {
	http.HandleFunc("/api/get", get)
	http.Handle("/", http.FileServer(http.Dir("/home/nikhilkumar/systemhealth/dist")))
	http.ListenAndServe(":8106", nil)
}
