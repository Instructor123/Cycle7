package main

import (
	"Cycle7-Client/wireGuard"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"time"
)

/*
	The intent is to create a a server and a client that communicate via wireguard.
	To start this will mean configuring a wireguard interface with all the necessary requirements.
	Eventually this could lead to a DNS query for a DGA-domain, but that'll be a stretch goal.
 */

/*
	Primary goals (must get done):
		1.) Creat an HTTP(S) client that communicates with an HTTP(S) server.
		2.) Client
			a.) The commands will be extremely simple for the sake of prototyping.
			b.) Create the wireguard interface and configure it to talk with the server
				I.) Could eventually get details from DNS, but for now it'll be hard coded
			c.) Send a Get request with a custom header indicating the client is alive
		3.) Server
			a.) Create the wireguard interface and configure it to talk with the client.
				I.) This will likely always be hardcoded as it is the server
			b.) Accept get requests and register clients.

	Secondary goals (really really want to get done):
		1.) Server
			a.) Allow the client to talk to up stream server(s), making a chain of communications.
				I.) The client A would connect to server B which would forward the request to server/client C
				II.) This could continue indefinatly
			b.) do additional tasking that would be expected of a remote application
		2.) Client
			a.) Request remote resources that are beyond the server.
			b.) Accept tasking
 */

func checkErr(e error){
	if e != nil{
		panic(e)
	}
}

//Retrieve the hardware address and the network address
func getHostInfo()(net.HardwareAddr, net.Addr){
	iFaces, err := net.InterfaceByName("ens33")
	checkErr(err)

	networkIPs, err := iFaces.Addrs()

	networkIP := networkIPs[0]
	macAddr := iFaces.HardwareAddr

	return macAddr, networkIP
}

//commands mostly taken from the demo on www.wireguard.com/quickstart
func createInterface(IP, Port string){

	//setup wg
	createWG(IP)

}

func createWG(IP string){
	//create the link:							ip link add wg0 type wireguard
	cmd := exec.Command("ip", "link", "add", "wg0", "type", "wireguard")
	_, err := cmd.Output()
	checkErr(err)

	//add the ip address:						ip addr add <address>/24 dev wg0
	cmd = exec.Command("ip", "addr", "add", IP+"/24", "dev", "wg0")
	_, err = cmd.Output()
	checkErr(err)

	//create the private key:					wg genkey > <location>
	//output, err := exec.Command("wg", "genkey").Output()
	//checkErr(err)

	//err = ioutil.WriteFile("wg_key",output,0744)
	//checkErr(err)

	//set the private key						wg set wg0 private-key <private key location>
	cmd = exec.Command("wg", "set", "wg0", "private-key", "wg_key")
	_, err = cmd.Output()
	checkErr(err)
	//set the link to up						ip link set wg0 up
	cmd = exec.Command("ip", "link", "set", "wg0", "up")
	_, err = cmd.Output()
	checkErr(err)
}

func wgConfigPeer(key, port string){
	cmd := exec.Command("wg", "set", "wg0", "listen-port", "12345", "peer", key,
		"allowed-ips", "172.16.0.2/32", "endpoint", "192.168.75.3:"+port)
	_, err := cmd.Output()
	checkErr(err)
}

//Must currently be run as root
func main(){

	//send beacon to server and wait for response
	for{
		resp, err := http.Get("http://192.168.75.3:8080/getInfo")
		if err != nil{
			fmt.Println("Couldn't connect, sleeping then trying again...")
			time.Sleep(5*time.Second)
		} else {

			port := resp.Header.Get("port")
			key := resp.Header.Get("key")
			vpnIP := resp.Header.Get("VPNIP")

			wireGuard.Initialize("172.16.0.3", "12345")
			wireGuard.ConfigPeer(key, port, "192.168.75.3", vpnIP, "12345")
			resp.Body.Close()
			break
		}
	}

	req, err := http.NewRequest(http.MethodGet, "http://192.168.75.3:8080/Data", nil)
	checkErr(err)
	req.Header.Set("Key", "ESRsw0jlzcFrpQsFRFyMSqVTqyhyvDfmnkCmu9uvYTs=")
	req.Header.Set("Port", "12345")
	req.Header.Set("VPNIP", "172.16.0.3")
	client := &http.Client{}

	resp, err := client.Do(req)
	checkErr(err)
	if resp.StatusCode != http.StatusOK{
		fmt.Println("Problem sending key...")
	}
	resp.Body.Close()
}
