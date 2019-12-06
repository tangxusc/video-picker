package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const debugArgName = "debug"

func InitLog() {
	if viper.GetBool(debugArgName) {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
		logrus.Debug("已开启debug模式...")
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	Instance.Debug = viper.GetBool(debugArgName)
}

func BindParameter(cmd *cobra.Command) {
	viper.SetEnvPrefix("picker")
	viper.AutomaticEnv()

	cmd.PersistentFlags().BoolVarP(&Instance.Debug, debugArgName, "v", false, "debug mod")
	cmd.PersistentFlags().IntVarP(&Instance.Downloader.TimeOut, "downloaderTimeout", "", 10, "下载器超时时间")
	cmd.PersistentFlags().StringVarP(&Instance.Downloader.OutPath, "downloaderOutPath", "", "./video", "输出目录")

	cmd.PersistentFlags().StringVarP(&Instance.Picker.OutTime, "pickerOutTime", "", "30", "智能输出时间")
	cmd.PersistentFlags().StringVarP(&Instance.Picker.OutWidth, "pickerOutWidth", "", "1080", "智能输出宽度")

	_ = viper.BindPFlag(debugArgName, cmd.PersistentFlags().Lookup(debugArgName))
	_ = viper.BindPFlag("downloaderTimeout", cmd.PersistentFlags().Lookup("downloaderTimeout"))
	_ = viper.BindPFlag("downloaderOutPath", cmd.PersistentFlags().Lookup("downloaderOutPath"))

	_ = viper.BindPFlag("pickerOutTime", cmd.PersistentFlags().Lookup("pickerOutTime"))
	_ = viper.BindPFlag("pickerOutWidth", cmd.PersistentFlags().Lookup("pickerOutWidth"))
}

type Config struct {
	Debug      bool
	Downloader *DownloaderConfig
	Picker     *PickerConfig
}
type DownloaderConfig struct {
	TimeOut int
	OutPath string
}
type PickerConfig struct {
	OutWidth string
	OutTime  string
}

var Instance = &Config{
	Downloader: &DownloaderConfig{},
	Picker:     &PickerConfig{},
}
