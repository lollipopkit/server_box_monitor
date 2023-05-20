[English](README.md) | 简体中文

## ServerBox 监测器
这个应用程序运行在服务器端, 监测服务器状态.  
这是 [ServerBox](https://github.com/lollipopkit/flutter_server_box) 项目的一部分.
**正处于活跃开发中，你可能需要在更新后重新配置.**

## 📖 使用方法
1. 这里有多种方式安装.
   - `Docker`:
     - (推荐) [Docker compose](docker-compose.yaml)
     - 或者 `docker run -d --name srvbox -v ./config:/root/.config/server_box lollipopkit/srvbox_monitor:latest`
     - 如果你要更新, 先执行 `docker rm srvbox -f && docker rmi lollipopkit/server_box_monitor:latest` 来删除旧的镜像.
   - 可执行文件.
     - 如果你有安装 `go`, `go install github.com/lollipopkit/server_box_monitor@latest`
     - 或者从 [发布](https://github.com/lollipopkit/server_box_monitor/releases) 下载
     - (推荐) 使用 `systemd` 来运行.
       - 示例文件在 [这里](doc/srvbox.service)，请阅读文件中的注释！
       - 非 root
         - 复制示例文件到 `~/.config/systemd/user/srvbox.service`
         - `systemctl --user enable --now srvbox`
         -  `sudo loginctl enable-linger $USER` 让服务在注销后继续运行.
       - root
         - 复制示例文件到 `/etc/systemd/system/srvbox.service`
         - `systemctl enable --now srvbox`
2. 修改配置.
   - 配置文件在
     - 二进制: `~/.config/server_box/config.json`
     - docker: `./config/config.json`
   - 完整配置模版在 [这里](doc/CONFIG.jsonc)

## 🔖 许可证
`GPL v3. lollipopkit 2023`