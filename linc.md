

# lindocker

## Linux基础概念

### Namespace

Linux Namespace是Kernel的一个功能，它可以隔离一系列的系统资源，如PID、User ID、Network

UTS Namespace：允许拥有自己的hostname

hostname -b



IPC Namespace：  用来隔离System V IPC 和POSIX message queues

 ipcs -q (查看现有的message queues);		ipcmk -Q （创建一个 message queue）



PID Namespace：用来隔离进程ID。同一个进程在不同的PID Namespace里可以拥有不同的PID。

echo $$（查看当前Namespace的PID）



Mount Namespace：用来隔离各个进程看到的挂载点视图。

User Namespace

Network Namespace：隔离网络设备、IP地址

ifconfig （查看自己的网络设备）





### Cgroups

Linux Cgroups (Control Groups)提供了对 一 组进程及将来子进程的资源限制、控制和统 计的能力，这些资源包括 CPU、内存、存储、网络等 。 通过 Cgroups，可以方便地限制某个进 程的资源占用，并且可以实时地监控进程的监控和统计信息 。

Cgroups 中的 3 个组件
。 cgroup 是对进程分组管理的一种机制， 一个 cgroup 包含 一 组进程，井可以在这个 cgroup 上增加 Linux subsystem 的各种参数配置，将一组进程和一组 subsystem 的系统参数关联 起来。

。 subsystem 是一组资源控制的模块， 

。hierarchy 的功能是把 一 组 cgroup 串成 一 个树状的结构，一个这样的树便是 一 个 hierarchy，通过这种树状结构， Cgroups 可以做到继承 。 比如，系统对 一组定 时的任务 进程通过 cgroupl 限制了 CPU 的使用率，然后其中有一个定时 dump 日志的进程还需要 限制磁盘 IO，为了避免限制了磁盘 IO 之后影响到其他进程，就可以创建 cgroup2，使 其继承于 cgroupl 井限制磁盘的 IO，这样 cgroup2 便继承了 cgroupl 中对 CPU 使用率的 限制，并且增加了磁盘 IO 的限制而不影响到 cgroupl 中的其他进程。

#### Kernel 接口

前面介绍了那么多 Cgroups 结构的内容，那么到底要怎么调用 Kernel 才能配置 Cgroups 呢?通过前面的介绍了解到， Cgroups 中的 hierarchy 是一种树状的组织结构， Kernel 为了使对 Cgroups 的配置更直观，是通过一个虚拟的树状文件系统配置 Cgroups 的，通过层级的目录虚 拟出 cgroup树。下面，就以一个配置的例子来了解一下如何操作 Cgroups。

1. 首先， 要创建并挂载一个 hierarchy Ccgroup树)，如下。
   • ~ mkdir cgroup-test #创建一个 hierarchy 挂载点
   • ~ sudo mount -t cgroup -o none,name=cgroup-test cgroup-test ./cgroup-test #挂载 一个 hierarchy
   • ~ ls ./cgroup test #挂载后我们就可以看到系统在这个目录下生成了一些默认文件
   cgroup. clone_children 	cgroup. procs	 cgroup. sane_behavior 	notify_on_release 	release_agent 	 tasks
   这些文件就是这个 hierarchy 中 cgroup 根节点的配置项，上面这些文件的含义分别如下。

