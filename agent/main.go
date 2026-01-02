package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var startTime = time.Now()

// -------- HEALTH --------
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// -------- MEMÓRIA (/proc/meminfo) --------
func getMemoryUsage() (total uint64, available uint64) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			total, _ = strconv.ParseUint(fields[1], 10, 64)
		case "MemAvailable:":
			available, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}
	return
}

// -------- CPU (/proc/stat) --------
type cpuTimes struct {
	idle  uint64
	total uint64
}

func readCPUTimes() cpuTimes {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return cpuTimes{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		var idle, total uint64
		for i := 1; i < len(fields); i++ {
			v, _ := strconv.ParseUint(fields[i], 10, 64)
			total += v
			if i == 4 || i == 5 { // idle + iowait
				idle += v
			}
		}
		return cpuTimes{idle: idle, total: total}
	}
	return cpuTimes{}
}

func getCPUUsagePercent() float64 {
	t1 := readCPUTimes()
	time.Sleep(500 * time.Millisecond)
	t2 := readCPUTimes()

	idleDelta := float64(t2.idle - t1.idle)
	totalDelta := float64(t2.total - t1.total)
	if totalDelta == 0 {
		return 0
	}
	return 100.0 * (1.0 - idleDelta/totalDelta)
}

// -------- METRICS --------
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime).Seconds()

	totalMem, availMem := getMemoryUsage()
	usedMem := totalMem - availMem
	cpuUsage := getCPUUsagePercent()

	fmt.Fprintf(w, "agent_uptime_seconds %.0f\n", uptime)
	fmt.Fprintf(w, "go_goroutines %d\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "process_pid %d\n", os.Getpid())

	fmt.Fprintf(w, "node_memory_total_kb %d\n", totalMem)
	fmt.Fprintf(w, "node_memory_used_kb %d\n", usedMem)
	fmt.Fprintf(w, "node_memory_available_kb %d\n", availMem)

	fmt.Fprintf(w, "node_cpu_usage_percent %.2f\n", cpuUsage)
}

// -------- PORT (ENV) --------
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9200"
	}
	return ":" + port
}

// -------- MAIN --------
func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/metrics", metricsHandler)

	addr := getPort()
	log.Println("linux-monitoring-agent listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
