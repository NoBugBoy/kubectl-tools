# kubectl-tools

这是一款kubectl的工具集合，未来会集成一些好用的功能

## 源代码

https://github.com/NoBugBoy/kubectl-tools  点个star不过分

## 使用方法

kubernetes版本 > 1.12 + 

直接从release中下载 https://github.com/NoBugBoy/kubectl-tools/releases/tag/1.0

进入kubernetes目录 使用交叉编译打包为linux平台的包，将生成的kubectl-tools可执行文件放入k8s集群master节点的/root/bin目录下（kubectl plugin list）
```shell
cd kubernetes
CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o kubectl-tools
```
使用 kubectl tools -h 即可查看帮助


![img_1.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/21f43e80ef9340e785319bcd9ef10643~tplv-k3u1fbpfcp-watermark.image)

目前提供的功能

1. debug,提供一个带工具的容器，加入到目标容器的namespace中，在相同的视图下进行诊断，debug的目标节点如果是第一次操作，则需要多等待一些时间，等待拉取debug-k8s的镜像，还有指定的debug container的镜像，重复操作如果只更改了debug container的话就只拉取新的工具镜像，否则就不需要太多等待时间

![img0.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/beeaa0311d4d4d68a26ec131ac8938f5~tplv-k3u1fbpfcp-watermark.image)

![img.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/8903e05d171442e49cf2bba77e05f8ef~tplv-k3u1fbpfcp-watermark.image)


2. etcd集群备份 仅支持kubeadm安装的集群，使用前可以先使用下面的命令查看对应cm中是否存在etcd配置信息，包含ca证书和节点等
```shell
kubectl describe cm kubeadm-config -n kube-system
```

![1.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/386be153296642f8821efc6ee9ace992~tplv-k3u1fbpfcp-watermark.image)
```shell
[root@node0 ~]# kubectl tools backup -o /usr
Snapshot saved at /usr/2021.09.03-04:37:56-etcd-snapshot.db
```




debug实现思路参考 https://aleiwu.com/post/kubectl-debug-intro/
