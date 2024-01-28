package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/kleczynski/go-microservices-k8s/details"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthHanler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Health Handler returned wrong code status: got %v want: %v", status, http.StatusOK)
	}
	expected := HealthResponse{
		Status:    "UP",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Error("Error marshaling expected response")
	}

	if rr.Body.String() != string(expectedBytes) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RootHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Root Handler returned wrong code status: got %v want: %v", rr.Code, http.StatusOK)
	}
	exptected := "Application is up and running"
	if rr.Body.String() != exptected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), exptected)
	}
}

func TestDetailsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/details", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DetailsHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Details Handler returned wrong code status: got %v want: %v", rr.Code, http.StatusOK)
	}
	hostname, err := details.GetHostName()
	if err != nil {
		t.Errorf("GetHostname returned error: got %v", err)
	}
	ip, err := details.GetIP()
	if err != nil {
		t.Errorf("GetIP returned error: got %v", err)
	}
	operationSystem := details.GetOperatingSystem()
	CPUCount := details.GetCPUCount()
	strCPUCount := strconv.Itoa(CPUCount)
	expected := DetailsResponse{
		Hostname:        hostname,
		IP:              ip.String(),
		OperationSystem: operationSystem,
		CPUCount:        strCPUCount,
	}
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Error("Error marshaling expected response")
	}
	if rr.Body.String() != string(expectedBytes) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
