# uit_cloud_exporter部署与功能介绍
###镜像制作
- 准备项目代码
  ```
  cd ~/uit_cloud_exporter/
  ```
- 构建镜像
  ```
  docker build -t uit-cloud-exporter:v0.1.0 .
  ```
- 保存镜像到本地
  ```
  docker save -0 uit-cloud-exporter.tar uit-cloud-exporter:v0.1.0
  ```
###部署vipbind
  ```
  kubectl apply -f uitCon-deploy.yaml
  ```
###配合promtheus监控需要添加servicemonitor配置
  ```
  kubectl apply -f uitCon-conf.yaml
  ```

###功能及原理

- 针对uit云平台的服务，开发的exporter
- 运用成熟的docker的sdk，client-go，go系统调用等方式

### indexs of the exporter supplied
uit_iface_up{iface_ip="fe80::20ef:f85a:a31d:b8a3",iface_name="utun4"} 0/1 \
uit_iface_up{iface_ip="172.18.50.191",iface_name="en0"} 0/1 \
uit_keepalive_up{keepalive="keepalive"} 0/1 \
uit_container_up{container="minikube"} 0/1 

####UIT platform container names
host-agent\
ws-endpoint\
search-api\
k8s-result-processor\
web-console\
novnc\
uit-haproxy\
uit-keepalived\
timer-task\
event-alarm\
resource-alarm\
k8s-scheduler\
k8s-events-listener\
k8s-api


