package envoyfilter

import (
	"fmt"

	meshconfig "istio.io/api/mesh/v1alpha1"
	networking "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/config/memory"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/mesh"
	"istio.io/istio/pkg/config/schema/collections"
	"istio.io/istio/pkg/config/schema/gvk"
	"istio.io/istio/pkg/log"
)

func newTestEnvironment(serviceDiscovery model.ServiceDiscovery, meshConfig *meshconfig.MeshConfig,
	configStore model.ConfigStore,
) *model.Environment {
	e := &model.Environment{
		ServiceDiscovery: serviceDiscovery,
		ConfigStore:      configStore,
		Watcher:          mesh.NewFixedWatcher(meshConfig),
	}

	pushContext := model.NewPushContext()
	e.Init()
	_ = pushContext.InitContext(e, nil, nil)
	e.SetPushContext(pushContext)
	return e
}

func buildEnvoyFilterConfigStore(configPatches []*networking.EnvoyFilter_EnvoyConfigObjectPatch) model.ConfigStore {
	store := memory.Make(collections.Pilot)

	for i, cp := range configPatches {
		_, err := store.Create(config.Config{
			Meta: config.Meta{
				Name:             fmt.Sprintf("test-envoyfilter-%d", i),
				Namespace:        "not-default",
				GroupVersionKind: gvk.EnvoyFilter,
			},
			Spec: &networking.EnvoyFilter{
				ConfigPatches: []*networking.EnvoyFilter_EnvoyConfigObjectPatch{cp},
			},
		})
		if err != nil {
			log.Errorf("create envoyfilter failed %v", err)
		}
	}
	return store
}
