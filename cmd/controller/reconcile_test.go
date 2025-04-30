package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVectorConfig(t *testing.T) {
	expectedConfig := `{
		"sources": {
			"host_metrics": {
				"type": "host_metrics",
				"filesystem": {
					"devices": {
						"exclude": ["binfmt_misc"]
					},
					"filesystems": {
						"exclude": ["binfmt_misc"]
					},
					"mountpoints": {
						"exclude": ["*/proc/sys/fs/binfmt_misc"]
					}
				}
			},
			"internal_metrics": {
				"type": "internal_metrics"
			},
			"kubernetes_logs": {
				"type": "kubernetes_logs"
			}
		},
		"sinks": {
			"prometheus_exporter": {
				"type": "prometheus_exporter",
				"inputs": ["host_metrics", "internal_metrics"],
				"address": "0.0.0.0:9090"
			},
			"loki_logs": {
				"type": "loki",
				"inputs": ["kubernetes_logs"],
				"endpoint": "http://loki-distributor.loki.svc.cluster.local:3100",
				"out_of_order_action": "accept",
				"acknowledgements": {
					"enabled": true
				},
				"encoding": {
					"codec": "json"
				},
				"request": {
					"concurrency": "adaptive"
				},
				"labels": {
					"pod_labels_*": "{{ kubernetes.pod_labels }}",
					"*": "{{ metadata }}",
					"source": "vector",
					"vector_instance": "inf-${HOSTNAME}",
					"tenant_id": "vector"
				}
			}
		}
	}`

	config, err := vectorAgentConfig()
	require.NoError(t, err)
	require.JSONEq(t, expectedConfig, config)
}
