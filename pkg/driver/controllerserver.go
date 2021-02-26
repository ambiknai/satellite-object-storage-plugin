/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2019, 2020 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/


package driver

import (
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/container-storage-interface/spec/lib/go/csi"
		"golang.org/x/net/context"
		"google.golang.org/grpc/status"
		"google.golang.org/grpc/codes"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
}

func (cs *controllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return &csi.ControllerExpandVolumeResponse{}, status.Error(codes.Unimplemented, "ControllerExpandVolume is not implemented")
}
