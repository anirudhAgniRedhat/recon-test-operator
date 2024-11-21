package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ReconTestReconciler reconciles a ReconTest object
type ReconTestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=create;get;list;update;watch;delete

// Reconcile is part of the main kubernetes reconciliation loop
func (r *ReconTestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Generate and reconcile CRDs
	return r.reconcileAllCRDs(ctx, logger)
}

// reconcileAllCRDs generates and ensures all CRDs are created
func (r *ReconTestReconciler) reconcileAllCRDs(ctx context.Context, logger logr.Logger) (ctrl.Result, error) {
	// Number of CRDs to generate
	numCRDs := 100

	for i := 1; i <= numCRDs; i++ {
		// Construct CRD name
		crdName := fmt.Sprintf("recontests%d.example.anirudh.io", i)

		// Create CRD object
		crd := r.generateCRD(i, crdName)

		// Check if CRD exists
		existingCRD := &v1.CustomResourceDefinition{}
		err := r.Get(ctx, types.NamespacedName{Name: crdName}, existingCRD)
		if err != nil {
			if errors.IsNotFound(err) {
				// CRD doesn't exist, create it
				logger.Info(fmt.Sprintf("Creating CRD: %s", crdName))
				if createErr := r.Create(ctx, crd); createErr != nil {
					logger.Error(createErr, fmt.Sprintf("Failed to create CRD %s", crdName))
					return ctrl.Result{}, createErr
				}
				logger.Info(fmt.Sprintf("Successfully created CRD: %s", crdName))
			} else {
				// Other error occurred
				logger.Error(err, fmt.Sprintf("Error checking CRD %s", crdName))
				return ctrl.Result{}, err
			}
		} else {
			// CRD exists, update if necessary
			logger.Info(fmt.Sprintf("CRD %s already exists", crdName))

			// Optionally, you can add update logic here if needed
			// For example, comparing spec and updating if different
		}
	}

	return ctrl.Result{}, nil
}

// generateCRD creates a CustomResourceDefinition for a specific index
func (r *ReconTestReconciler) generateCRD(index int, crdName string) *v1.CustomResourceDefinition {
	return &v1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: crdName,
			Labels: map[string]string{
				"generated-by": "recontest-controller",
				"index":        fmt.Sprintf("%d", index),
			},
		},
		Spec: v1.CustomResourceDefinitionSpec{
			Group: "example.anirudh.io",
			Names: v1.CustomResourceDefinitionNames{
				Kind:     fmt.Sprintf("Recontest%d", index),
				ListKind: fmt.Sprintf("Recontest%dList", index),
				Plural:   fmt.Sprintf("recontests%d", index),
				Singular: fmt.Sprintf("recontest%d", index),
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
									Properties: map[string]v1.JSONSchemaProps{
										// Add any specific validation you need
										"description": {
											Type: "string",
										},
										"count": {
											Type: "integer",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReconTestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.CustomResourceDefinition{}).
		Watches(
			&source.Kind{Type: &v1.CustomResourceDefinition{}},
			&handler.EnqueueRequestForObject{},
		).
		Complete(r)
}
