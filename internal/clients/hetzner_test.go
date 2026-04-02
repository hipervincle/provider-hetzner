package clients

import (
	"context"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	xpv2 "github.com/crossplane/crossplane-runtime/v2/apis/common/v2"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ktypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	clusterserverv1alpha1 "github.com/miaits/provider-hetzner/apis/cluster/server/v1alpha1"
	clusterv1beta1 "github.com/miaits/provider-hetzner/apis/cluster/v1beta1"
	namespacedserverv1alpha1 "github.com/miaits/provider-hetzner/apis/namespaced/server/v1alpha1"
	namespacedv1beta1 "github.com/miaits/provider-hetzner/apis/namespaced/v1beta1"
)

func TestResolveProviderConfigModernNamespacedProviderConfig(t *testing.T) {
	t.Parallel()

	s := runtime.NewScheme()
	mustAddToScheme(t, s, namespacedv1beta1.SchemeBuilder.AddToScheme, namespacedserverv1alpha1.AddToScheme)

	mg := &namespacedserverv1alpha1.Server{
		TypeMeta: metav1.TypeMeta{
			APIVersion: namespacedserverv1alpha1.CRDGroupVersion.String(),
			Kind:       "Server",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-server",
			Namespace: "tenant-a",
			UID:       ktypes.UID("modern-namespaced-uid"),
		},
		Spec: namespacedserverv1alpha1.ServerSpec{
			ManagedResourceSpec: xpv2.ManagedResourceSpec{
				ProviderConfigReference: &xpv1.ProviderConfigReference{
					Name: "default",
					Kind: "ProviderConfig",
				},
			},
		},
	}
	pc := &namespacedv1beta1.ProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: namespacedv1beta1.SchemeGroupVersion.String(),
			Kind:       "ProviderConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "tenant-a",
		},
		Spec: namespacedProviderConfigSpec("shared-creds"),
	}

	baseClient := fake.NewClientBuilder().WithScheme(s).WithRESTMapper(newModernRESTMapper()).WithObjects(pc).Build()
	cl := modernClusterScopeClient{Client: baseClient}

	got, err := resolveProviderConfig(context.Background(), cl, mg)
	if err != nil {
		t.Fatalf("resolveProviderConfig() unexpected error: %v", err)
	}
	if got.Credentials.SecretRef == nil {
		t.Fatal("resolveProviderConfig() returned nil SecretRef")
	}
	if got.Credentials.SecretRef.Namespace != "tenant-a" {
		t.Fatalf("resolveProviderConfig() secret namespace = %q, want %q", got.Credentials.SecretRef.Namespace, "tenant-a")
	}

	usage := &namespacedv1beta1.ProviderConfigUsage{}
	if err := cl.Get(context.Background(), client.ObjectKey{Name: string(mg.GetUID()), Namespace: mg.GetNamespace()}, usage); err != nil {
		t.Fatalf("Get(ProviderConfigUsage) unexpected error: %v", err)
	}
	if usage.ProviderConfigReference.Kind != "ProviderConfig" {
		t.Fatalf("ProviderConfigUsage kind = %q, want %q", usage.ProviderConfigReference.Kind, "ProviderConfig")
	}
	if usage.ProviderConfigReference.Name != "default" {
		t.Fatalf("ProviderConfigUsage name = %q, want %q", usage.ProviderConfigReference.Name, "default")
	}
}

func TestResolveProviderConfigModernClusterProviderConfig(t *testing.T) {
	t.Parallel()

	s := runtime.NewScheme()
	mustAddToScheme(t, s, namespacedv1beta1.SchemeBuilder.AddToScheme, namespacedserverv1alpha1.AddToScheme)

	mg := &namespacedserverv1alpha1.Server{
		TypeMeta: metav1.TypeMeta{
			APIVersion: namespacedserverv1alpha1.CRDGroupVersion.String(),
			Kind:       "Server",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-server",
			Namespace: "tenant-a",
			UID:       ktypes.UID("modern-cluster-uid"),
		},
		Spec: namespacedserverv1alpha1.ServerSpec{
			ManagedResourceSpec: xpv2.ManagedResourceSpec{
				ProviderConfigReference: &xpv1.ProviderConfigReference{
					Name: "default",
					Kind: "ClusterProviderConfig",
				},
			},
		},
	}
	// controller-runtime's fake client does not ignore the namespace segment for
	// cluster-scoped objects during Get, so mirror the managed namespace here.
	pc := &namespacedv1beta1.ClusterProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: namespacedv1beta1.SchemeGroupVersion.String(),
			Kind:       "ClusterProviderConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "tenant-a",
		},
		Spec: namespacedProviderConfigSpec("crossplane-system"),
	}

	cl := fake.NewClientBuilder().WithScheme(s).WithRESTMapper(newModernRESTMapper()).WithObjects(pc).Build()

	got, err := resolveProviderConfig(context.Background(), cl, mg)
	if err != nil {
		t.Fatalf("resolveProviderConfig() unexpected error: %v", err)
	}
	if got.Credentials.SecretRef == nil {
		t.Fatal("resolveProviderConfig() returned nil SecretRef")
	}
	if got.Credentials.SecretRef.Namespace != "crossplane-system" {
		t.Fatalf("resolveProviderConfig() secret namespace = %q, want %q", got.Credentials.SecretRef.Namespace, "crossplane-system")
	}

	usage := &namespacedv1beta1.ProviderConfigUsage{}
	if err := cl.Get(context.Background(), client.ObjectKey{Name: string(mg.GetUID()), Namespace: mg.GetNamespace()}, usage); err != nil {
		t.Fatalf("Get(ProviderConfigUsage) unexpected error: %v", err)
	}
	if usage.ProviderConfigReference.Kind != "ClusterProviderConfig" {
		t.Fatalf("ProviderConfigUsage kind = %q, want %q", usage.ProviderConfigReference.Kind, "ClusterProviderConfig")
	}
}

