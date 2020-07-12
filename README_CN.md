# inlets

<!--翻译贡献者请注意译文三大原则：信、达、雅。格式和 Markdown 风格请与英文原文保持一致。-->

将你的本地服务暴露到公网。

[![Build Status](https://travis-ci.org/inlets/inlets.svg?branch=master)](https://travis-ci.org/inlets/inlets) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/inlets/inlets)](https://goreportcard.com/report/github.com/inlets/inlets) [![Documentation](https://godoc.org/github.com/inlets/inlets?status.svg)](http://godoc.org/github.com/inlets/inlets) [![Derek App](https://alexellis.o6s.io/badge?owner=inlets&repo=inlets)](https://github.com/alexellis/derek/)
[![Setup Automated](https://img.shields.io/badge/setup-automated-blue?logo=gitpod)](https://gitpod.io/from-referrer/)

[English](./README.md) | [中文文档](./README_CN.md)

## 简介

inlets 利用反向代理和 Websocket 隧道，将内部、或是开发中的服务通过「出口节点」暴露到公网。出口节点可以是几块钱一个月的 VPS，也可以是任何带有公网 IPv4 的电脑。

为什么需要这个项目？类似的工具例如 [ngrok](https://ngrok.com/) 和由 [Cloudflare](https://www.cloudflare.com/) 开发的 [Argo Tunnel](https://developers.cloudflare.com/argo-tunnel/) 皆为闭源，内置了一些限制，并且价格不菲，以及对 arm/arm64 的支持很有限。Ngrok 还经常会被公司防火墙策略拦截而导致无法使用。而其它开源的隧道工具，基本只考虑到静态地配置单个隧道。inlets 旨在动态地发现本地服务，通过 Websocket 隧道将它们暴露到公网 IP 或域名，并自动化配置 TLS 证书。

当开启 SSL 时，inlets 可以通过任何支持 `CONNECT` 方法的 HTTP 代理服务器。

![](docs/inlets.png)

*inlets 概念示意图*

## 协议与条款

**重要**

如您需要在企业网络中使用 inlets，建议先征求 IT 管理员的同意。下载、使用或分发 inlets 前，您必须同意 [协议](./LICENSE) 条款与限制。本项目不提供任何担保，亦不承担任何责任。

### 幕后的开发者是谁？

inlets 由 [Alex Ellis](https://twitter.com/alexellisuk) 开发。Alex 是一名 [CNCF 大使](https://www.cncf.io/people/ambassadors/)，同时是 [OpenFaaS](https://github.com/openfaas/faas/) 的创始人。

> [OpenFaaS&reg;](https://github.com/openfaas/faas) 使得开发者将由事件驱动的函数和微服务部署到 Kubernetes 更加容易，而无需编写重复的样板代码。把代码或现成的二进制文件打包进 Docker 镜像，即可获得带有自动扩容和监控指标的服务入口。该项目目前已有接近 19k GitHub stars，超过 240 名贡献者；越来越多的用户已将它应用到生产环境。

### 待办事项与目标

### 已完成

* 基于客户端的定义，自动在出口节点创建服务入口
  * 通过 DNS / 域名实现单端口、单 Websocket 承载多站点
* 利用 SSL over Websockets 实现链路加密（`wss://`）
* 自动重连
* 通过 Service Account 或 HTTP Basic Auth 实现权限认证
  * 通过 HTTP01 challenge 使用 LetsEncrypt Staging 或 Production 签发证书
* 原生跨平台支持，包括 ARMHF 和 ARM64 架构
* 提供 Dockerfile 以及 Kubernetes YAML 文件
* 自动发现并实例化 Kubernetes 集群内 `LoadBalancer` 类型的 `Service` - [inlets-operator](https://github.com/inlets/inlets-operator)
* 除 HTTP(s) 以外，还支持在隧道内传输 Websocket 流量
* [为该项目制作一枚 Logo](https://github.com/inlets/inlets/issues/46)

#### 延伸目标

* 自动配置 DNS / A 记录。
* 基于 Azure ACI 和 AWS Fargate，以 Serverless 容器的方式运行「出口节点」。
* 通过 DNS01 challenge 使用 LetsEncrypt Staging 或 Production 签发证书

#### 非本项目的目标

* 通过 Websocket 隧道传输原始 TCP 流量。

  inlets-pro 涵盖了该使用场景，您可以向我咨询 [inlets-pro](mailto:alex@openfaas.com) 的内测事宜（English Only）。

### 项目状态

与 HTTP 1.1 遵循同步的请求/响应模型不同，Websocket 使用异步的发布/订阅模型来发送和接收消息。这带来了一些挑战 —— 通过 *异步总线* 隧道化传输 *同步协议*。

inlets 2.0 带来了性能上的提升，以及调用部分 Kubernetes 和 Rancher API 的能力。本项目使用了 [Rancher 的 K3s 项目](https://k3s.io) 实现节点间通讯同样的隧道依赖包。它非常适用于开发，在生产环境中也很实用。不过在部署 `inlets` 到生产环境中之前，建议先做好充足的测试。

如果您有任何评论、建议或是贡献想法，欢迎提交 Issue 讨论。

* 隧道链路通过 `--token` 选项指定的共享密钥保证安全
* 默认配置使用不带 SSL 的 Websocket `ws://`，但支持开启加密，即启用 SSL `wss://`
* 可通过服务器端选项设定请求超时时间
* ~~服务发现机制完成前，在服务端和客户端都必须配置上游 URL~~ 客户端可发布其可提供服务的上游 URLs
* 默认情况下，隧道传输会移除响应内的 CORS 头，但你可以在服务端使用 `--disable-transport-wrapping` 关闭该特性

### 相关项目

Inlets 作为服务代理 [已被列入 Cloud Native Landscape](https://landscape.cncf.io/category=service-proxy&format=card-mode&grouping=category&sort=stars)

* [inlets](https://github.com/inlets/inlets) - 开源的七层 HTTP 隧道和反向代理
* [inlets-pro](https://github.com/inlets/inlets-pro-pkg) - 四层 TCP 负载均衡
* [inlets-operator](https://github.com/inlets/inlets-operator) - 深度集成 Inlets 和 Kubernetes，实现 LoadBalancer 类型的 Service
* [inletsctl](https://github.com/inlets/inletsctl) - 配置出口节点的 CLI 工具，配合 inlets 和 inlets-pro 使用

### 大家如何评论 inlets？

> 你可以在社交媒体使用 `@inletsdev`、`#inletsdev` 和 `https://inlets.dev` 分享 inlets 的相关内容。

inlets 曾两次登上 Hacker News 首页推荐：

* [inlets 1.0](https://news.ycombinator.com/item?id=19189455) - 146 points, 48 评论
* [inlets 2.0](https://news.ycombinator.com/item?id=20410552) - 218 points, 66 评论

相关教程（英文）：

* [Get a LoadBalancer for your private Kubernetes cluster with inlets-operator by Alex Ellis](https://blog.alexellis.io/ingress-for-your-local-kubernetes-cluster/)
* [Blog post - webhooks, great when you can get them by Alex Ellis](https://blog.alexellis.io/webhooks-are-great-when-you-can-get-them/)
* [Micro-tutorial inlets with KinD by Alex Ellis](https://gist.github.com/alexellis/c29dd9f1e1326618f723970185195963)
* [The Awesomeness of Inlets by Ruan Bekker](https://sysadmins.co.za/the-awesomeness-of-inlets/)
* [K8Spin - What does fit in a low resources namespace? Inlets](https://medium.com/k8spin/what-does-fit-in-a-low-resources-namespace-3rd-part-inlets-6cc278835e57)
* [Exposing Magnificent Image Classifier with inlets](https://blog.baeke.info/2019/07/17/exposing-a-local-endpoint-with-inlets/)
* ["Securely access external applications as Kubernetes Services, from your laptop or from any other host, using inlets"](https://twitter.com/BanzaiCloud/status/1164168218954670080)
* [Using local services in Gitpod with inlets](https://www.gitpod.io/blog/local-services-in-gitpod/)

推文：

* ["I just transferred a 70Gb disk image from a NATed NAS to a remote NATed server with @alexellisuk inlets tunnels and a one-liner python web server" by Roman Dodin](https://twitter.com/ntdvps/status/1143071544203186176)
* ["Really amazed by inlets by @alexellisuk - "Up and running in 15min - I will be able to watch my #RaspberryPi servers running at home while staying on the beach 🏄‍♂️🌴🍸👏👏👏" by Florian Dambrine](https://twitter.com/DambrineF/status/1158364581624012802?s=20)
* [Testing an OAuth proxy by Vivek Singh](https://twitter.com/viveksyngh/status/1142054203478564864)
* [inlets used at KubeCon to power a live IoT demo at a booth](https://twitter.com/tobruzh/status/1130421702914129921)
* [PR to support Risc-V by Carlos Eduardo](https://twitter.com/carlosedp/status/1140740494617645061)
* [Recommended by Michael Hausenblas for use with local Kubernetes](https://twitter.com/mhausenblas/status/1143020953380753409)
* [5 top facts about inlets by Alex Ellis](https://twitter.com/alexellisuk/status/1140552115204608001)
* ["Cool! I hadn't heard of inlets until now, but I love the idea of exposing internal services this way. I've been using TOR to do this!" by Stephen Doskett, Tech Field Day](https://twitter.com/SFoskett/status/1108989190912524288)
* ["Learn how to set up HTTPS for your local endpoints with inlets, Caddy, and DigitalOcean thanks to @alexellisuk!" by @DigitalOcean](https://twitter.com/digitalocean/status/1113440166310502400)
* ["See how Inlets helped me to expose my local endpoints for my homelab that sits behind a Carrier-Grade NAT"](https://twitter.com/ruanbekker/status/1161399537417801728)

> 提示：欢迎提交 PR 添加你的故事或是使用场景，很期待听到你的声音！

阅读 [ADOPTERS.md](./ADOPTERS.md) 查看哪些公司正在使用 inlets。

## 开始使用

你可以使用 `curl` 下载安装脚本，或是用 `brew` 安装，或者直接在 Releases 页面直接下载二进制文件。安装完成后即可使用 `inlets` 命令。

### 安装 CLI

> 提示：虽然 `inlets` 是一款免费工具，但你也可以在 [GitHub Sponsors](https://insiders.openfaas.io/) 页面支持后续的开发 💪

使用 `curl` 和辅助脚本：

```bash
# 安装到当前目录
curl -sLS https://get.inlets.dev | sh

# 安装到 /usr/local/bin/
curl -sLS https://get.inlets.dev | sudo sh
```

使用 `brew`：

```bash
brew install inlets
```

> 提示：`brew` 分发的版本由 Homebrew 团队维护，因此可能会与 GitHub releases 存在一定延迟。

二进制文件可在 [Releases 页面](https://github.com/inlets/inlets/releases) 找到；包含 Linux（x86_64、armhf、arm64），Windows（实验性）以及 Darwin（MacOS）版本。如果你想要验证你的下载，也可以查看 SHA 校验值。

鼓励 Windows 用户使用 [Git bash](https://git-scm.com/downloads) 来安装 inlets。

### 入门教程

你可以在任何两台互相连接的「电脑」之间运行 inlets，「电脑」可以是两个容器，虚拟机，物理机，甚至你笔记本的环回网络也可以。

推荐阅读 [how to provision an "exit-node" with a public IPv4 address using a VPS](./docs/vps.md)。

* 以下步骤在 *出口节点*（又称服务端）执行。

首先在任何有公网 IP 的机器上（例如 VPS）启动隧道服务端。

例子如下，生成客户端认证的 Token 并启动服务端：

```bash
export token=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)
inlets server --port=8090 --token="$token"
```

> 提示：同时在服务端和客户端配置 `--token` 选项和密钥，可避免未授权地连接到隧道。


```bash
inlets server --port=8090
```

也可以像上面这样完全无保护地运行，但是并不推荐。

随后记下你的公网 IP。

* 接下来到运行 HTTP 服务的机器。

你可以使用我开发的 hash-browns 服务作为测试，该服务可生成哈希值。

```sh
export GO111MODULE=off
export GOPATH=$HOME/go/

go get -u github.com/alexellis/hash-browns
cd $GOPATH/src/github.com/alexellis/hash-browns

port=3000 go run server.go
```

如果你没安装 Golang，也可以运行 [Python 内置的 HTTP 服务](https://docs.python.org/2/library/simplehttpserver.html)：

```sh
mkdir -p /tmp/inlets-test/
cd /tmp/inlets-test/
touch hello-world
python -m SimpleHTTPServer 3000
```

* 在同一台机器上，启动 inlets 客户端。

启动隧道客户端：

```sh
export REMOTE="127.0.0.1:8090"    # 替换成刚刚记下的公网 IP
export TOKEN="CLIENT-TOKEN-HERE"  # Token 的值可在刚刚启动 "inlets server" 时找到
inlets client \
 --remote=$REMOTE \
 --upstream=http://127.0.0.1:3000 \
 --token $TOKEN
```

* 务必替换 `--remote` 的值为运行 `inlets server` （即出口节点）的 IP。
* 务必将 `--token` 的值与服务端保持一致。

我们现在总开启了三个进程：
* 用于测试的 HTTP 服务（运行 hash-browns 或是 Python Web 服务器）
* 出口节点运行着的隧道服务（`inlets server`）
* 隧道客户端（`inlets client`）

接下来是时候给 inlets 服务端发请求了，用指向它的域名或 IP 均可：

假设你的服务端位于 `127.0.0.1`，使用 `/etc/hosts` 文件或是 DNS 服务将域名 `gateway.mydomain.tk` 指向 `127.0.0.1`。

```sh
curl -d "hash this" http://127.0.0.1:8090/hash -H "Host: gateway.mydomain.tk"
# 或
curl -d "hash this" http://127.0.0.1:8090/hash
# 或
curl -d "hash this" http://gateway.mydomain.tk/hash
```

你会看到有流量通过隧道客户端到出口节点，如果你运行的是 hash-browns 服务，会出现类似下面的日志：

```sh
~/go/src/github.com/alexellis/hash-browns$ port=3000 go run main.go
2018/12/23 20:15:00 Listening on port: 3000
"hash this"
```

顺便还可以看看 hash-browns 服务内置的 Metrics 数据：

```sh
curl $REMOTE/metrics | grep hash
```

此外你还可以使用多个域名，并将它们分别绑定到多个内网服务。

这里我们在两个端口上启动 Python Web 服务，分别将两个本地目录作为服务内容，并将它们映射到不同的 Host 头，也就是域名：

```sh
mkdir -p /tmp/store1
cd /tmp/store1/
touch hello-store-1
python -m SimpleHTTPServer 8001 &


mkdir -p /tmp/store2
cd /tmp/store2/
touch hello-store-2
python -m SimpleHTTPServer 8002 &
```

```sh
export REMOTE="127.0.0.1:8090"    # 替换成刚刚记下的公网 IP
export TOKEN="CLIENT-TOKEN-HERE"  # Token 的值可在刚刚启动 "inlets server" 时找到
inlets client \
 --remote=$REMOTE \
 --token $TOKEN \
 --upstream="store1.example.com=http://127.0.0.1:8001,store2.example.com=http://127.0.0.1:8002"
```

随后修改 `store1.example.com` 和 `store2.example.com` 的 DNS 指向或设置 `/etc/hosts` 文件，即可通过浏览器访问了。

## 继续深入

### 文档与特色教程

教程：[HTTPS for your local endpoints with inlets and Caddy](https://blog.alexellis.io/https-inlets-local-endpoints/)

文档：[Inlets & Kubernetes recipes](./docs/kubernetes.md)

文档：[Run Inlets on a VPS](./docs/vps.md)

教程：[Get a LoadBalancer for your private Kubernetes cluster with inlets-operator](https://blog.alexellis.io/ingress-for-your-local-kubernetes-cluster/)

### 视频 Demo

使用 inlets 实现为我的 JavaScript & Webpack 应用配置公开的服务入口，以及自定义的域名：[Create React App](https://github.com/facebook/create-react-app)。

[![https://img.youtube.com/vi/jrAqqe8N3q4/hqdefault.jpg](https://img.youtube.com/vi/jrAqqe8N3q4/maxresdefault.jpg)](https://youtu.be/jrAqqe8N3q4)

### Docker

适用于多种架构的 Docker 镜像已发布，支持 `x86_64`, `arm64` and `armhf`。

* `inlets/inlets:2.6.3`

### 单出口节点多服务

你可以通过 inlets 暴露 OpenFaaS 或 OpenFaaS Cloud deployment，只需要将 `--upstream=http://127.0.0.1:3000` 改为 `--upstream=http://127.0.0.1:8080` 或是 `--upstream=http://127.0.0.1:31112` 即可。甚至可以指向任何内网或是外网 IP 地址，例如：`--upstream=http://192.168.0.101:8080`。

### 为控制平面设定独立端口

你可以为用户访问和隧道传输分别指定不同的端口。

* `--port` - 指定用户访问、提供对外服务的端口，又称 *数据平面*
* `--control-port` - 指定底层 Websocket 隧道连接的端口，又称 *控制平面*

### 开发指引

首先需要在出口节点和客户端都安装 Golang 1.10 或 1.11。

使用类似如下命令获取代码：

```bash
go get -u github.com/inlets/inlets
cd $GOPATH/src/github.com/inlets/inlets
```

另外，你也可以使用 [Gitpod](https://gitpod.io) 一键在浏览器中配置好开发环境：

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/inlets/inlets)

### 附录

其它 Kubernetes 端口转发工具：

* [`kubectl port-forward`](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) - built into the Kubernetes CLI, forwards a single port to the local computer.
* [kubefwd](https://github.com/txn2/kubefwd) - Kubernetes utility to port-forward multiple services to your local computer.
* [kurun](https://github.com/banzaicloud/kurun) - Run main.go in Kubernetes with one command, also port-forward your app into Kubernetes.
