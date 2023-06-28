/*
 (c) Copyright [2021-2023] Open Text.
 Licensed under the Apache License, Version 2.0 (the "License");
 You may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package vadmin

import (
	"context"
	"fmt"

	"github.com/icza/gog"
	vops "github.com/vertica/vcluster/vclusterops"
	vapi "github.com/vertica/vertica-kubernetes/api/v1beta1"
	"github.com/vertica/vertica-kubernetes/pkg/net"
	"github.com/vertica/vertica-kubernetes/pkg/vadmin/opts/addnode"
)

// AddNode will add a new vertica node to the cluster
func (v *VClusterOps) AddNode(ctx context.Context, opts ...addnode.Option) error {
	v.Log.Info("Starting vcluster AddNode")
	s := addnode.Parms{}
	s.Make(opts...)
	vaddNodeOptions, err := v.buildAddNodeOptions(ctx, &s)
	if err != nil {
		return fmt.Errorf("failed to build add node options: %w", err)
	}
	newNodeInfo, err := vops.VAddNode(vaddNodeOptions)
	v.Log.Info("vcluster AddNode is done", "newNodeInfo", newNodeInfo, "err", err)
	return err
}

// buildAddNodeOptions will build up the options struct to be passed into VAddNode
func (v *VClusterOps) buildAddNodeOptions(ctx context.Context, s *addnode.Parms) (*vops.VAddNodeOptions, error) {
	// get the certs
	certs, err := v.retrieveHTTPSCerts(ctx)
	if err != nil {
		return nil, err
	}

	opts := &vops.VAddNodeOptions{
		BootstrapHost:  s.InitiatorIP,
		SubclusterName: s.Subcluster,
		DatabaseOptions: vops.DatabaseOptions{
			Name:            &v.VDB.Spec.DBName,
			RawHosts:        s.Hosts,
			Ipv6:            gog.Ptr(net.IsIPv6(s.InitiatorIP)),
			CatalogPrefix:   gog.Ptr(v.VDB.Spec.Local.GetCatalogPath()),
			DataPrefix:      &v.VDB.Spec.Local.DataPath,
			ConfigDirectory: nil,
			DepotPrefix:     &v.VDB.Spec.Local.DepotPath,
			IsEon:           gog.Ptr(v.VDB.IsEON()),
			UserName:        gog.Ptr(vapi.SuperUser),
			Password:        &v.Password,
			Key:             certs.Key,
			Cert:            certs.Cert,
			CaCert:          certs.CaCert,
			LogPath:         nil,
			HonorUserInput:  gog.Ptr(true),
		},
	}
	return opts, nil
}
