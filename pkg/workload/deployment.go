package workload

import (
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
)

type Deployment struct {
	request *v1.AdmissionRequest
}

func NewDeployment(request *v1.AdmissionRequest) Operator {
	return &Deployment{
		request: request,
	}
}

func (d *Deployment) Create() error {
	return nil
}

func (d *Deployment) Update() error {
	return nil
}

func (d *Deployment) Delete() error {
	deployment, err := decodeToDeployment(d.request.OldObject.Raw)
	if err != nil {
		return err
	}

	if deployment.Labels["can-delete"] == "not-allow" {
		return fmt.Errorf("can not delete")
	}

	return nil
}

func decodeToDeployment(raw []byte) (*appsv1.Deployment, error) {
	var deployment *appsv1.Deployment
	if err := json.Unmarshal(raw, &deployment); err != nil {
		return nil, err
	}

	return deployment, nil
}
