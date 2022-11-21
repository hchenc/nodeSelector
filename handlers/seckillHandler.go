package handlers

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func seckillHandler(sts *appsv1.StatefulSet) (string, bool) {
	sts.Spec.Template.Spec.PriorityClassName = "high-priority"
	//sts.Spec.Template.Spec.Containers

	if redisNode != -1 {
		sts.Spec.Template.Spec.NodeName = NodeArray[redisNode]
		return fmt.Sprintf("sts %s scheduled to  %s", sts.Name, NodeArray[redisNode]), true
	}
	return "", true
}

func seckillRedisHandler(sts *appsv1.StatefulSet) (string, bool) {
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

	if redisNode != -1 {
		sts.Spec.Template.Spec.NodeName = NodeArray[redisNode]
		return fmt.Sprintf("sts %s scheduled to  %s", sts.Name, NodeArray[redisNode]), true
	}
	logger.Info(sts)
	return "", true
}

func seckillSeckillAppHandler(dp *appsv1.Deployment) (string, bool) {
	dp.Spec.Template.Spec.PriorityClassName = "medium-priority"
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
	logger.Info(dp)
	return "", true
}

func seckillStressBossHandler(dp *appsv1.Deployment) (string, bool) {
	logger.Info(dp.Name)
	dp.Spec.Template.Spec.PriorityClassName = "low-priority"
	logger.Info(dp)
	return "", true
}

func seckillStressWorkerHandler(dp *appsv1.Deployment) (string, bool) {
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
	logger.Info(dp)
	return "", true
}
