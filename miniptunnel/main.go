// Package ip_tunnel
// @Author: Jianming Que
// @Description:
// @Version: 1.0.0
// @Date: 2021/9/15 9:00 下午
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
package main

import "C"
import (
	"github.com/urfave/cli/v2"
	"ip-tunnel/iptun"
	"ip-tunnel/iptun/min"
	"minlib/common"
	"os"
)

// StartIPTunnel 启动 IP Tunnel
//
// @Description:
// @param ipTunnelConfig
//
func StartIPTunnel(ipTunnelConfig *iptun.IPTunnelConfig) {
	ipTun, err := iptun.NewIPTun(iptun.Config{
		InterfaceName: ipTunnelConfig.TunConfig.InterfaceName,
		IPv4Addr:      ipTunnelConfig.TunConfig.IPv4Addr,
		MTU:           ipTunnelConfig.TunConfig.Mtu,
		Mask:          ipTunnelConfig.TunConfig.Mask,
		Route:         ipTunnelConfig.TunConfig.Route,
	})

	// 创建一个 MIN 隧道适配器
	minAdapter, err := min.NewMINTunnelAdapter(ipTunnelConfig)
	if err != nil {
		common.LogFatal(err)
	}

	// 启动IP隧道程序
	ipTun.StartTunnel(minAdapter)
}

const defaultConfigFilePath = "/usr/local/etc/mir/iptunconf.ini"

func main() {
	var configFilePath string
	mirApp := cli.NewApp()
	mirApp.Name = "miniptunnel"
	mirApp.Usage = " A IP tunnel over MIN "
	mirApp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "f",
			Value:       defaultConfigFilePath,
			Usage:       "Config file path for MIR",
			Destination: &configFilePath,
			Required:    true,
		},
	}
	mirApp.Action = func(context *cli.Context) error {
		ipTunnelConfig, err := iptun.ParseConfig(configFilePath)
		if err != nil {
			common.LogFatal(err)
		}

		// 初始化日志模块
		var loggerParameters common.LoggerParameters
		loggerParameters.LogLevel = ipTunnelConfig.LogConfig.LogLevel
		loggerParameters.ReportCaller = true
		common.InitLogger(&loggerParameters)
		StartIPTunnel(ipTunnelConfig)
		return nil
	}

	if err := mirApp.Run(os.Args); err != nil {
		return
	}
}
