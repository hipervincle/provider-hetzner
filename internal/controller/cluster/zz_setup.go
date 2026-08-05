// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	zone "github.com/miaits/provider-hetzner/internal/controller/cluster/dns/zone"
	zonerrset "github.com/miaits/provider-hetzner/internal/controller/cluster/dns/zonerrset"
	loadbalancer "github.com/miaits/provider-hetzner/internal/controller/cluster/loadbalancer/loadbalancer"
	loadbalancerservice "github.com/miaits/provider-hetzner/internal/controller/cluster/loadbalancer/loadbalancerservice"
	loadbalancertarget "github.com/miaits/provider-hetzner/internal/controller/cluster/loadbalancer/loadbalancertarget"
	networkattachment "github.com/miaits/provider-hetzner/internal/controller/cluster/loadbalancer/networkattachment"
	certificate "github.com/miaits/provider-hetzner/internal/controller/cluster/network/certificate"
	firewall "github.com/miaits/provider-hetzner/internal/controller/cluster/network/firewall"
	managedcertificate "github.com/miaits/provider-hetzner/internal/controller/cluster/network/managedcertificate"
	network "github.com/miaits/provider-hetzner/internal/controller/cluster/network/network"
	route "github.com/miaits/provider-hetzner/internal/controller/cluster/network/route"
	subnet "github.com/miaits/provider-hetzner/internal/controller/cluster/network/subnet"
	providerconfig "github.com/miaits/provider-hetzner/internal/controller/cluster/providerconfig"
	firewallattachment "github.com/miaits/provider-hetzner/internal/controller/cluster/server/firewallattachment"
	floatingip "github.com/miaits/provider-hetzner/internal/controller/cluster/server/floatingip"
	floatingipassignment "github.com/miaits/provider-hetzner/internal/controller/cluster/server/floatingipassignment"
	networkattachmentserver "github.com/miaits/provider-hetzner/internal/controller/cluster/server/networkattachment"
	placementgroup "github.com/miaits/provider-hetzner/internal/controller/cluster/server/placementgroup"
	server "github.com/miaits/provider-hetzner/internal/controller/cluster/server/server"
	snapshot "github.com/miaits/provider-hetzner/internal/controller/cluster/server/snapshot"
	sshkey "github.com/miaits/provider-hetzner/internal/controller/cluster/server/sshkey"
	volume "github.com/miaits/provider-hetzner/internal/controller/cluster/server/volume"
	volumeattachment "github.com/miaits/provider-hetzner/internal/controller/cluster/server/volumeattachment"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		zone.Setup,
		zonerrset.Setup,
		loadbalancer.Setup,
		loadbalancerservice.Setup,
		loadbalancertarget.Setup,
		networkattachment.Setup,
		certificate.Setup,
		firewall.Setup,
		managedcertificate.Setup,
		network.Setup,
		route.Setup,
		subnet.Setup,
		providerconfig.Setup,
		firewallattachment.Setup,
		floatingip.Setup,
		floatingipassignment.Setup,
		networkattachmentserver.Setup,
		placementgroup.Setup,
		server.Setup,
		snapshot.Setup,
		sshkey.Setup,
		volume.Setup,
		volumeattachment.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		zone.SetupGated,
		zonerrset.SetupGated,
		loadbalancer.SetupGated,
		loadbalancerservice.SetupGated,
		loadbalancertarget.SetupGated,
		networkattachment.SetupGated,
		certificate.SetupGated,
		firewall.SetupGated,
		managedcertificate.SetupGated,
		network.SetupGated,
		route.SetupGated,
		subnet.SetupGated,
		providerconfig.SetupGated,
		firewallattachment.SetupGated,
		floatingip.SetupGated,
		floatingipassignment.SetupGated,
		networkattachmentserver.SetupGated,
		placementgroup.SetupGated,
		server.SetupGated,
		snapshot.SetupGated,
		sshkey.SetupGated,
		volume.SetupGated,
		volumeattachment.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
