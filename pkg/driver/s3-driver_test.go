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
	"testing"

	"github.com/stretchr/testify/assert"
	cloudProvider "github.ibm.com/alchemy-containers/ibm-csi-common/pkg/ibmcloudprovider"
)

func inits3Driver(t *testing.T) *s3Driver {
	vendorVersion := "test-vendor-version-1.1.2"
	driver := "mydriver"

	// Creating test logger
	logger, teardown := cloudProvider.GetTestLogger(t)
	defer teardown()
	icDriver := GetS3CSIDriver()

	// Setup the IBM CSI driver
	err := icDriver.Setups3Driver(logger, driver, vendorVersion)
	if err != nil {
		t.Fatalf("Failed to setup IBM CSI Driver: %v", err)
	}

	return icDriver
}

func TestSetups3Driver(t *testing.T) {
	// success setting up driver
	driver := inits3Driver(t)
	assert.NotNil(t, driver)

	// common code
	// Creating test logger
	vendorVersion := "test-vendor-version-1.1.2"
	name := ""
	logger, teardown := cloudProvider.GetTestLogger(t)
	defer teardown()
	icDriver := GetS3CSIDriver()

	// Failed setting up driver, name  nil
	err := icDriver.Setups3Driver(logger, name, vendorVersion)
	assert.NotNil(t, err)
}
