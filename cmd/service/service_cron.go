package service

import (
	"context"
	"fmt"
	"github.com/IceBearAI/aigc/src/pkg/files"
	"github.com/IceBearAI/aigc/src/pkg/finetuning"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/oklog/pkg/group"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"os"
)

var (
	cronJobNames = []string{
		"finetuning.run-waiting-train",
		"finetuning.running-log",
		"deployment.status",
	}

	cronJobCmd = &cobra.Command{
		Use:               "cronjob command <args> [flags]",
		Short:             "定时任务",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
// ┌─────────────second 范围 (0 - 60)
// │ ┌───────────── min (0 - 59)
// │ │ ┌────────────── hour (0 - 23)
// │ │ │ ┌─────────────── day of month (1 - 31)
// │ │ │ │ ┌──────────────── month (1 - 12)
// │ │ │ │ │ ┌───────────────── day of week (0 - 6) (0 to 6 are Sunday to
// │ │ │ │ │ │                  Saturday)
// │ │ │ │ │ │
// │ │ │ │ │ │
// * * * * * *
// 每隔5秒执行一次：*/5 * * * * ?
// 每隔1分钟执行一次：0 */1 * * * ?
// 每天23点执行一次：0 0 23 * * ?
// 每天凌晨1点执行一次：0 0 1 * * ?
// 每月1号凌晨1点执行一次：0 0 1 1 * ?
// 每周一和周三晚上22:30: 00 30 22 * * 1,3
// 在26分、29分、33分执行一次：0 26,29,33 * * * ?
// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
// 每年三月的星期四的下午14:10和14:40:  00 10,40 14 ? 3 4

可用的配置类型：
[start]
aigc-server cronjob -h`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err = prepare(cmd.Context()); err != nil {
				return errors.Wrap(err, "prepare")
			}
			return nil
		},
	}

	cronJobStartCmd = &cobra.Command{
		Use:               "start <args> [flags]",
		Short:             "定时任务启动",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `如果 cronjob.auto 设置为 true 并且没有传入相应用的任务名称，则将自动运行所有的任务

aigc-server cronjob start -h`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cronJobAuto && len(args) == 0 {
				args = cronJobNames
			}
			if len(args) == 0 {
				return fmt.Errorf("unknown command: %s", args)
			}
			return cronStart(cmd.Context(), args)
		},
	}
)

func cronStart(ctx context.Context, args []string) (err error) {
	fileSvc = files.NewService(logger, traceId, store, apiSvc, []files.CreationOption{
		files.WithLocalDataPath(serverStoragePath),
		files.WithServerUrl(fmt.Sprintf("%s/storage", serverDomain)),
		files.WithStorageType("local"),
	}...)
	fineTuningSvc = finetuning.New(traceId, logger, store, fileSvc, apiSvc,
		finetuning.WithGpuTolerationValue(datasetsGpuToleration),
		finetuning.WithCallbackHost(serverDomain),
	)

	crontab := cron.New(cron.WithSeconds()) //精确到秒

	for _, v := range args {
		switch v {
		case "finetuning.run-waiting-train":
			entryID, err := crontab.AddJob("0 0/1 * * * *", &fineTuningRunWaitingTrainCronJob{
				fineTuning: fineTuningSvc,
				logger:     log.With(logger, "cron", "finetuning.run-waiting-train"),
				Name:       "finetuning.run-waiting-train",
				ctx:        ctx,
			})
			if err != nil {
				_ = level.Error(logger).Log("msg", "add cron job failed", "err", err.Error())
				return err
			}
			_ = level.Info(logger).Log("msg", "add cron job success", "entryID", entryID, "name", "finetuning.run-waiting-train")
		case "finetuning.running-log":
			entryID, err := crontab.AddJob("0/30 * * * * *", &fineTuningRunningLogCronJob{
				fineTuning: fineTuningSvc,
				logger:     log.With(logger, "cron", "finetuning.running-log"),
				Name:       "finetuning.running-log",
				ctx:        ctx,
				store:      store,
				runFun:     fineTuningRunningJobLog,
			})
			if err != nil {
				_ = level.Error(logger).Log("msg", "add cron job failed", "err", err.Error())
				return err
			}
			_ = level.Info(logger).Log("msg", "add cron job success", "entryID", entryID, "name", "finetuning.running-log")
		case "deployment.status":
			entryID, err := crontab.AddJob("0 0/1 * * * *", &deploymentStatusCronJob{
				logger: log.With(logger, "cron", "deployment.status"),
				Name:   "deployment.status",
				ctx:    ctx,
				store:  store,
				apiSvc: apiSvc,
			})
			if err != nil {
				_ = level.Error(logger).Log("msg", "add cron job failed", "err", err.Error())
				return err
			}
			_ = level.Info(logger).Log("msg", "add cron job success", "entryID", entryID, "name", "deployment.status")
		default:
			return fmt.Errorf("unknown command: %s", v)
		}
	}

	crontab.Start()
	g := &group.Group{}

	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer crontab.Stop()

	initCronHandler(ctx, g)
	initCancelInterrupt(ctx, g)
	_ = level.Error(logger).Log("cron server exit", g.Run())
	return nil
}

func initCronHandler(ctx context.Context, g *group.Group) {
	g.Add(func() error {
		_ = level.Info(logger).Log("msg", "cron server start")
		select {}
	}, func(err error) {
		_ = level.Info(logger).Log("msg", "cron server stop", "err", err)
		// TODO： 退出的时候需要判断是否有正在评估的任务，如果有则把评估的任务数据更新为取消状态
		os.Exit(0)
	})
}
