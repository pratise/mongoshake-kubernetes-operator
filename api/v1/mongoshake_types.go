/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MongoShakeSpec defines the desired state of MongoShake
type MongoShakeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Image           string                 `json:"image,omitempty"`
	BackoffLimit    *int32                 `json:"backoffLimit,omitempty"`
	Pause           bool                   `json:"pause,omitempty"`
	ImagePullPolicy corev1.PullPolicy      `json:"imagePullPolicy,omitempty"`
	Resources       *ResourcesSpec         `json:"resources,omitempty"`
	ReadinessProbe  *corev1.Probe          `json:"readinessProbe,omitempty"`
	LivenessProbe   *LivenessProbeExtended `json:"livenessProbe,omitempty"`
	Configuration   string                 `json:"configuration,omitempty"`
	Annotations     map[string]string      `json:"annotations,omitempty"`
	NodeSelector    map[string]string      `json:"nodeSelector,omitempty"`
	VolumeSpec      *VolumeSpec            `json:"volumeSpec,omitempty"`
}

type JobStatus struct {
	Size    int      `json:"size"`
	Ready   int      `json:"ready"`
	Status  AppState `json:"status,omitempty"`
	Message string   `json:"message,omitempty"`
}

// MongoShakeStatus defines the observed state of MongoShake
type MongoShakeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	State      AppState           `json:"state,omitempty"`
	JobStatus  *JobStatus         `json:"job_status,omitempty"`
	Message    string             `json:"message,omitempty"`
	Conditions []ClusterCondition `json:"conditions,omitempty"`
	Host       string             `json:"host,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName="ms"
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`,description="The Actual status of Mongoshake"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// MongoShake is the Schema for the mongoshakes API
type MongoShake struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MongoShakeSpec   `json:"spec,omitempty"`
	Status MongoShakeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MongoShakeList contains a list of MongoShake
type MongoShakeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoShake `json:"items"`
}

type VolumeSpec struct {
	// EmptyDir represents a temporary directory that shares a pod's lifetime.
	EmptyDir *corev1.EmptyDirVolumeSource `json:"emptyDir,omitempty"`

	// HostPath represents a pre-existing file or directory on the host machine
	// that is directly exposed to the container.
	HostPath *corev1.HostPathVolumeSource `json:"hostPath,omitempty"`

	// PersistentVolumeClaim represents a reference to a PersistentVolumeClaim.
	// It has the highest level of precedence, followed by HostPath and
	// EmptyDir. And represents the PVC specification.
	PersistentVolumeClaim *corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaim,omitempty"`
}

type ResourceSpecRequirements struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type ResourcesSpec struct {
	Limits   *ResourceSpecRequirements `json:"limits,omitempty"`
	Requests *ResourceSpecRequirements `json:"requests,omitempty"`
}
type LivenessProbeExtended struct {
	corev1.Probe        `json:",inline"`
	StartupDelaySeconds int `json:"startupDelaySeconds,omitempty"`
}

type AppState string

const (
	AppStateInit     AppState = "initializing"
	AppStateRunning  AppState = "running"
	AppStateStopping AppState = "stopping"
	AppStatePaused   AppState = "paused"
	AppStateComplete AppState = "complete"
	AppStateError    AppState = "error"
)

func (l LivenessProbeExtended) CommandHas(flag string) bool {
	if l.Handler.Exec == nil {
		return false
	}

	for _, v := range l.Handler.Exec.Command {
		if v == flag {
			return true
		}
	}

	return false
}

type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

type ClusterCondition struct {
	Status             ConditionStatus `json:"status"`
	Type               AppState        `json:"type"`
	LastTransitionTime metav1.Time     `json:"lastTransitionTime,omitempty"`
	Reason             string          `json:"reason,omitempty"`
	Message            string          `json:"message,omitempty"`
}

const maxStatusesQuantity = 20

func (s *MongoShakeStatus) AddCondition(c ClusterCondition) {
	if len(s.Conditions) == 0 {
		s.Conditions = append(s.Conditions, c)
		return
	}

	if s.Conditions[len(s.Conditions)-1].Type != c.Type {
		s.Conditions = append(s.Conditions, c)
	}

	if len(s.Conditions) > maxStatusesQuantity {
		s.Conditions = s.Conditions[len(s.Conditions)-maxStatusesQuantity:]
	}
}

func init() {
	SchemeBuilder.Register(&MongoShake{}, &MongoShakeList{})
}
