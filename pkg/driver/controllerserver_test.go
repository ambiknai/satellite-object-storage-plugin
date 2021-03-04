/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2019, 2020 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

//Package driver ...
package driver

import (
	//"fmt"
	//"net/http"
	//"reflect"
	//"strconv"
	//"testing"

	//"github.com/google/uuid"

	//"github.com/IBM/ibmcloud-volume-interface/lib/provider"
	//providerError "github.com/IBM/ibmcloud-volume-interface/lib/utils"
	csi "github.com/container-storage-interface/spec/lib/go/csi"
	//"github.com/stretchr/testify/assert"
	"github.ibm.com/alchemy-containers/ibm-csi-common/pkg/utils"

	//"github.com/IBM/ibmcloud-volume-interface/lib/provider/fake"
//	cloudProvider "github.ibm.com/alchemy-containers/ibm-csi-common/pkg/ibmcloudprovider"
//	"golang.org/x/net/context"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
)

var (
	// Define "normal" parameters
	stdVolCap = []*csi.VolumeCapability{
		{
			AccessType: &csi.VolumeCapability_Mount{
				Mount: &csi.VolumeCapability_MountVolume{FsType: "ext2"},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
		},
	}
	stdVolCapNotSupported = []*csi.VolumeCapability{
		{
			AccessType: &csi.VolumeCapability_Mount{
				Mount: &csi.VolumeCapability_MountVolume{FsType: "ext2"},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
			},
		},
	}
	stdBlockVolCap = []*csi.VolumeCapability{
		{
			AccessType: &csi.VolumeCapability_Block{
				Block: &csi.VolumeCapability_BlockVolume{},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
		},
	}
	stdCapRange = &csi.CapacityRange{
		RequiredBytes: 20 * 1024 * 1024 * 1024,
	}
	stdTopology = []*csi.Topology{
		{
			Segments: map[string]string{utils.NodeZoneLabel: "myzone", utils.NodeRegionLabel: "myregion"},
		},
	}
)

