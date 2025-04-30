package main

import "github.com/jacobbrewer1/vector-config-controller/pkg/vector"

func configForLogs(vCfg *vector.Config) {
	vCfg.AddSourceUntyped("kubernetes_logs", map[string]any{
		"type": "kubernetes_logs",
	})

	vCfg.AddSinkUntyped("loki_logs", map[string]any{
		"type": "loki",
		"inputs": []string{
			"kubernetes_logs",
		},
		"endpoint":            "http://loki-distributor.loki.svc.cluster.local:3100",
		"out_of_order_action": "accept",
		"acknowledgements": map[string]any{
			"enabled": true,
		},
		"encoding": map[string]any{
			"codec": "json",
		},
		"request": map[string]any{
			"concurrency": "adaptive",
		},
		"labels": map[string]any{
			"pod_labels_*":    "{{ kubernetes.pod_labels }}",
			"*":               "{{ metadata }}",
			"source":          "vector",
			"vector_instance": "inf-${HOSTNAME}",
			"tenant_id":       "vector",
		},
	})
}
