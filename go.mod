module github.com/IBM/satellite-object-storage-plugin

go 1.15

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190516230258-a675ac48af67
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190313205120-8b27c41bdbb1
)

require (
	github.com/container-storage-interface/spec v1.2.0
	github.com/ctrox/csi-s3 v1.1.1 // indirect
	github.com/kubernetes-csi/drivers v1.0.2
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.ibm.com/alchemy-containers/ibm-csi-common v1.0.0-beta08
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/grpc v1.27.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)
