// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CheckPointSpec) DeepCopyInto(out *CheckPointSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheckPointSpec.
func (in *CheckPointSpec) DeepCopy() *CheckPointSpec {
	if in == nil {
		return nil
	}
	out := new(CheckPointSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterCondition) DeepCopyInto(out *ClusterCondition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterCondition.
func (in *ClusterCondition) DeepCopy() *ClusterCondition {
	if in == nil {
		return nil
	}
	out := new(ClusterCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CollectorSpec) DeepCopyInto(out *CollectorSpec) {
	*out = *in
	if in.Log != nil {
		in, out := &in.Log, &out.Log
		*out = new(LogSpec)
		**out = **in
	}
	if in.Mongo != nil {
		in, out := &in.Mongo, &out.Mongo
		*out = new(MongoSpec)
		**out = **in
	}
	if in.Tunnel != nil {
		in, out := &in.Tunnel, &out.Tunnel
		*out = new(TunnelSpec)
		**out = **in
	}
	if in.Filter != nil {
		in, out := &in.Filter, &out.Filter
		*out = new(FilterSpec)
		**out = **in
	}
	if in.CheckPoint != nil {
		in, out := &in.CheckPoint, &out.CheckPoint
		*out = new(CheckPointSpec)
		**out = **in
	}
	if in.FullSync != nil {
		in, out := &in.FullSync, &out.FullSync
		*out = new(FullSyncSpec)
		**out = **in
	}
	if in.IncrSync != nil {
		in, out := &in.IncrSync, &out.IncrSync
		*out = new(IncrSyncSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CollectorSpec.
func (in *CollectorSpec) DeepCopy() *CollectorSpec {
	if in == nil {
		return nil
	}
	out := new(CollectorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FilterSpec) DeepCopyInto(out *FilterSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FilterSpec.
func (in *FilterSpec) DeepCopy() *FilterSpec {
	if in == nil {
		return nil
	}
	out := new(FilterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FullSyncSpec) DeepCopyInto(out *FullSyncSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FullSyncSpec.
func (in *FullSyncSpec) DeepCopy() *FullSyncSpec {
	if in == nil {
		return nil
	}
	out := new(FullSyncSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IncrSyncSpec) DeepCopyInto(out *IncrSyncSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IncrSyncSpec.
func (in *IncrSyncSpec) DeepCopy() *IncrSyncSpec {
	if in == nil {
		return nil
	}
	out := new(IncrSyncSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobStatus) DeepCopyInto(out *JobStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobStatus.
func (in *JobStatus) DeepCopy() *JobStatus {
	if in == nil {
		return nil
	}
	out := new(JobStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LivenessProbeExtended) DeepCopyInto(out *LivenessProbeExtended) {
	*out = *in
	in.Probe.DeepCopyInto(&out.Probe)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LivenessProbeExtended.
func (in *LivenessProbeExtended) DeepCopy() *LivenessProbeExtended {
	if in == nil {
		return nil
	}
	out := new(LivenessProbeExtended)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogSpec) DeepCopyInto(out *LogSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogSpec.
func (in *LogSpec) DeepCopy() *LogSpec {
	if in == nil {
		return nil
	}
	out := new(LogSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoShake) DeepCopyInto(out *MongoShake) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoShake.
func (in *MongoShake) DeepCopy() *MongoShake {
	if in == nil {
		return nil
	}
	out := new(MongoShake)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MongoShake) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoShakeList) DeepCopyInto(out *MongoShakeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MongoShake, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoShakeList.
func (in *MongoShakeList) DeepCopy() *MongoShakeList {
	if in == nil {
		return nil
	}
	out := new(MongoShakeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MongoShakeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoShakeSpec) DeepCopyInto(out *MongoShakeSpec) {
	*out = *in
	if in.BackoffLimit != nil {
		in, out := &in.BackoffLimit, &out.BackoffLimit
		*out = new(int32)
		**out = **in
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(ResourcesSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.ReadinessProbe != nil {
		in, out := &in.ReadinessProbe, &out.ReadinessProbe
		*out = new(corev1.Probe)
		(*in).DeepCopyInto(*out)
	}
	if in.LivenessProbe != nil {
		in, out := &in.LivenessProbe, &out.LivenessProbe
		*out = new(LivenessProbeExtended)
		(*in).DeepCopyInto(*out)
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]corev1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeSpec != nil {
		in, out := &in.VolumeSpec, &out.VolumeSpec
		*out = new(VolumeSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Collector != nil {
		in, out := &in.Collector, &out.Collector
		*out = new(CollectorSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoShakeSpec.
func (in *MongoShakeSpec) DeepCopy() *MongoShakeSpec {
	if in == nil {
		return nil
	}
	out := new(MongoShakeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoShakeStatus) DeepCopyInto(out *MongoShakeStatus) {
	*out = *in
	if in.JobStatus != nil {
		in, out := &in.JobStatus, &out.JobStatus
		*out = new(JobStatus)
		**out = **in
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]ClusterCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoShakeStatus.
func (in *MongoShakeStatus) DeepCopy() *MongoShakeStatus {
	if in == nil {
		return nil
	}
	out := new(MongoShakeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoSpec) DeepCopyInto(out *MongoSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoSpec.
func (in *MongoSpec) DeepCopy() *MongoSpec {
	if in == nil {
		return nil
	}
	out := new(MongoSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSpecRequirements) DeepCopyInto(out *ResourceSpecRequirements) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSpecRequirements.
func (in *ResourceSpecRequirements) DeepCopy() *ResourceSpecRequirements {
	if in == nil {
		return nil
	}
	out := new(ResourceSpecRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourcesSpec) DeepCopyInto(out *ResourcesSpec) {
	*out = *in
	if in.Limits != nil {
		in, out := &in.Limits, &out.Limits
		*out = new(ResourceSpecRequirements)
		**out = **in
	}
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = new(ResourceSpecRequirements)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourcesSpec.
func (in *ResourcesSpec) DeepCopy() *ResourcesSpec {
	if in == nil {
		return nil
	}
	out := new(ResourcesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TunnelSpec) DeepCopyInto(out *TunnelSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TunnelSpec.
func (in *TunnelSpec) DeepCopy() *TunnelSpec {
	if in == nil {
		return nil
	}
	out := new(TunnelSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSpec) DeepCopyInto(out *VolumeSpec) {
	*out = *in
	if in.EmptyDir != nil {
		in, out := &in.EmptyDir, &out.EmptyDir
		*out = new(corev1.EmptyDirVolumeSource)
		(*in).DeepCopyInto(*out)
	}
	if in.HostPath != nil {
		in, out := &in.HostPath, &out.HostPath
		*out = new(corev1.HostPathVolumeSource)
		(*in).DeepCopyInto(*out)
	}
	if in.PersistentVolumeClaim != nil {
		in, out := &in.PersistentVolumeClaim, &out.PersistentVolumeClaim
		*out = new(corev1.PersistentVolumeClaimSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSpec.
func (in *VolumeSpec) DeepCopy() *VolumeSpec {
	if in == nil {
		return nil
	}
	out := new(VolumeSpec)
	in.DeepCopyInto(out)
	return out
}
