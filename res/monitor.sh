# Script for app `ServerBox`
# Delete this file while app is running will cause app crash

export LANG=en_US.utf-8
echo SrvBox
cat /proc/net/dev && date +%s
echo SrvBox
cat /etc/os-release | grep PRETTY_NAME
echo SrvBox
cat /proc/stat | grep cpu
echo SrvBox
uptime
echo SrvBox
cat /proc/net/snmp
echo SrvBox
df -h
echo SrvBox
cat /proc/meminfo
echo SrvBox
cat /sys/class/thermal/thermal_zone*/type
echo SrvBox
cat /sys/class/thermal/thermal_zone*/temp