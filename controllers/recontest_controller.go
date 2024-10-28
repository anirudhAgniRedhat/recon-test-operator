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

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ReconTestReconciler reconciles a ReconTest object
type ReconTestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=example.anirudh.io,resources=recontests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.anirudh.io,resources=recontests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.anirudh.io,resources=recontests/finalizers,verbs=update
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=create;get;list;update;watch

// Reconcile is part of the main kubernetes reconciliation loop
func (r *ReconTestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Check if any of the CRDs have been deleted
	for i := 1; i <= 100; i++ {
		crdName := fmt.Sprintf("recontests%d.example.anirudh.io", i)
		crd := &v1.CustomResourceDefinition{}

		err := r.Get(ctx, client.ObjectKey{Name: crdName}, crd)
		if err != nil {
			if errors.IsNotFound(err) {
				logger.Info(fmt.Sprintf("CRD %s is missing, recreating all CRDs...", crdName))
				// If any CRD is missing, recreate all CRDs
				return r.reapplyAllCRDs(ctx, logger)
			}
			logger.Error(err, fmt.Sprintf("Failed to fetch CRD %s", crdName))
			return ctrl.Result{}, err
		}
	}

	// All CRDs are present, no further action needed
	logger.Info("All CRDs are present.")
	return ctrl.Result{}, nil
}

// reapplyAllCRDs will reapply all the CRDs defined in the Python script
func (r *ReconTestReconciler) reapplyAllCRDs(ctx context.Context, logger logr.Logger) (ctrl.Result, error) {
	for i := 1; i <= 100; i++ {
		crdName := fmt.Sprintf("recontests%d.example.anirudh.io", i)

		newCRD := &v1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: crdName,
			},
			Spec: v1.CustomResourceDefinitionSpec{
				Group: "example.anirudh.io",
				Names: v1.CustomResourceDefinitionNames{
					Kind:     fmt.Sprintf("Recontest%d", i),
					ListKind: fmt.Sprintf("Recontest%dList", i),
					Plural:   fmt.Sprintf("recontests%d", i),
					Singular: fmt.Sprintf("recontest%d", i),
				},
				Scope: v1.NamespaceScoped,
				Versions: []v1.CustomResourceDefinitionVersion{
					{
						Name:    "v1alpha1",
						Served:  true,
						Storage: true,
						Schema: &v1.CustomResourceValidation{
							OpenAPIV3Schema: &v1.JSONSchemaProps{
								Type: "object",
								Properties: map[string]v1.JSONSchemaProps{
									"spec": {
										Type: "object",
									},
								},
							},
						},
					},
				},
			},
		}

		// Recreate the CRD
		if err := r.Create(ctx, newCRD); err != nil {
			if errors.IsAlreadyExists(err) {
				logger.Info(fmt.Sprintf("CRD %s already exists", crdName))
				continue
			}
			logger.Error(err, fmt.Sprintf("Failed to recreate CRD %s", crdName))
			return ctrl.Result{}, err
		}
		logger.Info(fmt.Sprintf("Successfully recreated CRD: %s", crdName))
	}
	return ctrl.Result{Requeue: true}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReconTestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.CustomResourceDefinition{}).
		Complete(r)
}
