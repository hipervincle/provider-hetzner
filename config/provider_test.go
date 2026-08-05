package config

import (
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

type expectedResource struct {
	shortGroup string
	kind       string
}

var expectedResources = map[string]expectedResource{
	"hcloud_firewall":               {shortGroup: "network"},
	"hcloud_firewall_attachment":    {shortGroup: "server", kind: "FirewallAttachment"},
	"hcloud_floating_ip":            {shortGroup: "server", kind: "FloatingIP"},
	"hcloud_floating_ip_assignment": {shortGroup: "server", kind: "FloatingIPAssignment"},
	"hcloud_load_balancer":          {shortGroup: "loadbalancer", kind: "LoadBalancer"},
	"hcloud_load_balancer_network":  {shortGroup: "loadbalancer", kind: "NetworkAttachment"},
	"hcloud_load_balancer_service":  {shortGroup: "loadbalancer", kind: "LoadBalancerService"},
	"hcloud_load_balancer_target":   {shortGroup: "loadbalancer", kind: "LoadBalancerTarget"},
	"hcloud_managed_certificate":    {shortGroup: "network", kind: "ManagedCertificate"},
	"hcloud_network":                {shortGroup: "network"},
	"hcloud_network_route":          {shortGroup: "network"},
	"hcloud_network_subnet":         {shortGroup: "network"},
	"hcloud_placement_group":        {shortGroup: "server", kind: "PlacementGroup"},
	"hcloud_server":                 {shortGroup: "server"},
	"hcloud_server_network":         {shortGroup: "server", kind: "NetworkAttachment"},
	"hcloud_snapshot":               {shortGroup: "server"},
	"hcloud_ssh_key":                {shortGroup: "server", kind: "SSHKey"},
	"hcloud_uploaded_certificate":   {shortGroup: "network"},
	"hcloud_volume":                 {shortGroup: "server"},
	"hcloud_volume_attachment":      {shortGroup: "server", kind: "VolumeAttachment"},
	"hcloud_zone":                   {shortGroup: "dns"},
	"hcloud_zone_rrset":             {shortGroup: "dns", kind: "ZoneRRSet"},
}

func TestProviderConfigurationConsistency(t *testing.T) {
	t.Parallel()

	t.Run("ClusterProvider", func(t *testing.T) {
		t.Parallel()
		assertProviderConsistency(t, GetProvider())
	})

	t.Run("NamespacedProvider", func(t *testing.T) {
		t.Parallel()
		assertProviderConsistency(t, GetProviderNamespaced())
	})
}

func assertProviderConsistency(t *testing.T, pc *ujconfig.Provider) {
	t.Helper()

	for name := range expectedResources {
		if _, ok := ExternalNameConfigs[name]; !ok {
			t.Fatalf("ExternalNameConfigs[%q] is missing", name)
		}
	}

	for name, expected := range expectedResources {
		r, ok := pc.Resources[name]
		if !ok {
			t.Fatalf("provider is missing resource %q", name)
		}
		if r.ShortGroup != expected.shortGroup {
			t.Fatalf("resource %q ShortGroup = %q, want %q", name, r.ShortGroup, expected.shortGroup)
		}
		if expected.kind != "" && r.Kind != expected.kind {
			t.Fatalf("resource %q Kind = %q, want %q", name, r.Kind, expected.kind)
		}
	}

	assertReference(t, pc.Resources["hcloud_firewall_attachment"], "firewall_id", "hcloud_firewall")
	assertReference(t, pc.Resources["hcloud_firewall_attachment"], "server_ids", "hcloud_server")
	assertReference(t, pc.Resources["hcloud_server"], "network.network_id", "hcloud_network")
	assertReference(t, pc.Resources["hcloud_server"], "network.*.network_id", "hcloud_network")
	assertReference(t, pc.Resources["hcloud_load_balancer_service"], "http.*.certificates", "hcloud_managed_certificate")
	assertReference(t, pc.Resources["hcloud_volume_attachment"], "server_id", "hcloud_server")
	assertReference(t, pc.Resources["hcloud_volume_attachment"], "volume_id", "hcloud_volume")
	assertReference(t, pc.Resources["hcloud_zone_rrset"], "zone", "hcloud_zone")
	assertIgnoredLateInitField(t, pc.Resources["hcloud_server"], "datacenter")
}

func assertReference(t *testing.T, r *ujconfig.Resource, path, terraformName string) {
	t.Helper()

	ref, ok := r.References[path]
	if !ok {
		t.Fatalf("resource %q is missing reference path %q", r.Name, path)
	}
	if ref.TerraformName != terraformName {
		t.Fatalf("resource %q reference %q targets %q, want %q", r.Name, path, ref.TerraformName, terraformName)
	}
}

func assertIgnoredLateInitField(t *testing.T, r *ujconfig.Resource, field string) {
	t.Helper()

	for _, ignored := range r.LateInitializer.IgnoredFields {
		if ignored == field {
			return
		}
	}

	t.Fatalf("resource %q late initializer ignored fields do not include %q", r.Name, field)
}
