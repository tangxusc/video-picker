package picker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"path/filepath"
	"testing"
)

func TestNewAiPicker(t *testing.T) {
	picker := NewAiPicker(context.TODO(), 1)
	picker.Pick(`../../video/武林萌主唐小姐/2019-12-06-13-40.mp4`)
	select {}
}

func TestDockerClient(t *testing.T) {
	cli, e := client.NewEnvClient()
	if e != nil {
		panic(e)
	}
	abs, e := filepath.Abs(".")
	fmt.Println(abs, e)
	return
	config := &container.Config{
		Cmd: strslice.StrSlice{"hecate", "-i", "/hecate/examples/2019-12-06-13-40.mp4",
			`-o`, `/hecate/examples/`, `--generate_mov`, `--lmov`, `30`, `-w`, `1080`},
		Image: "creactiviti/hecate",
	}

	createResp, e := cli.ContainerCreate(context.TODO(), config, &container.HostConfig{
		Binds: []string{
			//`/home/tangxu/openProject/video-picker/video/武林萌主唐小姐:/hecate/examples/:rw`,
			`./video-picker/video/武林萌主唐小姐:/hecate/examples/:rw`,
		},
		AutoRemove: true,
	}, nil, "")
	if e != nil {
		panic(e)
	}
	e = cli.ContainerStart(context.TODO(), createResp.ID, types.ContainerStartOptions{})
	if e != nil {
		panic(e)
	}

	statusCh, errCh := cli.ContainerWait(context.TODO(), createResp.ID)
	fmt.Println(statusCh, errCh)
}
