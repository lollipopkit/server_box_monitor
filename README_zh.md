[English](README.md) | 简体中文

## ServerBox 监测器
这个应用程序运行在服务器端, 监测服务器状态.  
这是 [ServerBox](https://github.com/lollipopkit/flutter_server_box) 项目的一部分.

## 📖 使用方法
1. 在服务器上安装此应用程序
    - 如果你安装了 `go`, 你可以运行 `go install github.com/lollipopkit/server_box_monitor` 来安装
    - 如果你没有安装 `go`, 你可以从 [发布](https://github.com/lollipopkit/server_box_monitor/releases) 下载二进制文件
2. 编辑配置文件
    - 配置文件保存在 `~/.server_box/config.json`
    - 完整的配置示例 [在这里](CONFIG_zh.jsonc)
3. 执行 `server_box_monitor` 来运行
    - 如果你是下载的二进制, 你需要执行 `./server_box_monitor`

## 🔖 License
`GPL v3. lollipopkit 2023`