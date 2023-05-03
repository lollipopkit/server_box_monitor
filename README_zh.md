[English](README.md) | 简体中文

## ServerBox 监测器
这个应用程序运行在服务器端, 监测服务器状态.  
这是 [ServerBox](https://github.com/lollipopkit/flutter_server_box) 项目的一部分.
**正处于活跃开发中，你可能需要在更新后重新配置.**

## 📖 使用方法
1. 在服务器上安装此应用程序
    - 如果你安装了 `go`, 你可以运行 `go install github.com/lollipopkit/server_box_monitor@latest` 来安装
    - 如果你没有安装 `go`, 你可以从 [发布](https://github.com/lollipopkit/server_box_monitor/releases) 下载二进制文件
2. 编辑配置文件
    - 配置文件保存在 `~/.config/server_box/config.json`
    - 完整的配置示例 [在这里](doc/CONFIG_zh.jsonc)
3. 运行.
    - 注意: 如果是下载的二进制文件，命令为 `./server_box_monitor`
    - 有多种方式运行
        1. (推荐) 配置 `systemd`
            - 示例配置文件 [这里](doc/srvbox.service)
            - Rootless
                - 复制文件到 `~/.config/systemd/user/srvbox.service`
                - `systemctl --user enable --now srvbox`
                - 你可以执行 `sudo loginctl enable-linger $USER` 让服务在注销后继续运行
            - Rootful
                - 复制文件到 `/etc/systemd/system/srvbox.service`
                - `systemctl enable --now srvbox`
        2. 使用 `screen`
            - 运行: `screen -S srvbox`, 然后 `server_box_monitor`
            - 移至后台（Detach）: `Ctrl + A + D`
            - 移至前台（Attach）: `screen -r srvbox`
        3. 直接运行 `server_box_monitor`
            - 这会在前台运行, 你可以使用 `Ctrl + C` 来停止.

## 🔖 License
`GPL v3. lollipopkit 2023`