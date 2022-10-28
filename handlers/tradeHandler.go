package handlers

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"nodeSelector/utils"
	"os"
	"strconv"
	"strings"
)

const (
	Slave  = "ftdb-slave1-bobft"
	Mysql  = "ftdb-mysql-bobft"
	Master = "ftdb-master-bobft"
)

var (
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
	NodeArray = []string{
		"worker1", "worker2", "worker3", "worker4", "worker5", "worker6", "worker7", "worker8", "worker9",
		"worker1", "worker2", "worker3", "worker4", "worker5", "worker6", "worker7", "worker8", "worker9",
		"worker1", "worker2", "worker3", "worker4", "worker5", "worker6", "worker7", "worker8", "worker9",
		"worker1", "worker2", "worker3", "worker4", "worker5", "worker6", "worker7", "worker8",
	}

	NodeCount = len(NodeArray)

	logger.Info("Node Count: ", NodeCount)
	logger.Info("Node Array: ", NodeArray)

	logger.Info("SlaveRequestCpu: ", SlaveRequestCpu)
	logger.Info("SlaveRequestMem: ", SlaveRequestMem)
	logger.Info("SlaveLimitCpu: ", SlaveLimitCpu)
	logger.Info("SlaveLimitMem: ", SlaveLimitMem)

	logger.Info("MySQLRequestCpu: ", MySQLRequestCpu)
	logger.Info("MySQLRequestMem: ", MySQLRequestMem)
	logger.Info("MySQLLimitCpu: ", MySQLLimitCpu)
	logger.Info("MySQLLimitMem: ", MySQLLimitMem)
}

func tradHandler(deployment *appsv1.Deployment) (string, bool) {
	logger.Info(deployment.Name)
	arrayNum := 0
	if strings.Contains(deployment.Name, Master) {
		arrayNum, _ = strconv.Atoi(strings.Split(deployment.Name, Master)[1])
	} else if strings.Contains(deployment.Name, Slave) {
		arrayNum, _ = strconv.Atoi(strings.Split(deployment.Name, Slave)[1])
		if SlaveRequestCpu != "" {
			deployment.Spec.Template.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(SlaveRequestCpu),
			}
		}
		if SlaveRequestMem != "" {
			deployment.Spec.Template.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(SlaveRequestMem),
			}
		}

		if SlaveLimitCpu != "" {
			deployment.Spec.Template.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(SlaveLimitCpu),
			}
		}
		if SlaveLimitMem != "" {
			deployment.Spec.Template.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(SlaveLimitMem),
			}
		}

	} else if strings.Contains(deployment.Name, Mysql) {
		arrayNum, _ = strconv.Atoi(strings.Split(deployment.Name, Mysql)[1])
		if MySQLRequestCpu != "" {
			deployment.Spec.Template.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(MySQLRequestCpu),
			}
		}
		if MySQLRequestMem != "" {
			deployment.Spec.Template.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(MySQLRequestMem),
			}
		}

		if MySQLLimitCpu != "" {
			deployment.Spec.Template.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(MySQLLimitCpu),
			}
		}
		if MySQLLimitMem != "" {
			deployment.Spec.Template.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(MySQLLimitMem),
			}
		}
	} else {
		return "deployment name not format", false
	}
	targetNode := NodeArray[arrayNum%NodeCount]
	deployment.Spec.Template.Spec.NodeName = targetNode
	logger.Info("Requests: ", deployment.Spec.Template.Spec.Containers[0].Resources.Requests)
	logger.Info("Limits", deployment.Spec.Template.Spec.Containers[0].Resources.Limits)

	return fmt.Sprintf("deployment %s scheduled to node %s", deployment.Name, targetNode), true
}
