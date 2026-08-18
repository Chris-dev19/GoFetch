package main

// Needed Imports

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/disk"
)

// Core Function
func main() {
	operating_System := GetOSRelease()
	uptime, _ := GetUptime()
	cpu := GetCPU()
	totalKB, availableKB := GetMem()
	usedKB := totalKB - availableKB
	totalGB := float64(totalKB) / 1024 / 1024
	usedGB := float64(usedKB) / 1024 / 1024
	gpu := GetGPU()
	shell := GetShell()
	battery_level := GetBatteryLevel()
	battery_status := GetBatteryStatus()
	lang := GetLocale()
	isgoInstalled := IsGoInstalled()
	term := GetTerm()

	// 1. Initialize your info lines slice with the standard stats
	infoLines := []string{
		fmt.Sprintf("OS: %s", operating_System),
		fmt.Sprintf("Uptime: %s", uptime),
		fmt.Sprintf("CPU: %s", cpu),
		fmt.Sprintf("RAM: %.1f/%.1fGB", usedGB, totalGB),
	}

	// 2. Fetch all disks and dynamically append each one to infoLines
	disks, err := GetDisks()
	if err == nil {
		for _, d := range disks {
			// Formats each drive nicely (e.g., "Disk [/]: 45.2/128.0GB (35.3%)")
			diskLine := fmt.Sprintf("Disk [%s]: %.1f/%.1fGB (%.1f%%)",
				d.Mountpoint, d.UsedGB, d.TotalGB, d.UsedPercent)
			infoLines = append(infoLines, diskLine)
		}
	} else {
		infoLines = append(infoLines, "Disk: Unknown")
	}

	// 3. Append the remaining system items
	infoLines = append(infoLines,
		fmt.Sprintf("GPU: %s", gpu),
		fmt.Sprintf("Shell: %s", shell),
		fmt.Sprintf("Battery Level: %s%%", battery_level),
		fmt.Sprintf("Status: %s", battery_status),
		fmt.Sprintf("Lang: %s", lang),
		fmt.Sprintf("Go Installed: %v", isgoInstalled),
		fmt.Sprintf("Terminal: %s", term),
	)

	// 4. Print everything out
	PrintLogo(infoLines)
}

func GetOSRelease() string {
	// OS Variable declaration

	osname, err := os.ReadFile("/etc/os-release")
	// Error handling

	if err != nil {
		return "Unknown"
	}
	// Triming OS name variables

	osname_new_lines := strings.Split(string(osname), "\n")

	for _, osname_new_line := range osname_new_lines {
		if strings.HasPrefix(string(osname_new_line), "PRETTY_NAME=") {
			return strings.Trim(strings.TrimPrefix(string(osname_new_line), "PRETTY_NAME="), `"`)
		}
	}
	return "Unknown"
}

// GetUptime function declaration

func GetUptime() (time.Duration, error) {
	// Uptime variables declarations

	timeup, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	// Parsing timeup to Float64

	timeup_improved := strings.Fields(string(timeup))
	seconds, err := strconv.ParseFloat(timeup_improved[0], 64)
	if err != nil {
		return 0, err
	}
	// Return Result

	return time.Duration(seconds) * time.Second, nil
}

// GetCPU function declarations

func GetCPU() string {
	// CpuiInfo variable declaration

	cpuinfo, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "Unknown"
	}
	// Triming Cpuinfos

	cpu_infos_lines := strings.Split(string(cpuinfo), "\n")

	for _, cpu_infos_line := range cpu_infos_lines {
		if strings.HasPrefix(cpu_infos_line, "model name") {
			model_parts := strings.SplitN(cpu_infos_line, ":", 2)
			return strings.TrimSpace(model_parts[1])
		}
	}
	// Final Return

	return "Unknown"
}

// GetMem function declaration

func GetMem() (uint64, uint64) {
	// Memory variable declaration

	memory, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0
	}
	// Parsing memory

	memory_lines := strings.Split(string(memory), "\n")
	var TotalKB, availableKB uint64

	for _, memory_line := range memory_lines {
		mem_fields := strings.Fields(string(memory_line))
		if len(mem_fields) < 2 {
			continue
		}
		mem_value, err := strconv.ParseUint(mem_fields[1], 10, 64)
		if err != nil {
			continue
		}
		switch mem_fields[0] {
		case "MemTotal:":
			TotalKB = mem_value
		case "MemAvailable:":
			availableKB = mem_value
		}
	}
	// Final return

	return TotalKB, availableKB
}

