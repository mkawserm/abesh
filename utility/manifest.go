package utility

import (
	"github.com/mkawserm/abesh/model"
)

// MergeManifest merges `from` manifest to `to` manifest
func MergeManifest(to, from *model.Manifest) *model.Manifest {
	if from == nil {
		return to
	}

	fromCapabilities := make(map[string]*model.CapabilityManifest)

	for _, capabilityManifest := range from.Capabilities {
		fromCapabilities[capabilityManifest.ContractId] = capabilityManifest
	}

	for _, capabilityManifest := range to.Capabilities {
		fromCapability, found := fromCapabilities[capabilityManifest.ContractId]
		if found {
			for k, v := range fromCapability.Values {
				capabilityManifest.Values[k] = v
			}
		}
	}

	return to
}
