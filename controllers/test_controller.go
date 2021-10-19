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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	testv1alpha1 "github.com/kgibm/testoperator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// TestReconciler reconciles a Test object
type TestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const FinalizerName = "example.com/finalizer"

//+kubebuilder:rbac:groups=test.example.com,resources=tests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.example.com,resources=tests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.example.com,resources=tests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Test object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *TestReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconcile started")

	test := &testv1alpha1.Test{}
	err := r.Get(ctx, request.NamespacedName, test)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Resource not found")
			return ctrl.Result{}, nil
		}
		logger.Info("Error retrieving resource")
		return ctrl.Result{}, err
	}

	// Check if we are finalizing
	isMarkedToBeDeleted := test.GetDeletionTimestamp() != nil
	if isMarkedToBeDeleted {
		logger.Info("Marked to be deleted")
		if controllerutil.ContainsFinalizer(test, FinalizerName) {
			logger.Info("Contains finalizer")

			// Remove finalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(test, FinalizerName)
			err := r.Update(ctx, test)
			if err != nil {
				logger.Info("Removing finalizer error")
				return ctrl.Result{}, err
			}
			logger.Info("Removed finalizer")
		}
		logger.Info("Returning successfully within mark check")
		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR
	if !controllerutil.ContainsFinalizer(test, FinalizerName) {
		logger.Info("Adding finalizer")
		controllerutil.AddFinalizer(test, FinalizerName)
		err = r.Update(ctx, test)
		if err != nil {
			logger.Info("Failed to add finalizer")
			return ctrl.Result{}, err
		} else {
			logger.Info("Added finalizer")
			return ctrl.Result{}, nil
		}
	}

	logger.Info("Started normal processing")

	logger.Info("Final successful return")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1alpha1.Test{}).
		Complete(r)
}
