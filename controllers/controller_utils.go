package controllers

import (
	"reflect"

	operatorsv2 "convect.ai/notebook-crd/api/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func generateStatefulSet(instance *operatorsv2.Jupyter) *appsv1.StatefulSet {
	replicas := int32(1)

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"statefulset": instance.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"statefulset":   instance.Name,
						"notebook-name": instance.Name,
					},
				},
				Spec: instance.Spec.Template.Spec,
			},
		},
	}
	// copy all of the Notebook labels to the pod including poddefault related labels
	labels := &statefulSet.Spec.Template.ObjectMeta.Labels
	for k, v := range instance.ObjectMeta.Labels {
		(*labels)[k] = v
	}

	podSpec := &statefulSet.Spec.Template.Spec
	container := &podSpec.Containers[0]

	if container.WorkingDir == "" {
		container.WorkingDir = "/home/jovyan"
	}

	if container.Ports == nil {
		container.Ports = []corev1.ContainerPort{
			{
				ContainerPort: 8888,
				Protocol:      "TCP",
				Name:          "notebook-port",
			},
		}
	}

	return statefulSet
}

func generateService(instance *operatorsv2.Jupyter) *corev1.Service {
	port := 8888

	containerPorts := instance.Spec.Template.Spec.Containers[0].Ports

	if containerPorts != nil {
		port = int(containerPorts[0].ContainerPort)
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: "ClusterIP",
			Selector: map[string]string{
				"statefulset": instance.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http-" + instance.Name,
					Port:       80,
					TargetPort: intstr.FromInt(port),
					Protocol:   "TCP",
				},
			},
		},
	}

	return svc

}

// CopyStatefulSetFields copies the owned fields from one StatefulSet to another
// Returns true if the fields copied from don't match to.
func copyStatefulSetFields(from, to *appsv1.StatefulSet) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			requireUpdate = true
		}
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			requireUpdate = true
		}
	}
	to.Annotations = from.Annotations

	if from.Spec.Replicas != to.Spec.Replicas {
		to.Spec.Replicas = from.Spec.Replicas
		requireUpdate = true
	}

	if !reflect.DeepEqual(to.Spec.Template.Spec, from.Spec.Template.Spec) {
		requireUpdate = true
	}
	to.Spec.Template.Spec = from.Spec.Template.Spec

	return requireUpdate
}

// CopyServiceFields copies the owned fields from one Service to another
func copyServiceFields(from, to *corev1.Service) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			requireUpdate = true
		}
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			requireUpdate = true
		}
	}
	to.Annotations = from.Annotations

	// Don't copy the entire Spec, because we can't overwrite the clusterIp field

	if !reflect.DeepEqual(to.Spec.Selector, from.Spec.Selector) {
		requireUpdate = true
	}
	to.Spec.Selector = from.Spec.Selector

	if !reflect.DeepEqual(to.Spec.Ports, from.Spec.Ports) {
		requireUpdate = true
	}
	to.Spec.Ports = from.Spec.Ports

	return requireUpdate
}
