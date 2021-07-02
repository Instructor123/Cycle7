package main

import (
	"Cycle7-Server/wireGuard"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func checkErr(e error){
	if e != nil{
		panic(e)
	}
}

func getHostInfo()(net.HardwareAddr, net.Addr){
	iFaces, err := net.InterfaceByName("ens33")
	checkErr(err)

	networkIPs, err := iFaces.Addrs()

	networkIP := networkIPs[0]
	macAddr := iFaces.HardwareAddr

	return macAddr, networkIP
}

func handleRequest(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	switch r.Method {
	case http.MethodGet:
		if r.RequestURI == "/getInfo" {
			w.Header().Set("port", "34221")
			w.Header().Set("key", "IGH0Pu19GFb7/q/eAjM34ILgtZnKU3aTVmf8iEIjDmA=")
			w.Header().Set("VPNIP", "172.16.0.2")
			fmt.Fprint(w, "")
		}
		if r.RequestURI == "/Data" {
			peerKey := r.Header.Get("Key")
			peerPort := r.Header.Get("Port")
			peerVPN := r.Header.Get("VPNIP")
			ip := strings.Split(r.RemoteAddr, ":")[0]
			wireGuard.ConfigPeer(peerKey, peerPort, ip, peerVPN, "34221")
		}
	}
}

func main(){

	_, CIDR := getHostInfo()
	myIP, _, err := net.ParseCIDR(CIDR.String())
	checkErr(err)

	wireGuard.Initialize("172.16.0.2", "34221")

	quitChan := make(chan int)

	myServer := &http.Server{
		Addr: myIP.String()+":8080",
		Handler: http.HandlerFunc(handleRequest),
	}

	go func(){
		log.Fatal(myServer.ListenAndServe())
	}()

	<-quitChan

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()
	myServer.Shutdown(ctx)

}
