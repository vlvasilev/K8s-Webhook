package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// replaceImage is a container for the mutation replacing container image
type replaceImage struct {
	Logger    logrus.FieldLogger
	PodConfig *Pod
}

// replaceImage implements the podMutator interface
var _ podMutator = (*replaceImage)(nil)

// Name returns the struct name
func (se replaceImage) Name() string {
	return "inject_env"
}

// Mutate returns a new mutated pod according to set containerImage rules
func (se replaceImage) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	se.Logger = se.Logger.WithField("mutation", se.Name())
	mpod := pod.DeepCopy()

	// containerImage := map[string]string{
	// 	"sapcc-operator": "columbus.common.repositories.cloud.sap/com.sap.edgelm/sapcc-operator:0.51.0-test.i330716-latest",
	// }
	if se.PodConfig == nil || se.PodConfig.ContainerReplacements == nil {
		se.Logger.Debugf("Pod NAME %s has no container replacements", mpod.Name)
		return mpod, nil
	}
	se.replaceImage(mpod, se.PodConfig.ContainerReplacements)

	return mpod, nil
}

func (se replaceImage) replaceImage(pod *corev1.Pod, containerImage map[string]string) {
	for containerName, image := range containerImage {
		for i, container := range pod.Spec.Containers {
			if container.Name == containerName {
				se.Logger.Debugf("pod image injected %s", image)
				pod.Spec.Containers[i].Image = image
			}
		}
		for i, initContainer := range pod.Spec.InitContainers {
			if initContainer.Name == containerName {
				se.Logger.Debugf("pod image injected %s", image)
				pod.Spec.InitContainers[i].Image = image
			}
		}
	}
}
