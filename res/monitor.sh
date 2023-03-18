# Script for `ServerBoxMonitor`
# Please do not edit it

export LANG=en_US.utf-8
echo SrvBox
cat /proc/net/dev
echo SrvBox
cat /proc/stat | grep cpu
echo SrvBox
df -h
echo SrvBox
cat /proc/meminfo
echo SrvBox
cat /sys/class/thermal/thermal_zone*/type
echo SrvBox
cat /sys/class/thermal/thermal_zone*/temp