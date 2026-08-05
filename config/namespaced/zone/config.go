package zone

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure adds configurations for DNS zones.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("hcloud_zone", func(r *config.Resource) {
		r.ShortGroup = "dns"
	})
}
