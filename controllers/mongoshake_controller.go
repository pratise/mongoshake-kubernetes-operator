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

package controllers

import (
	"context"
	"encoding/base64"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	api "github.com/pratise/mongoshake-kubernetes-operator/api/v1"
	"github.com/pratise/mongoshake-kubernetes-operator/controllers/ms"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sync"
	"sync/atomic"
	"time"
)

var log = logf.Log.WithName("controller_mongoshake")

// MongoShakeReconciler reconciles a MongoShake object
type MongoShakeReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Lockers lockStore
}

type lockStore struct {
	store *sync.Map
}

func NewLockStore() lockStore {
	return lockStore{
		store: new(sync.Map),
	}
}

func (l lockStore) LoadOrCreate(key string) lock {
	val, _ := l.store.LoadOrStore(key, lock{
		statusMutex: new(sync.Mutex),
		updateSync:  new(int32),
	})

	return val.(lock)
}

type lock struct {
	statusMutex *sync.Mutex
	updateSync  *int32
}

const (
	updateDone = 0
	updateWait = 1
)

// +kubebuilder:rbac:groups=pratise.github.com,resources=mongoshakes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=pratise.github.com,resources=mongoshakes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=pratise.github.com,resources=mongoshakes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MongoShake object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *MongoShakeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)

	// 需要Result添加RequeueAfter，保证失败如果失败以后几秒钟后再次进行协调
	rr := ctrl.Result{
		RequeueAfter: time.Second * 5,
	}
	// 如果一个operator需要控制多个命名空间的话，应该为每个集群创建锁，以避免锁定其他集群的cron job
	l := r.Lockers.LoadOrCreate(req.NamespacedName.String())
	l.statusMutex.Lock()
	defer l.statusMutex.Unlock()
	defer atomic.StoreInt32(l.updateSync, updateDone)

	// /获取MongoShake实例
	cr := &api.MongoShake{}
	err := r.Get(context.TODO(), req.NamespacedName, cr)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// 如果请求对象不存的，可能在协调请求以后已经被删除。
			// 自身的对象会自动被垃圾收集。对于其他清理逻辑使用finalizers。
			// 返回，不重新请求
			return reconcile.Result{}, nil
		}
		return rr, err
	}

	mongoshakeStatus := api.AppStateInit
	defer func() {
		err = r.updateStatus(cr, err, mongoshakeStatus)
		if err != nil {
			logger.Error(err, "failed to update mongoshake status", "mongoshake", cr.Name)
		}
	}()

	err = r.checkNSetDefaults(log)
	if err != nil {
		err = errors.Wrap(err, "wrong mongoshake options")
		return reconcile.Result{}, err
	}

	err = r.reconcileMongoshakeConfigMaps(cr)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "reconcile mongoshake configmaps")
	}
	err = r.reconcileMongoshakeLogPVC(cr)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "reconcile mongoshake log pvc")
	}
	err = r.reconcileJob(cr, logger)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "reconcile deployments")
	}

	return rr, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MongoShakeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.MongoShake{}).
		Complete(r)
}

func (r *MongoShakeReconciler) reconcileMongoshakeConfigMaps(cr *api.MongoShake) error {
	name := ms.MongoshakeCustomConfigName(cr.Name)

	err := r.createOrUpdateConfigMap(cr, &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cr.Namespace,
		},
		Data: map[string]string{
			"collector.conf": cr.Spec.Configuration,
		},
	})
	if err != nil {
		return errors.Wrapf(err, "update or create configMap: %s", name)
	}

	return nil
}

// 创建配置文件configmap
func (r *MongoShakeReconciler) createOrUpdateConfigMap(cr *api.MongoShake, configMap *corev1.ConfigMap) error {
	err := setControllerReference(cr, configMap, r.Scheme)
	if err != nil {
		return errors.Wrapf(err, "failed to set controller ref for config map %s", configMap.Name)
	}

	currMap := &corev1.ConfigMap{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Namespace: configMap.Namespace,
		Name:      configMap.Name,
	}, currMap)
	if err != nil && !k8serrors.IsNotFound(err) {
		return errors.Wrap(err, "get current configmap")
	}

	if k8serrors.IsNotFound(err) {
		log.Info("create mongoshake configmap", "mongoshake", cr.Name, "configmap", configMap.Name)
		return r.Create(context.TODO(), configMap)
	}

	if !mapsEqual(currMap.Data, configMap.Data) {
		log.Info("update mongoshake configmap", "mongoshake", cr.Name, "configmap", configMap.Name)
		return r.Update(context.TODO(), configMap)
	}

	return nil
}

