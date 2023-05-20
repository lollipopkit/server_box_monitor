English | [简体中文](README_zh.md)

## ServerBox Monitor
This app runs on server end and monitors the server status.  
It is a part of [ServerBox](https://github.com/lollipopkit/flutter_server_box) project.  
**It's under active development, you may need to reconfig it after upgrading.**

## 📖 Usage
1. There are serveral ways to install it.
   - `Docker`:
     - (Recommonded) [Docker compose](docker-compose.yaml)
     - Or `docker run -d --name srvbox -v ./config:/root/.config/server_box lollipopkit/srvbox_monitor:latest`
     - (Optional) If you need to update it, `docker rm srvbox -f && docker rmi lollipopkit/server_box_monitor:latest` to delete old image.
   - Use binary.
     - If you have `go` installed, you can run `go install github.com/lollipopkit/server_box_monitor@latest`
     - If you don't have `go` installed, you can download the binary from [release page](https://github.com/lollipopkit/server_box_monitor/releases)
     - (Recommended) Config `systemd` to run the app.
       - Example service file [here](doc/srvbox.service)
       - Rootless
         - Copy file to `~/.config/systemd/user/srvbox.service`
         - Run `systemctl --user enable --now srvbox`
         - You can run `sudo loginctl enable-linger $USER` to make the servicerun   after logout
       - Rootful
         - Copy file to `/etc/systemd/system/srvbox.service`
         - Run `systemctl enable --now srvbox`
2. Edit the config file.
   - The config file is located at
     - binary: `~/.config/server_box/config.json`
     - docker: `./config/config.json`
   - Fully example is [here](doc/CONFIG.jsonc)
    

## 🔖 License
`GPL v3. lollipopkit 2023`