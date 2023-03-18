English | [简体中文](README_zh.md)

## ServerBox Monitor
This app runs on server end and monitors the server status.  
It is a part of [ServerBox](https://github.com/lollipopkit/flutter_server_box) project.  

## 📖 Usage
1. Install the app on your server.
    - If you have `go` installed, you can run `go install github.com/lollipopkit/server_box_monitor`
    - If you don't have `go` installed, you can download the binary from [release page](https://github.com/lollipopkit/server_box_monitor/releases)
2. Edit the config file.
    - The config file is located at `~/.server_box/config.json`
    - Fully example is [here](doc/CONFIG.jsonc)
3. Run the app with `server_box_monitor` command.
    - If you download the binary, you can run `./server_box_monitor`

## 🔖 License
`GPL v3. lollipopkit 2023`