package details

import (
	"log"
	"net"
	"os"
	"runtime"
)

func GetHostName() (string, error) {
	hostname, err := os.Hostname()
	return hostname, err
}

func GetIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, err
}

func GetOperatingSystem() string {
	return runtime.GOOS
}

func GetCPUCount() int {
	return runtime.NumCPU()
}

func GetMemoryUsage() runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem
}
