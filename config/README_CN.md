## 配置文件
命令:
```bash
glider -config glider.conf
```
配置文件，**只需使用命令行标志名作为键名**：
```bash
  # 注释行
  键=值
  键=值
  # 键等同于命令行标志名: listen, forward, strategy...
```

示例:
```bash
### glider 配置文件

# 详细模式，打印日志
verbose

# 监听8443端口，同时作为http/socks5代理服务器。
listen=:8443

# 上游转发代理
forward=socks5://192.168.1.10:1080

# 上游转发代理
forward=ss://method:pass@1.1.1.1:8443

# 上游转发代理（转发链）
forward=http://1.1.1.1:8080,socks5://2.2.2.2:1080

# 多个上游代理的转发策略
strategy=rr

# 转发器健康检查
check=http://www.msftconnecttest.com/connecttest.txt#expect=200

# 检查间隔
checkinterval=30

# 设置一个DNS转发服务器
dns=:53
# 全局远程DNS服务器（你可以在规则文件中指定不同的DNS服务器）
dnsserver=8.8.8.8:53

# 规则文件
rules-dir=rules.d
#rulefile=office.rule
#rulefile=home.rule

# 包含更多配置文件
#include=dnsrecord.inc.conf
#include=more.inc.conf
```

参见:
- [glider.conf.example](glider.conf.example)
- [examples](examples)

## 规则文件
规则文件，**与配置文件类似，但基于目标指定转发器**：
```bash
# 你可以使用全局配置文件中的所有键，除了"listen"和"rulefile"
forward=socks5://192.168.1.10:1080
forward=ss://method:pass@1.1.1.1:8443
forward=http://192.168.2.1:8080,socks5://192.168.2.2:1080
strategy=rr
check=http://www.msftconnecttest.com/connecttest.txt#expect=200
checkinterval=30

# 此规则文件中域名的DNS服务器
dnsserver=208.67.222.222:53

# IPSET管理
# ---------
# 基于规则文件中的目标在Linux上创建和管理ipset
#   - 启动时添加规则文件中的IP/CIDR
#   - 通过DNS转发服务器为规则文件中的域名添加解析后的IP
# 通常在Linux的透明代理模式下使用
ipset=glider

# 你可以指定目标以使用上述转发器
# 匹配abc.com及*.abc.com
domain=abc.com

# 匹配1.1.1.1
ip=1.1.1.1

# 匹配192.168.100.0/24
cidr=192.168.100.0/24

# 我们可以包含一个仅包含目标设置的列表文件
include=office.list.example
```

参见:
- [office.rule.example](rules.d/office.rule.example)
- [examples](examples)