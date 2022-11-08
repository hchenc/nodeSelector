package handlers

import (
	"context"
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"nodeSelector/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	Slave  = "ftdb-slave1-bobft"
	Mysql  = "ftdb-mysql-bobft"
	Master = "ftdb-master-bobft"
)

var (
	config *rest.Config

	limitCPU string

	NodeCount int
	NodeArray []string
	logger    = utils.GetLogger()

	ResourceArray []string

	MasterRequestCpu = os.Getenv("MasterRequestCpu")
	MasterRequestMem = os.Getenv("MasterRequestMem")
	MasterLimitCpu   = os.Getenv("MasterLimitCpu")
	MasterLimitMem   = os.Getenv("MasterLimitMem")

	SlaveRequestCpu = os.Getenv("SlaveRequestCpu")
	SlaveRequestMem = os.Getenv("SlaveRequestMem")
	SlaveLimitCpu   = os.Getenv("SlaveLimitCpu")
	SlaveLimitMem   = os.Getenv("SlaveLimitMem")

	MySQLRequestCpu = os.Getenv("MySQLRequestCpu")
	MySQLRequestMem = os.Getenv("MySQLRequestMem")
	MySQLLimitCpu   = os.Getenv("MySQLLimitCpu")
	MySQLLimitMem   = os.Getenv("MySQLLimitMem")
)

func init() {
	if cpu := os.Getenv("Limit_CPU"); cpu != "" {
		limitCPU = cpu
	} else {
		limitCPU = "2"
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodes, _ := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	for _, node := range nodes.Items {
		if len(node.Spec.Taints) != 0 {
			continue
		}
		NodeArray = append(NodeArray, node.Name)
	}

	//NodeArray = []string{
	//	"worker1", "worker2", "worker3", "worker4", "worker5", "worker6", "worker7", "worker8", "worker9",
	//	"worker1", "worker2", "worker3", "worker4", "worker5", "worker6", "worker7", "worker8", "worker9",
	//	"worker1", "worker2", "worker3", "worker4", "worker5", "worker6", "worker7", "worker8", "worker9",
	//	"worker1", "worker2", "worker3", "worker4", "worker5", "worker6", "worker7", "worker8",
	//}

	NodeCount = len(NodeArray)

	logger.Info("Node Count: ", NodeCount)
	logger.Info("Node Array: ", NodeArray)

	//logger.Info("SlaveRequestCpu: ", SlaveRequestCpu)
	//logger.Info("SlaveRequestMem: ", SlaveRequestMem)
	//logger.Info("SlaveLimitCpu: ", SlaveLimitCpu)
	//logger.Info("SlaveLimitMem: ", SlaveLimitMem)
	//
	//logger.Info("MySQLRequestCpu: ", MySQLRequestCpu)
	//logger.Info("MySQLRequestMem: ", MySQLRequestMem)
	//logger.Info("MySQLLimitCpu: ", MySQLLimitCpu)
	//logger.Info("MySQLLimitMem: ", MySQLLimitMem)
}

func oltpHandler(deployment *appsv1.Deployment) (string, bool) {
	logger.Info(deployment.Name)
	arrayNum := 0
	if strings.Contains(deployment.Name, Master) {
		arrayNum, _ = strconv.Atoi(strings.Split(deployment.Name, Master)[1])
	} else if strings.Contains(deployment.Name, Slave) {
		deployment.Spec.Template.Spec.Containers[0].Resources.Requests = v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("250m"),
			v1.ResourceMemory: resource.MustParse("2Gi"),
		}
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits = v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse(limitCPU),
			v1.ResourceMemory: resource.MustParse("2Gi"),
		}
		arrayNum, _ = strconv.Atoi(strings.Split(deployment.Name, Slave)[1])
	} else if strings.Contains(deployment.Name, Mysql) {
		deployment.Spec.Template.Spec.Containers[0].Resources.Requests = v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("250m"),
			v1.ResourceMemory: resource.MustParse("2Gi"),
		}
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits = v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse(limitCPU),
			v1.ResourceMemory: resource.MustParse("2Gi"),
		}
		arrayNum, _ = strconv.Atoi(strings.Split(deployment.Name, Mysql)[1])
	} else {
		return "deployment name not format", false
	}
	targetNode := NodeArray[arrayNum%NodeCount]
	deployment.Spec.Template.Spec.NodeName = targetNode

	return fmt.Sprintf("deployment %s scheduled to node %s", deployment.Name, targetNode), true
}
