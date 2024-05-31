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

	webappv1alpha1 "dmpe.github.io/DeclarativeLabels/api/v1alpha1"
	"golang.org/x/exp/maps"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// DeclarativeLabelsReconciler reconciles a DeclarativeLabels object
type DeclarativeLabelsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.dmpe.github.io,resources=declarativelabels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.dmpe.github.io,resources=declarativelabels/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=webapp.dmpe.github.io,resources=declarativelabels/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DeclarativeLabels object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *DeclarativeLabelsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Starting DeclarativeLabels reconciler/controller...")

	var nodes corev1.Node
	declareLabelsCRD := &webappv1alpha1.DeclarativeLabels{}
	desiredLabels := declareLabelsCRD.Spec.NodeLabels
	clusterNodes := &corev1.NodeList{}

	// Get all nodes
	if err := r.Get(ctx, req.NamespacedName, &nodes); err != nil {
		log.Error(err, "Unable to fetch any nodes")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// skip those with non-worker
	// if we have a specific worker label than use this below
	// workerSelector := client.MatchingLabels{"node-role.kubernetes.io/control-plane": ""}
	// otherwise use inverse
	labelSelector := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "node-role.kubernetes.io/control-plane",
				Operator: metav1.LabelSelectorOpNotIn,
				Values:   []string{""},
			},
		},
	}
	// Convert LabelSelector to Selector
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		log.Error(err, "Failed to create selector from label selector")
		return ctrl.Result{}, err
	}
	listOptions := &client.ListOptions{
		LabelSelector: selector,
	}

	// List all nodes
	if err := r.List(ctx, clusterNodes, listOptions); err != nil {
		log.Error(err, "Unable to fetch any nodes")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// print all worker nodes
	// iterate over each worker node and get its labels
	for i, node := range clusterNodes.Items {
		fmt.Println("Iteration number: ", i, " == Node name == ", node.Name)
		for key, value := range node.GetLabels() {
			fmt.Println("Iteration number: ", i, "inner loop - key: ", key, "value:", value)
		}
		fmt.Println("Updating labels now...")

		allLabels := node.GetLabels()
		actualLabelsKeys := maps.Keys(allLabels)
		desiredLabelsKeys := maps.Keys(desiredLabels)

		// both are []string
		fmt.Println(actualLabelsKeys)
		fmt.Println(desiredLabelsKeys)

		for myKey, myValue := range desiredLabels {
			if value, ok := allLabels[myKey]; ok {
				fmt.Printf("Desired %s-%s already exists on node: %s\n", myKey, value, node.Name)
			} else {
				fmt.Println("Desired label not found - adding it now...")
				node.Labels[myKey] = myValue
				if err := r.Update(ctx, &node); err != nil {
					log.Error(err, "Failed to patch node", "node", node.Name)
					return ctrl.Result{}, err
				}
			}
		}
	}

	return ctrl.Result{}, nil
}

// Some new method which works with minNodes value

func IsMaster(node *corev1.Node) bool {
	if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
		return true
	}

	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeclarativeLabelsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1alpha1.DeclarativeLabels{}).
		// additionally watch for all nodes changes
		Watches(&corev1.Node{}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
