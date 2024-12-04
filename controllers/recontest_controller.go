package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
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

	// Generate and create CRDs
	return r.createAllCRDs(ctx, logger)
}

// createAllCRDs generates and creates all CRDs
func (r *ReconTestReconciler) createAllCRDs(ctx context.Context, logger logr.Logger) (ctrl.Result, error) {
	// Number of CRDs to generate
	numCRDs := 100

	// Track if any creation errors occurred
	var lastErr error

	for i := 1; i <= numCRDs; i++ {
		// Construct CRD name
		crdName := fmt.Sprintf("complexrecontests%d.example.anirudh.io", i)

		logger.Info(fmt.Sprintf("Creating CRD %s", crdName))

		// Create CRD object
		crd := r.generateComplexCRD(i, crdName)

		// Attempt to create CRD
		if err := r.Create(ctx, crd); err != nil {
			// Check if the error is due to the CRD already existing
			if apierrors.IsAlreadyExists(err) {
				// Log that the CRD already exists
				logger.Info(fmt.Sprintf("CRD already exists: %s", crdName))
				continue // Move to the next CRD
			}

			// Log other errors
			logger.Error(err, fmt.Sprintf("Failed to create complex CRD: %s", crdName))
			lastErr = err
		} else {
			logger.Info(fmt.Sprintf("Successfully created complex CRD: %s", crdName))
		}
	}

	// If there were any errors, requeue with a delay
	if lastErr != nil {
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second * 30, // Requeue after 30 seconds to retry failed CRD creations
		}, lastErr
	}

	// Continuously requeue to keep trying to create CRDs
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: time.Minute * 1, // Requeue every 5 minutes
	}, nil
}

// generateComplexCRD creates a highly nested CustomResourceDefinition
func (r *ReconTestReconciler) generateComplexCRD(index int, crdName string) *v1.CustomResourceDefinition {
	plural := fmt.Sprintf("complexrecontests%d", index) // spec.names.plural
	group := "example.anirudh.io"                       // spec.group

	return &v1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", plural, group), // Ensure name is in the correct format
			Labels: map[string]string{
				"generated-by": "complex-recontest-controller",
				"complexity":   "high",
				"index":        fmt.Sprintf("%d", index),
				"timestamp":    fmt.Sprintf("%d", time.Now().Unix()),
			},
		},
		Spec: v1.CustomResourceDefinitionSpec{
			Group: group,
			Names: v1.CustomResourceDefinitionNames{
				Kind:     fmt.Sprintf("ComplexRecontest%d", index),
				ListKind: fmt.Sprintf("ComplexRecontest%dList", index),
				Plural:   plural,
				Singular: fmt.Sprintf("complexrecontest%d", index),
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
										"organization": {
											Type: "object",
											Properties: map[string]v1.JSONSchemaProps{
												"name": {
													Type: "string",
												},
												"foundedYear": {
													Type:    "integer",
													Minimum: float64Ptr(1800.0),
													Maximum: float64Ptr(2100.0),
												},
												"address": {
													Type: "object",
													Properties: map[string]v1.JSONSchemaProps{
														"street": {
															Type: "string",
														},
														"city": {
															Type: "string",
														},
														"state": {
															Type: "string",
														},
														"postalCode": {
															Type:    "string",
															Pattern: `^\d{5}(-\d{4})?$`, // US Zip code pattern
														},
														"geoLocation": {
															Type: "object",
															Properties: map[string]v1.JSONSchemaProps{
																"latitude": {
																	Type:    "number",
																	Minimum: float64Ptr(-90.0),
																	Maximum: float64Ptr(90.0),
																},
																"longitude": {
																	Type:    "number",
																	Minimum: float64Ptr(-180.0),
																	Maximum: float64Ptr(180.0),
																},
															},
															Required: []string{"latitude", "longitude"},
														},
													},
													Required: []string{"street", "city", "state"},
												},
												"departments": {
													Type: "array",
													Items: &v1.JSONSchemaPropsOrArray{
														Schema: &v1.JSONSchemaProps{
															Type: "object",
															Properties: map[string]v1.JSONSchemaProps{
																"name": {
																	Type: "string",
																},
																"headCount": {
																	Type:    "integer",
																	Minimum: float64Ptr(0.0),
																},
																"budget": {
																	Type: "object",
																	Properties: map[string]v1.JSONSchemaProps{
																		"annual": {
																			Type:    "number",
																			Minimum: float64Ptr(0.0),
																		},
																		"currency": {
																			Type: "string",
																			Enum: []v1.JSON{
																				{Raw: []byte(`"USD"`)},
																				{Raw: []byte(`"EUR"`)},
																				{Raw: []byte(`"GBP"`)},
																				{Raw: []byte(`"JPY"`)},
																			},
																		},
																	},
																	Required: []string{"annual", "currency"},
																},
															},
															Required: []string{"name", "headCount"},
														},
													},
												},
											},
											Required: []string{"name", "foundedYear"},
										},
										"projectMetadata": {
											Type: "object",
											Properties: map[string]v1.JSONSchemaProps{
												"projectId": {
													Type: "string",
												},
												"status": {
													Type: "string",
													Enum: []v1.JSON{
														{Raw: []byte(`"planning"`)},
														{Raw: []byte(`"in-progress"`)},
														{Raw: []byte(`"completed"`)},
														{Raw: []byte(`"on-hold"`)},
													},
												},
												"resources": {
													Type: "array",
													Items: &v1.JSONSchemaPropsOrArray{
														Schema: &v1.JSONSchemaProps{
															Type: "object",
															Properties: map[string]v1.JSONSchemaProps{
																"type": {
																	Type: "string",
																},
																"quantity": {
																	Type:    "integer",
																	Minimum: float64Ptr(0.0),
																},
																"details": {
																	Type: "object",
																	AdditionalProperties: &v1.JSONSchemaPropsOrBool{
																		Allows: true,
																	},
																},
															},
															Required: []string{"type", "quantity"},
														},
													},
												},
												"timeline": {
													Type: "object",
													Properties: map[string]v1.JSONSchemaProps{
														"startDate": {
															Type:   "string",
															Format: "date",
														},
														"endDate": {
															Type:   "string",
															Format: "date",
														},
														"milestones": {
															Type: "array",
															Items: &v1.JSONSchemaPropsOrArray{
																Schema: &v1.JSONSchemaProps{
																	Type: "object",
																	Properties: map[string]v1.JSONSchemaProps{
																		"name": {
																			Type: "string",
																		},
																		"completionDate": {
																			Type:   "string",
																			Format: "date",
																		},
																		"dependencies": {
																			Type: "array",
																			Items: &v1.JSONSchemaPropsOrArray{
																				Schema: &v1.JSONSchemaProps{
																					Type: "string",
																				},
																			},
																		},
																	},
																	Required: []string{"name", "completionDate"},
																},
															},
														},
													},
													Required: []string{"startDate", "endDate"},
												},
											},
											Required: []string{"projectId", "status"},
										},
									},
									Required: []string{"organization", "projectMetadata"},
								},
							},
							Required: []string{"spec"},
						},
					},
				},
			},
		},
	}
}

// Helper function to create float64 pointers
func float64Ptr(f float64) *float64 {
	return &f
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
