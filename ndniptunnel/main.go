// Package main
// @Author: Jianming Que
// @Description:
// @Version: 1.0.0
// @Date: 2021/9/21 5:58 下午
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
package main

/*
#cgo CXXFLAGS: -std=c++14
#cgo LDFLAGS: -lpthread -ljsoncpp -lpcap -lndn-cxx -lboost_system -lboost_thread -lcrypto
#include "bridge.h"
#include "common.h"
*/
import "C"
import (
	"github.com/songgao/water/waterutil"
	"github.com/urfave/cli/v2"
	"ip-tunnel/iptun"
	"ip-tunnel/iptun/ndn"
	"minlib/common"
	"os"
	"sync/atomic"
	"time"
	"unsafe"
)

var ipTun iptun.IPTun
var adapter ndn.NDNTunnelAdapter
var ipTunnelConfig *iptun.IPTunnelConfig

// 一个队列，缓存TUN网卡收到的包
var pktChan chan *iptun.IPPacket = make(chan *iptun.IPPacket, 10)
var unSatisfiedInterestCount int64 = 0

//export GoOnData
func GoOnData(cstr *C.char, size C.int) {
	//fmt.Println("GoOnData")
	// 收到 CPacket
	data := C.GoBytes(unsafe.Pointer(cstr), size)
	adapter.OnReceivePktFromNDN(&iptun.IPPacket{
		Src:        waterutil.IPv4Source(data),
		Dst:        waterutil.IPv4Destination(data),
		RawPackets: data,
	})
	sendInterest(ipTunnelConfig.TargetIdentifier)
	//atomic.AddInt64(&unSatisfiedInterestCount, -1)
}

//export GoOnInterest
func GoOnInterest(cstr *C.char) {
	name := C.GoString(cstr)
	//common.LogWarn("onInterest: ", name)
	// 接收到兴趣包，从缓存中取出一个IP包，封装到一个Data中发出
	if len(pktChan) > 0 {
		ipPacket := <-pktChan
		sendData(ipPacket.RawPackets, name)
	}
}

//export GoOnNack
func GoOnNack() {
	//atomic.AddInt64(&unSatisfiedInterestCount, -1)
	sendInterest(ipTunnelConfig.TargetIdentifier)
}

//export GoOnTimeout
func GoOnTimeout() {
	//atomic.AddInt64(&unSatisfiedInterestCount, -1)
	sendInterest(ipTunnelConfig.TargetIdentifier)
}

////int sendInterest(char *buf, int size, char *name) ;
//func sendPacket(pkt []byte, name string) error {
//	cstr := C.CString(name)
//	C.sendInterest((*C.char)(unsafe.Pointer(&pkt[0])), C.int(len(pkt)), cstr)
//	C.free(unsafe.Pointer(cstr))
//	return nil
//}

//int sendInterest(char *buf, int size, char *name) ;
func sendInterest(name string) error {
	cstr := C.CString(name)
	C.sendInterest(cstr)
	C.free(unsafe.Pointer(cstr))
	return nil
}

//int sendData(char *buf, int size, char *name) ;
func sendData(pkt []byte, name string) error {
	cstr := C.CString(name)
	if pkt == nil {
		C.sendData((*C.char)(unsafe.Pointer(nil)), C.int(0), cstr)
	} else {
		C.sendData((*C.char)(unsafe.Pointer(&pkt[0])), C.int(len(pkt)), cstr)
	}
	C.free(unsafe.Pointer(cstr))
	return nil
}

// StartIPTunnel 启动 IP Tunnel
//
// @Description:
// @param ipTunnelConfig
//
func StartIPTunnel(config *iptun.IPTunnelConfig) error {
	ipTunnelConfig = config
	if err := ipTun.Init(iptun.Config{
		InterfaceName: ipTunnelConfig.TunConfig.InterfaceName,
		IPv4Addr:      ipTunnelConfig.TunConfig.IPv4Addr,
		MTU:           ipTunnelConfig.TunConfig.Mtu,
		Mask:          ipTunnelConfig.TunConfig.Mask,
		Route:         ipTunnelConfig.TunConfig.Route,
	}); err != nil {
		return err
	}
	adapter.Init()

	// 在单独的协程里面处理从TUN读取IP包，构造一个Data并发出
	go func() {
		//count := 0
		for {
			// 从 TUN 中读取IP包，并通过CPacket发出
			ipPacket := adapter.GetPktFromTun()
			pktChan <- ipPacket
			//sendPacket(ipPacket.RawPackets, ipTunnelConfig.TargetIdentifier+"/"+strconv.Itoa(count))
			//count++
		}
	}()

	// 在单独的协程里面不停的发送 Interest，保持已发出未满足的兴趣包数量在30个左右
	go func() {
		for {
			if atomic.LoadInt64(&unSatisfiedInterestCount) < 30 {
				sendInterest(ipTunnelConfig.TargetIdentifier)
				atomic.AddInt64(&unSatisfiedInterestCount, 1)
			} else {
				break
				// 睡1ms
				time.Sleep(time.Millisecond)
			}
		}
	}()

	// 在单独的协程里面处理接收兴趣包
	go func() {
		cstr := C.CString(ipTunnelConfig.MirConfig.ListenIdentifier)
		ret := C.start(cstr, C.int(len(ipTunnelConfig.MirConfig.ListenIdentifier)))
		C.free(unsafe.Pointer(cstr))
		if ret < 0 {
			common.LogFatal("start face fail")
		}
	}()
	// 启动IP隧道程序
	ipTun.StartTunnel(&adapter)
	return nil
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
		return StartIPTunnel(ipTunnelConfig)
	}

	if err := mirApp.Run(os.Args); err != nil {
		return
	}
}
