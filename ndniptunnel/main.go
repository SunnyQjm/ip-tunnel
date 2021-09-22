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
	"strconv"
	"unsafe"
)

var ipTun iptun.IPTun
var adapter ndn.NDNTunnelAdapter

//export GoOnData
func GoOnData(cstr *C.char, size C.int) {
	//fmt.Println("GoOnData")
}

//export GoOnInterest
func GoOnInterest(cstr *C.char, size C.int) {
	// 接收到兴趣包，提取出其中携带的 IP 并传递给 adapter
	data := C.GoBytes(unsafe.Pointer(cstr), size)
	adapter.OnReceivePktFromNDN(&iptun.IPPacket{
		Src:        waterutil.IPv4Source(data),
		Dst:        waterutil.IPv4Destination(data),
		RawPackets: data,
	})
}

//int sendInterest(char *buf, int size, char *name) ;
func sendPacket(pkt []byte, name string) error {
	cstr := C.CString(name)
	C.sendInterest((*C.char)(unsafe.Pointer(&pkt[0])), C.int(len(pkt)), cstr)
	C.free(unsafe.Pointer(cstr))
	return nil
}

// StartIPTunnel 启动 IP Tunnel
//
// @Description:
// @param ipTunnelConfig
//
func StartIPTunnel(ipTunnelConfig *iptun.IPTunnelConfig) error {
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

	// 在单独的协程里面处理从TUN读取IP包，构造一个Interest并发出
	go func() {
		count := 0
		for {
			// 从 TUN 中读取IP包，并通过CPacket发出
			ipPacket := adapter.GetPktFromTun()
			sendPacket(ipPacket.RawPackets, ipTunnelConfig.TargetIdentifier+"/"+strconv.Itoa(count))
			count++
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
