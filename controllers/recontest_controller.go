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
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplev1alpha1 "github.com/anirudhAgniRedhat/recon-test-operator/api/v1alpha1"
)

// ReconTestReconciler reconciles a ReconTest object
type ReconTestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=example.anirudh.io,resources=recontests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.anirudh.io,resources=recontests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.anirudh.io,resources=recontests/finalizers,verbs=update
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=create;get;list;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ReconTest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ReconTestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	crdName := "recontests.example.anirudh.io"
	crd := &apiextensionsv1.CustomResourceDefinition{}
	logger.Info(fmt.Sprintf("reconciling the loop for CRD %s", crdName))

	// Check if the CRD exists
	err := r.Get(ctx, client.ObjectKey{Name: crdName}, crd)
	if err != nil {
		if errors.IsNotFound(err) {
			// CRD is missing, recreate it
			logger.Info(fmt.Sprintf("CRD %s is missing, recreating...", crdName))
			// Define the CRD (match this with your actual CRD definition)
			newCRD := &apiextensionsv1.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name: crdName,
				},
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Group: "example.anirudh.io",
					Names: apiextensionsv1.CustomResourceDefinitionNames{
						Kind:     "ReconTest",
						ListKind: "ReconTestList",
						Plural:   "recontests",
						Singular: "recontest",
					},
					Scope: apiextensionsv1.NamespaceScoped,
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name:    "v1alpha1",
							Served:  true,
							Storage: true,
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type:       "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{}, // define your schema properties
								},
							},
						},
					},
				},
			}
			// Create the CRD
			if createErr := r.Create(ctx, newCRD); createErr != nil {
				logger.Error(createErr, "Failed to recreate the CRD")
				return ctrl.Result{}, createErr
			}
			logger.Info(fmt.Sprintf("Successfully recreated CRD: %s", crdName))
			return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
		}
		logger.Error(err, "Failed to fetch CRD")
		return ctrl.Result{}, err
	}
	// Continue with the regular reconciliation logic
	logger.Info(fmt.Sprintf("CRD %s already exists.", crdName))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReconTestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1alpha1.ReconTest{}).
		Complete(r)
}
