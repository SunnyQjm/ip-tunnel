// Package iptun
// @Author: Jianming Que
// @Description:
// @Version: 1.0.0
// @Date: 2021/9/20 11:42 上午
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
package iptun

import "gopkg.in/ini.v1"

type IPTunnelConfig struct {
	TunConfig `ini:"Tun"`
	MirConfig `ini:"Mir"`
	LogConfig `ini:"Log"`
}

func (i *IPTunnelConfig) Init() {
	i.TunConfig.InterfaceName = "tun0"
	i.TunConfig.IPv4Addr = "192.168.88.1"
	i.TunConfig.Mtu = 1500
	i.TunConfig.Mask = "255.255.255.0"
	i.TunConfig.Route = make([]string, 0)

	i.MirConfig.UnixSocketPath = "/tmp/mir.sock"
	i.MirConfig.ListenIdentifier = ""
	i.MirConfig.TargetIdentifier = ""
}

type TunConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// Tun
	////////////////////////////////////////////////////////////////////////////////////////////////
	InterfaceName string   `ini:"InterfaceName"`
	IPv4Addr      string   `ini:"IPv4Addr"`
	Mtu           int      `ini:"Mtu"`
	Mask          string   `ini:"Mask"`
	Route         []string `ini:"Route"`
}

type MirConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// Mir
	////////////////////////////////////////////////////////////////////////////////////////////////
	UnixSocketPath   string `ini:"UnixSocketPath"`
	ListenIdentifier string `ini:"ListenIdentifier"`
	TargetIdentifier string `ini:"TargetIdentifier"`
}

type LogConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// Log
	////////////////////////////////////////////////////////////////////////////////////////////////
	LogLevel string `ini:"LogLevel"`
}

// ParseConfig
// 解析配置文件
//
// @Description:
// @receiver m
// @param configPath
// @return error
//
func ParseConfig(configPath string) (*IPTunnelConfig, error) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, err
	}
	mirConfig := new(IPTunnelConfig)
	// 初始化配置，给所有的配置项设置默认值
	mirConfig.Init()
	// 加载配置文件中的配置
	if err = cfg.MapTo(&mirConfig); err != nil {
		return nil, err
	}
	return mirConfig, nil
}
