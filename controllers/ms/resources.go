package ms

import (
	"fmt"
	api "github.com/pratise/mongoshake-kubernetes-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func CreateResources(r *api.ResourcesSpec) (rr corev1.ResourceRequirements, err error) {
	if r == nil {
		return rr, nil
	}
	if r.Requests != nil {
		rList, err := CreateResourcesList(r.Requests)
		if err != nil {
			return rr, nil
		}
		rr.Requests = rList
	}
	if r.Limits != nil {
		rList, err := CreateResourcesList(r.Limits)
		if err != nil {
			return rr, nil
		}
		rr.Limits = rList
	}
	return rr, nil
}

func CreateResourcesList(r *api.ResourceSpecRequirements) (rlist corev1.ResourceList, err error) {
	rlist = make(corev1.ResourceList)
	if len(r.CPU) > 0 {
		rlist[corev1.ResourceCPU], err = resource.ParseQuantity(r.CPU)
		if err != nil {
			return nil, fmt.Errorf("malformed CPU resources: %v", err)
		}
	}
	if len(r.Memory) > 0 {
		rlist[corev1.ResourceMemory], err = resource.ParseQuantity(r.Memory)
		if err != nil {
			return nil, fmt.Errorf("malformed Memory resources: %v", err)
		}
	}
	return rlist, nil
}
