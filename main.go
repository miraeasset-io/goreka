package main

import (
	goreka "goreka/pkg"
)

func main() {
	goreka.Init(goreka.RegistrationForm{
		ServiceName: "test-service",
		ServiceHost: "localhost",
		ServicePort: 8080,
		InstanceId:  "test-service-1",
		ServiceUrl:  "https://admin:M1rAeA553t2910@dev-jhipster.miraeasset.io/eureka/apps/",
	})
}
