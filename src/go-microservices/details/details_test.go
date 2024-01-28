package details

import "testing"

func TestGetHostname(t *testing.T) {
	hostname, err := GetHostName()
	if err != nil {
		t.Errorf("GetHostname returned an unexpected error: %v", err)
	}
	if hostname == "" {
		t.Errorf("GetHostname returned an empty hostname %v", hostname)
	}
}

func TestGetIP(t *testing.T) {
	ip, err := GetIP()
	if err != nil {
		t.Errorf("GetIP returned an unexpected error: %v", err)
	}
	if ip.String() == "" {
		t.Errorf("GetIP returned empty IP address: %v", ip.String())
	}
}

func TestGetOperationSystem(t *testing.T) {
	operationSystem := GetOperatingSystem()
	if operationSystem == "" {
		t.Errorf("GetOparationSystem returned an unexpected empty value: %v", operationSystem)
	}
}

func TestGetCPUCount(t *testing.T) {
	CPUCount := GetCPUCount()
	if CPUCount <= 0 {
		t.Errorf("GetCPUCount has to return positive which is at least 1 got: %v", CPUCount)
	}
}
