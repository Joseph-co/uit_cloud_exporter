package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"strings"
	"sync"
	"uit_cloud_exporter/docker"
	syst "uit_cloud_exporter/syst"
)

const (
	namespace = "uit"
)

var (
	val int
)

type NodeInfo struct {
	HostName string
	VIPAddr  string
}

type Container struct {
	ID      string
	Name    string
	Running bool
}

type Exporter struct {
	NodeInfo
	Container
	metrix map[string]* prometheus.Desc
	mutex sync.Mutex
	clientSet *kubernetes.Clientset
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
	}
    strc,kubeconfig := docker.GetK8sConf()
	//	create config
	//config, err := clientcmd.BuildConfigFromFlags("https://172.18.70.241:6443", clientcmd.RecommendedHomeFile)
	config, err := clientcmd.BuildConfigFromFlags(strc, *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	//	create client
	clientset, err := kubernetes.NewForConfig(config)
	//config, err := clientcmd.BuildConfigFromFlags(master, "config")
	if err != nil {
		log.Fatal(err)
	}
	return &Exporter{
		clientSet: clientset,
		NodeInfo: nodeinfo,
		metrix: map[string]*prometheus.Desc{
			"container_status" : prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "container_up"),
				"Was the last Mirth query successful.",
				[]string{"container"}, nil,
			),
			"keepalive" : prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "keepalive_up"),
				"Was the last Mirth query successful.",
				[]string{"keepalive"}, nil,
			),
			"net_ifaces" : prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "iface_up"),
				"Was the last Mirth query successful.",
				[]string{"iface_name","iface_ip"}, nil,
			),
			"haproxy" : prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "haproxy_up"),
				"Was the last Mirth query successful.",
				[]string{"ha"}, nil,
			),
			"promMetrics" : prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "promMetrics_up"),
				"Was the last Mirth query successful.",
				[]string{"metrics"}, nil,
			),
		},
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // 加锁
	defer e.mutex.Unlock()

	for _, i := range docker.GetContainerIDs() {

		cInfo := docker.GetContainerInspect(i)
		if cInfo.State.Running {
			val = 1
		}else {
			val = 0
		}
		e.Container.Name = cInfo.Name[1:]
		if strings.HasPrefix(e.Container.Name,"k8s_") {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			e.metrix["container_status"], prometheus.GaugeValue, float64(val), e.Container.Name,
		)
	}
	if  syst.GetKeepAlive() {
		ch <- prometheus.MustNewConstMetric(
			e.metrix["keepalive"], prometheus.GaugeValue, 1, "keepalive",
		)
	}else {
		ch <- prometheus.MustNewConstMetric(
			e.metrix["keepalive"], prometheus.GaugeValue, 0, "keepalive",
		)
	}
	if syst.GetHaproxy() {
		ch <- prometheus.MustNewConstMetric(
			e.metrix["haproxy"], prometheus.GaugeValue, 1, "haproxy",
		)
	}else {
		ch <- prometheus.MustNewConstMetric(
			e.metrix["haproxy"], prometheus.GaugeValue, 0, "haproxy",
		)
	}

	ipMap,ifCheck := syst.GetIpMap()
    Iface := syst.GetSpecIface()
	for _,iname := range Iface {
		var iface_ip string
		nodeInfoNet := []string{
			iname,
			iface_ip,
		}
		mi,ok := ipMap[iname]
			if ok {
				for _,iface := range mi{
					nodeInfoNet[1] = iface
					if ifCheck[iname] {
						ch <- prometheus.MustNewConstMetric(
							e.metrix["net_ifaces"], prometheus.GaugeValue, 1, nodeInfoNet...,
						)
					}else{
						ch <- prometheus.MustNewConstMetric(
							e.metrix["net_ifaces"], prometheus.GaugeValue, 0, nodeInfoNet...,
						)
					}
				}
			}else {
				nodeInfoNet[1] = "none"
				if ifCheck[iname] {
					ch <- prometheus.MustNewConstMetric(
						e.metrix["net_ifaces"], prometheus.GaugeValue, 1, nodeInfoNet...,
					)
				}else{
					ch <- prometheus.MustNewConstMetric(
						e.metrix["net_ifaces"], prometheus.GaugeValue, 0, nodeInfoNet...,
					)
				}
			}
	}
	//for iname := range ipMap {
	//	ic, ok := ifCheck[iname]
	//	if ok {
	//		for _, ip := range ipMap[iname] {
	//			nodeInfoNet := []string{
	//				e.HostName,
	//				iname,
	//				ip,
	//			}
	//			if ic {
	//				ch <- prometheus.MustNewConstMetric(
	//					e.metrix["net_ifaces"], prometheus.GaugeValue, 1, nodeInfoNet...,
	//				)
	//			} else {
	//				ch <- prometheus.MustNewConstMetric(
	//					e.metrix["net_ifaces"], prometheus.GaugeValue, 0, nodeInfoNet...,
	//				)
	//			}
	//		}
	//	}
	//}
	if docker.GetDeployStatus(e.clientSet) {
		ch <- prometheus.MustNewConstMetric(
			e.metrix["promMetrics"], prometheus.GaugeValue, 1, "metrics",
		)
	}else{
		ch <- prometheus.MustNewConstMetric(
			e.metrix["promMetrics"], prometheus.GaugeValue, 0, "metrics",
		)
	}
}
func main() {
	exporter := NewExporter()
	prometheus.MustRegister(exporter)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
