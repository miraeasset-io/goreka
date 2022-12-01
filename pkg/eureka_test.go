package goreka

import "testing"

func TestEurekaRegistration(t *testing.T) {
	form := RegistrationForm{
		ServiceName: "test-service",
		ServiceHost: "localhost",
		ServicePort: 8080,
		InstanceId:  "test-service-1",
		ServiceUrl:  "https://admin:M1rAeA553t2910@dev-jhipster.miraeasset.io/eureka/apps/",
	}

	err := form.RegisterService()
	if err != nil {
		t.Error(err)
	}
}

func TestEurekaHeartbeat(t *testing.T) {
	form := RegistrationForm{
		ServiceName: "test-service",
		ServiceHost: "localhost",
		ServicePort: 8080,
		InstanceId:  "test-service-1",
		ServiceUrl:  "https://admin:M1rAeA553t2910@dev-jhipster.miraeasset.io/eureka/apps/",
	}

	err := form.RegisterService()
	if err != nil {
		t.Error(err)
	}

	err = form.Heartbeat()
	if err != nil {
		t.Error(err)
	}
}
