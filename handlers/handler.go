package handlers

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strconv"
)

//type WebHookHandler struct {
//	Client  client.Client
//	Decoder *admission.Decoder
//}
//
//func (handler *WebHookHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
//	sts := &appsv1.StatefulSet{}
//	if err := handler.Decoder.Decode(req, sts); err == nil {
//		logger.Info("start to handle sts: ", sts.GetName())
//		if sts.GetName() == "seckill-master" {
//			sts.Spec.Template.Spec.PriorityClassName = "system-cluster-critical"
//			for _, container := range sts.Spec.Template.Spec.Containers {
//				if container.Name != "bobft-seckill" {
//					continue
//				}
//				container.Resources.Requests = corev1.ResourceList{
//					corev1.ResourceCPU:    resource.MustParse("4"),
//					corev1.ResourceMemory: resource.MustParse("8Gi"),
//				}
//				container.Resources.Limits = corev1.ResourceList{
//					corev1.ResourceCPU:    resource.MustParse("4"),
//					corev1.ResourceMemory: resource.MustParse("8Gi"),
//				}
//			}
//		} else if sts.GetName() == "seckill-slave" {
//			sts.Spec.Template.Spec.PriorityClassName = "high-priority"
//			for _, container := range sts.Spec.Template.Spec.Containers {
//				if container.Name != "bobft-seckill" {
//					continue
//				}
//				container.Resources.Requests = corev1.ResourceList{
//					corev1.ResourceCPU:    resource.MustParse("4"),
//					corev1.ResourceMemory: resource.MustParse("8Gi"),
//				}
//				container.Resources.Limits = corev1.ResourceList{
//					corev1.ResourceCPU:    resource.MustParse("4"),
//					corev1.ResourceMemory: resource.MustParse("8Gi"),
//				}
//			}
//		} else {
//			return admission.Allowed("no target StatefulSet")
//		}
//
//		if redisNode != -1 {
//			sts.Spec.Template.Spec.NodeName = NodeArray[redisNode]
//		}
//		//logger.Info(sts.Spec.Template.Spec.Containers[0].Resources.Requests)
//		//logger.Info(sts.Spec.Template.Spec.Containers[0].Resources.Limits)
//		marshaledSts, err := json.Marshal(sts)
//		if err != nil {
//			logger.Info(err.Error())
//			return admission.Errored(http.StatusInternalServerError, err)
//		}
//		logger.Info("end to handle sts: ", sts.GetName())
//		return admission.PatchResponseFromRaw(req.Object.Raw, marshaledSts)
//	}
//
//	dp := &appsv1.Deployment{}
//	if err := handler.Decoder.Decode(req, dp); err == nil {
//		logger.Info("start to handle dp: ", dp.GetName())
//		if dp.GetName() == "stress-boss" {
//			dp.Spec.Template.Spec.PriorityClassName = "low-priority"
//		} else if dp.GetName() == "stress-worker" {
//			dp.Spec.Template.Spec.PriorityClassName = "low-priority"
//			dp.Spec.Template.Spec.Affinity.PodAffinity = &corev1.PodAffinity{
//				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
//					{
//						Weight: 100,
//						PodAffinityTerm: corev1.PodAffinityTerm{
//							TopologyKey: "kubernetes.io/hostname",
//							LabelSelector: &metav1.LabelSelector{
//								MatchLabels: map[string]string{
//									"io.kompose.service": "stress-worker",
//								},
//							},
//						},
//					},
//				},
//			}
//		} else if dp.GetName() == "seckill-app" {
//			dp.Spec.Template.Spec.PriorityClassName = "high-priority"
//			dp.Spec.Template.Spec.Affinity.PodAffinity = &corev1.PodAffinity{
//				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
//					{
//						Weight: 100,
//						PodAffinityTerm: corev1.PodAffinityTerm{
//							TopologyKey: "kubernetes.io/hostname",
//							LabelSelector: &metav1.LabelSelector{
//								MatchLabels: map[string]string{
//									"app": "seckill-app",
//								},
//							},
//						},
//					},
//				},
//			}
//		} else {
//			return admission.Allowed("no target deployment")
//		}
//		marshaledDeploy, err := json.Marshal(dp)
//		if err != nil {
//			logger.Info(err.Error())
//			return admission.Errored(http.StatusInternalServerError, err)
//		}
//		logger.Info("end to handle dp: ", dp.GetName())
//
//		return admission.PatchResponseFromRaw(req.Object.Raw, marshaledDeploy)
//	}
//	return admission.Allowed("No Resource Found")
//}

var (
	redisNode     int
	appNode       int
	BossNode      = ""
	resourceOn    bool
	cmdOn         bool
	appRequestCPU string
	appRequestMEM string
	appLimitCPU   string
	appLimitMEM   string

	workerRequestCPU string
	workerRequestMEM string
	workerLimitCPU   string
	workerLimitMEM   string
	xms              string
	xmx              string
	xmn              string
)

