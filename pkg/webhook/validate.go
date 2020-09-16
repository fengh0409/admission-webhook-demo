package webhook

import (
	"admission-webhook-demo/pkg/workload"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	v1 "k8s.io/api/admission/v1"
)

func Validate(w http.ResponseWriter, r *http.Request) {
	serve(w, r, serveValidate)
}

func serveValidate(ar v1.AdmissionReview) *v1.AdmissionResponse {
	glog.Infof("admit %v %v", ar.Request.Kind.Kind, ar.Request.Operation)

	var err error
	switch ar.Request.Kind.Kind {
	case "Deployment":
		err = workload.GetResource(workload.NewDeployment, ar.Request)
	case "StatefulSet":
	case "DaemonSet":
	case "CronJob":
	case "Job":
	case "Scale":
		err = workload.GetResource(workload.NewScale, ar.Request)
	default:
		err = fmt.Errorf("unsupported workload type: %v", ar.Request.Kind.Kind)
	}

	if err != nil {
		glog.Error(err)
		return toAdmissionResponse(err)
	}

	return &v1.AdmissionResponse{
		Allowed: true,
	}
}
