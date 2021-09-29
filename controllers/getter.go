package controllers

import (
	"context"
	api "github.com/pratise/mongoshake-kubernetes-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func mongoshakeLables(cr *api.MongoShake) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "mongodshake-job",
		"app.kubernetes.io/instance":   cr.Name,
		"app.kubernetes.io/component":  "mongoshake",
		"app.kubernetes.io/managed-by": "mongoshake-operator",
	}
}

func (r *MongoShakeReconciler) getJobPods(cr *api.MongoShake) (corev1.PodList, error) {
	job := corev1.PodList{}
	err := r.Client.List(context.TODO(), &job, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(mongoshakeLables(cr)),
		Namespace:     cr.Namespace,
	})
	return job, err
}
