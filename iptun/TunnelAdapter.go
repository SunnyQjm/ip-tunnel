// Package iptun
// @Author: Jianming Que
// @Description:
// @Version: 1.0.0
// @Date: 2021/9/21 10:18 上午
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
package iptun

// TunnelAdapter IP隧道适配器，定义了实现IP隧道所需要实现的方法
//  @Description:
//
type TunnelAdapter interface {
	// OnReceiveIPPktFromTun 每当从 TUN 收到一个IP包时，会调用这个回调进行处理
	//
	// @Description:
	//  1. 基于MIN或者NDN实现隧道时，要在这个回调里面对生成一个新的网络包，携带这个网络分组并通过 MIN/NDN 网络隧道到另一个边界路由器
	// @param packet
	//
	OnReceiveIPPktFromTun(packet *IPPacket) error

	// ReadIPPkt 从 MIN/NDN 中读取一个隧道过来的 IP 包
	//
	// @Description:
	//  1. 基于MIN或者NDN实现隧道时，需要在其中实现从 MIN/NDN 网络中接收一个携带有IP包的网络分组，并提取出其中的IP包
	// @return *IPPacket
	// @return error
	//
	ReadIPPkt() (*IPPacket, error)
}
