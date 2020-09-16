package workload

import (
	"encoding/json"

	v1 "k8s.io/api/admission/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Scale struct {
	request *v1.AdmissionRequest
}

func NewScale(request *v1.AdmissionRequest) Operator {
	return &Scale{
		request: request,
	}
}

func (d *Scale) Create() error {
	return nil
}

func (d *Scale) Update() error {
	scale, err := decodeToScale(d.request.Object.Raw)
	if err != nil {
		return err
	}

	// do something
	_ = scale

	return nil
}

func (d *Scale) Delete() error {
	return nil
}

// CustomScale define scale object for json decode
// for dealing with different apiVersion like that
// `kubectl scale --replicas=3 deploy.apps mydeploy` which the Scale.Object.Status.Selector is a map type
// and `kubectl scale --replicas=3 deploy.extensions mydeploy` which the Scale.Object.Status.Selector is a string type
type CustomScale struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spect             autoscalingv1.ScaleSpec `json:"spec,omitempty"`
}

func decodeToScale(raw []byte) (*CustomScale, error) {
	var scale CustomScale
	if err := json.Unmarshal(raw, &scale); err != nil {
		return nil, err
	}

	return &scale, nil
}
