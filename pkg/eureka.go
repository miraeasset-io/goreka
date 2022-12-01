package goreka

import (
	"fmt"
	"goreka/tools"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/carlescere/scheduler"
	log "github.com/sirupsen/logrus"
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

type RegistrationManager struct {
	ServiceName string
	ServiceHost string
	ServicePort int
	InstanceId  string
	AppProfile  string
	ServiceUrl  string
}

func RegisterService(erm RegistrationManager) {
	log.Println("Registering service with status: STARTING")
	body := ConstructRegistrationBody(erm, "STARTING")

	serviceName := strings.ToUpper(erm.ServiceName)

	postUrl := erm.ServiceUrl + serviceName
	log.Debugln(postUrl)

	_, err := tools.HttpPostReq(postUrl, body, nil)
	if err != nil {
		panic(err)
	}

	waitSec, _ := strconv.Atoi(fmt.Sprintf("%v", viper.Get("STARTUP_WAIT_SEC")))
	d := time.Duration(waitSec)

	log.Printf("Waiting for %d seconds for application to start properly ...", waitSec)
	time.Sleep(d * time.Second)

	log.Print("Updating the status to: UP")
	bodyUP := ConstructRegistrationBody(erm, "UP")

	_, err = tools.HttpPostReq(erm.ServiceUrl+serviceName, bodyUP, nil)
	if err != nil {
		panic(err)
	}
}

func SendHeartBeat(erm RegistrationManager) {
	serviceName := strings.ToUpper(erm.ServiceName)
	heartBeat := func() {
		putUrl := erm.ServiceUrl + serviceName + "/" + erm.ServiceName + ":" + erm.InstanceId
		res, err := tools.HttpPutReq(putUrl, nil, nil)
		if err != nil {
			log.Errorln(err)
		}
		log.Debugln(res)
		log.Debugln("Heartbeat sent ...")
	}

	// Run every 25 seconds but not now.
	res, err := scheduler.Every(25).Seconds().Run(heartBeat)
	if err != nil {
		log.Errorln(err)
	}
	log.Debugln(res)

	runtime.Goexit()
}

func InitServiceDiscovery(Config Config) *RegistrationManager {
	log.Infoln("Initializing service discovery ...")

	manager := new(RegistrationManager)
	manager.AppProfile = Config.AppProfile
	manager.InstanceId = Config.InstanceId
	manager.ServiceName = Config.ServiceName
	manager.ServiceHost = Config.ServiceHost
	manager.ServicePort = Config.ServicePort
	manager.ServiceUrl = Config.EurekaUrl

	RegisterService(*manager)
	go SendHeartBeat(*manager)

	return manager
}

func UnRegisterEurekaService(rm RegistrationManager) {
	log.Warningln("UnRegistering service from eureka ...")
	res, err := tools.HttpPostReq(rm.ServiceUrl, nil, nil)
	if err != nil {
		log.Errorln(err)
	}
	log.Debugln(res)
}

func ConstructRegistrationBody(erm RegistrationManager, status string) *AppRegistrationBody {
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
