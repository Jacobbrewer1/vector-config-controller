package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/jacobbrewer1/vector-config-controller/pkg/vector"
)

var iterationsHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "vector_config_controller_iterations_total",
	Help: "The seconds taken to process iterations of this reconciler.",
})

func configForMetrics(vCfg *vector.Config) {
	vCfg.AddSourceUntyped("host_metrics", map[string]any{
		"type": "host_metrics",
		"filesystem": map[string]any{
			"devices": map[string]any{
				"exclude": []string{
					"binfmt_misc",
				},
			},
			"filesystems": map[string]any{
				"exclude": []string{
					"binfmt_misc",
				},
			},
			"mountpoints": map[string]any{
				"exclude": []string{
					"*/proc/sys/fs/binfmt_misc",
				},
			},
		},
	})

	vCfg.AddSourceUntyped("internal_metrics", map[string]any{
		"type": "internal_metrics",
	})

	vCfg.AddSinkUntyped("prometheus_exporter", map[string]any{
		"type": "prometheus_exporter",
		"inputs": []string{
			"host_metrics",
			"internal_metrics",
		},
		"address": "0.0.0.0:9090",
	})
}
