package scheduled

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
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
				log.Printf("Error in job %s: %v", job.Name(), err)
			} else {
				log.Printf("Job %s executed successfully", job.Name())
			}
		}); err != nil {
			log.Printf("Failed to schedule job %s: %v", job.Name(), err)
		}
	}
	jm.cron.Start()
}