func init() {

	if num := os.Getenv("ResourceOn"); num != "" {
		resourceOn = true
	} else {
		resourceOn = false
	}

	if num := os.Getenv("CmdOn"); num != "" {
		cmdOn = true
	} else {
		cmdOn = false
	}

	if num := os.Getenv("RedisNode"); num != "" {
		redisNode, _ = strconv.Atoi(num)
	} else {
		redisNode = -1
	}

	if num := os.Getenv("AppNode"); num != "" {
		appNode, _ = strconv.Atoi(num)
	} else {
		appNode = -1
	}

	if num := os.Getenv("XMS"); num != "" {
		xms = num
	} else {
		xms = "4096m"
	}

	if num := os.Getenv("XMX"); num != "" {
		xmx = num
	} else {
		xmx = "4096m"
	}

	if num := os.Getenv("XMN"); num != "" {
		xmn = num
	} else {
		xmn = "3072m"
	}

	if num := os.Getenv("RequestCPU"); num != "" {
		appRequestCPU = num
	} else {
		appRequestCPU = "4"
	}

	if num := os.Getenv("RequestMEM"); num != "" {
		appRequestMEM = num
	} else {
		appRequestMEM = "16Gi"
	}

	if num := os.Getenv("LimitCPU"); num != "" {
		appLimitCPU = num
	} else {
		appLimitCPU = "4"
	}

	if num := os.Getenv("LimitMEM"); num != "" {
		appLimitMEM = num
	} else {
		appLimitMEM = "16Gi"
	}

	if node := os.Getenv("BossNode"); node != "" {
		BossNode = node
	}
}

type WebHookDeploymentHandler struct {
	Client  client.Client
	Decoder *admission.Decoder
}

func (handler *WebHookDeploymentHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	logger.Info("Req Method: ", req.Operation)
	dp := &appsv1.Deployment{}
	if err := handler.Decoder.Decode(req, dp); err == nil {
		logger.Info("start to handle dp: ", dp.GetName())
		if dp.GetName() == "seckill-app" {
			logger.Info("start to handle dp: ", dp.GetName())
			dp.Spec.Template.Spec.PriorityClassName = "high-priority"
			if resourceOn {
				containers := dp.Spec.Template.Spec.Containers
				container := corev1.Container{}
				index := 0
				if containers[0].Name == "mysql" {
					container = containers[0]
					index = 0
				} else {
					container = containers[1]
					index = 1
				}
				container.Resources = corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("6"),
						corev1.ResourceMemory: resource.MustParse("64Gi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("6"),
						corev1.ResourceMemory: resource.MustParse("64Gi"),
					},
				}
				if cmdOn {
					container.Command = []string{
						"java",
						"-server",
						"-Xms32768m",
						"-Xmx32768m",
						"-Xmn16384m",
						"-Xss1024k",
						"-jar",
						"/opt/seckill.jar",
						"--spring.profiles.active=container",
						"--server.port=30080",
					}
				}
				dp.Spec.Template.Spec.Containers[index] = container
			}
			dp.Spec.Template.Spec.Affinity = &corev1.Affinity{
				PodAntiAffinity: &corev1.PodAntiAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
						{
							Weight: 100,
							PodAffinityTerm: corev1.PodAffinityTerm{
								TopologyKey: "kubernetes.io/hostname",
								LabelSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"app": "seckill-app",
									},
								},
							},
						},
					},
				},
			}
		} else if dp.GetName() == "stress-boss" {
			dp.Spec.Template.Spec.PriorityClassName = "low-priority"
			if BossNode != "" {
				dp.Spec.Template.Spec.NodeName = BossNode
			}
		} else if dp.GetName() == "stress-worker" {
			dp.Spec.Template.Spec.PriorityClassName = "low-priority"
			if resourceOn {
				containers := dp.Spec.Template.Spec.Containers
				container := corev1.Container{}
				index := 0
				if containers[0].Name == "stress-worker" {
					container = containers[0]
					index = 0
				} else {
					container = containers[1]
					index = 1
				}
				container.Resources = corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("1"),
						corev1.ResourceMemory: resource.MustParse("8Gi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("1"),
						corev1.ResourceMemory: resource.MustParse("8Gi"),
					},
				}
				if cmdOn {
					container.Command = []string{
						"java",
						"-server",
						"-Xms4096m",
						"-Xmx4096m",
						"-Xmn2048m",
						"-Xss512k",
						"-jar",
						"/opt/stress-worker.jar",
						"--spring.profiles.active=container",
						"--server.port=30070",
					}
				}
				dp.Spec.Template.Spec.Containers[index] = container
			}
			dp.Spec.Template.Spec.Affinity = &corev1.Affinity{
				PodAntiAffinity: &corev1.PodAntiAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
						{
							Weight: 100,
							PodAffinityTerm: corev1.PodAffinityTerm{
								TopologyKey: "kubernetes.io/hostname",
								LabelSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"io.kompose.service": "stress-worker",
									},
								},
							},
						},
					},
				},
			}
		} else if dp.GetName() == "seckill-mysql" {
			dp.Spec.Template.Spec.PriorityClassName = "low-priority"

			//containers := dp.Spec.Template.Spec.Containers
			//container := corev1.Container{}
			//index := 0
			//if containers[0].Name == "mysql" {
			//	container = containers[0]
			//	index = 0
			//} else {
			//	container = containers[1]
			//	index = 1
			//}
			//container.Resources = corev1.ResourceRequirements{
			//	Limits: corev1.ResourceList{
			//		corev1.ResourceCPU:    resource.MustParse("2"),
			//		corev1.ResourceMemory: resource.MustParse("2Gi"),
			//	},
			//	Requests: corev1.ResourceList{
			//		corev1.ResourceCPU:    resource.MustParse("2"),
			//		corev1.ResourceMemory: resource.MustParse("2Gi"),
			//	},
			//}
			//dp.Spec.Template.Spec.Containers[index] = container
		} else if dp.GetName() == "seckill-statis" {
			dp.Spec.Template.Spec.PriorityClassName = "low-priority"

			//containers := dp.Spec.Template.Spec.Containers
			//container := corev1.Container{}
			//index := 0
			//if containers[0].Name == "mysql" {
			//	container = containers[0]
			//	index = 0
			//} else {
			//	container = containers[1]
			//	index = 1
			//}
			//container.Resources = corev1.ResourceRequirements{
			//	Limits: corev1.ResourceList{
			//		corev1.ResourceCPU:    resource.MustParse("1"),
			//		corev1.ResourceMemory: resource.MustParse("2Gi"),
			//	},
			//	Requests: corev1.ResourceList{
			//		corev1.ResourceCPU:    resource.MustParse("1"),
			//		corev1.ResourceMemory: resource.MustParse("2Gi"),
			//	},
			//}
			//dp.Spec.Template.Spec.Containers[index] = container
		} else {
			dp.Spec.Template.Spec.PriorityClassName = "low-priority"
			logger.Info("no target deployment")
		}
		marshaledDeploy, err := json.Marshal(dp)
		if err != nil {
			logger.Info(err.Error())
			return admission.Errored(http.StatusInternalServerError, err)
		}
		logger.Info("end to handle dp: ", dp.GetName())

		return admission.PatchResponseFromRaw(req.Object.Raw, marshaledDeploy)
	}
	return admission.Allowed("No Resource Found")
}