> cgroup.clone_children, cpuset 的 subsystem 会读取这个配置文件，如果这个值是 I C默 认是 0)，子 cgroup 才会继承父 cgroup 的 cpuset 的配置 。
> 》 cgroup.procs 是树中 当前节点 cgroup 中的进程组 ID，现在的位置是在根节点，这个文 件中会有现在系统中所有进程组的 ID。
> notify_on_release和 release_agent会一起使用。 notify_on_release标识当这个 cgroup最 后一个进程退出的时候是否执行了 release_agent; release_agent则是一个路径，通常 用作进程退出之后自动清理掉不再使用的 cgroup。
> 》 tasks标识该 cgroup下面的进程 ID，如果把一个进程 ID写到 tasks文件中，便会将相 应的进程加入到这个 cgroup 中 。
2. 然后， 创建刚刚创 建好 的 hierarchy 上 cgroup 根节 点中扩展出的两个子 cgroup。
可以看到，在 一个 cgroup 的目录下创建文件夹时，
Kernel 会把文件夹标记为这个 cgroup
的子 cgroup，它们会继承父 cgroup 的属性。
2.  在 cgroup 中添加和移动进程 。
一个进程在一个 Cgroups 的 hierarchy 中，只能在一个 cgroup节点上存在，系统的所有 进程都会默认在根节点上存在 ，可以将进程移动到其他 cgroup 节点，只需要将进程 ID
写到移动到的 cgroup节点的 tasks文件中即可。
• [ cgroup-1] echo $$
7475 •[cgroup-1]sudosh-c”echo时》tasks”#将我所在的终端进程移动到cgroup-1中 •[ cgroup-1] cat /proc/7475/cgroup
13:name=cgroup-test:/cgroup-1
2. 通过 subsystem 限制 cgroup 中进程的资源 。
在上面创建 hierarchy 的时候，这个 hierarchy 并没有关联到任何的 subsystem，所以没办 法通过那个 hierarchy 中的 cgroup 节点限制进程的资源占用，其实系统默认已经为每个 subsystem 创建了 一个默认的 hierarchy，比如 memory 的 hierarchyo
第 2章基础技术 23


•~ mount I grep memory
cgroup on /sys/fs/cgroup/memory type cgroup
(rw, nosuid,nodev, noexec,relatime,memory, nsroot=/)
可以看到，/sys/fs/cgroup/memory 目录便是挂在了 memory subsystem 的 hierarchy上。



### AUFS

一种存储驱动类型，快速启动、高效利用内存

在 **mount aufs** 的命令中，没有指定待挂 载的 5个文件夹的权限， 默认的行为是， dirs指定的左边起第一个目录是 read-write权限， 后续 的都是 read-only权限。
$ sudo mount -t aufs -o dirs=./container-layer:./image-layer4:./image-layer3:./image- layer2 :. /image-layerl none . /mnt

#### image layer 和AUFS

image layer 的 内容存在host 的/var/lib/docker/aufs/**diff**目录下

在/var/lib/docker/aufs/**layers**目录存储者image layer 如何堆栈这些layer 的metadata

#### container layer 和 AUFS

Docker使用AUFS的CoW技术实现image layer 共享和**减少磁盘空间占用**。

启动一个container的时候，Docker会为其创建一个read-only的init layer，用来存储与这个容器内环境相关的内容；Docker还会为其创建一个read-write的layer来执行所有写操作。放在diff目录下



container layer 的**mount目录也是/var/lib/docker/aufs/mnt**。

container 的 metadata 和配置文 件都存放在/var/lib/docker/containers/<container-id>目录中。 container 的 read-write layer存储在 /var/lib/docker/aufs/diff/ 目录下。即使容器停止，这个可读写层仍然存在，因而重启容器不会丢 失数据，只有当 一个容器被删除的时候，这个可读写层才会一起删除。



eg：changed-ubuntu

**diff目录多了2个文件夹**，f9cc...-init,f9cc...,在layers目录cat可以看到f9cc...依赖与f9cc...-init,

同时/var/lib/docker/container/目录下多了一个与containerid相同的文件夹，存放着容器的metadata和config文件

从系统AUFS来看mount情况，在**/sys/fs/afus**目录下多了一个si_fe6d...文件夹，可以看到容器的layer权限：f9cc...=rw,f9cc...init（和其他镜像层）=ro+wh



AUFS删除文件file1，在container的rw层生成一个.wh.file1的文件夹来隐藏所有ro层的file1文件。



## 构造容器

