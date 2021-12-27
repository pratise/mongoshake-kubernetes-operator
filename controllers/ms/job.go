package ms

import (
	"crypto/md5"
	"fmt"
	"github.com/go-logr/logr"
	api "github.com/pratise/mongoshake-kubernetes-operator/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MongoshakeCustomConfigName(clusterName string) string {
	return fmt.Sprintf("%s-%s", MongoshakeConfigName, clusterName)
}
func MongoshakeCustomPersistentVolumeClaimLogName(clusterName string) string {
	return fmt.Sprintf("%s-%s", MongoshakeLogVolClaimName, clusterName)
}

func MongoshakeJob(cr *api.MongoShake) *batchv1.Job {
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
	}
}

func MongoshakeJobSpec(cr *api.MongoShake, log logr.Logger) (batchv1.JobSpec, error) {
	ls := map[string]string{
		"app.kubernetes.io/name":       "mongodshake-job",
		"app.kubernetes.io/instance":   cr.Name,
		"app.kubernetes.io/component":  "mongoshake",
		"app.kubernetes.io/managed-by": "mongoshake-operator",
	}

	container, err := mongoshakeContainer(cr)
	if err != nil {
		return batchv1.JobSpec{}, fmt.Errorf("failed to create container %v", err)
	}

	annotations := cr.Spec.Annotations
	if annotations == nil {
		annotations = make(map[string]string)
	}

	if c := cr.Spec.Configuration; c != "" {
		annotations["mongoshake/configuration-hash"] = fmt.Sprintf("%x", md5.Sum([]byte(c)))
	}

	manualSelector := true
	return batchv1.JobSpec{
		// 存活时间
		ActiveDeadlineSeconds: nil,
		// 重试次数
		BackoffLimit:   cr.Spec.BackoffLimit,
		ManualSelector: &manualSelector,
		Selector: &metav1.LabelSelector{
			MatchLabels: ls,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      ls,
				Annotations: annotations,
			},
			Spec: corev1.PodSpec{
				NodeSelector:  cr.Spec.NodeSelector,
				Tolerations:   cr.Spec.Tolerations,
				RestartPolicy: corev1.RestartPolicyNever,
				Containers:    []corev1.Container{container},
				Volumes:       volumes(cr),
			},
		},
	}, nil
}

func volumes(cr *api.MongoShake) []corev1.Volume {
	hostPathType := corev1.HostPathFile
	volumes := []corev1.Volume{
		{
			Name: MongoshakeConfigVolClaimName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: MongoshakeCustomConfigName(cr.Name),
					},
				},
			},
		},
		{
			Name: MongoshakeLogVolClaimName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: MongoshakeCustomPersistentVolumeClaimLogName(cr.Name),
					ReadOnly:  false,
				},
			},
		},
		{
			Name: timezone,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/usr/share/zoneinfo/Asia/Shanghai",
					Type: &hostPathType,
				},
			},
		},
	}
	return volumes
}

func mongoshakeContainer(cr *api.MongoShake) (corev1.Container, error) {
	resources, err := CreateResources(cr.Spec.Resources)
	if err != nil {
		return corev1.Container{}, fmt.Errorf("resource creation: %v", err)
	}

	volumes := []corev1.VolumeMount{
		{
			Name:      MongoshakeConfigVolClaimName,
			ReadOnly:  true,
			MountPath: MongoshakeConfigDir,
		},
		{
			Name:      timezone,
			MountPath: timezoneDir,
		},
		{
			Name:      MongoshakeLogVolClaimName,
			MountPath: MongoshakeLogDir,
		},
	}

	container := corev1.Container{
		Name:  "mongoshake",
		Image: cr.Spec.Image,
		// Command:    []string{"/bin/sh"},
		// Args: []string{"-c","`while true; do echo hello world; sleep 30; done`"},
		Command:    []string{"./collector.linux"},
		Args:       []string{"--conf", "./config/collector.conf", "--verbose", "1"},
		WorkingDir: WorkingDir,
		Ports: []corev1.ContainerPort{
			{
				Name:          FullSyncPortName,
				ContainerPort: 9101,
			}, {
				Name:          IncrSyncPortName,
				ContainerPort: 9100,
			},
		},
		Resources:    resources,
		VolumeMounts: volumes,
		// LivenessProbe:            nil,
		// ReadinessProbe:           nil,
		// Lifecycle:                nil,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
	}
	return container, nil
}
