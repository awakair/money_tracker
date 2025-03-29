package cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

type CronService struct{}

func NewCronService() *CronService {
	return &CronService{}
}

func (s *CronService) GetNextGenerationTime(expression string, after time.Time) (time.Time, error) {
	schedule, err := cron.Parse(expression)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse cron expression: %w", err)
	}

	nextTime := schedule.Next(after)
	return nextTime, nil
}
