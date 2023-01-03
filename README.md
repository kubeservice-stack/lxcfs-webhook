# lxcfs-webhook
[![lxcfs docker publish](https://github.com/kubeservice-stack/lxcfs-webhook/actions/workflows/lxcfs.yml/badge.svg?branch=main)](https://github.com/kubeservice-stack/lxcfs-webhook/actions/workflows/lxcfs.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkubeservice-stack%2Flxcfs-webhook.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkubeservice-stack%2Flxcfs-webhook?ref=badge_shield)


Automatically deploy LXCFS while mounted to the container

## 动机
Pod 容器内资源可见性：让Pod的资源视角真实、准确

❓**是否有个发现：`Pod`中限定了`CPU`、`MEM`等资源大小，然而登入的`POD`中查询资源，却还是`Node`总的资源大小？**

对于**业务上云**, java（识别`内存资源`开辟`堆`大小）、golang(识别`CPU个数`开启`runtime线程个数`) 等语言，在`OOM`、`GC`方面的问题，有时常发生的原因

利用lxcfs将`容器`中读取出来的`CPU`、`MEM`、`disk`、`swaps`的信息是宿主机的信息，与容器实际分配和限制的资源量相同。 解决低层通过`os.syscall`获得的资源信息一致

## 依赖

Kubernetes: >= `1.16.0`




## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkubeservice-stack%2Flxcfs-webhook.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkubeservice-stack%2Flxcfs-webhook?ref=badge_large)