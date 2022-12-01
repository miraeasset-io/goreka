package goreka

import (
	"fmt"
	"goreka/tools"
	"strings"
	"time"
)

type AppRegistrationBody struct {
	Instance InstanceDetails `json:"instance"`
}

type InstanceDetails struct {
	InstanceId       string         `json:"instanceId"`
	HostName         string         `json:"hostName"`
	App              string         `json:"app"`
	VipAddress       string         `json:"vipAddress"`
	SecureVipAddress string         `json:"secureVipAddress"`
	IpAddr           string         `json:"ipAddr"`
	Status           string         `json:"status"`
	Port             Port           `json:"port"`
	SecurePort       Port           `json:"securePort"`
	HealthCheckUrl   string         `json:"healthCheckUrl"`
	StatusPageUrl    string         `json:"statusPageUrl"`
	HomePageUrl      string         `json:"homePageUrl"`
	DataCenterInfo   DataCenterInfo `json:"dataCenterInfo"`
	Metadata         Metadata       `json:"metadata"`
}

type Port struct {
	Port    string `json:"$"`
	Enabled string `json:"@enabled"`
}

type DataCenterInfo struct {
	Class string `json:"@class"`
	Name  string `json:"name"`
}

type Metadata struct {
	Zone    string  `json:"zone"`
	Profile string  `json:"profile"`
	Port    int     `json:"management.port"`
	Version float32 `json:"version"`
}

type RegistrationForm struct {
	ServiceName string
	ServiceHost string
	ServicePort int
	InstanceId  string
	ServiceUrl  string
}

func (form RegistrationForm) RegisterService() error {
	fmt.Println("Registering service with status: STARTING")
	body := ConstructRegistrationBody(form, "STARTING")

	serviceName := strings.ToUpper(form.ServiceName)

	postUrl := form.ServiceUrl + serviceName
	fmt.Println(postUrl)

	_, err := tools.HttpPostReq(postUrl, body, nil)
	if err != nil {
		return err
	}

	fmt.Println("Updating the status to: UP")
	bodyUP := ConstructRegistrationBody(form, "UP")

	_, err = tools.HttpPostReq(form.ServiceUrl+serviceName, bodyUP, nil)
	if err != nil {
		return err
	}

	return nil
}

func (form RegistrationForm) UnRegisterEurekaService() {
	fmt.Println("UnRegistering service from eureka ...")
	res, err := tools.HttpPostReq(form.ServiceUrl, nil, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func (form RegistrationForm) Heartbeat() error {
	serviceName := strings.ToUpper(form.ServiceName)
	putUrl := form.ServiceUrl + serviceName + "/" + form.ServiceName + ":" + form.InstanceId
	_, err := tools.HttpPutReq(putUrl, nil, nil)
	if err != nil {
		return err
	}
	fmt.Println("Heartbeat sent ...")
	return nil
}

func (form RegistrationForm) SendHeartBeat() {
	for {
		err := form.Heartbeat()
		if err != nil {
			fmt.Println("Error!: ", err)
		}
		time.Sleep(20 * time.Second)
	}
}

func ConstructRegistrationBody(erm RegistrationForm, status string) *AppRegistrationBody {
	instanceId := erm.ServiceName + ":" + erm.InstanceId
	servicePort := fmt.Sprintf("%d", erm.ServicePort)
	hostAddress := erm.ServiceHost
	statusPageUrl := hostAddress + ":" + servicePort + "/health"
	healthCheckUrl := hostAddress + ":" + servicePort + "/health"
	homePageUrl := hostAddress + ":" + servicePort

	port := Port{
		Port:    servicePort,
		Enabled: "true",
	}

	securePort := Port{
		Port:    "443",
		Enabled: "false",
	}

	dataCenterInfo := DataCenterInfo{
		Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
		Name:  "MyOwn",
	}

	metadata := Metadata{
		Zone:    "primary",
		Profile: "dev",
		Port:    erm.ServicePort,
		Version: 1.01,
	}

	instance := InstanceDetails{
		InstanceId:       instanceId,
		HostName:         hostAddress,
		App:              erm.ServiceName,
		VipAddress:       erm.ServiceName,
		SecureVipAddress: erm.ServiceName,
		IpAddr:           hostAddress,
		Status:           status,
		Port:             port,
		SecurePort:       securePort,
		HomePageUrl:      homePageUrl,
		HealthCheckUrl:   healthCheckUrl,
		StatusPageUrl:    statusPageUrl,
		DataCenterInfo:   dataCenterInfo,
		Metadata:         metadata,
	}

	body := &AppRegistrationBody{
		Instance: instance,
	}

	return body
}
