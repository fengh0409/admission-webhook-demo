package workload

import (
	"fmt"

	v1 "k8s.io/api/admission/v1"
)

type Operator interface {
	Create() error
	Update() error
	Delete() error
}

type ResourceFunc func(*v1.AdmissionRequest) Operator

func GetResource(newResource ResourceFunc, request *v1.AdmissionRequest) error {
	operator := newResource(request)

	var err error
	switch request.Operation {
	case v1.Create:
		err = operator.Create()
	case v1.Update:
		err = operator.Update()
	case v1.Delete:
		err = operator.Delete()
	default:
		err = fmt.Errorf("unsupported operation: %v", request.Operation)
	}

	if err != nil {
		return err
	}

	return nil
}
