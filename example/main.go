package main

import (
	"fmt"
	"log"

	omada "github.com/dougbw/go-omada"
)

func main() {

	// variables
	controllerUrl := ""
	user := ""
	pass := ""

	// setup
	omada := omada.New(controllerUrl)
	err := omada.GetControllerInfo()
	if err != nil {
		log.Fatal(err)
	}

	// login
	err = omada.Login(user, pass)
	if err != nil {
		log.Fatal(err)
	}

	// get clients
	clients := omada.GetClients()
	for _, client := range clients {
		fmt.Printf("client ip: %s, dnsName: %s, name: %s\n", client.Ip, client.DnsName, client.Name)
	}

	// get devices
	devices := omada.GetDevices()
	for _, device := range devices {
		fmt.Printf("device name: %s, dnsName: %s,  ip: %s\n", device.Name, device.DnsName, device.IP)
	}

	// get networks
	networks := omada.GetNetworks()
	for _, network := range networks {
		fmt.Printf("network name: %s, subnet: %s, domain: %s\n", network.Name, network.Subnet, network.Domain)
	}

}
