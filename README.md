# IP tunnel 

This project implements NDN-based IP tunnel and MIN-based IP tunnel to compare the performance of the two.

## 1. Experiment intro

 ![多标识路由器设计](https://gitee.com/quejianming/pic-bed/raw/master/uPic/2021/10/05/%E5%A4%9A%E6%A0%87%E8%AF%86%E8%B7%AF%E7%94%B1%E5%99%A8%E8%AE%BE%E8%AE%A1-1633441220.svg) 

We implement an IP tunnel based on MIN and NDN, respectively, and compare the performance of the two network architectures in this scenario. The experimental topology is shown in the above figure, with R1, R2, and R3 running on three Linux servers whose configuration information is shown in the following table. Each server runs a MIN/NDN as a network router, and they are directly connected via an ethernet link. We implement a TUN tunnel on the edge router, which collects IP packets sent to the TUN interface and encapsulates them into MIN/NDN packets. Then, IP packets are transmitted to the edge router on the other side through MIN/NDN communication. Considering that MIN/NDN packets have some header overhead, to avoid fragmentation, we set the MTU of the TUN interface to 1500 and the MTU of the interface connected between routers to 2000. We will use *iperf3* to measure the network's performance using MIN/NDN as an IP tunnel.

| CPU                                                 | RAM   | Operating System | Software Version      |
| --------------------------------------------------- | ----- | ---------------- | --------------------- |
| Intel\(R\) Xeon\(R\) Gold 5118 CPU @ 2\.30GHz \* 96 | 128GB | Ubuntu 20.04.3   | MIR v0.1.0;NFD v0.7.1 |

## 2. Install

### 2.1 Running MIN IP Tunnel

Since MIN is not open source yet, only the experimental code is given here, and no operation guidance is provided.

### 2.2 Runnin NDN IP Tunnel

We should config three Linux servers, each server running a Ubuntu 20.04 operating system.

#### 2.2.1 Install NFD environment

- Linux Server 1:

  ```bash
  git clone https://gitea.qjm253.cn/SunnyQjm/NDNInstaller.git
  cd NDNInstaller
  ./install_ndn-cxx.sh
  ./install_nfd.sh
  ./install_ndn-tools.sh
  ```

- Linux Server 2:

  ```bash
  git clone https://gitea.qjm253.cn/SunnyQjm/NDNInstaller.git
  cd NDNInstaller
  ./install_ndn-cxx.sh
  ./install_nfd.sh
  ./install_ndn-tools.sh
  ```

- Linux Server 3:

  ```bash
  git clone https://gitea.qjm253.cn/SunnyQjm/NDNInstaller.git
  cd NDNInstaller
  ./install_ndn-cxx.sh
  ./install_nfd.sh
  ./install_ndn-tools.sh
  ```

#### 2.2.2 Config and running NFD

- First running NFD in all Linux Server

  ```bash
  # Linux Server 1, execute follow command
  nfd-start
  
  # Linux Server 2, execute follow command
  nfd-start
  
  # Linux Server 3, execute follow command
  nfd-start
  ```

- Then create an ethernet face 

  - Linux Server 1:

    ```bash
    nfdc face create remote ether://[76:8f:3c:bb:65:31] local dev://eno1
    # after execute above command, you will get a face id 
    nfdc route add prefix /pkusz/C nexthop <face-id>
    ```

  - Linux Server 2:

    ```bash
    nfdc face create remote ether://[76:8f:3c:bb:64:31] local dev://eno1
    # after execute above command, you will get a face id 
    nfdc route add prefix /pkusz/B nexthop <face-id>
    
    nfdc face create remote ether://[76:8f:3c:bb:66:31] local dev://eno2
    # after execute above command, you will get a face id 
    nfdc route add prefix /pkusz/C nexthop <face-id>
    ```

  - Linux Server 3:

    ```bash
    nfdc face create remote ether://[76:8f:3c:bb:65:32] local dev://eno1
    # after execute above command, you will get a face id 
    nfdc route add prefix /pkusz/B nexthop <face-id>
    ```

#### 2.2.3 Config and running NDN-based IP Tunnel

- modify config file

  - Linux Server 1:

    Edit iptunconf.ini，replace the config file value as follow:

    ```ini
    [Tun]
    InterfaceName = tun0
    IPv4Addr = 192.168.88.1
    Mtu = 1500
    Mask = 255.255.255.0
    Route = 192.168.90.0/24
    
    [Mir]
    UnixSocketPath = /tmp/mir.sock
    ListenIdentifier = /pkusz/B
    TargetIdentifier = /pkusz/C
    
    [Log]
    LogLevel = WARN
    ```

  - Linux Server 3:

    Edit iptunconf.ini，replace the config file value as follow:

    ```bash
    [Tun]
    InterfaceName = tun0
    IPv4Addr = 192.168.88.2
    Mtu = 1500
    Mask = 255.255.255.0
    Route = 192.168.89.0/24
    
    [Mir]
    UnixSocketPath = /tmp/mir.sock
    ListenIdentifier = /pkusz/C
    TargetIdentifier = /pkusz/B
    
    [Log]
    LogLevel = WARN
    ```

- running IP Tunnel

  - Linux Server 1:

    ```bash
    # cd to ip-tunnel/ndniptunnel, then execute follow command
    go run .
    ```

  - Linux Server 3:

    ```bash
    # cd to ip-tunnel/ndniptunnel, then execute follow command
    go run .
    ```

- Test whether the tunnel is established successfully

  - Linux Server 1:

    ```bash
    ping 192.168.90.1
    ```

#### 2.2.4 Running experiment test

- Linux Server 3:

  ```bash
  iperf3 -s -B 192.168.90.1
  ```

- Linux Server 1:

  - Test TCP:

    ```bash
    # 1 Client
    iperf3 -c 192.168.90.1 -B 192.168.89.1 -t 60
    # 2 Client
    iperf3 -c 192.168.90.1 -B 192.168.89.1 -t 60 -P 2
    # 5 Client
    iperf3 -c 192.168.90.1 -B 192.168.89.1 -t 60 -P 5
    # 10 Client 
    iperf3 -c 192.168.90.1 -B 192.168.89.1 -t 60 -P 10
    ```

  - Test UDP:

    ```bash
    # 100Mbps
    iperf3 -c 192.168.90.1 -B 192.168.89.1 -u -l 1400 -t 60 -b 100m -R
    # 200Mbps
    iperf3 -c 192.168.90.1 -B 192.168.89.1 -u -l 1400 -t 60 -b 200m -R
    # 300Mbps
    iperf3 -c 192.168.90.1 -B 192.168.89.1 -u -l 1400 -t 60 -b 300m -R
    # 400Mbps
    iperf3 -c 192.168.90.1 -B 192.168.89.1 -u -l 1400 -t 60 -b 400m -R
    # 500Mbps
    iperf3 -c 192.168.90.1 -B 192.168.89.1 -u -l 1400 -t 60 -b 500m -R
    ```

## 3. Result

The example of test data you can get from [min-graph/result_1500]([min-graph/result_1500 at main · SunnyQjm/min-graph (github.com)](https://github.com/SunnyQjm/min-graph/tree/main/result_1500))

The analysis code you can get from [min-graph]([SunnyQjm/min-graph (github.com)](https://github.com/SunnyQjm/min-graph))