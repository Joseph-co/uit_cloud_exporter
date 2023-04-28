package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"sync"
	"uit_cloud_exporter/docker"
	syst "uit_cloud_exporter/syst"
)

const (
	namespace = "uit"
)

var (
	val int
	vip string = "127.0.0.1"
	netdev string = "en0"
)

type NodeInfo struct {
	HostName string
	IPAddr   string
	VIPAddr  string
}

type Container struct {
	ID      string
	Name    string
	Image   string
	Running bool
}

type Exporter struct {
	NodeInfo
	Container
	metrix map[string]* prometheus.Desc
	mutex sync.Mutex
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.metrix {
		ch <- m
	}
}

func NewExporter() *Exporter {
	hostname, _ := syst.GetHostName()
	nodeinfo := NodeInfo{
		HostName: hostname,
		IPAddr:   "127.0.0.1",
		VIPAddr: vip,
	}
	return &Exporter{
		NodeInfo: nodeinfo,
		metrix: map[string]*prometheus.Desc{
			"container_status" : prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "container_up"),
				"Was the last Mirth query successful.",
				[]string{"hostname","host_ip","container"}, nil,
			),
			"keepalive_vip" : prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "vip_up"),
				"Was the last Mirth query successful.",
				[]string{"hostname","host_ip","vip"}, nil,
			),
			"net_ifaces" : prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "iface_up"),
				"Was the last Mirth query successful.",
				[]string{"hostname","iface_name","iface_ip"}, nil,
			),
			"ha_proxy" : prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "ha_proxy_up"),
				"Was the last Mirth query successful.",
				[]string{"hostname","host_ip","container"}, nil,
			),
		},
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // 加锁
	defer e.mutex.Unlock()
	nodeIn := []string{
		e.HostName,
		e.IPAddr,
		" ",
	}
	for _, i := range docker.GetContainerIDs() {
		cInfo := docker.GetContainerInspect(i)
		if cInfo.State.Running {
			val = 1
		}else {
			val = 0
		}
		e.Container.Name = cInfo.Name
		nodeIn[2] = e.Container.Name
		ch <- prometheus.MustNewConstMetric(
			e.metrix["container_status"], prometheus.GaugeValue, float64(val), nodeIn...,
		)
	}
	nodeInfoVip := []string{
		e.HostName,
		e.IPAddr,
		e.VIPAddr,
	}
	if  syst.VipCheck(vip) {
		ch <- prometheus.MustNewConstMetric(
			e.metrix["keepalive_vip"], prometheus.GaugeValue, 1, nodeInfoVip...,
		)
	}else {
		ch <- prometheus.MustNewConstMetric(
			e.metrix["keepalive_vip"], prometheus.GaugeValue, 0, nodeInfoVip...,
		)
	}
	ipMap,ifCheck := syst.GetIpMap()
	//for iname := range ifCheck {
	//	var iface_ip string
	//	miname,ok := ipMap[iname]
	//		if ok {
	//			iface_ip = miname[0]
	//		}else {
	//			iface_ip = "none"
	//		}
	//	nodeInfoNet := []string{
	//		e.HostName,
	//		iname,
	//		iface_ip,
	//	}
	//	if ifCheck[iname] {
	//		ch <- prometheus.MustNewConstMetric(
	//			e.metrix["net_ifaces"], prometheus.GaugeValue, 1, nodeInfoNet...,
	//		)
	//	}else{
	//		ch <- prometheus.MustNewConstMetric(
	//			e.metrix["net_ifaces"], prometheus.GaugeValue, 0, nodeInfoNet...,
	//		)
	//	}
	//}
	for iname := range ipMap {
		ic, ok := ifCheck[iname]
		if ok {
			for _, ip := range ipMap[iname] {
				nodeInfoNet := []string{
					e.HostName,
					iname,
					ip,
				}
				if ic {
					ch <- prometheus.MustNewConstMetric(
						e.metrix["net_ifaces"], prometheus.GaugeValue, 1, nodeInfoNet...,
					)
				} else {
					ch <- prometheus.MustNewConstMetric(
						e.metrix["net_ifaces"], prometheus.GaugeValue, 0, nodeInfoNet...,
					)
				}
			}
		}
	}
}

func main() {
	exporter := NewExporter()
	prometheus.MustRegister(exporter)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9100", nil))
}
