package handlers

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

type WebHookHandler struct {
	Client  client.Client
	Decoder *admission.Decoder
}

func (handler *WebHookHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	sts := &appsv1.StatefulSet{}
	if err := handler.Decoder.Decode(req, sts); err == nil {
		logger.Info("start to handle sts: ", sts.GetName())
		if sts.GetName() == "seckill-master" {
			sts.Spec.Template.Spec.PriorityClassName = "system-cluster-critical"
			for _, container := range sts.Spec.Template.Spec.Containers {
				if container.Name != "bobft-seckill" {
					continue
				}
				container.Resources.Requests = corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("4"),
					corev1.ResourceMemory: resource.MustParse("8Gi"),
				}
				container.Resources.Limits = corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("4"),
					corev1.ResourceMemory: resource.MustParse("8Gi"),
				}
			}
		} else if sts.GetName() == "seckill-slave" {
			sts.Spec.Template.Spec.PriorityClassName = "high-priority"
			for _, container := range sts.Spec.Template.Spec.Containers {
				if container.Name != "bobft-seckill" {
					continue
				}
				container.Resources.Requests = corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("4"),
					corev1.ResourceMemory: resource.MustParse("8Gi"),
				}
				container.Resources.Limits = corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("4"),
					corev1.ResourceMemory: resource.MustParse("8Gi"),
				}
			}
		} else {
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

	dp := &appsv1.Deployment{}
	if err := handler.Decoder.Decode(req, dp); err == nil {
		logger.Info("start to handle dp: ", dp.GetName())
		if dp.GetName() == "stress-boss" {
			dp.Spec.Template.Spec.PriorityClassName = "low-priority"
		} else if dp.GetName() == "stress-worker" {
			dp.Spec.Template.Spec.PriorityClassName = "low-priority"
			dp.Spec.Template.Spec.Affinity.PodAffinity = &corev1.PodAffinity{
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
			}
		} else if dp.GetName() == "seckill-app" {
			dp.Spec.Template.Spec.PriorityClassName = "high-priority"
			dp.Spec.Template.Spec.Affinity.PodAffinity = &corev1.PodAffinity{
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
			}
		} else {
			return admission.Allowed("no target deployment")
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

type WebHookDeploymentHandler struct {
	Client  client.Client
	Decoder *admission.Decoder
}

func (handler *WebHookDeploymentHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	logger.Info("Req Method: ", req.Operation)
	dp := &appsv1.Deployment{}
	if err := handler.Decoder.Decode(req, dp); err == nil {
		if dp.GetName() == "seckill-app" {
			logger.Info("start to handle dp: ", dp.GetName())
			dp.Spec.Template.Spec.PriorityClassName = "high-priority"
		} else if dp.GetName() == "stress-boss" {
			logger.Info("start to handle dp: ", dp.GetName())
			dp.Spec.Template.Spec.PriorityClassName = "low-priority"
		} else if dp.GetName() == "stress-worker" {
			logger.Info("start to handle dp: ", dp.GetName())
			dp.Spec.Template.Spec.PriorityClassName = "low-priority"
			//dp.Spec.Template.Spec.Affinity = &corev1.Affinity{
			//	PodAntiAffinity: &corev1.PodAntiAffinity{
			//		PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
			//			{
			//				Weight: 100,
			//				PodAffinityTerm: corev1.PodAffinityTerm{
			//					TopologyKey: "kubernetes.io/hostname",
			//					LabelSelector: &metav1.LabelSelector{
			//						MatchLabels: map[string]string{
			//							"io.kompose.service": "stress-worker",
			//						},
			//					},
			//				},
			//			},
			//		},
			//	},
			//}
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
			for _, container := range sts.Spec.Template.Spec.Containers {
				if container.Name != "bobft-seckill" {
					continue
				}
				container.ReadinessProbe = &corev1.Probe{
					Handler: corev1.Handler{
						TCPSocket: &corev1.TCPSocketAction{
							Port: intstr.FromInt(6379),
						},
					},
					InitialDelaySeconds: 0,
					TimeoutSeconds:      1,
					PeriodSeconds:       3,
					SuccessThreshold:    1,
					FailureThreshold:    3,
				}
				//container.Resources = corev1.ResourceRequirements{
				//	Limits: corev1.ResourceList{
				//		corev1.ResourceCPU:    resource.MustParse("4"),
				//		corev1.ResourceMemory: resource.MustParse("8Gi"),
				//	},
				//	Requests: corev1.ResourceList{
				//		corev1.ResourceCPU:    resource.MustParse("4"),
				//		corev1.ResourceMemory: resource.MustParse("8Gi"),
				//	},
				//}
			}
		} else if sts.GetName() == "seckill-slave" {
			logger.Info("start to handle sts: ", sts.GetName())
			sts.Spec.Template.Spec.PriorityClassName = "high-priority"
			for _, container := range sts.Spec.Template.Spec.Containers {
				if container.Name != "bobft-seckill" {
					continue
				}
				container.ReadinessProbe = &corev1.Probe{
					Handler: corev1.Handler{
						TCPSocket: &corev1.TCPSocketAction{
							Port: intstr.FromInt(6379),
						},
					},
					InitialDelaySeconds: 0,
					TimeoutSeconds:      1,
					PeriodSeconds:       3,
					SuccessThreshold:    1,
					FailureThreshold:    3,
				}
				//container.Resources = corev1.ResourceRequirements{
				//	Limits: corev1.ResourceList{
				//		corev1.ResourceCPU:    resource.MustParse("4"),
				//		corev1.ResourceMemory: resource.MustParse("8Gi"),
				//	},
				//	Requests: corev1.ResourceList{
				//		corev1.ResourceCPU:    resource.MustParse("4"),
				//		corev1.ResourceMemory: resource.MustParse("8Gi"),
				//	},
				//}
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

type NodeSelectorHandler struct {
	Client  client.Client
	Decoder *admission.Decoder
}

func (handler *NodeSelectorHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	logger.Info("start!!!!!!!!")
	resourceName, resourceObject := resourceType(req.Name)
	switch resourceName {
	case "oltp":
		deployment, _ := resourceObject.(*appsv1.Deployment)

		err := handler.Decoder.Decode(req, deployment)
		if err != nil {
			logger.Info(err.Error())
			return admission.Errored(http.StatusBadRequest, err)
		}
		msg, success := oltpHandler(deployment)
		logger.Info(msg)
		if !success {
			admission.Allowed(msg)
		} else {
			marshaledDeploy, err := json.Marshal(deployment)
			if err != nil {
				logger.Info(err.Error())
				return admission.Errored(http.StatusInternalServerError, err)
			}
			logger.Info("end!!!!!!!!")

			return admission.PatchResponseFromRaw(req.Object.Raw, marshaledDeploy)
		}
	case "olap":
		return admission.Allowed("olap no handler")
	case "seckill-sts":
		sts, _ := resourceObject.(*appsv1.StatefulSet)
		logger.Info("start to handle: ", sts.Name)

		err := handler.Decoder.Decode(req, sts)
		if err != nil {
			logger.Info(err.Error())
			return admission.Errored(http.StatusBadRequest, err)
		}
		msg, success := seckillRedisHandler(sts)
		logger.Info(msg)
		if !success {
			admission.Allowed(msg)
		} else {
			marshaledSts, err := json.Marshal(sts)
			if err != nil {
				logger.Info(err.Error())
				return admission.Errored(http.StatusInternalServerError, err)
			}

			response := admission.PatchResponseFromRaw(req.Object.Raw, marshaledSts)
			logger.Info(response.Result)
			logger.Info("end!!!!!!!!")
			return response
		}
	case "seckill-app":
		deployment, _ := resourceObject.(*appsv1.Deployment)
		logger.Info("start to handle: ", deployment.Name)
		err := handler.Decoder.Decode(req, deployment)
		if err != nil {
			logger.Info(err.Error())
			return admission.Errored(http.StatusBadRequest, err)
		}
		switch deployment.Name {
		case "stress-boss":
			msg, success := seckillStressBossHandler(deployment)
			logger.Info(msg)
			if !success {
				admission.Allowed(msg)
			} else {
				marshaledDeploy, err := json.Marshal(deployment)
				if err != nil {
					logger.Info(err.Error())
					return admission.Errored(http.StatusInternalServerError, err)
				}
				response := admission.PatchResponseFromRaw(req.Object.Raw, marshaledDeploy)
				logger.Info(response.Result)
				logger.Info("end!!!!!!!!")
				return response
			}
		case "stress-worker":
			msg, success := seckillStressWorkerHandler(deployment)
			logger.Info(msg)
			if !success {
				admission.Allowed(msg)
			} else {
				marshaledDeploy, err := json.Marshal(deployment)
				if err != nil {
					logger.Info(err.Error())
					return admission.Errored(http.StatusInternalServerError, err)
				}
				response := admission.PatchResponseFromRaw(req.Object.Raw, marshaledDeploy)
				logger.Info(response.Result)
				logger.Info("end!!!!!!!!")
				return response
			}
		case "seckill-app":
			msg, success := seckillSeckillAppHandler(deployment)
			logger.Info(msg)
			if !success {
				admission.Allowed(msg)
			} else {
				marshaledDeploy, err := json.Marshal(deployment)
				if err != nil {
					logger.Info(err.Error())
					return admission.Errored(http.StatusInternalServerError, err)
				}
				response := admission.PatchResponseFromRaw(req.Object.Raw, marshaledDeploy)
				logger.Info(response.Result)
				logger.Info("end!!!!!!!!")
				return response
			}
		}
	}
	return admission.Allowed("no resource found")
}

func resourceType(resourceName string) (name string, resourceObject client.Object) {
	if resourceName == "seckill-slave" || resourceName == "seckill-master" {
		resourceObject = &appsv1.StatefulSet{}
		return "seckill-sts", resourceObject
	} else if resourceName == "seckill-app" || resourceName == "stress-worker" || resourceName == "stress-boss" {
		resourceObject = &appsv1.Deployment{}
		return "seckill-dp", resourceObject
	} else if strings.Contains(resourceName, "ftdb") && strings.Contains(resourceName, "bobft") {
		return "oltp", &appsv1.Deployment{}
	} else if strings.Contains(resourceName, "flink-taskmanager") {
		return "olap", &appsv1.Deployment{}
	} else {
		return "", nil
	}
}
