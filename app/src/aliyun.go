package main

//curl -X POST  127.0.0.1:8040/aliyun/slbChange -d"machine=stg&weight=1"
import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"net/http"
)

const (
	accessKeyId     = "LTAI8NnspfitFU3q"
	accessKeySecret = "VLRNtYfzzeLjGz4lkSz8khKER6wbv0"
)

var server = map[string]string{"stg": "i-j5efdcn2conrdv3gnddi", "p1": "i-j5e0y1jwzpb1ylt9z7nu", "p2": "i-j5e5z9q2vuv7z3wlq4pl", "p3": "i-j5efdcn2conrdv3gnddj"}

func slbChangeModCall(mod string) string {
	client, err := slb.NewClientWithAccessKey("cn-shenzhen-finance-1", accessKeyId, accessKeySecret)

	request := slb.CreateSetBackendServersRequest()
	if mod == "alone" {
		request.BackendServers = "[{\"ServerId\":\"" + server["stg"] + "\",\"Weight\":\"1\"},{\"ServerId\":\"" + server["p1"] + "\",\"Weight\":\"0\"},{\"ServerId\":\"" + server["p2"] + "\",\"Weight\":\"0\"},{\"ServerId\":\"" + server["p3"] + "\",\"Weight\":\"0\"}]"
	}
	if mod == "normal" {
		request.BackendServers = "[{\"ServerId\":\"" + server["stg"] + "\",\"Weight\":\"0\"},{\"ServerId\":\"" + server["p1"] + "\",\"Weight\":\"1\"},{\"ServerId\":\"" + server["p2"] + "\",\"Weight\":\"1\"},{\"ServerId\":\"" + server["p3"] + "\",\"Weight\":\"1\"}]"
	}

	request.LoadBalancerId = "lb-j5eb66q79lhgdvkseetgp"

	response, err := client.SetBackendServers(request)
	if err != nil {
		fmt.Print(err.Error())
		return "failed"
	}
	fmt.Printf("response is %#v\n", response)
	return "ok"
}

func slbChangeCall(machine string, weight string) string {
	client, err := slb.NewClientWithAccessKey("cn-shenzhen-finance-1", accessKeyId, accessKeySecret)

	request := slb.CreateSetBackendServersRequest()

	var serverId string
	_, ok := server[machine]
	if ok {
		serverId = server[machine]
	} else {
		serverId = machine
	}
	request.BackendServers = "[{\"ServerId\":\"" + serverId + "\",\"Weight\":\"" + weight + "\"}]"
	request.LoadBalancerId = "lb-j5eb66q79lhgdvkseetgp"

	response, err := client.SetBackendServers(request)
	if err != nil {
		fmt.Print(err.Error())
		return "failed"
	}
	fmt.Printf("response is %#v\n", response)
	return "ok"
}

func slbChange(w http.ResponseWriter, r *http.Request) string {
	machine := r.FormValue("machine")
	weight := r.FormValue("weight")
	return slbChangeCall(machine, weight)
}

func getServerCall() string {
	client, err := slb.NewClientWithAccessKey("cn-shenzhen-finance-1", accessKeyId, accessKeySecret)

	request := slb.CreateAddBackendServersRequest()

	request.BackendServers = "[{\"ServerId\":\"i-j5efdcn2conrdv3gnddi\",\"Weight\":\"100\",\"Type\":\"ecs\"}]"
	request.LoadBalancerId = "lb-j5eb66q79lhgdvkseetgp"

	response, err := client.AddBackendServers(request)
	if err != nil {
		fmt.Print(err.Error())
		return "failed"
	}
	fmt.Printf("response is %#v\n", response)
	res, _ := json.Marshal(response.BackendServers)
	return string(res)
}

func getServer(w http.ResponseWriter, r *http.Request) string {
	return getServerCall()
}

func getStatus(w http.ResponseWriter, r *http.Request) string {
	return getStatusCall()
}

func getStatusCall() string {
	client, err := slb.NewClientWithAccessKey("cn-shenzhen-finance-1", accessKeyId, accessKeySecret)

	request := slb.CreateDescribeHealthStatusRequest()
	request.LoadBalancerId = "lb-j5eb66q79lhgdvkseetgp"
	request.ListenerProtocol = "https"
	request.ListenerPort = requests.NewInteger(443)

	response, err := client.DescribeHealthStatus(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
	res, _ := json.Marshal(response.BackendServers)
	return string(res)
}
