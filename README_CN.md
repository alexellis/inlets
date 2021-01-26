# inlets

<!--翻译贡献者请注意译文三大原则：信、达、雅。格式和 Markdown 风格请与英文原文保持一致。-->

轻松将服务暴露到公网或其它网络，穿透内网、代理和 NAT。

[![Build Status](https://travis-ci.com/inlets/inlets.svg?branch=master)](https://travis-ci.com/inlets/inlets)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/inlets/inlets)](https://goreportcard.com/report/github.com/inlets/inlets)
[![Documentation](https://godoc.org/github.com/inlets/inlets?status.svg)](http://godoc.org/github.com/inlets/inlets)
![GitHub All Releases](https://img.shields.io/github/downloads/inlets/inlets/total)

[Follow @inletsdev on Twitter](https://twitter.com/inletsdev)

[English](./README.md) | [中文文档](./README_CN.md)

## 简介

inlets&reg; 利用反向代理和 Websocket 隧道，将内部、或是开发中的服务通过「出口节点」暴露到公网。出口节点可以是几块钱一个月的 VPS，也可以是任何带有公网 IPv4 的电脑。
如果无需公开服务到公网，可以只建立服务隧道，使 inlets 成为 VPN 的替代品。

为什么需要这个项目？类似的工具例如 [ngrok](https://ngrok.com/) 和由 [Cloudflare](https://www.cloudflare.com/) 开发的 [Argo Tunnel](https://developers.cloudflare.com/argo-tunnel/) 皆为闭源，内置了一些限制，并且价格不菲，以及对 arm/arm64 的支持很有限。Ngrok 还经常会被公司防火墙策略拦截而导致无法使用。而其它开源的隧道工具，基本只考虑到静态地配置单个隧道。inlets 旨在动态地发现本地服务，通过 Websocket 隧道将它们暴露到公网 IP 或域名，并自动化配置 TLS 证书。

当开启 SSL 时，inlets 可以通过任何支持 `CONNECT` 方法的 HTTP 代理服务器。

![概念示意图](docs/inlets.png)

*inlets 概念示意图*

## Built for developers by developers

<a href="https://github.com/sponsors/inlets/">
<img alt="Sponsor this project" src="https://github.com/alexellis/alexellis/blob/master/sponsor-today.png" width="90%">
</a>

## 协议与条款

**重要**

如您需要在企业网络中使用 inlets，建议先征求 IT 管理员的同意。下载、使用或分发 inlets 前，您必须同意 [协议](./LICENSE) 条款与限制。本项目不提供任何担保，亦不承担任何责任。

### 待办事项和目标

#### 已完成

* 基于客户端的定义，自动在出口节点创建服务入口
  * 通过 DNS / 域名实现单端口、单 Websocket 承载多站点
* 利用 SSL over Websockets 实现链路加密（`wss://`）
* 服务端和客户端验权
* 自动重连
* 原生多架构（ARMHF / ARM64）支持
* 提供 Dockerfile 以及 Kubernetes YAML 文件
* 自动发现并实例化 Kubernetes 集群内 `LoadBalancer` 类型的 `Service` - [inlets-operator](https://github.com/inlets/inlets-operator)
* 传输 Websocket 流量
* [为该项目制作一枚 Logo](https://github.com/inlets/inlets/issues/46)
* 和反向代理或 inlets PRO 搭配时支持配置 TLS 证书

#### inlets PRO

以下是 [inlets PRO](https://inlets.dev) 的特性和用例：

* 传输 4 层的 TCP 流量，例如 Websockets、数据库流量、反向代理、远程桌面和 SSH
* 暴露同一个客户端的多个端口 - 例如 80 和 443
* 作为反向代理或 Kubernetes IngressController 运行
* 自动为控制平面部署 TLS 证书
* 商业服务和客户支持
* 文档、博客、教程、视频等

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

Inlets 作为 *服务代理* [已被列入 Cloud Native Landscape](https://landscape.cncf.io/category=service-proxy&format=card-mode&grouping=category&sort=stars)

* [inlets PRO](https://inlets.dev) - 云原生的隧道服务 - TCP、HTTP 以及 Websockets，支持全自动 TLS 加密
* [inlets](https://github.com/inlets/inlets) - 云原生的隧道服务，只支持 HTTP - TLS 需要单独配置
* [inlets-operator](https://github.com/inlets/inlets-operator) - 给私有 Kubernetes Services 实现公网 IP，支持 CRD
* [inletsctl](https://github.com/inlets/inletsctl) - 搭建出口节点的最快方法

## 安装 inlets

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

Windows 用户建议使用 [Git bash](https://git-scm.com/downloads) 来安装 inlets。

## 使用 inlets

### 视频 Demo

使用 inlets 实现为我的 JavaScript & Webpack 应用配置公开的服务入口，以及自定义的域名：[Create React App](https://github.com/facebook/create-react-app)。

[![https://img.youtube.com/vi/jrAqqe8N3q4/hqdefault.jpg](https://img.youtube.com/vi/jrAqqe8N3q4/maxresdefault.jpg)](https://youtu.be/jrAqqe8N3q4)

### 快速上手

你可以在任何两台互相连接的「电脑」之间运行 inlets，「电脑」可以是两个容器，虚拟机，物理机，甚至单台机器的环回网络也可以。

可在本地尝试 [quickstart tutorial now](./docs/quickstart.md)（英文）。

### 文档和教程

inlets 和 inlets PRO 有了独立的文档站点（英文）：

官方文档：[docs.inlets.dev](https://docs.inlets.dev)

* 文档：[Quickstart tutorial on your laptop](./docs/quickstart.md)
* 文档：[Inlets & Kubernetes recipes](./docs/kubernetes.md)
* 教程：[Get a LoadBalancer for your private Kubernetes cluster with inlets-operator](https://blog.alexellis.io/ingress-for-your-local-kubernetes-cluster/)

延伸阅读: [advanced usage of inlets including Docker, Kubernetes, multiple-services, and binding to private IPs](./docs/advanced.md)

### 大家如何评价 inlets？

读一读社区教程、Hacker News 推送，如果你有编写关于 inlets 和 inlets PRO 的内容，欢迎提交 PR：

* [社区教程](docs/community.md)

> 你可以使用这些关键词在社交媒体分享 inlets：`@inletsdev`、`#inletsdev` 和 `https://inlets.dev`。

### 工作用途或是在生产环境使用 inlets？

查看 [ADOPTERS.md](./ADOPTERS.md) 来了解如今有哪些公司在使用 inlets 了。

### 周边商品

查看 [OpenFaaS Ltd SWAG store](https://store.openfaas.com/) 获得属于你自己的 inlets 卫衣、T 恤和水杯。

<img src="https://pbs.twimg.com/media/EQuxmEJWoAAP0Ga?format=jpg&name=small" width=300>

### 开发指引

查看 [CONTRIBUTING.md](./CONTRIBUTING.md)

### 其它 Kubernetes 端口转发工具

* [`kubectl port-forward`](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) - built into the Kubernetes CLI, forwards a single port to the local computer.
* [kubefwd](https://github.com/txn2/kubefwd) - Kubernetes utility to port-forward multiple services to your local computer.
* [kurun](https://github.com/banzaicloud/kurun) - Run main.go in Kubernetes with one command, also port-forward your app into Kubernetes.

inlets&reg; is a registered trademark of OpenFaaS Ltd. All rights reserved, registered company in the UK: 11076587
