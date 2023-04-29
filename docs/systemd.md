#### systemd主要概念和术语

- 系统初始化需要启动后台服务，需要完成一系列配置工作（如挂载文件系统），其中每一步或每一项任务都被systemd抽象为一个单元（unit），一个服务、一个挂载点、一个文件路径都可以被视为单元。
- systemd单元类型

| 单元类型            | 配置文件扩展名    | 说明                                                   |
|-----------------|------------|------------------------------------------------------|
| service(服务)     | .service   | 定义系统服务。这是最常用的一类，与早期Linux版本的/etc/init.d/目录下的服务脚本的作用相同 |
| device（设备）      | .device    | 定义内核识别的设备。每一个使用udev规则标识的设备都会在systemd中作为一个设备单元出现      |
| mount(挂载)       | .mount     | 定义文件挂载点                                              |
| automount(自动挂载) | .automount | 用于文件系统自动挂载设备                                         |
| socket(套接字)     | .socket    | 定义系统和互联网中的一个套接字，标识进程间通信用到的socket文件                   |
| swap(交换空间)      | .swap      | 标识管理用于交换空间的设备                                        |
| path(路径)        | .path      | 定义文件系统中的文件或目录                                        |
| timer(定时器)      | .timer     | 用来定时触发用户定义的操作，以取代atd、crond等传统的定时服务                   |
| target(目标)      | .target    | 用于对其他单元进行逻辑分组，主要用于模拟实现运行级别的概念                        |
| snapshot(快照)    | .snapshot  | 快照是一组配置单元，保存了系统当前的运行状态                               |

- systemd事物
- systemd能够保证事务的完整性。次事务与数据库中有所不同，旨在保证多个依赖单元之间没有循环引用。你如单元A、B、C之间存在循环依赖，systemd将无法启动任意一个服务
- systemd将单元之间的依赖关系分为两种
    - required 强依赖
    - wants 弱依赖，systemd将去除wants关键字指定的弱依赖，以打破循环。若无法修复，则systemd会报错。

#### systemd单元文件

- systemd对服务、设备、套接字和挂载点进行控制管理，都是有单元文件实现的。
- 单元文件分为若干节（Section，也可以译为段）
- 第一节用来定义单元的通用选项、配置与其他单元的关系，常用的字段指定如下：
    - Description：提供简短的描述信息
    - Requires：指定当前单元所依赖的其他单元. 这是强依赖，被依赖的单元无法启动时，当前单元也无法启动
    - Wants：指定与当前单元配合的其他单元。这是弱依赖，被依赖的单元无法启动时，当前单元可以被激活
    - Before和After：指定当前单元启动的前后单元。
    - Conflicts: 定义单元之间的冲突关系。列如此字段中的单元如果正在运行，此单元就不能运行，反之亦然。
- [Install]通常是配置文件的最后一节，用来定义如何启动，以及是否开机启动。常用的字段（指令）如下：
    - Alias：当前单元的别名
    - Also：与当前单元一起安装或被协助的单元。
    - RequiredBy: 指定被哪些单元所依赖，这是强依赖
    - WantedBy:指定被哪些单元所依赖，这是弱依赖。
- 对于新创建或修改过的单元文件，必须让systemd重新识别下此配置文件，通常执行``` systemctl daemon-reload ``` 命令重载配置。
- 建议将手动创建的单元文件放在/etc/systemd/system/目录下
- 开机启动``` systemctl enable cpus.service ``` 执行该命令后会看到在/etc/systemd/system/multi-user.target.wants/创建了链接文件
- 禁止开机启动``` systemctl disable cpus.service ``` 执行该命令就是删除/etc/systemd/system/目录下相应的链接文件

#### target单元文件

- 启动目标使用target单元文件描述，target单元文件的扩展名是.target，目的是将其他systemd单元文件通过一连串的依赖关系组织在一起。

#### 创建自定义服务

- Type 配置单元进程启动时的类型，simple为默认值，
- ExecStart 指定启动单元的命令或脚本，ExecStartPre、ExecStartPost字段指定在ExecStart之前或者之后用户自定义执行脚本。
- ExecStop 指定启动单元停止时执行的脚本
- ExecReload 指定单元重新装载时执行的脚本
- Restart 如果设置为always，服务重启时进程会退出，会通过systemctl命令清除并重启的操作。

#### systemd单元管理

- 单元的活动状态
    - active 正在运行
    - inactive 没有运行
    - failed 表示运行不成功
    - running 表示一次或多次持续运行
    - exited 表示成功完成一次性配置，仅运行一次就正常结束，目前已经没有运行该进程
    - waiting 表示正在运行中，不过还需要再等待其他事件才能再继续处理
    - dead 表示没有运行
    - failed 表示运行失败
    - mounted 表示成功挂载（文件系统）
    - plugged 表示已接入（设备）

| 功能               | 命令                                                                                                                         |
|------------------|----------------------------------------------------------------------------------------------------------------------------|
| 启动服务             | systemctl start 服务名                                                                                                        |
| 停止服务             | systemctl stop 服务名                                                                                                         |
| 重启服务             | systemctl restart 服务名                                                                                                      |
| 查看服务运行状态         | systemctl status 服务名                                                                                                       |
| 重载服务的配置文件而不重启服务  | systemctl reload 服务名                                                                                                       |
| 条件式重启服务          | systemctl tryrestart 服务名                                                                                                   |
| 重载或重启服务          | systemctl reload-or-restart 服务名                                                                                            |
| 重载或条件式重启         | systemctl reload-or-restart 服务名                                                                                            |
| 查看服务是否激活（正在运行）   | systemctl is-active 服务名                                                                                                    |
| 查看服务启动是否失败       | systemctl is-failed 服务名                                                                                                    |
| 杀死服务             | systemctl kill 服务名                                                                                                         |
| 列出所有已装载的单元       | systemctl list-units 或者 systemctl list-units --all 或者 systemctl list-units --all --state=not-found(此处状态有dead、active等，参考上面) |                                                                                  |
| 根据服务类型列出所有已装载的单元 | systemctl list-units -t service 、systemctl list-units --type=device                                                        |
| 列出指定单元的所有依赖      | systemctl list-dependencies cpus                                                                                           |
| 列出单元文件（可用的）      | systemctl list-unit-files                                                                                                  |
| 禁止某服务设为开机启动      | systemctl mask 服务名                                                                                                         |
| 取消禁止某服务设为开机启动    | systemctl unmask 服务名                                                                                                       |
| 查看某服务是否能够开机启动    | systemctl is-enabled 服务名                                                                                                   |

#### 常用例子

- systemd.timer定时任务 [systemd.timer定时任务.pdf](..%2Fstatic%2Fsystemd.timer%E5%AE%9A%E6%97%B6%E4%BB%BB%E5%8A%A1.pdf)
- systemd.service 服务 [risk-control.service](..%2Fdemo%2Frisk-control.service)