package scheduled

import (
	"context"
	"log/slog"

	"github.com/robfig/cron/v3"
)

type JobManager struct {
	cron    *cron.Cron
	jobs    []Job
	context context.Context
}

func NewJobManager(ctx context.Context) *JobManager {
	return &JobManager{
		cron:    cron.New(cron.WithSeconds()),
		jobs:    []Job{},
		context: ctx,
	}
}

func (jm *JobManager) RegisterJob(job Job) {
	jm.jobs = append(jm.jobs, job)
}

func (jm *JobManager) StartScheduler() {
	for _, job := range jm.jobs {
		schedule := job.Schedule()
		if _, err := jm.cron.AddFunc(schedule, func() {
			if err := job.Run(jm.context); err != nil {
				slog.Error("Error in job", "job", job.Name(), "error", err)
			} else {
				slog.Info("Job executed successfully", "job", job.Name())
			}
		}); err != nil {
			slog.Error("Failed to schedule job", "job", job.Name(), "error", err)
		}
	}
	jm.cron.Start()
}

func (jm *JobManager) Stop() {
	slog.Info("Stopping scheduled jobs...")
	jm.cron.Stop()
	slog.Info("Scheduled jobs stopped")
}
