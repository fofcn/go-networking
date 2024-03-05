快速搭建一个k8s环境
# 概述
我们在VirutBox安装了一个Ubuntu系统，这次我们在这个Ubuntu系统上安装一个k8s集群。
k8s的安装步骤繁多似乎有点复杂，今天我来更新一篇k8s的安装过程。至于我们为什么使用k8s我觉得我不用说太多，我觉得两个关键点：自动化和保证应用健康。

在本地搭建k8s集群有很多种方式：
1. minikube(最简方式：DockerDesktop+minikube)
2. Linux + k8s
3. Kind
3. MicroK8s
4. k3s
5. DockerDesktop

我们这里选择Linux + k8s官方的方式，是想让大家感受一下k8s安装的复杂过程，然后通过安装过程可以很直观的了解到k8s的一些关键点，比如：容器和流量分发。

在做k8s安装时我们还是要做很多选择：
1. 容器运行时(CRI)使用Docker,Containerd还是其他的
2. 网络转发用iptables还是ipvs
3. k8s的网络组件用什么？flannel还是canal还是Calico还是什么？
4. ...


我们这儿选择了containerd,ipvs和Calico（tigera是公司的名字，这个公司开源了Calico）。

如果你想选择其他比如：docker,iptable, cannel也是可以，根据自己的喜好而定，毕竟是学习类。

我们在安装一些组件使用了Helm，Helm是k8s yaml生成工具，大家先了解下就可以了。

安装过程的组织：
1. Linux环境配置
2. Containerd安装
3. k8s安装
4. k8s集群初始化
5. k8s服务部署


# 1. 安装K8s
## 1.1 配置containerd运行环境
创建/etc/modules-load.d/containerd.conf配置文件，确保在系统启动时自动加载所需的内核模块，以满足容器运行时的要求:

```shell
cat << EOF > /etc/modules-load.d/containerd.conf
overlay
br_netfilter
EOF
```
使配置生效
```shell

modprobe overlay
modprobe br_netfilter
```

## 1.2 创建/etc/sysctl.d/99-kubernetes-cri.conf
```shell

cat << EOF > /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
user.max_user_namespaces=28633
EOF
```

配置生效
```shell

sysctl -p /etc/sysctl.d/99-kubernetes-cri.conf

```

## 1.3 开启ipvs
```shell

cat > /etc/modules-load.d/ipvs.conf <<EOF
ip_vs
ip_vs_rr
ip_vs_wrr
ip_vs_sh
EOF
```
配置生效
```shell

modprobe ip_vs
modprobe ip_vs_rr
modprobe ip_vs_wrr
modprobe ip_vs_sh
```

安装ipvsadm
```shell

apt install -y ipset ipvsadm

```

## 1.5 安装containerd
```shell

wget https://github.com/containerd/containerd/releases/download/v1.7.3/containerd-1.7.3-linux-amd64.tar.gz

```
解压缩
```shell

tar Cxzvf /usr/local containerd-1.7.3-linux-amd64.tar.gz
```

## 1.6 安装runc:
```shell

wget https://github.com/opencontainers/runc/releases/download/v1.1.9/runc.amd64
install -m 755 runc.amd64 /usr/local/sbin/runc
```

## 1.7 生成containerd配置
```shell

mkdir -p /etc/containerd
containerd config default > /etc/containerd/config.toml
```

配置containerd使用systemd作为容器cgroup driver(这个一定要找对地方)
```shell

[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
  ...
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
    SystemdCgroup = true
```
国内已经修改sanbox image为阿里云的Registry，这样就不会有网络访问问题了。
```shell
[plugins."io.containerd.grpc.v1.cri"]
  ...
  # sandbox_image = "registry.k8s.io/pause:3.8"
  sandbox_image = "registry.aliyuncs.com/google_containers/pause:3.9"
```

## 1.8 下载containerd.service
链接：https://raw.githubusercontent.com/containerd/containerd/main/containerd.service,如果无法下载可以直接复制下面的内容
```shell

cat << EOF > /etc/systemd/system/containerd.service
[Unit]
Description=containerd container runtime
Documentation=https://containerd.io
After=network.target local-fs.target

[Service]
#uncomment to enable the experimental sbservice (sandboxed) version of containerd/cri integration
#Environment="ENABLE_CRI_SANDBOXES=sandboxed"
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/local/bin/containerd

Type=notify
Delegate=yes
KillMode=process
Restart=always
RestartSec=5
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNPROC=infinity
LimitCORE=infinity
LimitNOFILE=infinity
# Comment TasksMax if your systemd version does not supports it.
# Only systemd 226 and above support this version.
TasksMax=infinity
OOMScoreAdjust=-999

[Install]
WantedBy=multi-user.target
EOF
```

