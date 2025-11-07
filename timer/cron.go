package timer

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"ppt/log"
	"ppt/util"
	"sync"
	"time"
)

var (
	cronInstance *cron.Cron
	cronMgr      = sync.Map{}
	once         sync.Once
)

// CreateCron 创建定时任务
func CreateCron(key, spec, timeZone string, task func()) error {
	once.Do(func() {
		cronInstance = cron.New(
			cron.WithSeconds(),
			cron.WithLocation(util.GetTz(timeZone)),
		)
		cronInstance.Start()
	})

	jobID, err := cronInstance.AddJob(spec, cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(cron.FuncJob(func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Cron job run error", zap.Any("err", err), zap.Stack("stack"))
				return
			}
		}()
		task()
	})))
	if err != nil {
		log.Error("Cron add cron job error", zap.Any("err", err), zap.String("key", key), zap.String("spec", spec))
		return err
	}

	cronMgr.Store(key, jobID)
	return nil
}

func CloseCron() {
	cronMgr.Range(func(key, jobID any) bool {
		cronInstance.Remove(jobID.(cron.EntryID))
		return true
	})
	ctx := cronInstance.Stop()
	select {
	case <-ctx.Done():
		return
	case <-time.After(30 * time.Second):
		log.Info("Cron stop timeout")
		return
	}
}
