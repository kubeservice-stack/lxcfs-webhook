#!/bin/bash +x

LXCFS="/var/lib/lxc/lxcfs"

containers=$(crictl ps | grep -v pause  | grep -v calico | grep -v cilium  |awk '{print $1}' | grep -v CONTAINER)
for container in $containers; do
    # 获取挂载点信息
    mounts=$(crictl inspect -o go-template --template='{{range .info.config.mounts}}{{.container_path}} -> {{.host_path}}{{println}}{{end}}' $container | grep "$LXCFS/")
    
    echo "Mounts for container $container:"
    echo "$mounts"
    
    # 检查是否有挂载到 LXCFS 路径
    count=$(echo "$mounts" | grep  "$LXCFS/" | wc -l)
    if [ "$count" != "0" ]; then
        echo "remount $container"
        PID=$(crictl inspect --output go-template --template '{{- .info.pid -}}' $container)
        # mount /proc
        for file in meminfo cpuinfo loadavg stat diskstats swaps uptime; do
            echo nsenter --target $PID --mount -- mount -B "$LXCFS/proc/$file" "/proc/$file"
            nsenter --target $PID --mount -- mount -B "$LXCFS/proc/$file" "/proc/$file"
        done

        # mount /sys
        for file in online; do
            echo nsenter --target $PID --mount -- mount -B "$LXCFS/sys/devices/system/cpu/$file" "/sys/devices/system/cpu/$file"
            nsenter --target $PID --mount -- mount -B "$LXCFS/sys/devices/system/cpu/$file" "/sys/devices/system/cpu/$file"
    else
        echo "容器 $container 没有挂载 /var/lib/lxc/lxcfs"
    fi
done
exit 0