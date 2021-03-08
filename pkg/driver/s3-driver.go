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
	"errors"
	"github.com/IBM/satellite-object-storage-plugin/pkg/s3client"
	csi "github.com/container-storage-interface/spec/lib/go/csi"
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
	commonError "github.ibm.com/alchemy-containers/ibm-csi-common/pkg/messages"
	"go.uber.org/zap"
	"github.ibm.com/alchemy-containers/ibm-csi-common/pkg/utils"
)

const (
	kib    int64 = 1024
	mib    int64 = kib * 1024
	gib    int64 = mib * 1024
	gib10  int64 = gib * 10
	gib100 int64 = gib * 100
	tib    int64 = gib * 1024
	tib100 int64 = tib * 100
)

type s3Driver struct {
	name     string
	driver   *csicommon.CSIDriver
	s3client s3client.S3Client
	endpoint string

	ids           *identityServer
	ns            *nodeServer
	cs            *controllerServer
	vendorVersion string
	logger        *zap.Logger

	cap     []*csi.ControllerServiceCapability
	vc      []*csi.VolumeCapability_AccessMode
	nscap []*csi.NodeServiceCapability
}

type s3Volume struct {
	VolName string `json:"volName"`
	VolID   string `json:"volID"`
	VolSize int64  `json:"volSize"`
	VolPath string `json:"volPath"`
}

type s3VolumeSnapshot struct {
	Name      string `json:"name"`
	Id        string `json:"id"`
	VolID     string `json:"volID"`
	Path      string `json:"path"`
	CreateAt  int64  `json:"createAt"`
	SizeBytes int64  `json:"sizeBytes"`
}

var (
	s3CosVolumes         map[string]s3Volume
	s3CosVolumeSnapshots map[string]s3VolumeSnapshot
)

// GetS3CSIDriver ...
func GetS3CSIDriver() *s3Driver {
	return &s3Driver{}
}

func (csiDriver *s3Driver) Setups3Driver(lgr *zap.Logger, name, vendorVersion string) (*s3Driver, error) {
	csiDriver.logger = lgr
	csiDriver.logger.Info("S3CSIDriver-SetupS3CSIDriver setting up S3 CSI Driver...")

	if name == "" {
		return nil, fmt.Errorf("driver name missing")
	}

	// Setup messaging
	commonError.MessagesEn = commonError.InitMessages()

	csiDriver.name = name
	csiDriver.vendorVersion = vendorVersion

	csiDriver.logger.Info("Successfully setup IBM CSI driver")
	return csiDriver, nil
}

func (s3 *s3Driver) newIdentityServer(d *csicommon.CSIDriver) *identityServer {
	s3.logger.Info("-newIdentityServer-")
	return &identityServer{
		DefaultIdentityServer: csicommon.NewDefaultIdentityServer(d),
		s3Driver:                s3,
	}
}

func (s3 *s3Driver) newControllerServer(d *csicommon.CSIDriver) *controllerServer {
	s3.logger.Info("-newControllerServer-")
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(d),
		s3Driver: s3,
	}
}

func (s3 *s3Driver) newNodeServer(d *csicommon.CSIDriver) *nodeServer {
	s3.logger.Info("-newNodeServer-")
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
		s3Driver: s3,
	}
}

func  (csiDriver *s3Driver) NewS3CosDriver(nodeID string, endpoint string) (*s3Driver, error) {
	driver := csicommon.NewCSIDriver(csiDriver.name, csiDriver.vendorVersion, nodeID)
	if driver == nil {
		csiDriver.logger.Error("Failed to initialize CSI Driver")
		return nil, errors.New("failed to initialize CSI Driver")
	}
	s3client, err := s3client.NewS3Client("awss3")
	if err != nil {
		return nil, err
	}

	csiDriver.endpoint = endpoint
	csiDriver.driver = driver
	csiDriver.s3client =  s3client

	csiDriver.driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME})
	csiDriver.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER})

	/*csc := []csi.ControllerServiceCapability_RPC_Type{
                csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
        }
        _ = csiDriver.AddControllerServiceCapabilities(csc)

        ns := []csi.NodeServiceCapability_RPC_Type{
                csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
        }

        _ = csiDriver.AddNodeServiceCapabilities(ns)

	vcam := []csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
	}
	_ = csiDriver.AddVolumeCapabilityAccessModes(vcam)*/

        // Create GRPC servers
        csiDriver.ids = csiDriver.newIdentityServer(csiDriver.driver)
        csiDriver.ns = csiDriver.newNodeServer(csiDriver.driver)
        csiDriver.cs = csiDriver.newControllerServer(csiDriver.driver)

	return csiDriver, nil
}

func (s3 *s3Driver) Run() {
	s3.logger.Info("-S3CSIDriver Run-")
	s3.logger.Info("Driver:", zap.Reflect("Driver Name", s3.name))
	s3.logger.Info("Version:", zap.Reflect("Driver Version", s3.vendorVersion))
	// Initialize default library driver

	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(s3.endpoint, s3.ids, s3.cs, s3.ns)
	s.Wait()
}

func getVolumeByName(volName string) (s3Volume, error) {
	for _, s3CosVol := range s3CosVolumes {
		if s3CosVol.VolName == volName {
			return s3CosVol, nil
		}
	}
	return s3Volume{}, fmt.Errorf("volume name %s does not exit in the volumes list", volName)
}


// AddControllerServiceCapabilities ...
func (s3 *s3Driver) AddControllerServiceCapabilities(cl []csi.ControllerServiceCapability_RPC_Type) error {
	s3.logger.Info("IBMCSIDriver-AddControllerServiceCapabilities...", zap.Reflect("ControllerServiceCapabilities", cl))
	var csc []*csi.ControllerServiceCapability
	for _, c := range cl {
		s3.logger.Info("Adding controller service capability", zap.Reflect("Capability", c.String()))
		csc = append(csc, utils.NewControllerServiceCapability(c))
	}
	s3.cap = csc
	s3.logger.Info("Successfully added Controller Service Capabilities")
	return nil
}

// AddNodeServiceCapabilities ...
func (s3 *s3Driver) AddNodeServiceCapabilities(nl []csi.NodeServiceCapability_RPC_Type) error {
	s3.logger.Info("IBMCSIDriver-AddNodeServiceCapabilities...", zap.Reflect("NodeServiceCapabilities", nl))
	var nsc []*csi.NodeServiceCapability
	for _, n := range nl {
		s3.logger.Info("Adding node service capability", zap.Reflect("NodeServiceCapabilities", n.String()))
		nsc = append(nsc, utils.NewNodeServiceCapability(n))
	}
	s3.nscap = nsc
	s3.logger.Info("Successfully added Node Service Capabilities")
	return nil
}

// AddVolumeCapabilityAccessModes ...
func (s3 *s3Driver) AddVolumeCapabilityAccessModes(vcapa []csi.VolumeCapability_AccessMode_Mode) error {
	s3.logger.Info("IBMCSIDriver-AddVolumeCapabilityAccessModes...", zap.Reflect("VolumeCapabilityAccessModes", vcapa))
	var vca []*csi.VolumeCapability_AccessMode
	for _, c := range vcapa {
		s3.logger.Info("Enabling volume access mode", zap.Reflect("Mode", c.String()))
		vca = append(vca, utils.NewVolumeCapabilityAccessMode(c))
	}
	s3.vc = vca
	s3.logger.Info("Successfully enabled Volume Capability Access Modes")
	return nil
}
