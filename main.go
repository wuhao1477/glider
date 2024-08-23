package main

import (
	"context"
	"net"       // 网络连接相关的包
	"os"        // 操作系统相关的包
	"os/signal" // 信号处理包
	"syscall"   // 系统调用接口包
	"time"      // 时间处理包

	"github.com/nadoo/glider/dns"     // DNS 服务器相关的包
	"github.com/nadoo/glider/ipset"   // IP 集管理相关的包
	"github.com/nadoo/glider/pkg/log" // 日志记录相关的包
	"github.com/nadoo/glider/proxy"   // 代理相关的包
	"github.com/nadoo/glider/rule"    // 规则管理相关的包
	"github.com/nadoo/glider/service" // 服务管理相关的包
)

var (
	version = "0.17.0"      // 版本号
	config  = parseConfig() // 解析配置文件
)

func main() {
	// 全局规则代理
	pxy := rule.NewProxy(config.Forwards, &config.Strategy, config.rules)

	// ipset 管理器
	ipsetM, _ := ipset.NewManager(config.rules)

	// 检查并设置 DNS 服务器
	if config.DNS != "" {
		d, err := dns.NewServer(config.DNS, pxy, &config.DNSConfig)
		if err != nil {
			// 如果 DNS 服务器创建失败，记录错误并退出程序
			log.Fatal(err)
		}

		// 设置规则中的 DNS 服务器
		for _, r := range config.rules {
			if len(r.DNSServers) > 0 {
				for _, domain := range r.Domain {
					// 为每个域名设置 DNS 服务器
					d.SetServers(domain, r.DNSServers)
				}
			}
		}

		// 添加处理程序，当域名解析时更新代理规则
		d.AddHandler(pxy.AddDomainIP)
		if ipsetM != nil {
			d.AddHandler(ipsetM.AddDomainIP)
		}

		// 启动 DNS 服务器
		d.Start()

		// 自定义解析器
		net.DefaultResolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: time.Second * 3}
				return d.DialContext(ctx, "udp", config.DNS)
			},
		}
	}

	// 清空规则中的 IP, CIDR 和域名信息
	for _, r := range config.rules {
		r.IP, r.CIDR, r.Domain = nil, nil, nil
	}

	// 启用检查器
	pxy.Check()

	// 启动代理服务器
	for _, listen := range config.Listens {
		local, err := proxy.ServerFromURL(listen, pxy)
		if err != nil {
			// 如果代理服务器创建失败，记录错误并退出程序
			log.Fatal(err)
		}
		// 以协程方式运行代理服务器
		go local.ListenAndServe()
	}

	// 启动服务
	for _, s := range config.Services {
		service, err := service.New(s)
		if err != nil {
			// 如果服务创建失败，记录错误并退出程序
			log.Fatal(err)
		}
		// 以协程方式运行服务
		go service.Run()
	}

	// 处理系统信号，用于优雅地退出程序
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// 等待系统信号
	<-sigCh
}
