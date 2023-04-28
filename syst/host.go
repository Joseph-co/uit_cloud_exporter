package syst

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func VipCheck(vip string) bool {
	interface_list, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	var byName *net.Interface
	var addrList []net.Addr
	var oneAddrs []string
	for _, i := range interface_list {
		byName, err = net.InterfaceByName(i.Name)
		if err != nil {
			log.Fatal(err)
		}
		addrList, err = byName.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, oneAddr := range addrList {
			oneAddrs = strings.SplitN(oneAddr.String(), "/", 2)
			if oneAddrs[0] == vip {
				return true
			}
		}
	}
	return false
}

func GetIpMap() (ipMap map[string][]string,ifCheck map[string]bool) {
	ipmap := make(map[string][]string)
	ifcheck := make(map[string]bool)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp != 0 {
			fmt.Printf("Interface %v is up\n", iface.Name)
			ifcheck[iface.Name] = true
		} else {
			fmt.Printf("Interface %v is down\n", iface.Name)
			ifcheck[iface.Name] = false
		}
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			ip := strings.SplitN(addr.String(), "/", 2)
			ipmap[iface.Name] = append(ipMap[iface.Name], ip[0])
		}
	}
	return ipmap, ifcheck
}



func GetHostName() (string, error) {
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return hostName, nil
}




