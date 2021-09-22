// Package min
// @Author: Jianming Que
// @Description:
// @Version: 1.0.0
// @Date: 2021/9/21 10:27 上午
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
package min

import (
	"fmt"
	"github.com/songgao/water/waterutil"
	"ip-tunnel/iptun"
	"minlib/common"
	"minlib/component"
	"minlib/logicface"
	"minlib/packet"
)

type MINTunnelAdapter struct {
	logicFace logicface.LogicFace
	config    *iptun.IPTunnelConfig
}

// NewMINTunnelAdapter 创建一个 MINTunnelAdapter
//
// @Description:
// @param config
// @return *MINTunnelAdapter
// @return error
//
func NewMINTunnelAdapter(config *iptun.IPTunnelConfig) (*MINTunnelAdapter, error) {
	adapter := new(MINTunnelAdapter)
	adapter.config = config
	return adapter, adapter.Init(config)
}

// Init 初始化 MINTunnelAdapter
//
// @Description:
// @receiver M
// @param config
// @return error
//
func (M *MINTunnelAdapter) Init(config *iptun.IPTunnelConfig) error {
	// 连接上 MIR
	if err := M.logicFace.InitWithUnixSocket(config.MirConfig.UnixSocketPath); err != nil {
		return err
	}

	// 注册前缀监听
	if listenIdentifier, err := component.CreateIdentifierByString(config.MirConfig.ListenIdentifier); err != nil {
		return err
	} else {
		if err := M.logicFace.RegisterIdentifier(listenIdentifier, 1000); err != nil {
			return err
		}
	}
	return nil
}

// OnReceiveIPPktFromTun 处理从TUN网卡接收到的IP包
//
// @Description:
//	1. 生成一个 CPacket，将IP包放进去
//	2. 然后将 CPacket 发出即可
// @receiver M
// @param packet
//
func (M *MINTunnelAdapter) OnReceiveIPPktFromTun(ipPacket *iptun.IPPacket) error {
	srcIdentifier, err := component.CreateIdentifierByString(M.config.ListenIdentifier)
	if err != nil {
		return err
	}
	dstIdentifier, err := component.CreateIdentifierByString(M.config.TargetIdentifier)
	if err != nil {
		return err
	}

	// 在这边构造CPacket发出
	cPacket := new(packet.CPacket)
	cPacket.SetTtl(5)
	cPacket.SetSrcIdentifier(srcIdentifier)
	cPacket.SetDstIdentifier(dstIdentifier)
	cPacket.Payload.SetValue(ipPacket.RawPackets)

	common.LogDebug(fmt.Sprintf("Packet Received: %v -> %v \t %x\n", ipPacket.Src.String(), ipPacket.Dst.String(),
		ipPacket.RawPackets))
	if err := M.logicFace.SendCPacket(cPacket); err != nil {
		return err
	}
	return nil
}

// ReadIPPkt 从MIN网络中接收携带 IP 包的 CPacket，并
//
// @Description:
// @receiver M
// @return *IPPacket
// @return error
//
func (M *MINTunnelAdapter) ReadIPPkt() (*iptun.IPPacket, error) {
	cPacket, err := M.logicFace.ReceiveCPacket(4000)
	if err != nil {
		return nil, err
	} else {
		common.LogDebug(fmt.Sprintf("Write %d bytes, %s -> %s, %x", len(cPacket.Payload.GetValue()),
			waterutil.IPv4Source(cPacket.Payload.GetValue()),
			waterutil.IPv4Destination(cPacket.Payload.GetValue()), cPacket.Payload.GetValue()))
	}
	return &iptun.IPPacket{
		Src:        waterutil.IPv4Source(cPacket.Payload.GetValue()),
		Dst:        waterutil.IPv4Destination(cPacket.Payload.GetValue()),
		RawPackets: cPacket.Payload.GetValue(),
	}, nil
}