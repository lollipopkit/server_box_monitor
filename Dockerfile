FROM alpine:3.18

# Install binary from Github
RUN apk update
ARG APPVER
RUN arch="$(apk --print-arch)"; case "$arch" in 'x86_64') arch="amd64"; ;; 'armhf') arch="armv6"; ;; 'armv7') arch="armv7"; ;; 'aarch64') arch="arm64"; ;; *) echo >&2 "error: unsupported architecture '$arch' (likely packaging update needed)"; exit 1 ;; esac \
    && wget "https://github.com/lollipopkit/server_box_monitor/releases/download/v${APPVER}/server_box_monitor_${APPVER}_linux_$arch.tar.gz" \
    && tar -xvf "server_box_monitor_${APPVER}_linux_$arch.tar.gz" \
    && rm "server_box_monitor_${APPVER}_linux_$arch.tar.gz" \
    && mv server_box_monitor /usr/bin \
    && chmod +x /usr/bin/server_box_monitor

# Chdir
WORKDIR /root/.config/server_box

ENTRYPOINT ["/usr/bin/server_box_monitor", "serve"]