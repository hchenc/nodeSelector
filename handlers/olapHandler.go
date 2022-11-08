package handlers

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Flink = "flink-taskmanager-bobft"
)

func olapHandler(deployment *appsv1.Deployment) (string, bool) {
	logger.Info(deployment.Name)

	deployment.Spec.Template.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse("1"),
		corev1.ResourceMemory: resource.MustParse("2Gi"),
	}
	deployment.Spec.Template.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse("1"),
		corev1.ResourceMemory: resource.MustParse("2Gi"),
	}

	deployment.Spec.Template.Spec.Affinity.PodAffinity = &corev1.PodAffinity{
		PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
			{
				Weight: 100,
				PodAffinityTerm: corev1.PodAffinityTerm{
					TopologyKey: "kubernetes.io/hostname",
					LabelSelector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "baz",
								Operator: metav1.LabelSelectorOpIn,
								Values:   []string{"qux", "norf"},
							},
						},
						MatchLabels: map[string]string{
							"io.kompose.service": "stress-worker",
						},
					},
				},
			},
		},
	}
	return fmt.Sprintf("request: %v, limit: %v", deployment.Spec.Template.Spec.Containers[0].Resources.Requests, deployment.Spec.Template.Spec.Containers[0].Resources.Limits), true
}
