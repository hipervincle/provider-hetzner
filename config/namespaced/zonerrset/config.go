package zonerrset

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure adds configurations for DNS zone record sets.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("hcloud_zone_rrset", func(r *config.Resource) {
		r.Kind = "ZoneRRSet"
		r.ShortGroup = "dns"
		r.References["zone"] = config.Reference{TerraformName: "hcloud_zone"}
	})
}
