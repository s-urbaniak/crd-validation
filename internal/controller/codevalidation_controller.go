/*
Copyright 2024.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	validationv1 "my.domain/validation/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// CodeValidationReconciler reconciles a Guestbook object
type CodeValidationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=my.domain,resources=guestbooks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=my.domain,resources=guestbooks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=my.domain,resources=guestbooks/finalizers,verbs=update
//
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Guestbook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *CodeValidationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var customResource validationv1.CodeValidation
	if err := r.Get(ctx, req.NamespacedName, &customResource); err != nil {
		logger.Error(err, "error reconciling")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if customResource.Spec.Foo == "invalid" {
		meta.SetStatusCondition(&customResource.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             "False",
			LastTransitionTime: metav1.Now(),
			Reason:             "InvalidField",
			Message:            "foo has invalid value",
		})
	} else {
		meta.SetStatusCondition(&customResource.Status.Conditions, metav1.Condition{
			Type:   "Ready",
			Status: "True",
			Reason: "Valid",
		})
	}

	err := r.Status().Update(ctx, &customResource)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CodeValidationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&validationv1.CodeValidation{}).
		Complete(r)
}
