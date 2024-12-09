package svc

import (
	"camp/internal/config"
	"camp/internal/job"
)

type ServiceContext struct {
	Config *config.Config
	Job    *job.Job
}

func NewServiceContext(c *config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Job:    job.NewJob(c),
	}
}
