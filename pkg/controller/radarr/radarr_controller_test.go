package radarr

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"

	radarrv1alpha1 "github.com/parflesh/radarr-operator/pkg/apis/radarr/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestRadarrController(t *testing.T) {
	var (
		name      = "radarr-operator"
		namespace = "radarr"
	)
	// A Radarr object with metadata and spec.
	cr := &radarrv1alpha1.Radarr{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: radarrv1alpha1.RadarrSpec{},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{cr}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(radarrv1alpha1.SchemeGroupVersion, cr)

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClientWithScheme(s, objs...)

	// Create a ReconcileRadarr object with the scheme and fake client.
	r := &ReconcileRadarr{
		client: cl,
		scheme: s,
		/*imageInspector: &image_inspect.MockImageInspector{
			GetImageLabelsOutput: &imagetypes.ImageInspectInfo{
				Tag: "test",
			},
			GetImageLabelsError: nil,
		},*/
	}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	// Check spec updates
	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, cr)
	if cr.Spec.Image == "" {
		t.Error("Image spec not updated")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, cr)
	if cr.Spec.WatchFrequency == "" {
		t.Error("Watch Frequency not updated")
	}

	depDep := &appsv1.Deployment{}
	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, depDep)
	if err != nil {
		t.Error("Deployment not created")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, cr)
	if cr.Status.Image != cr.Spec.Image {
		t.Error("status image mismatch")
	}

	depSvc := &corev1.Service{}
	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, depSvc)
	if err != nil {
		t.Error("Service not created")
	}

	// Everything should be good.  Lets check
	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if res.Requeue {
		t.Error("reconcile requeued even though all should be good")
	}
}
