# [glider](https://github.com/nadoo/glider)

[![Go Version](https://img.shields.io/github/go-mod/go-version/nadoo/glider?style=flat-square)](https://go.dev/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/nadoo/glider?style=flat-square)](https://goreportcard.com/report/github.com/nadoo/glider)
[![GitHub release](https://img.shields.io/github/v/release/nadoo/glider.svg?style=flat-square&include_prereleases)](https://github.com/nadoo/glider/releases)
[![Actions Status](https://img.shields.io/github/actions/workflow/status/nadoo/glider/build.yml?branch=dev&style=flat-square)](https://github.com/nadoo/glider/actions)
[![DockerHub](https://img.shields.io/docker/image-size/nadoo/glider?color=blue&label=docker&style=flat-square)](https://hub.docker.com/r/nadoo/glider)

glider 是一个支持多种协议的正向代理，同时也是一个具备 ipset 管理功能的 DNS/DHCP 服务器（类似 dnsmasq）。

我们可以设置本地监听器作为代理服务器，并通过转发器将请求转发到互联网。

```bash
                |转发器 -------------------->|
   监听器 -----> |                            | 互联网
                |转发器 --> 转发器 -> ...     |
```

## 特性
- 作为代理客户端和代理服务器（协议转换器）
- 灵活的代理和协议链
- 负载均衡，支持以下调度算法：
  - rr: 轮询
  - ha: 高可用
  - lha: 基于延迟的高可用
  - dh: 目标哈希
- 基于规则和优先级的转发器选择：[配置示例](config/examples)
- DNS 转发服务器：
  - 通过代理的 DNS
  - 强制通过 TCP 查询上游
  - DNS 与转发器选择之间的关联规则
  - DNS 与 ipset 之间的关联规则
  - 支持 DNS 缓存
  - 自定义 DNS 记录
- IPSet 管理（Linux 内核版本 >= 2.6.32）：
  - 启动时从规则文件中添加 IP/CIDR
  - 通过 DNS 转发服务器从规则文件中解析域名的 IP
- 在同一个端口上提供 HTTP 和 SOCKS5 服务
- 定期检查转发器的可用性
- 从指定的本地 IP/接口发送请求
- 服务：
  - dhcpd: 可以运行在故障转移模式下的简单 DHCP 服务器

## 支持的协议

<details>
<summary>点击查看详情</summary>

|协议          | 监听/TCP | 监听/UDP | 转发/TCP | 转发/UDP | 描述
|:-:           |:-:|:-:|:-:|:-:|:-
|Mixed         |√|√| | |http+socks5 服务器
|HTTP          |√| |√| |客户端和服务器
|SOCKS5        |√|√|√|√|客户端和服务器
|SS            |√|√|√|√|客户端和服务器
|Trojan        |√|√|√|√|客户端和服务器
|Trojanc       |√|√|√|√|Trojan 明文（无 TLS）
|VLESS         |√|√|√|√|客户端和服务器
|VMess         | | |√|√|仅客户端
|SSR           | | |√| |仅客户端
|SSH           | | |√| |仅客户端
|SOCKS4        | | |√| |仅客户端
|SOCKS4A       | | |√| |仅客户端
|TCP           |√| |√| |TCP 隧道客户端和服务器
|UDP           | |√| |√|UDP 隧道客户端和服务器
|TLS           |√| |√| |传输客户端和服务器
|KCP           | |√|√| |传输客户端和服务器
|Unix          |√|√|√|√|传输客户端和服务器
|VSOCK         |√| |√| |传输客户端和服务器
|Smux          |√| |√| |传输客户端和服务器
|Websocket(WS) |√| |√| |传输客户端和服务器
|WS Secure     |√| |√| |WebSocket 安全传输 (WSS)
|Proxy Protocol|√| | | |仅版本 1 服务器
|Simple-Obfs   | | |√| |仅传输客户端
|Redir         |√| | | |Linux 重定向代理
|Redir6        |√| | | |Linux 重定向代理 (IPv6)
|TProxy        | |√| | |Linux TProxy（仅 UDP）
|Reject        | | |√|√|拒绝所有请求

</details>

## 安装

- 二进制文件: [https://github.com/nadoo/glider/releases](https://github.com/nadoo/glider/releases)
- Docker: `docker pull nadoo/glider`
- Manjaro: `pamac install glider`
- ArchLinux: `sudo pacman -S glider`
- Homebrew: `brew install glider`
- MacPorts: `sudo port install glider`
- 源代码: `go install github.com/nadoo/glider@latest`

## 使用方法

#### 运行

```bash
glider -verbose -listen :8443
# docker run --rm -it nadoo/glider -verbose -listen :8443
```

#### 帮助

<details>
<summary><code>glider -help</code></summary>

```bash
用法: glider [-listen URL]... [-forward URL]... [选项]...

  例如: glider -config /etc/glider/glider.conf
        glider -listen :8443 -forward socks5://serverA:1080 -forward socks5://serverB:1080 -verbose

选项:
  -check string
        check=tcp[://HOST:PORT]: TCP 端口连接检查
        check=http://HOST[:PORT][/URI][#expect=REGEX_MATCH_IN_RESP_LINE]
        check=https://HOST[:PORT][/URI][#expect=REGEX_MATCH_IN_RESP_LINE]
        check=file://SCRIPT_PATH: 运行检查脚本，退出码为 0 时表示健康，环境变量：FORWARDER_ADDR,FORWARDER_URL
        check=disable: 禁用健康检查（默认值 "http://www.msftconnecttest.com/connecttest.txt#expect=200"）
  -checkdisabledonly
        仅检查禁用的转发器
  -checkinterval int
        转发器检查间隔（秒）（默认值 30）
  -checklatencysamples int
        使用最近 N 次检查的平均延迟（默认值 10）
  -checktimeout int
        转发器检查超时（秒）（默认值 10）
  -checktolerance int
        转发器检查容忍度（毫秒），仅在 lha 模式下使用
  -config string
        配置文件路径
  -dialtimeout int
        拨号超时（秒）（默认值 3）
  -dns string
        本地 DNS 服务器监听地址
  -dnsalwaystcp
        无论是否有转发器，始终使用 TCP 查询上游 DNS 服务器
  -dnscachelog
        显示 DNS 缓存查询日志
  -dnscachesize int
        缓存中的 DNS 响应最大数量（默认值 4096）
  -dnsmaxttl int
        缓存中条目的最大 TTL 值（秒）（默认值 1800）
  -dnsminttl int
        缓存中条目的最小 TTL 值（秒）
  -dnsnoaaaa
        禁用 AAAA 查询
  -dnsrecord value
        自定义 DNS 记录，格式：domain/ip
  -dnsserver value
        远程 DNS 服务器地址
  -dnstimeout int
        多个 DNS 服务器切换时的超时值（秒）（默认值 3）
  -example
        显示使用示例
  -forward value
        转发 URL，详见下文 URL 部分
  -include value
        包含文件
  -interface string
        源 IP 或源接口
  -listen value
        监听 URL，详见下文 URL 部分
  -logflags int
        如果不了解它，请不要更改，参考：https://pkg.go.dev/log#pkg-constants （默认值 19）
  -maxfailures int
        更改转发器状态为禁用所需的最大失败次数（默认值 3）
  -relaytimeout int
        中继超时（秒）
  -rulefile value
        规则文件路径
  -rules-dir string
        规则文件夹路径
  -scheme string
        显示代理方案的帮助信息，使用 'all' 查看所有方案
  -service value
        运行指定

的服务 (dhcpd)
  -tun string
        Linux tun 设备名称
  -verbose
        详细日志
  -vscok
        使用 Virtio VSOCK 协议

有关详细信息，请访问：https://github.com/nadoo/glider

```

</details>

#### URL 格式

##### `scheme://[user[:pass]@]host[:port]`

|协议|描述|
|---|---|
|http|HTTP 代理服务器（仅限监听）|
|https|HTTP over TLS（监听）|
|ws|WebSocket（监听）|
|wss|WebSocket Secure（监听）|
|socks5|SOCKS5 代理服务器|
|socks5h|SOCKS5 代理服务器（通过 SOCKS5 的 DNS 解析）|
|ss|Shadowsocks（监听和转发）|
|trojan|Trojan（监听和转发）|
|trojanc|Trojan 明文（无 TLS）|
|vless|VLESS（监听和转发）|
|vmess|VMess（仅限转发）|
|redir|Linux 重定向代理|
|tun|Linux tun 设备|
|unix|Unix 套接字（监听和转发）|
|vsock|Virtio VSOCK 协议（监听和转发）|

## 示例配置

#### HTTP 代理监听

```bash
glider -listen http://0.0.0.0:8443
```

#### HTTP 代理转发

```bash
glider -forward http://serverA:1080 -forward http://serverB:1080
```

#### HTTP over TLS 代理转发

```bash
glider -forward https://server:8443
```

#### SOCKS5 代理监听和转发

```bash
glider -listen socks5://0.0.0.0:1080 -forward socks5://serverA:1080 -forward socks5://serverB:1080
```

#### Shadowsocks (ss) 代理监听

```bash
glider -listen ss://AEAD_CHACHA20_POLY1305:password@:8443
```

#### Trojan (trojan) 代理监听

```bash
glider -listen trojan://password@:8443
```

#### VLESS 代理监听

```bash
glider -listen vless://server:443?encryption=none
```

#### WebSocket Secure (wss) 代理转发

```bash
glider -listen socks5://0.0.0.0:1080 -forward wss://server:443
```

## Docker

#### 创建 glider 容器

```bash
docker run --name=glider -d \
    -p 1080:1080 \
    --restart unless-stopped \
    -v /etc/glider:/etc/glider \
    nadoo/glider \
    -config /etc/glider/glider.conf
```

## 许可

本项目采用 MIT 许可。请查阅 [LICENSE](LICENSE) 文件获取更多信息。

## 贡献

欢迎提交问题和 PR。如果你有任何建议或意见，欢迎加入讨论。

