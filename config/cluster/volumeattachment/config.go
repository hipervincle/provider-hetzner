package volumeattachment

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure adds configurations for volume attachment cluster.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("hcloud_volume_attachment", func(r *config.Resource) {
		r.Kind = "VolumeAttachment"
		r.ShortGroup = "server"
		r.References["server_id"] = config.Reference{
			TerraformName: "hcloud_server",
		}
		r.References["volume_id"] = config.Reference{
			TerraformName: "hcloud_volume",
		}
	})
}
