package syst

import (
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func VipCheck(vip string) bool {
	interfaceList, err := net.Interfaces()
	if err != nil {
		log.Fatalf("errors when getting interfacelist: %v\n",err)
	}
	var byName *net.Interface
	var addrList []net.Addr
	var oneAddrs []string
	for _, i := range interfaceList {
		byName, err = net.InterfaceByName(i.Name)
		if err != nil {
			log.Fatalf("error when geting interface addr by name: %v\n",err)
		}
		addrList, err = byName.Addrs()
		if err != nil {
			log.Fatalf("error when geting addrList: %v\n",err)
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
		log.Fatalf("error when get interfaces: %v\n",err)
	}
	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp != 0 {
			ifcheck[iface.Name] = true
		} else {
			ifcheck[iface.Name] = false
		}
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatalf("error when geting interfaces addr:%v\n",err)
		}
		for _, addr := range addrs {
			ip := strings.SplitN(addr.String(), "/", 2)
			ipmap[iface.Name] = append(ipmap[iface.Name], ip[0])
		}
	}
	return ipmap, ifcheck
}

func GetSpecIface()[]string {
	ifaces, err := net.Interfaces()
	Iface := make([]string,0,len(ifaces))
	if err != nil {
		log.Fatalf("error when geting ifaces:%v\n",err)
	}
	for _, iface := range ifaces {
		if strings.HasPrefix(iface.Name, "net") || strings.HasPrefix(iface.Name, "eth") {
			Iface = append(Iface,iface.Name)
		}
	}

	path,err :=  PathExists("/proc/net/bonding")
	if err != nil {
		log.Fatalf("error bonding path exit or not:%v\n",err)
	}
	if path {
		cmd := exec.Command("ls", "/proc/net/bonding")
		output, err := cmd.Output()
		if err != nil {
			log.Fatalf("Error executing command:%v\n", err)
		}
		for _,i := range strings.Split(string(output), "\n"){
			Iface = append(Iface,i)
		}
	}else {
		return Iface
	}
	return Iface
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


func GetHostName() (string, error) {
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatalf("error when geting hostname:%v\n",err)
		return "", err
	}
	return hostName, nil
}


func GetKeepAlive() bool {
	cmd := exec.Command("ps", "-ef")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error executing command:%v\n", err)
		return false
	}
	//substr := "/bin/zsh"
	substr := "/etc/cz-ha/keepalived.conf"
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line,substr) {
			return true
		}
	}
	return false
}

func GetHaproxy() bool {
	cmd := exec.Command("ps", "-ef")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error executing command:%v\n", err)
		return false
	}
	//substr := "/bin/zsh"
	substr := "/etc/haproxy/haproxy.cfg"
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line,substr) {
			return true
		}
	}
	return false
}







