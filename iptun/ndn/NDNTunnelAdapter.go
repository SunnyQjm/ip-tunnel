// Package ndn
// @Author: Jianming Que
// @Description:
// @Version: 1.0.0
// @Date: 2021/9/21 5:23 下午
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
package ndn

import (
	"ip-tunnel/iptun"
)

type NDNTunnelAdapter struct {
	receivePktFromTunChannel chan *iptun.IPPacket
	receivePktFromNDN        chan *iptun.IPPacket
}

func (N *NDNTunnelAdapter) Init() {
	N.receivePktFromTunChannel = make(chan *iptun.IPPacket, 10)
	N.receivePktFromNDN = make(chan *iptun.IPPacket, 10)
}

func (N *NDNTunnelAdapter) GetPktFromTun() *iptun.IPPacket {
	return <-N.receivePktFromTunChannel
}

func (N *NDNTunnelAdapter) OnReceivePktFromNDN(packet *iptun.IPPacket) {
	N.receivePktFromNDN <- packet
}

// OnReceiveIPPktFromTun 从 TUN 接收到IP包时，调用本回调
//
// @Description:
// @receiver N
// @param packet
// @return error
//
func (N *NDNTunnelAdapter) OnReceiveIPPktFromTun(packet *iptun.IPPacket) error {
	// 通过 chan 传递
	N.receivePktFromTunChannel <- packet
	return nil
}

// ReadIPPkt 从NDN接收隧道过来的IP包
//
// @Description:
// @receiver N
// @return *iptun.IPPacket
// @return error
//
func (N *NDNTunnelAdapter) ReadIPPkt() (*iptun.IPPacket, error) {
	ipPacket := <-N.receivePktFromNDN
	return ipPacket, nil
}
