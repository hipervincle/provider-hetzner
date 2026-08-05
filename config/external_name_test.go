package config

import (
	"context"
	"testing"
)

func TestLoadBalancerNetworkExternalName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		parameters map[string]any
		want       string
		wantErr    string
	}{
		"UsesExplicitNetworkID": {
			parameters: map[string]any{
				"load_balancer_id": float64(101),
				"network_id":       float64(202),
			},
			want: "101-202",
		},
		"DerivesNetworkIDFromSubnetID": {
			parameters: map[string]any{
				"load_balancer_id": "101",
				"subnet_id":        "202-10.0.1.0/24",
			},
			want: "101-202",
		},
		"FailsWithoutNetworkOrSubnetID": {
			parameters: map[string]any{
				"load_balancer_id": 101,
			},
			wantErr: "network_id or subnet_id must be set",
		},
	}

	en := loadBalancerNetworkExternalName()
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := en.GetIDFn(context.Background(), "", tc.parameters, nil)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("GetIDFn() error = %v, want %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetIDFn() unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("GetIDFn() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestLoadBalancerServiceExternalName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		parameters map[string]any
		want       string
		wantErr    string
	}{
		"UsesExplicitListenPort": {
			parameters: map[string]any{
				"load_balancer_id": 101,
				"listen_port":      8443,
			},
			want: "101__8443",
		},
		"DefaultsHTTPSPort": {
			parameters: map[string]any{
				"load_balancer_id": float64(101),
				"protocol":         "https",
			},
			want: "101__443",
		},
		"FailsForTCPWithoutListenPort": {
			parameters: map[string]any{
				"load_balancer_id": 101,
				"protocol":         "tcp",
			},
			wantErr: "listen_port not set",
		},
	}

	en := loadBalancerServiceExternalName()
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := en.GetIDFn(context.Background(), "", tc.parameters, nil)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("GetIDFn() error = %v, want %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetIDFn() unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("GetIDFn() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestLoadBalancerTargetExternalName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		parameters map[string]any
		want       string
		wantErr    string
	}{
		"BuildsServerTargetID": {
			parameters: map[string]any{
				"load_balancer_id": 11,
				"type":             "server",
				"server_id":        float64(22),
			},
			want: "11__server__22",
		},
		"BuildsLabelSelectorTargetID": {
			parameters: map[string]any{
				"load_balancer_id": 11,
				"type":             "label_selector",
				"label_selector":   "role=web",
			},
			want: "11__label_selector__role=web",
		},
		"BuildsIPTargetID": {
			parameters: map[string]any{
				"load_balancer_id": 11,
				"type":             "ip",
				"ip":               "10.0.0.2",
			},
			want: "11__ip__10.0.0.2",
		},
	}

	en := loadBalancerTargetExternalName()
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := en.GetIDFn(context.Background(), "", tc.parameters, nil)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("GetIDFn() error = %v, want %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetIDFn() unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("GetIDFn() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestServerNetworkExternalName(t *testing.T) {
	t.Parallel()

	en := serverNetworkExternalName()
	got, err := en.GetIDFn(context.Background(), "", map[string]any{
		"server_id": float64(301),
		"subnet_id": "404-10.0.2.0/24",
	}, nil)
	if err != nil {
		t.Fatalf("GetIDFn() unexpected error: %v", err)
	}
	if got != "301-404" {
		t.Fatalf("GetIDFn() = %q, want %q", got, "301-404")
	}
}

func TestVolumeExternalName(t *testing.T) {
	t.Parallel()

	en := volumeExternalName()

	t.Run("UsesPlaceholderIDForUnnamedVolume", func(t *testing.T) {
		t.Parallel()

		got, err := en.GetIDFn(context.Background(), "", map[string]any{}, nil)
		if err != nil {
			t.Fatalf("GetIDFn() unexpected error: %v", err)
		}
		if got != "0" {
			t.Fatalf("GetIDFn() = %q, want %q", got, "0")
		}
	})

	t.Run("ConvertsNumericTerraformStateID", func(t *testing.T) {
		t.Parallel()

		got, err := en.GetExternalNameFn(map[string]any{"id": float64(55)})
		if err != nil {
			t.Fatalf("GetExternalNameFn() unexpected error: %v", err)
		}
		if got != "55" {
			t.Fatalf("GetExternalNameFn() = %q, want %q", got, "55")
		}
	})
}

func TestZoneRRSetExternalName(t *testing.T) {
	t.Parallel()

	en := ExternalNameConfigs["hcloud_zone_rrset"]
	got, err := en.GetIDFn(context.Background(), "", map[string]any{
		"zone": "example.com",
		"name": "www",
		"type": "A",
	}, nil)
	if err != nil {
		t.Fatalf("GetIDFn() unexpected error: %v", err)
	}
	if got != "example.com/www/A" {
		t.Fatalf("GetIDFn() = %q, want %q", got, "example.com/www/A")
	}
}
