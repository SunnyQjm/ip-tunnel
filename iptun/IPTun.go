// Package iptun
// @Author: Jianming Que
// @Description:
// @Version: 1.0.0
// @Date: 2021/9/19 9:56 下午
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
package iptun

import (
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
	"log"
	"minlib/common"
	"net"
	"os/exec"
	"strconv"
	"sync"
)

type Config struct {
	InterfaceName string // The interface name want to create
	IPv4Addr      string // The IPv4 address assign to created interface
	MTU           int    // The MTU assign to created interface
	Mask          string // The network mask assign to created interface
	Route         []string
}

type IPPacket struct {
	Src        net.IP // 源ip地址
	Dst        net.IP // 目的ip地址
	RawPackets []byte // 原始IP包
}

type IPTun struct {
	iface   *water.Interface // TUN网卡对象，用于收发包
	packets []byte           // 缓存队列，用于保存从TUN收到的IP包
}

// NewIPTun 创建一个 IPTun
//
// @Description:
// @param config
//
func NewIPTun(config Config) (*IPTun, error) {
	ipTun := new(IPTun)
	if err := ipTun.Init(config); err != nil {
		return nil, err
	}
	return ipTun, nil
}

// Init 初始化 IPTun
//
// @Description:
// @receiver i
// @param config
// @return error
//
func (i *IPTun) Init(config Config) error {
	// 根据 MTU 设置抓包大小
	i.packets = make([]byte, config.MTU+200)

	// 创建一个 TUN 网卡
	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: config.InterfaceName,
		},
	})
	if err != nil {
		return err
	}
	i.iface = iface

	// 设置IP地址
	cmd := exec.Command("ifconfig", config.InterfaceName,
		config.IPv4Addr, "netmask", config.Mask, "mtu", strconv.Itoa(config.MTU))
	err = cmd.Run()
	if err != nil {
		log.Fatalln("Assign ip addr failed: ", err)
	}

	// 添加路由
	for i := 0; i < len(config.Route); i++ {
		cmd = exec.Command("route", "add", "-net", config.Route[i], "dev", config.InterfaceName)
		err = cmd.Run()
		if err != nil {
			common.LogFatal("add route error. cmd: ", cmd.String(), " ", err)
			return err
		}
	}
	return nil
}

// ReadIPPacketFromTun 从TUN中读取一个IP包
//
// @Description:
// @receiver i
// @param needCopy
// @return *IPPacket
// @return error
//
func (i *IPTun) ReadIPPacketFromTun(needCopy bool) (*IPPacket, error) {
	n, err := i.iface.Read(i.packets)
	if err != nil {
		return nil, err
	}
	ipPacket := new(IPPacket)
	ipPacket.Src = waterutil.IPv4Source(i.packets)
	ipPacket.Dst = waterutil.IPv4Destination(i.packets)
	if needCopy {
		ipPacket.RawPackets = make([]byte, n)
		copy(ipPacket.RawPackets, i.packets)
	} else {
		ipPacket.RawPackets = i.packets[:n]
	}
	return ipPacket, nil
}

// WriteIPPacketToTun 将一个 IP 包写入 TUN
//
// @Description:
// @receiver i
// @param packet
// @return error
//
func (i *IPTun) WriteIPPacketToTun(packet []byte) (int, error) {
	return i.iface.Write(packet)
}

// StartTunnel 启动隧道
//
// @Description:
// @receiver i
// @param adapter
//
func (i *IPTun) StartTunnel(adapter TunnelAdapter) {
	var wg sync.WaitGroup
	wg.Add(2)

	// 在一个单独的协程里面接收来自 TUN 的IP包
	go func(done func()) {
		defer done()
		for {
			ipPacket, err := i.ReadIPPacketFromTun(false)
			if err != nil {
				common.LogError(err)
			}
			if err := adapter.OnReceiveIPPktFromTun(ipPacket); err != nil {
				common.LogError(err)
			}
		}
	}(wg.Done)

	// 在单独的协程里面接收CPacket，并从中提取出 IP 包写入到TUN当中
	go func(done func()) {
		defer wg.Done()
		for {
			if ipPacket, err := adapter.ReadIPPkt(); err != nil {
				common.LogError(err)
			} else {
				if _, err := i.WriteIPPacketToTun(ipPacket.RawPackets); err != nil {
					common.LogError(err)
				}
			}

		}
	}(wg.Done)

	// 等待两个协程结束
	wg.Wait()
}
