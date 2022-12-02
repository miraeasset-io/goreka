package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type OperationBody interface {
}

func HttpPostReq(urlToPost string, body OperationBody, headers map[string]string) (*http.Response, error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.Encode(body)

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(http.MethodPost, urlToPost, &buffer)
	if err != nil {
		return nil, errors.New("Error while creating http request  " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Error while making POST call " + err.Error())
	}

	defer closeHttpIo(resp.Body)

	return resp, err
}

func HttpPutReq(urlToPost string, body OperationBody, headers map[string]string) (*http.Response, error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	var buffer bytes.Buffer
	if body != nil {
		encoder := json.NewEncoder(&buffer)
		err := encoder.Encode(body)
		if err != nil {
			return nil, errors.New("Error while encoding body " + err.Error())
		}
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(http.MethodPut, urlToPost, &buffer)
	if err != nil {
		fmt.Println("Error while creating the http request  " + err.Error())
		return nil, errors.New("Error while creating http request " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while making PUT call " + err.Error())
		return nil, errors.New("Error while making PUT call " + err.Error())
	}

	defer closeHttpIo(resp.Body)

	return resp, err
}

func closeHttpIo(Body io.ReadCloser) {
	Body.Close()
}

func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
