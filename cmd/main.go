package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tangxusc/video-picker/pkg/config"
	"github.com/tangxusc/video-picker/pkg/downloader"
	"github.com/tangxusc/video-picker/pkg/picker"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	newCommand := NewCommand(ctx)
	HandlerNotify(cancel)

	_ = newCommand.Execute()
	cancel()
}

func NewCommand(ctx context.Context) *cobra.Command {
	var command = &cobra.Command{
		Use:   "start",
		Short: "start picker",
		RunE: func(cmd *cobra.Command, args []string) error {
			rand.Seed(time.Now().Unix())
			config.InitLog()

			cancel, _ := context.WithCancel(ctx)
			aiPicker := picker.NewAiPicker(cancel, 1)
			d := downloader.NewHuyaDownloader(cancel, 1, aiPicker)

			_, cancelFunc := d.Download(`11336726`)
			time.Sleep(time.Second * 6)
			cancelFunc()

			<-ctx.Done()
			return nil
		},
	}
	logrus.SetFormatter(&logrus.TextFormatter{})
	config.BindParameter(command)

	return command
}

func HandlerNotify(cancel context.CancelFunc) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		cancel()
	}()
}
