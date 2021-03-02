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
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
)

type identityServer struct {
	*csicommon.DefaultIdentityServer
}
