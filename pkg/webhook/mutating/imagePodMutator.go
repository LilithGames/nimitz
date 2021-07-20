package mutating

import (
	"context"
	"github.com/getsentry/sentry-go"
	"log"
	"regexp"
	"solarland/infra/annunciation/nimitz/pkg/config"
	"sync"

	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var imageOnce sync.Once
var imageMap []interface{}

// ImagePodMutator mutator interceptor to help match images address
func ImagePodMutator(_ context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
	imageOnce.Do(func() {
		imageMap, _ = config.ImageMap()
	})

	pod, ok := obj.(*corev1.Pod)
	if !ok {
		// If not a pod just continue the mutation chain(if there is one) and don't do nothing.
		return &kwhmutating.MutatorResult{}, nil
	}

	// Mutate our object with the required annotations.
	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string)
	}

	// skip static pod
	for k := range pod.Annotations {
		if k == "kubernetes.io/config.mirror" {
			log.Printf("[route.Mutating] /mutating: pod %s has kubernetes.io/config.mirror annotation, skip image rename\n", pod.Name)
			return &kwhmutating.MutatorResult{
				MutatedObject: pod,
			}, nil
		}
	}

	// mutate image source by imageMap rules
	for i, c := range pod.Spec.Containers {
		for _, v := range imageMap {
			val := v.([]string)
			re, err := regexp.Compile(val[0])
			//reg, err := regexp.Compile(pattern.(string))
			if err != nil {
				sentry.CaptureException(err)
				log.Fatalf("regexp Compile err: %s", err)
			}
			match, err := regexp.MatchString(val[0], c.Image)
			if err != nil {
				sentry.CaptureException(err)
				log.Fatalf("regexp matchString err: %s", err)
			}
			if match {
				// add annotations and update image
				pod.Annotations["mutated"] = "true"
				pod.Annotations["mutator"] = "imageMutator"
				newName := re.ReplaceAllString(c.Image, val[1])
				pod.Spec.Containers[i].Image = newName
				log.Printf("[route.Mutating] /mutating: pod %s has been update image address from %s to %s \n", pod.Name, c.Image, newName)
				break
			}
		}
	}

	return &kwhmutating.MutatorResult{
		MutatedObject: pod,
	}, nil
}
