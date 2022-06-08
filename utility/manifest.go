package utility

import (
	"github.com/mkawserm/abesh/model"
)

// MergeManifest merges `from` manifest to `to` manifest
func MergeManifest(to, from *model.Manifest) *model.Manifest {
	if from == nil {
		return to
	}

	var fromCapabilities = make(map[string]*model.CapabilityManifest)
	for _, capabilityManifest := range from.Capabilities {
		if len(capabilityManifest.NewContractId) == 0 {
			fromCapabilities[capabilityManifest.ContractId] = capabilityManifest
		} else {
			fromCapabilities[capabilityManifest.NewContractId] = capabilityManifest
		}
	}

	for _, capabilityManifest := range to.Capabilities {
		var fromCapability *model.CapabilityManifest
		var found bool

		if len(capabilityManifest.NewContractId) == 0 {
			fromCapability, found = fromCapabilities[capabilityManifest.ContractId]
		} else {
			fromCapability, found = fromCapabilities[capabilityManifest.NewContractId]
		}

		if found {
			for k, v := range fromCapability.Values {
				if capabilityManifest.Values == nil {
					capabilityManifest.Values = make(map[string]string)
				}
				capabilityManifest.Values[k] = v
			}
		}
	}

	return to
}