func (r *MongoShakeReconciler) reconcileJob(cr *api.MongoShake, logger logr.Logger) error {

	msJob := ms.MongoshakeJob(cr)
	err := setControllerReference(cr, msJob, r.Scheme)
	if err != nil {
		return errors.Wrapf(err, "set owner ref for deployment %s", msJob.Name)
	}
	err = r.Get(context.TODO(), types.NamespacedName{Name: msJob.Name, Namespace: msJob.Namespace}, msJob)
	if err != nil && !k8serrors.IsNotFound(err) {
		return errors.Wrapf(err, "get job %s", msJob.Name)
	}

	jopSpec, err := ms.MongoshakeJobSpec(cr, logger)
	if err != nil {
		return errors.Wrapf(err, "create job spec %s", msJob.Name)
	}

	msJob.Spec = jopSpec
	err = r.createOrUpdateJob(msJob)
	if err != nil {
		return errors.Wrapf(err, "update or create job %s", msJob.Name)
	}

	return nil
}

func (r *MongoShakeReconciler) createOrUpdate(depl *appsv1.Deployment) error {
	objectMeta := depl.GetObjectMeta()
	if objectMeta.GetAnnotations() == nil {
		objectMeta.SetAnnotations(make(map[string]string))
	}
	objAnnotations := objectMeta.GetAnnotations()
	delete(objAnnotations, "mongoshake/last-config-hash")
	objectMeta.SetAnnotations(objAnnotations)

	hash, err := getObjectHash(depl)
	if err != nil {
		return errors.Wrap(err, "calculate object hash")
	}
	objAnnotations = objectMeta.GetAnnotations()
	objAnnotations["mongoshake/last-config-hash"] = hash
	objectMeta.SetAnnotations(objAnnotations)

	oldDeployment := &appsv1.Deployment{}
	err = r.Get(context.Background(), types.NamespacedName{
		Namespace: depl.Namespace,
		Name:      depl.Name,
	}, oldDeployment)
	if err != nil && !k8serrors.IsNotFound(err) {
		return errors.Wrapf(err, "get object")
	}
	if k8serrors.IsNotFound(err) {
		log.Info("create mongoshake deployment", "deployment", depl.Name)
		return r.Create(context.TODO(), depl)
	}

	oldObjectMeta := oldDeployment.GetObjectMeta()
	if oldObjectMeta.GetAnnotations()["mongoshake/last-config-hash"] != hash || !isObjectMetaEqual(depl, oldObjectMeta) {
		depl.SetResourceVersion(oldObjectMeta.GetResourceVersion())
		log.Info("update mongoshake deployment", "deployment", depl.Name)
		return r.Update(context.TODO(), depl)
	}
	return nil
}

func (r *MongoShakeReconciler) createOrUpdateJob(job *batchv1.Job) error {
	objectMeta := job.GetObjectMeta()
	if objectMeta.GetAnnotations() == nil {
		objectMeta.SetAnnotations(make(map[string]string))
	}
	objAnnotations := objectMeta.GetAnnotations()
	delete(objAnnotations, "mongoshake/last-config-hash")
	objectMeta.SetAnnotations(objAnnotations)

	hash, err := getObjectHash(job)
	if err != nil {
		return errors.Wrap(err, "calculate object hash")
	}
	objAnnotations = objectMeta.GetAnnotations()
	objAnnotations["mongoshake/last-config-hash"] = hash
	objectMeta.SetAnnotations(objAnnotations)

	oldJob := &batchv1.Job{}
	err = r.Get(context.Background(), types.NamespacedName{
		Namespace: job.Namespace,
		Name:      job.Name,
	}, oldJob)
	if err != nil && !k8serrors.IsNotFound(err) {
		return errors.Wrapf(err, "get object")
	}
	if k8serrors.IsNotFound(err) {
		log.Info("create mongoshake job", "job", job.Name)
		return r.Create(context.TODO(), job)
	}

	oldObjectMeta := oldJob.GetObjectMeta()
	if oldObjectMeta.GetAnnotations()["mongoshake/last-config-hash"] != hash || !isObjectMetaEqual(job, oldObjectMeta) {
		job.SetResourceVersion(oldObjectMeta.GetResourceVersion())
		log.Info("update mongoshake job", "job", job.Name)
		return r.Update(context.TODO(), job)
	}
	return nil
}