type WebHookStatefulSetHandler struct {
	Client  client.Client
	Decoder *admission.Decoder
}

func (handler *WebHookStatefulSetHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	logger.Info("Req Method: ", req.Operation)
	sts := &appsv1.StatefulSet{}
	if err := handler.Decoder.Decode(req, sts); err == nil {
		if sts.GetName() == "seckill-master" {
			logger.Info("start to handle sts: ", sts.GetName())
			sts.Spec.Template.Spec.PriorityClassName = "system-cluster-critical"
			//sts.Spec.Template.Spec.NodeName = "worker5"
			if resourceOn {
				containers := sts.Spec.Template.Spec.Containers
				container := corev1.Container{}
				index := 0
				if containers[0].Name == "bobft-seckill" {
					container = containers[0]
					index = 0
				} else {
					container = containers[1]
					index = 1
				}
				container.Resources = corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("4"),
						corev1.ResourceMemory: resource.MustParse("8Gi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("4"),
						corev1.ResourceMemory: resource.MustParse("8Gi"),
					},
				}
				sts.Spec.Template.Spec.Containers[index] = container
			}
		} else if sts.GetName() == "seckill-slave" {
			logger.Info("start to handle sts: ", sts.GetName())
			sts.Spec.Template.Spec.PriorityClassName = "high-priority"
			//sts.Spec.Template.Spec.NodeName = "worker5"
			if resourceOn {
				containers := sts.Spec.Template.Spec.Containers
				container := corev1.Container{}
				index := 0
				if containers[0].Name == "bobft-seckill" {
					container = containers[0]
					index = 0
				} else {
					container = containers[1]
					index = 1
				}
				container.Resources = corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("4"),
						corev1.ResourceMemory: resource.MustParse("8Gi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("4"),
						corev1.ResourceMemory: resource.MustParse("8Gi"),
					},
				}
				sts.Spec.Template.Spec.Containers[index] = container
			}
		} else {
			logger.Info("no target StatefulSet")
			return admission.Allowed("no target StatefulSet")
		}
		if redisNode != -1 {
			sts.Spec.Template.Spec.NodeName = NodeArray[redisNode]
		}
		//logger.Info(sts.Spec.Template.Spec.Containers[0].Resources.Requests)
		//logger.Info(sts.Spec.Template.Spec.Containers[0].Resources.Limits)
		marshaledSts, err := json.Marshal(sts)
		if err != nil {
			logger.Info(err.Error())
			return admission.Errored(http.StatusInternalServerError, err)
		}
		logger.Info("end to handle sts: ", sts.GetName())
		return admission.PatchResponseFromRaw(req.Object.Raw, marshaledSts)
	}
	return admission.Allowed("No Resource Found")
}
