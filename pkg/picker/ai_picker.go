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
	"path/filepath"
)

const ContainerInputFilePath = "/hecate/examples/"
const ImageName = `creactiviti/hecate`

type AiPicker struct {
}

func NewAiPicker() *AiPicker {
	return &AiPicker{}
}

func (p *AiPicker) Pick(ctx context.Context, values map[string]interface{}) error {
	file := values[`filepath`].(string)
	fmt.Println("=========开始智能剪辑文件:", file, "==========")
	cli, e := client.NewEnvClient()
	if e != nil {
		return e
	}
	file, _ = filepath.Abs(file)
	filename := filepath.Base(file)
	targetDir := filepath.Dir(file)
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
	fmt.Println("=========智能剪辑文件:", file, "结束==========")
	return nil
}