func (r *MongoShakeReconciler) reconcileMongoshakeLogPVC(cr *api.MongoShake) error {
	logName := ms.MongoshakeCustomPersistentVolumeClaimLogName(cr.Name)
	err := r.createOrUpdatePersistentVolumeClaim(cr, &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      logName,
			Namespace: cr.Namespace,
			Labels:    mongoshakeLables(cr),
		},
		Spec: *cr.Spec.VolumeSpec.PersistentVolumeClaim,
	})
	if err != nil {
		return errors.Wrapf(err, "update or create PersistentVolumeClaim: %s", logName)
	}
	return nil
}

func (r *MongoShakeReconciler) createOrUpdatePersistentVolumeClaim(cr *api.MongoShake, pvc *corev1.PersistentVolumeClaim) error {
	err := setControllerReference(cr, pvc, r.Scheme)
	if err != nil {
		return errors.Wrapf(err, "failed to set controller ref for PersistentVolumeClaim %s", pvc.Name)
	}

	currPvc := &corev1.PersistentVolumeClaim{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Namespace: pvc.Namespace,
		Name:      pvc.Name,
	}, currPvc)
	if err != nil && !k8serrors.IsNotFound(err) {
		return errors.Wrap(err, "get current pvc")
	}

	if k8serrors.IsNotFound(err) {
		log.Info("create mongoshake pvc", "mongoshake", cr.Name, "pvc", pvc.Name)
		return r.Create(context.TODO(), pvc)
	}

	if *currPvc.Spec.Resources.Requests.Storage() != *pvc.Spec.Resources.Requests.Storage() {
		log.Info("update mongoshake pvc", "mongoshake", cr.Name, "pvc", pvc.Name)
		return r.Update(context.TODO(), pvc)
	}

	return nil
}

func isObjectMetaEqual(old, new metav1.Object) bool {
	return compareMaps(old.GetAnnotations(), new.GetAnnotations()) && compareMaps(old.GetLabels(), new.GetLabels())
}

func compareMaps(x, y map[string]string) bool {
	if len(x) != len(y) {
		return false
	}
	for k, v := range x {
		yVal, ok := y[k]
		if !ok || yVal != v {
			return false
		}
	}
	return true
}

func getObjectHash(obj runtime.Object) (string, error) {
	var dataToMarshall interface{}
	switch object := obj.(type) {
	case *appsv1.StatefulSet:
		dataToMarshall = object.Spec
	case *appsv1.Deployment:
		dataToMarshall = object.Spec
	case *corev1.Service:
		dataToMarshall = object.Spec
	case *batchv1.Job:
		dataToMarshall = object.Spec
	default:
		dataToMarshall = obj
	}
	data, err := json.Marshal(dataToMarshall)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// 设置OwnerReferences
func setControllerReference(owner runtime.Object, obj metav1.Object, scheme *runtime.Scheme) error {
	ownerRef, err := OwnerRef(owner, scheme)
	if err != nil {
		return err
	}
	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
	return nil
}

// OwnerRef returns OwnerReference to object
func OwnerRef(ro runtime.Object, scheme *runtime.Scheme) (metav1.OwnerReference, error) {
	gvk, err := apiutil.GVKForObject(ro, scheme)
	if err != nil {
		return metav1.OwnerReference{}, err
	}

	trueVar := true

	ca, err := meta.Accessor(ro)
	if err != nil {
		return metav1.OwnerReference{}, err
	}

	return metav1.OwnerReference{
		APIVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
		Name:       ca.GetName(),
		UID:        ca.GetUID(),
		Controller: &trueVar,
	}, nil
}

// 数据比对
func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for ka, va := range a {
		if vb, ok := b[ka]; !ok || vb != va {
			return false
		}
	}

	return true
}
