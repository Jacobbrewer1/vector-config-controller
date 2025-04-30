package main

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/jacobbrewer1/vector-config-controller/pkg/k8s"
	"github.com/jacobbrewer1/vector-config-controller/pkg/vector"
	"github.com/jacobbrewer1/web/logging"
)

// Reconcile is the main reconciliation loop for the application.
func (a *App) Reconcile(ctx context.Context) error {
	l := logging.LoggerWithComponent(a.base.Logger(), "reconcile")

	tick := time.NewTicker(a.config.TickerInterval)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if !a.base.IsLeader() {
				l.Info("not leader, skipping reconciliation")
				continue
			}

			l.Info("reconciling")

			if err := reconcile(
				ctx,
				a.base.KubeClient(),
			); err != nil {
				l.Error("error reconciling", logging.KeyError, err)
			}
		case <-ctx.Done():
			l.Info("reconciler closing")
			return nil
		}
	}
}

// reconcile represents one iteration of the reconciliation process.
func reconcile(
	ctx context.Context,
	kubeClient kubernetes.Interface,
) error {
	t := prometheus.NewTimer(iterationsHistogram)
	defer t.ObserveDuration()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	agentConfig, err := vectorAgentConfig()
	if err != nil {
		return err
	}

	if err := k8s.UpsertResource(
		ctx,
		kubeClient,
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "vector-logs-config",
				Namespace: "vector",
				Labels: map[string]string{
					"owner": "vector-config-reconciler",
				},
			},
			Data: map[string]string{
				"config.json": agentConfig,
			},
		},
	); err != nil {
		return fmt.Errorf("failed to upsert resource: %w", err)
	}

	return nil
}

func vectorAgentConfig() (string, error) {
	vCfg := vector.NewConfig()

	configForMetrics(vCfg)
	configForLogs(vCfg)

	return vCfg.JSON()
}
