package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
)

type admitFunc func(v1.AdmissionReview) *v1.AdmissionResponse

func serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	var body []byte
	if r.Body == nil {
		glog.Error("body is nil")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body = data

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Error("content type should be application/json")
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		return
	}

	glog.Infof("get data:%s", body)
	requestAdmissionReview := v1.AdmissionReview{}
	responseAdmissionReview := v1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
	}

	var admissionResponse *v1.AdmissionResponse
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &requestAdmissionReview); err != nil {
		glog.Error(err)
		admissionResponse = toAdmissionResponse(err)
	} else {
		admissionResponse = admit(requestAdmissionReview)
	}

	if admissionResponse != nil {
		responseAdmissionReview.Response = admissionResponse
		if requestAdmissionReview.Request != nil {
			responseAdmissionReview.Response.UID = requestAdmissionReview.Request.UID
		}
	}

	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		glog.Error(err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}

	glog.Infof("response admission review:%s", respBytes)
	if _, err := w.Write(respBytes); err != nil {
		glog.Error(err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

func toAdmissionResponse(err error) *v1.AdmissionResponse {
	return &v1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}
