// Copyright 2018 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apiserver

import (
	// "strings"

	// openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/informers"

	"github.com/kubernetes-incubator/metrics-server/pkg/apiserver/generic"
	// generatedopenapi "github.com/kubernetes-incubator/metrics-server/pkg/generated/openapi"
	"github.com/kubernetes-incubator/metrics-server/pkg/version"
)

// Config contains configuration for launching an instance of metrics-server.
type Config struct {
	GenericConfig  *genericapiserver.Config
	ProviderConfig generic.ProviderConfig
}

type completedConfig struct {
	genericapiserver.CompletedConfig
	ProviderConfig *generic.ProviderConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete(informers informers.SharedInformerFactory) completedConfig {
	c.GenericConfig.Version = version.VersionInfo()

	// enable OpenAPI schemas
	// c.GenericConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(generatedopenapi.GetOpenAPIDefinitions, openapinamer.NewDefinitionNamer(generic.Scheme))
	// c.GenericConfig.OpenAPIConfig.Info.Title = "Kubernetes metrics-server"
	// c.GenericConfig.OpenAPIConfig.Info.Version = strings.Split(c.GenericConfig.Version.String(), "-")[0] // TODO(directxman12): remove this once autosetting this doesn't require security definitions

	return completedConfig{
		CompletedConfig: c.GenericConfig.Complete(informers),
		ProviderConfig:  &c.ProviderConfig,
	}
}

type MetricsServer struct {
	*genericapiserver.GenericAPIServer
}

// New returns a new instance of MetricsServer from the given config.
func (c completedConfig) New() (*MetricsServer, error) {
	genericServer, err := c.CompletedConfig.New("metrics-server", genericapiserver.NewEmptyDelegate()) // completion is done in Complete, no need for a second time
	if err != nil {
		return nil, err
	}

	if err := generic.InstallStorage(c.ProviderConfig, c.SharedInformerFactory.Core().V1(), genericServer); err != nil {
		return nil, err
	}

	return &MetricsServer{
		GenericAPIServer: genericServer,
	}, nil
}
