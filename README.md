English | [简体中文](README_zh.md)

## ServerBox Monitor
This app runs on server end and monitors the server status.  
It is a part of [ServerBox](https://github.com/lollipopkit/flutter_server_box) project.  

## 📖 Usage
1. Install the app on your server.
    - If you have `go` installed, you can run `go install github.com/lollipopkit/server_box_monitor`
    - If you don't have `go` installed, you can download the binary from [release page](https://github.com/lollipopkit/server_box_monitor/releases)
2. Edit the config file.
    - The config file is located at `~/.config/server_box/config.json`
    - Fully example is [here](doc/CONFIG.jsonc)
3. Run the app.
    - Note: If you download the binary, you need to run `./server_box_monitor`
    - There are several ways to run.
        1. (Recommended) Config `systemd` to run the app.
            - Example service file [here](doc/srvbox.service)
            - Rootless
                - Copy file to `~/.config/systemd/user/srvbox.service`
                - Run `systemctl --user enable --now srvbox`
                - You can run `sudo loginctl enable-linger $USER` to make the service run after logout
            - Rootful
                - Copy file to `/etc/systemd/system/srvbox.service`
                - Run `systemctl enable --now srvbox`
        2. Use `screen`
            - Run: `screen -S srvbox`, then `server_box_monitor`
            - Detach: `Ctrl + A`, then `D`
            - Reattach: `screen -r srvbox`
        3. Run `server_box_monitor` directly
            - It will run in foreground, you can use `Ctrl + C` to stop it.

## 🔖 License
`GPL v3. lollipopkit 2023`