// GPU function declaration

func GetGPU() string {
	// GPU variable declaration

	gpu, err := exec.Command("lspci").Output()
	if err != nil {
		return "Unknown"
	}
	// Searching for VGA Options...

	gpu_lines := strings.Split(string(gpu), "\n")

	for _, gpu_line := range gpu_lines {
		if strings.Contains(strings.ToLower(gpu_line), "vga") {

			gpu_parts := strings.SplitN(gpu_line, ": ", 2)
			if len(gpu_parts) == 2 {
				gpu_line := gpu_parts[1]
				// FIXED: Changed logic from idx != 1 to check for presence via != -1
				if idx := strings.Index(gpu_line, "(rev"); idx != -1 {
					gpu_line = gpu_line[:idx]
					return strings.TrimSpace(gpu_line)
				}
				return strings.TrimSpace(gpu_line)
			}
		}
	}
	// Final Return

	return "Unknown"
}

// GetShell function declaration

func GetShell() string {
	// Shell Path variable declaration

	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		return "Unknown"
	}
	// Final Return
	return filepath.Base(shellPath)
}

// GetBatteryLevel function declaration

func GetBatteryLevel() string {
	// Battery Level variable declaration

	battery_level, err := os.ReadFile("/sys/class/power_supply/BAT0/capacity")
	if err != nil {
		return "Unknown"
	}

	battery_level_lines := strings.Split(string(battery_level), "\n")
	// Final Return

	return battery_level_lines[0]
}

// GetBatteryStatus function declaration

func GetBatteryStatus() string {
	// Battery Status variable declaration

	battery_status, err := os.ReadFile("/sys/class/power_supply/BAT0/status")
	if err != nil {
		return "Unknown"
	}

	battery_status_lines := strings.Split(string(battery_status), "\n")
	// Final Return

	return battery_status_lines[0]
}

// GetLocale function declaration

func GetLocale() string {
	// Locale variable declaration

	locale, err := exec.Command("locale").Output()
	if err != nil {
		return "Unknown"
	}

	locale_lines := strings.Split(string(locale), "\n")
	for _, locale_line := range locale_lines {
		if strings.HasPrefix(locale_line, "LANG=") {
			return strings.Trim(strings.TrimPrefix(locale_line, "LANG="), `"`)
		}
	}
	// Final Return

	return "Unknown"
}

// IsGoInstalled function declaration

func IsGoInstalled() string {
	_, err := exec.LookPath("go")
	if err == nil {
		return "Yes"
	}
	return "No"
}

// GetTerm function declaration

func GetTerm() string {
	// Term variable declaration
	term := os.Getenv("TERM")
	if term == "" {
		return "Unknown"
	}
	return term
}

type DiskInfo struct {
	Mountpoint  string
	TotalGB     float64
	UsedGB      float64
	FreeGB      float64
	UsedPercent float64
}

func GetDisks() ([]DiskInfo, error) {
	const bytesInGB = 1024 * 1024 * 1024
	partitions, err := disk.Partitions(false)

	if err != nil {
		return nil, err
	}

	var diskList []DiskInfo

	for _, partition := range partitions {
		isReadOnly := false
		for _, opt := range strings.Split(partition.Opts, ",") {
			if opt == "ro" {
				isReadOnly = true
				break
			}
		}
		if isReadOnly {
			continue
		}

		if strings.HasPrefix(partition.Mountpoint, "/nix/store") {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		totalGB := float64(usage.Total) / bytesInGB
		usedGB := float64(usage.Used) / bytesInGB
		freeGB := float64(usage.Free) / bytesInGB

		if totalGB < 2.0 {
			continue
		}

		diskList = append(diskList, DiskInfo{
			Mountpoint:  partition.Mountpoint,
			TotalGB:     totalGB,
			UsedGB:      usedGB,
			FreeGB:      freeGB,
			UsedPercent: usage.UsedPercent,
		})
	}

	return diskList, nil
}
