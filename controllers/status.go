package controllers

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	api "github.com/pratise/mongoshake-kubernetes-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// 更新cr状态
func (r *MongoShakeReconciler) updateStatus(cr *api.MongoShake, reconcileErr error, appState api.AppState) error {
	clusterCondition := api.ClusterCondition{
		Status:             api.ConditionTrue,
		Type:               api.AppStateInit,
		LastTransitionTime: metav1.NewTime(time.Now()),
	}
	if reconcileErr != nil {
		if cr.Status.State != api.AppStateError {
			clusterCondition = api.ClusterCondition{
				Status:             api.ConditionTrue,
				Type:               api.AppStateError,
				Message:            reconcileErr.Error(),
				Reason:             "ErrorReconcile",
				LastTransitionTime: metav1.NewTime(time.Now()),
			}
			cr.Status.AddCondition(clusterCondition)
		}

		cr.Status.Message = "Error: " + reconcileErr.Error()
		cr.Status.State = api.AppStateError

		return r.writeStatus(cr)
	}

	jobStatus, err := r.jobStatus(cr)
	if err != nil {
		return errors.Wrap(err, "get mongoshake status")
	}

	if jobStatus.Status != cr.Status.State {
		mongoshakeCondition := api.ClusterCondition{
			Type:               jobStatus.Status,
			Status:             api.ConditionTrue,
			LastTransitionTime: metav1.NewTime(time.Now()),
		}

		switch jobStatus.Status {
		case api.AppStateRunning:
			cr.Status.State = api.AppStateRunning
			mongoshakeCondition.Reason = "Mongoshake Running"
		case api.AppStateComplete:
			cr.Status.State = api.AppStateComplete
			mongoshakeCondition.Reason = "Mongoshake Complete"
		case api.AppStateStopping:
			cr.Status.State = api.AppStateStopping
			mongoshakeCondition.Reason = "Mongoshake Stopping"
		case api.AppStatePaused:
			cr.Status.State = api.AppStatePaused
			mongoshakeCondition.Reason = "Mongoshake Paused"
		}
		cr.Status.AddCondition(mongoshakeCondition)
	}
	cr.Status.JobStatus = &jobStatus

	switch jobStatus.Status {
	case api.AppStateRunning:
		cr.Status.State = api.AppStateRunning
	case api.AppStateComplete:
		cr.Status.State = api.AppStateComplete
	case api.AppStateStopping:
		cr.Status.State = api.AppStateStopping
	case api.AppStatePaused:
		cr.Status.State = api.AppStatePaused
	default:
		cr.Status.State = appState
	}
	clusterCondition.Type = cr.Status.State
	cr.Status.AddCondition(clusterCondition)
	return r.writeStatus(cr)
}
func (r *MongoShakeReconciler) jobStatus(cr *api.MongoShake) (api.JobStatus, error) {
	jobPods, err := r.getJobPods(cr)
	if err != nil {
		return api.JobStatus{}, fmt.Errorf("get list: %v", err)
	}
	status := api.JobStatus{
		Size:   len(jobPods.Items),
		Status: api.AppStateInit,
	}
	if cr.Spec.Pause && status.Size == 0 {
		status.Status = api.AppStatePaused
	}
	for _, pod := range jobPods.Items {
		for _, cntr := range pod.Status.ContainerStatuses {
			if cntr.State.Waiting != nil && cntr.State.Waiting.Message != "" {
				status.Message += cntr.Name + ": " + cntr.State.Waiting.Message + "; "
			}
		}
		for _, cond := range pod.Status.Conditions {
			switch cond.Type {
			case corev1.ContainersReady:
				if cond.Status == corev1.ConditionTrue {
					status.Ready++
				}
			case corev1.PodScheduled:
				if cond.Reason == corev1.PodReasonUnschedulable &&
					cond.LastTransitionTime.Time.Before(time.Now().Add(-1*time.Minute)) {
					status.Message = cond.Message
				}
			}
		}
		switch {
		case cr.Spec.Pause && (cr.Spec.Collector.Mode == api.SyncModeAll || cr.Spec.Collector.Mode == api.SyncModelIncr):
			if err := r.Delete(context.TODO(), &pod); err != nil {
				log.Error(err, "dts-job is paused, sync mode is incr or all, delete job pod fail", "pod", pod.Name)
			}
			status.Status = api.AppStatePaused
		case cr.Spec.Pause && status.Ready > 0:
			status.Status = api.AppStateStopping
		case cr.Spec.Pause:
			if err := r.Delete(context.TODO(), &pod); err != nil {
				log.Error(err, "dts-job is paused,delete job pod fail", "pod", pod.Name)
			}
			status.Status = api.AppStatePaused
		case status.Size == status.Ready && pod.Status.Phase == corev1.PodRunning:
			status.Status = api.AppStateRunning
		case pod.Status.Phase == corev1.PodSucceeded:
			status.Status = api.AppStateComplete
		}
	}

	return status, nil
}

//
func (r *MongoShakeReconciler) writeStatus(cr *api.MongoShake) error {
	err := r.Status().Update(context.TODO(), cr)
	if err != nil {
		err := r.Update(context.TODO(), cr)
		if err != nil {
			return errors.Wrap(err, "send update")
		}
	}

	return nil
}
