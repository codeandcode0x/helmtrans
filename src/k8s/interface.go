package k8s

import (
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//interface
type K8sClientInf interface{
	//service
	UnmarshalService(bytes []byte) corev1.Service
	GetServices(apps ...string) []*corev1.Service
	DeployService(svc corev1.Service) error
	//deployment
	UnmarshalDeployment(bytes []byte) appsv1.Deployment
	DeployDeployment(deploy appsv1.Deployment) error
	GetDeployments(apps ...string) []*appsv1.Deployment
	//configmap
	UnmarshalConfigMap(bytes []byte) corev1.ConfigMap
	DeployConfigMap(cm corev1.ConfigMap) error
	GetConfigMaps(cms ...string) []*corev1.ConfigMap
	//statefulset
	//pod
	GetPodsByLabel(app string) *corev1.PodList
	PodExecCommand(namespace, podName, command, containerName string) (string, string, error)

	//resource delete
	ResDelete(resourceType string, bytes []byte)
}