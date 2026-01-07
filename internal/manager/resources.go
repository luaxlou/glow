package manager

import (
	"fmt"

	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/pkg/api"
)

func ListResources() ([]api.ResourceRef, error) {
	var hostCfg api.Host
	if err := configmanager.GetHostConfig(&hostCfg); err != nil {
		// No host config implies no managed resources yet?
		return []api.ResourceRef{}, nil
	}

	var resources []api.ResourceRef
	for name, svc := range hostCfg.Spec.Services {
		// Heuristic to determine Kind based on name or port?
		// Host spec map key is name (e.g. "mysql").
		// We treat key as Kind/Name combo?
		// Usually name is "mysql", "redis".
		kind := name // Simplified
		
		resources = append(resources, api.ResourceRef{
			Kind: kind,
			Name: fmt.Sprintf("%s-local", kind), // local instance
			Port: svc.Port,
		})
	}
	return resources, nil
}