## 1.9 配置Containerd开机启动
```shell

systemctl daemon-reload
systemctl enable containerd --now 
systemctl status containerd
```

## 1.10 安装Crictl
```shell

wget https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.28.0/crictl-v1.28.0-linux-amd64.tar.gz
tar -zxvf crictl-v1.28.0-linux-amd64.tar.gz
install -m 755 crictl /usr/local/bin/crictl
```

测试Crictl
```shell

crictl --runtime-endpoint=unix:///run/containerd/containerd.sock  version

Version:  0.1.0
RuntimeName:  containerd
RuntimeVersion:  v1.7.3
RuntimeApiVersion:  v1
```

## 1.11 更新apt仓库
```shell

sudo apt-get update
# apt-transport-https may be a dummy package; if so, you can skip that package
sudo apt-get install -y apt-transport-https ca-certificates curl
```

## 1.12 下载k8s包仓库公钥
```shell

curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.28/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

```

## 1.13 添加k8s apt 仓库
```shell

echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.28/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list

```

## 1.14 安装kubelet, kubeadm和kubelet
```shell

sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
```

## 1.15 关闭系统swap
```shell

swapoff -a
```

永久关闭
```shell

vim /etc/fstab
```
开机启动kubelet
```shell

systemctl enable kubelet

```



# 2. 初始化K8s集群
```shell

kubeadm init --apiserver-advertise-address=your_host-only-ip --pod-network-cidr=10.244.0.0/16
```

初始化完成后
```shell

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

# 3. 安装Helm
```shell

wget https://get.helm.sh/helm-v3.12.3-linux-amd64.tar.gz
tar -zxvf helm-v3.12.3-linux-amd64.tar.gz
install -m 755 linux-amd64/helm  /usr/local/bin/helm
```

# 4. 安装k8s网络插件
下载tigera-operator
```shell

wget https://github.com/projectcalico/calico/releases/download/v3.26.1/tigera-operator-v3.26.1.tgz

```

查看chart中可定制的配置
```shell

helm show values tigera-operator-v3.26.1.tgz
```
做点简单配置定制,保存为vlaues.yaml
```yaml

apiServer:
  enabled: false
installation:
  kubeletVolumePluginPath: None
```

Heml安装colico
```shell

helm install calico tigera-operator-v3.26.1.tgz -n kube-system  --create-namespace -f values.yaml

```
等待Pod处于Running
```shell

kubectl get pod -n kube-system | grep tigera-operator
```
安装kubectl插件
```shell

cd /usr/local/bin
curl -o kubectl-calico -O -L  "https://github.com/projectcalico/calicoctl/releases/download/v3.21.5/calicoctl-linux-amd64" 
chmod +x kubectl-calico
```

验证是否正常工作
```shell

kubectl calico -h

```

# 5. 测试
## 5.1 验证k8s DNS
```shell

kubectl run curl --image=radial/busyboxplus:curl -it
nslookup kubernetes.default
```

## 5.2 发布一个nginx
命令
```shell

kubectl apply -f nginx.yaml
```
yaml文件如下：
```yaml

apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  annotations:
    change-cause: "Rollout test"
spec:
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  selector:
    matchLabels:
      app: nginx
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
        resources:
          requests:
            cpu: 200m
          limits:
            cpu: 500m
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  type: NodePort
  selector:
    app: nginx
  ports:
    - protocol: TCP
      port: 40000
      targetPort: 80
      nodePort: 32000

```

# 6. FAQ
1. Nameserver limits exceeded" err="Nameserver limits were exceeded, some nameservers have been omitted, the applied nameserver line is: 10.22.16.2 10.22.16.254 10.245.0.10"

```shell

/run/systemd/resolve/resolv.conf
# 注释掉几个ip
```

2. Master节点作为Node

```shell

kubectl taint node k8s-master node-role.kubernetes.io/master:NoSchedule-
```

3. 忘记join集群命令

```shell

kubeadm token create --print-join-command
```

4. 部署curl测试

```shell

kubectl run curl --image=radial/busyboxplus:curl -it
```

5. kubelet日志查看

```shell

journal -xeu kubelet
journal -xeu kubelet > kubelet.log
```

6. kubeadm初始化有问题可以尝试reset后重新初始化

```shell

kubeadm reset
```

# 7. 总结
k8s安装步骤还是比较繁杂，不过通过k8s安装过程你可以学习到一些k8s底层原理。k8s是在容器技术之上，使用Linux的流量分发技术并将自己的流量路由规则应用到Linux的分发技术之上。

至于k8s的技术自动化就在于它对于你提出需求的严格保证，这里的需求就是我们发布到K8s应用中的资源需求（比如：几个Replica,多少CPU，多少内存，多少硬盘）。