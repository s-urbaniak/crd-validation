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
	"fmt"
	"k8s.io/apimachinery/pkg/util/json"

	"k8s.io/apimachinery/pkg/runtime"
	validationv1 "my.domain/validation/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type NoValidationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=my.domain,resources=guestbooks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=my.domain,resources=guestbooks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=my.domain,resources=guestbooks/finalizers,verbs=update
func (r *NoValidationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var customResource validationv1.NoValidation
	if err := r.Get(ctx, req.NamespacedName, &customResource); err != nil {
		logger.Error(err, "error reconciling")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}
	fmt.Printf("reconciling %v: ", req)

	isItThis := struct {
		Foo string `json:"foo"`
	}{}
	if err := json.Unmarshal(customResource.Spec.Raw, &isItThis); err == nil && isItThis.Foo != "" {
		fmt.Println(`it's a foo!`, isItThis.Foo)
		return reconcile.Result{}, nil
	}

	isItThat := struct {
		Bar string `json:"bar"`
	}{}
	if err := json.Unmarshal(customResource.Spec.Raw, &isItThat); err == nil && isItThat.Bar != "" {
		fmt.Println(`it's a bar!`, isItThat.Bar)
		return reconcile.Result{}, nil
	}

	anythingElse := map[string]interface{}{}
	if err := json.Unmarshal(customResource.Spec.Raw, &anythingElse); err == nil {
		fmt.Println(`it's JSON!`, anythingElse)
		return reconcile.Result{}, nil
	}

	fmt.Println(`it's an arbitrary value'`, string(customResource.Spec.Raw))

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NoValidationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&validationv1.NoValidation{}).
		Complete(r)
}
