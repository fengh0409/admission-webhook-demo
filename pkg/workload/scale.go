package workload

import (
	v1 "k8s.io/api/admission/v1"
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
	return nil
}

func (d *Scale) Delete() error {
	return nil
}
