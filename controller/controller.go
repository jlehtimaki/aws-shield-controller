package controller

import (
	"context"
	"github.com/jlehtimaki/aws-shield-controller/pkg/aws"
	"github.com/jlehtimaki/aws-shield-controller/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	"time"
)

type Controller struct {
	// The interval between individual synchronizations
	Interval time.Duration
	// The nextRunAt used for throttling and batching reconciliation
	nextRunAt time.Time
}

func (c *Controller) runReconcile() error {
	ingressList, err := kubernetes.GetIngresses()
	if err != nil {
		return err
	}

	err = aws.EnableAWSShield(ingressList)
	if err != nil {
		return err
	}
	return nil
}

// Run runs RunOnce in a loop with a delay until context is canceled
func (c *Controller) Run(ctx context.Context) {
	ticker := time.NewTicker(c.Interval)
	for {
		select {
		case <-ticker.C:
			err := c.runReconcile()
			if err != nil {
				log.Error(err)
			}
		case <-ctx.Done():
			log.Info("Terminating main controller loop")
			return
		}
	}
}
