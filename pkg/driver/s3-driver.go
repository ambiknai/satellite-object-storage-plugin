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
	"fmt"
	csi "github.com/container-storage-interface/spec/lib/go/csi"
	commonError "github.ibm.com/alchemy-containers/ibm-csi-common/pkg/messages"
	"go.uber.org/zap"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"errors"
)


type s3Driver struct {
	name          string
	driver   *csicommon.CSIDriver
	endpoint string

	ids *identityServer
	ns  *nodeServer
	cs  *controllerServer
	vendorVersion string
	logger        *zap.Logger

}


type s3Volume struct {
	VolName string `json:"volName"`
	VolID   string `json:"volID"`
	VolSize int64  `json:"volSize"`
	VolPath string `json:"volPath"`
}

// GetS3CSIDriver ...
func GetS3CSIDriver() *s3Driver {
	return &s3Driver{}
}

func (csiDriver *s3Driver) Setups3Driver(lgr *zap.Logger, name, vendorVersion string) error {
	csiDriver.logger = lgr
	csiDriver.logger.Info("S3CSIDriver-SetupS3CSIDriver setting up S3 CSI Driver...")

	if name == "" {
		return fmt.Errorf("driver name missing")
	}

	// Setup messaging
	commonError.MessagesEn = commonError.InitMessages()

	csiDriver.name = name
	csiDriver.vendorVersion = vendorVersion

	csiDriver.logger.Info("Successfully setup IBM CSI driver")

	return nil
}

// NewS3 initializes the driver
func (csiDriver *s3Driver) NewS3(nodeID string, endpoint string) (*s3Driver, error) {

	driver := csicommon.NewCSIDriver(csiDriver.name, csiDriver.vendorVersion, nodeID)
	if driver == nil {
		 csiDriver.logger.Error("Failed to initialize CSI Driver.")
		 return nil, errors.New("Failed to initialize CSI Driver")
	}

	csis3Driver := &s3Driver{
		endpoint: endpoint,
		driver:   driver,
	}
	return csis3Driver, nil
}

func (s3 *s3Driver) newIdentityServer(d *csicommon.CSIDriver) *identityServer {
	s3.logger.Info("-newIdentityServer-")
	return &identityServer{
		DefaultIdentityServer: csicommon.NewDefaultIdentityServer(d),
	}
}

func (s3 *s3Driver) newControllerServer(d *csicommon.CSIDriver) *controllerServer {
	s3.logger.Info("-newControllerServer-")
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(d),
	}
}

func (s3 *s3Driver) newNodeServer(d *csicommon.CSIDriver) *nodeServer {
	s3.logger.Info("-newNodeServer-")
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
	}
}

func (s3 *s3Driver) Run() {
	s3.logger.Info("-S3CSIDriver Run-")
	s3.logger.Info("Driver:", zap.Reflect("Driver Name", s3.name))
	s3.logger.Info("Version:", zap.Reflect("Driver Version", s3.vendorVersion))
	// Initialize default library driver

	s3.driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME})
	s3.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER})

	// Create GRPC servers
	s3.ids = s3.newIdentityServer(s3.driver)
	s3.ns = s3.newNodeServer(s3.driver)
	s3.cs = s3.newControllerServer(s3.driver)

	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(s3.endpoint, s3.ids, s3.cs, s3.ns)
	s.Wait()
}
