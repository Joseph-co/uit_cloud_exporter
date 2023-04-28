package syst

import (
	"log"
	"net"
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
			log.Println(oneAddrs)
			if oneAddrs[0] == vip {
				return true
			}
		}
	}
	return false
}




