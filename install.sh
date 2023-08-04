# (Un)Install script for installing ServerBoxMonitor

# Check root
if [ $(id -u) -ne 0 ]; then
    echo "Please run as root or use sudo"
    exit 1
fi

install() {
# Check arch: amd64 or arm64
arch=$(uname -m)
case $arch in
    x86_64)
        arch=amd64
        ;;
    aarch64)
        arch=arm64
        ;;
    *)
        echo "Not support arch: $arch"
        exit 1
        ;;
esac

# Check systemd
if [ ! -d /etc/systemd ]; then
    echo "Please install systemd"
    exit 1
fi

# Check wget
if ! command -v wget >/dev/null 2>&1; then
    echo "Please install wget"
    exit 1
fi

# Generate download url
newestTag=$(curl -s https://api.github.com/repos/lollipopkit/server_box_monitor/releases/latest | grep tag_name | cut -d '"' -f 4)
# Remove 'v' at the start -> "0.1.0"
newestTagLen=$(expr length $newestTag)
APPVER=$(expr substr $newestTag 2 $newestTagLen)
DOWNLOAD_URL="https://github.com/lollipopkit/server_box_monitor/releases/download/v${APPVER}/server_box_monitor_${APPVER}_linux_$arch.tar.gz"

# Download binary
echo "Download $DOWNLOAD_URL"
wget -q --show-progress $DOWNLOAD_URL -O /tmp/server_box_monitor.tar.gz
if [ ! -f /tmp/server_box_monitor.tar.gz ]; then
    echo "Download binary failed"
    exit 1
fi

# Extract binary
echo "Extracting binary..."
tar -xf /tmp/server_box_monitor.tar.gz -C /tmp
if [ ! -f /tmp/server_box_monitor ]; then
    echo "Extract binary failed"
    exit 1
fi

# Install binary
echo "Installing binary..."
mv /tmp/server_box_monitor /usr/local/bin/server_box_monitor
if [ $? -ne 0 ]; then
    echo "Install binary failed"
    exit 1
fi

# Install systemd service
echo "Installing systemd service..."
cat <<EOF > /etc/systemd/system/server_box_monitor.service
[Unit]
Description=Server Box Monitor
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/server_box_monitor serve
User=root
Restart=always

[Install]
WantedBy=multi-user.target
EOF

# Enable systemd service
echo "Enabling systemd service..."
systemctl enable server_box_monitor.service
if [ $? -ne 0 ]; then
    echo "Enable systemd service failed"
    exit 1
fi

# Start systemd service
echo "Starting systemd service..."
systemctl start server_box_monitor.service
if [ $? -ne 0 ]; then
    echo "Start systemd service failed"
    exit 1
fi

# Display systemd service status
echo "Displaying systemd service status..."
systemctl status server_box_monitor.service

# Clean up
echo "Cleaning up..."
rm /tmp/server_box_monitor.tar.gz
rm -rf /tmp/server_box_monitor

echo "Install success"
}

uninstall() {
# Stop systemd service
echo "Stopping systemd service..."
systemctl stop server_box_monitor.service
if [ $? -ne 0 ]; then
    echo "Stop systemd service failed"
    exit 1
fi

# Disable systemd service
echo "Disabling systemd service..."
systemctl disable server_box_monitor.service
if [ $? -ne 0 ]; then
    echo "Disable systemd service failed"
    exit 1
fi

# Remove systemd service
echo "Removing systemd service..."
rm /etc/systemd/system/server_box_monitor.service

# Remove binary
echo "Removing binary..."
rm /usr/local/bin/server_box_monitor

echo "Uninstall success"
}

# Request input
echo "ServerBoxMonitor (Un)Installer"
echo "----------------------------"
echo "1. Install"
echo "2. Uninstall"
echo "----------------------------"
echo -n "Please select: "
read action

case $action in
    1)
        install
        ;;
    2)
        uninstall
        ;;
    *)
        echo "Usage: $0 [install|uninstall]"
        exit 1
        ;;
esac