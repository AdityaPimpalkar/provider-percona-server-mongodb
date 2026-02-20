// Copyright (C) 2026 The OpenEverest Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package common defines shared metadata and topology definitions for the PSMDB provider.
package common

import (
	"github.com/openeverest/openeverest/v2/provider-runtime/controller"
	"github.com/openeverest/provider-percona-server-mongodb/types"
)

// Component name and type constants for PSMDB
const (
	ComponentEngine       = "engine"
	ComponentConfigServer = "configServer"
	ComponentProxy        = "proxy"
	ComponentBackupAgent  = "backupAgent"
	ComponentMonitoring   = "monitoring"

	ComponentTypeMongod = "mongod"
	ComponentTypeBackup = "backup"
	ComponentTypePMM    = "pmm"
)

// PSMDBMetadata returns the metadata for the PSMDB provider.
// This defines the component types, versions, components, and topologies
// that the provider supports.
//
// This metadata is use by runtime metadata accessed via c.Metadata(),
// and also for CLI generation of the provider manifest.
//
// The topologies are derived from the shared PSMDBTopologyDefinitions()
// to maintain a single source of truth across all provider implementations.
func PSMDBMetadata() *controller.ProviderMetadata {
	// Define component types and logical components
	metadata := &controller.ProviderMetadata{
		// ComponentTypes defines the available component types with their versions.
		// Each component type represents a different image/binary that can be deployed.
		ComponentTypes: map[string]controller.ComponentTypeMeta{
			// mongod is the main MongoDB server component
			ComponentTypeMongod: {
				Versions: []controller.ComponentVersionMeta{
					{Version: "6.0.19-16", Image: "percona/percona-server-mongodb:6.0.19-16-multi"},
					{Version: "6.0.21-18", Image: "percona/percona-server-mongodb:6.0.21-18"},
					{Version: "7.0.18-11", Image: "percona/percona-server-mongodb:7.0.18-11"},
					{Version: "8.0.4-1", Image: "percona/percona-server-mongodb:8.0.4-1-multi"},
					{Version: "8.0.8-3", Image: "percona/percona-server-mongodb:8.0.8-3", Default: true},
				},
			},
			// backup is the backup agent component
			ComponentTypeBackup: {
				Versions: []controller.ComponentVersionMeta{
					{Version: "2.9.1", Image: "percona/percona-backup-mongodb:2.9.1", Default: true},
				},
			},
			// pmm is the Percona Monitoring and Management component
			ComponentTypePMM: {
				Versions: []controller.ComponentVersionMeta{
					{Version: "2.44.1", Image: "percona/pmm-server:2.44.1", Default: true},
				},
			},
		},

		// Components defines the logical components that use the component types.
		// Multiple components can reference the same component type (e.g., engine and configServer both use mongod).
		Components: map[string]controller.ComponentMeta{
			ComponentEngine:       {Type: ComponentTypeMongod}, // Main database engine
			ComponentConfigServer: {Type: ComponentTypeMongod}, // Config server for sharded clusters
			ComponentProxy:        {Type: ComponentTypeMongod}, // Proxy/mongos for sharded clusters
			ComponentBackupAgent:  {Type: ComponentTypeBackup}, // Backup agent
			ComponentMonitoring:   {Type: ComponentTypePMM},    // Monitoring agent
		},
	}

	// Derive topologies from the shared topology definitions
	metadata.Topologies = controller.TopologiesFromSchemaProvider(PSMDBTopologyDefinitions())

	return metadata
}

// PSMDBTopologyDefinitions returns the topology definitions for PSMDB.
// This is shared by all provider implementations to maintain a single source of truth.
func PSMDBTopologyDefinitions() map[string]controller.TopologyDefinition {
	return map[string]controller.TopologyDefinition{
		string(types.TopologyTypeReplicaSet): {
			Schema: &types.ReplicaSetTopologyConfig{},
			Components: map[string]controller.TopologyComponentDefinition{
				ComponentEngine: {
					Optional: false,
					Defaults: map[string]any{"replicas": 3},
				},
				ComponentBackupAgent: {Optional: true},
				ComponentMonitoring:  {Optional: true},
			},
		},
		string(types.TopologyTypeSharded): {
			Schema: &types.ShardedTopologyConfig{},
			Components: map[string]controller.TopologyComponentDefinition{
				ComponentEngine: {
					Optional: false,
					Defaults: map[string]any{"replicas": 3},
				},
				ComponentProxy:        {Optional: false},
				ComponentConfigServer: {Optional: false},
				ComponentBackupAgent:  {Optional: true},
				ComponentMonitoring:   {Optional: true},
			},
		},
	}
}
