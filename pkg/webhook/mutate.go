package webhook

import (
	"admission-webhook-demo/pkg/config"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	admissionWebhookAnnotationInjectKey = "bazingafeng.com/inject"
)

var podResource = metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

var sidecarVolumeMounts = []corev1.VolumeMount{
	{
		Name:      "sidecar-volume",
		MountPath: "/opt/sidecar",
		ReadOnly:  true,
	},
}

func Mutate(w http.ResponseWriter, r *http.Request) {
	serve(w, r, serveInjection)
}

func serveInjection(ar v1.AdmissionReview) *v1.AdmissionResponse {
	glog.Info("admit injection")
	var req = ar.Request
	if req.Resource != podResource {
		err := fmt.Errorf("expect resource to be %s", podResource)
		glog.Error(err)
		return toAdmissionResponse(err)
	}

	if req.Operation != v1.Create {
		return &v1.AdmissionResponse{
			Allowed: true,
		}
	}

	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		glog.Error(err)
		return toAdmissionResponse(err)
	}

	// inject containers
	annotations := map[string]string{admissionWebhookAnnotationInjectKey: "injected"}
	patchBytes, err := createPatch(&pod, config.Conf.Sidecar, annotations)
	if err != nil {
		glog.Error(err)
		return toAdmissionResponse(err)
	}

	pt := v1.PatchTypeJSONPatch
	return &v1.AdmissionResponse{
		Allowed:   true,
		Patch:     patchBytes,
		PatchType: &pt,
	}
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func createPatch(pod *corev1.Pod, sidecarConfig *config.SidecarConfig, annotations map[string]string) ([]byte, error) {
	var patch []patchOperation

	patch = append(patch, addContainer(pod.Spec.Containers, sidecarConfig.Containers, "/spec/containers")...)
	patch = append(patch, addVolume(pod.Spec.Volumes, sidecarConfig.Volumes, "/spec/volumes")...)
	patch = append(patch, updateAnnotation(pod.Annotations, annotations)...)

	return json.Marshal(patch)
}

func addContainer(target, added []corev1.Container, basePath string) []patchOperation {
	var patch []patchOperation
	first := len(target) == 0
	var value interface{}
	for _, add := range added {
		value = add
		path := basePath
		if first {
			first = false
			value = []corev1.Container{add}
		} else {
			path = path + "/-"
		}
		patch = append(patch, patchOperation{
			Op:    "add",
			Path:  path,
			Value: value,
		})
	}

	patch = append(patch, addVolumeMounts(target)...)

	return patch
}

func addVolumeMounts(containers []corev1.Container) (patch []patchOperation) {
	basePath := "/spec/containers/%d/volumeMounts"
	for i := range containers {
		if len(containers[i].VolumeMounts) == 0 {
			path := fmt.Sprintf(basePath, i)
			patch = append(patch, patchOperation{
				Op:    "add",
				Path:  path,
				Value: sidecarVolumeMounts,
			})
		} else {
			path := fmt.Sprintf(basePath+"/-", i)
			for _, volumeMount := range sidecarVolumeMounts {
				patch = append(patch, patchOperation{
					Op:    "add",
					Path:  path,
					Value: volumeMount,
				})
			}
		}
	}

	return patch
}

func addVolume(target, added []corev1.Volume, basePath string) (patch []patchOperation) {
	first := len(target) == 0
	var value interface{}
	for _, add := range added {
		value = add
		path := basePath
		if first {
			first = false
			value = []corev1.Volume{add}
		} else {
			path = path + "/-"
		}
		patch = append(patch, patchOperation{
			Op:    "add",
			Path:  path,
			Value: value,
		})
	}
	return patch
}

func updateAnnotation(target map[string]string, added map[string]string) (patch []patchOperation) {
	for key, value := range added {
		if target == nil || target[key] == "" {
			target = map[string]string{}
			patch = append(patch, patchOperation{
				Op:   "add",
				Path: "/metadata/annotations",
				Value: map[string]string{
					key: value,
				},
			})
		} else {
			patch = append(patch, patchOperation{
				Op:    "replace",
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		}
	}
	return patch
}
