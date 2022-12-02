package goreka

import (
	"os"
	"testing"
)

func TestEurekaRegistration(t *testing.T) {
	username := os.Getenv("EUREKA_USERNAME")
	password := os.Getenv("EUREKA_PASSWORD")

	form := RegistrationForm{
		ServiceName: "test-service",
		ServiceHost: "localhost",
		ServicePort: 8080,
		EurekaUrl:   "https://" + username + ":" + password + "@dev-jhipster.miraeasset.io/eureka/apps/",
	}

	err := form.RegisterService()
	if err != nil {
		t.Error(err)
	}
}

func TestEurekaHeartbeat(t *testing.T) {
	username := os.Getenv("EUREKA_USERNAME")
	password := os.Getenv("EUREKA_PASSWORD")
	form := RegistrationForm{
		ServiceName: "test-service",
		ServiceHost: "localhost",
		ServicePort: 8080,
		EurekaUrl:   "https://" + username + ":" + password + "@dev-jhipster.miraeasset.io/eureka/apps/",
	}

	err := form.RegisterService()
	if err != nil {
		t.Error(err)
	}

	resp, err := form.Heartbeat()
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Error("Heartbeat failed")
	}
}
