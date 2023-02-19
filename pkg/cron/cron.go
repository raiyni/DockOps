package cron

import (
	"time"

	"github.com/go-co-op/gocron"
)

func cron(task interface{}) {
	s := gocron.NewScheduler(time.UTC)
	s.TagsUnique()

	s.Every(1).Minute().Tag("foo").Do(task)

	s.StartBlocking()
}