func TestResolveProviderConfigLegacyProviderConfig(t *testing.T) {
	t.Parallel()

	s := runtime.NewScheme()
	mustAddToScheme(t, s, clusterv1beta1.SchemeBuilder.AddToScheme, clusterserverv1alpha1.AddToScheme)

	mg := &clusterserverv1alpha1.Server{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clusterserverv1alpha1.CRDGroupVersion.String(),
			Kind:       "Server",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "legacy-server",
			UID:  ktypes.UID("legacy-uid"),
		},
		Spec: clusterserverv1alpha1.ServerSpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{
					Name: "default",
				},
			},
		},
	}
	pc := &clusterv1beta1.ProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clusterv1beta1.SchemeGroupVersion.String(),
			Kind:       "ProviderConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: clusterProviderConfigSpec("crossplane-system"),
	}

	cl := fake.NewClientBuilder().WithScheme(s).WithRESTMapper(newLegacyRESTMapper()).WithObjects(pc).Build()

	got, err := resolveProviderConfig(context.Background(), cl, mg)
	if err != nil {
		t.Fatalf("resolveProviderConfig() unexpected error: %v", err)
	}
	if got.Credentials.SecretRef == nil {
		t.Fatal("resolveProviderConfig() returned nil SecretRef")
	}
	if got.Credentials.SecretRef.Namespace != "crossplane-system" {
		t.Fatalf("resolveProviderConfig() secret namespace = %q, want %q", got.Credentials.SecretRef.Namespace, "crossplane-system")
	}

	usage := &clusterv1beta1.ProviderConfigUsage{}
	if err := cl.Get(context.Background(), client.ObjectKey{Name: string(mg.GetUID())}, usage); err != nil {
		t.Fatalf("Get(ProviderConfigUsage) unexpected error: %v", err)
	}
	if usage.ProviderConfigReference.Name != "default" {
		t.Fatalf("ProviderConfigUsage name = %q, want %q", usage.ProviderConfigReference.Name, "default")
	}
}

func mustAddToScheme(t *testing.T, s *runtime.Scheme, addFns ...func(*runtime.Scheme) error) {
	t.Helper()

	for _, addFn := range addFns {
		if err := addFn(s); err != nil {
			t.Fatalf("AddToScheme() unexpected error: %v", err)
		}
	}
}

func namespacedProviderConfigSpec(secretNamespace string) namespacedv1beta1.ProviderConfigSpec {
	return namespacedv1beta1.ProviderConfigSpec{
		Credentials: namespacedv1beta1.ProviderCredentials{
			Source: xpv1.CredentialsSourceSecret,
			CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
				SecretRef: &xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "provider-creds",
						Namespace: secretNamespace,
					},
					Key: "credentials.json",
				},
			},
		},
	}
}

func clusterProviderConfigSpec(secretNamespace string) clusterv1beta1.ProviderConfigSpec {
	return clusterv1beta1.ProviderConfigSpec{
		Credentials: clusterv1beta1.ProviderCredentials{
			Source: xpv1.CredentialsSourceSecret,
			CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
				SecretRef: &xpv1.SecretKeySelector{
					SecretReference: xpv1.SecretReference{
						Name:      "provider-creds",
						Namespace: secretNamespace,
					},
					Key: "credentials.json",
				},
			},
		},
	}
}

type modernClusterScopeClient struct {
	client.Client
}

func (c modernClusterScopeClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	if _, ok := obj.(*namespacedv1beta1.ClusterProviderConfig); ok {
		key.Namespace = ""
	}
	return c.Client.Get(ctx, key, obj, opts...)
}

func newModernRESTMapper() meta.RESTMapper {
	mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{
		namespacedv1beta1.SchemeGroupVersion,
		namespacedserverv1alpha1.CRDGroupVersion,
	})
	mapper.Add(namespacedv1beta1.ProviderConfigGroupVersionKind, meta.RESTScopeNamespace)
	mapper.Add(namespacedv1beta1.ClusterProviderConfigGroupVersionKind, meta.RESTScopeRoot)
	mapper.Add(namespacedv1beta1.ProviderConfigUsageGroupVersionKind, meta.RESTScopeNamespace)
	mapper.Add(namespacedserverv1alpha1.Server_GroupVersionKind, meta.RESTScopeNamespace)
	return mapper
}

func newLegacyRESTMapper() meta.RESTMapper {
	mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{
		clusterv1beta1.SchemeGroupVersion,
		clusterserverv1alpha1.CRDGroupVersion,
	})
	mapper.Add(clusterv1beta1.ProviderConfigGroupVersionKind, meta.RESTScopeRoot)
	mapper.Add(clusterv1beta1.ProviderConfigUsageGroupVersionKind, meta.RESTScopeRoot)
	mapper.Add(clusterserverv1alpha1.Server_GroupVersionKind, meta.RESTScopeRoot)
	return mapper
}
