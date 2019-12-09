package picker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/video-picker/pkg/config"
	"github.com/tangxusc/video-picker/pkg/dispatcher"
	"path/filepath"
)

const ContainerInputFilePath = "/hecate/examples/"
const ImageName = `creactiviti/hecate`

type AiPicker struct {
	dispatcher *dispatcher.Dispatcher
}

func NewAiPicker(ctx context.Context, maxCount int) *AiPicker {
	dis := dispatcher.NewDispatcher(ctx, maxCount)
	dis.Start()
	return &AiPicker{
		dispatcher: dis,
	}
}

func (a *AiPicker) Pick(target string) {
	job := dispatcher.NewJob(func(ctx context.Context) error {
		return pick(ctx, target)
	})
	a.dispatcher.Dispatch(job)
}

func pick(ctx context.Context, target string) error {
	fmt.Println("=========开始智能剪辑文件:", target, "==========")
	cli, e := client.NewEnvClient()
	if e != nil {
		return e
	}

	target, _ = filepath.Abs(target)
	filename := filepath.Base(target)
	targetDir := filepath.Dir(target)
	containerFilePath := ContainerInputFilePath

	logrus.Debugf("filename:%v,targetDir:%v,containerFilePath:%v", filename, targetDir, containerFilePath)

	createResp, e := cli.ContainerCreate(ctx,
		&container.Config{
			Cmd: strslice.StrSlice{"hecate", "-i", filepath.Join(ContainerInputFilePath, filename),
				`-o`, `/hecate/examples/`, `--generate_mov`, `--lmov`, config.Instance.Picker.OutTime, `-w`, config.Instance.Picker.OutWidth},
			Image: ImageName,
		},
		&container.HostConfig{
			Binds: []string{
				fmt.Sprintf(`%s:%s`, targetDir, containerFilePath),
			},
			AutoRemove: true,
		},
		nil,
		filename)
	if e != nil {
		return e
	}
	e = cli.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{})
	if e != nil {
		return e
	}

	status, e := cli.ContainerWait(context.TODO(), createResp.ID)
	if e != nil {
		return e
	}
	logrus.Infof(`容器运行结果:%v`, status)
	fmt.Println("=========智能剪辑文件:", target, "结束==========")
	return nil
}
