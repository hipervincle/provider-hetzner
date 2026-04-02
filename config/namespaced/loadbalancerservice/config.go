package loadbalancerservice

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure adds configurations for load balancer service namespaced.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("hcloud_load_balancer_service", func(r *config.Resource) {
		r.Kind = "LoadBalancerService"
		r.ShortGroup = "loadbalancer"
		r.References["load_balancer_id"] = config.Reference{
			TerraformName: "hcloud_load_balancer",
		}
		r.References["http.certificates"] = config.Reference{
			TerraformName: "hcloud_managed_certificate",
		}
		r.References["http.*.certificates"] = config.Reference{
			TerraformName: "hcloud_managed_certificate",
		}
	})
}
