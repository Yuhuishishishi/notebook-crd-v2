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

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	operatorsv2 "convect.ai/notebook-crd/api/v2"
)

// JupyterReconciler reconciles a Jupyter object
type JupyterReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=operators.convect.ai,resources=jupyters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operators.convect.ai,resources=jupyters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operators.convect.ai,resources=jupyters/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=services,verbs="*"
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs="*"

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Jupyter object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *JupyterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("jupyter", req.NamespacedName)

	instance := &operatorsv2.Jupyter{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		log.Error(err, "unable to fetch notebook")
		if apierrs.IsNotFound(err) {
			return ctrl.Result{}, nil // Ignore not found error
		}
		return ctrl.Result{}, err

	}

	// Reconcile statefulset
	ss := generateStatefulSet(instance)

	if err := ctrl.SetControllerReference(instance, ss, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if the statefulset already exists
	foundStateful := &appsv1.StatefulSet{}
	justCreate := false
	err := r.Get(ctx, types.NamespacedName{Name: ss.Name, Namespace: ss.Namespace}, foundStateful)

	if err != nil && apierrs.IsNotFound(err) {
		// Not found, create new
		log.Info("Creating StatefulSet", "namespace", ss.Namespace, "name", ss.Name)
		if err = r.Create(ctx, ss); err != nil {
			log.Error(err, "unable to create StatefulSet")
			return ctrl.Result{}, err
		}
		justCreate = true

	} else if err != nil {
		log.Error(err, "error getting StatefulSet")
		return ctrl.Result{}, err
	}
	// Update the foundStateful object and write the result back if there are any changes
	if !justCreate && copyStatefulSetFields(ss, foundStateful) {
		log.Info("Updating StatefulSet", "namespace", ss.Namespace, "name", ss.Name)
		if err = r.Update(ctx, foundStateful); err != nil {
			log.Error(err, "unable to update StatefulSet")
			return ctrl.Result{}, err
		}
	}

	// Reconcile service
	svc := generateService(instance)

	foundService := &corev1.Service{}
	justCreate = false
	err = r.Get(ctx, types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, foundService)

	if err != nil && apierrs.IsNotFound(err) {
		log.Info("Creating service", "namespace", svc.Namespace, "name", svc.Name)
		if err = r.Create(ctx, svc); err != nil {
			log.Error(err, "unable to create service")
			return ctrl.Result{}, err
		}
		justCreate = true
	} else if err != nil {
		log.Error(err, "error getting service")
		return ctrl.Result{}, err
	}

	// Update the service object if needed
	if !justCreate && copyServiceFields(svc, foundService) {
		log.Info("Updating service", "namespace", svc.Namespace, "name", svc.Name)
		if err = r.Update(ctx, foundService); err != nil {
			log.Error(err, "unable to update service")
			return ctrl.Result{}, err
		}
	}

	// Update the status

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JupyterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorsv2.Jupyter{}).
		Complete(r)
}
