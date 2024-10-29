#!/bin/bash +x

LXCFS="/var/lib/lxc/lxcfs"

containers=$(crictl ps | grep -v pause  | grep -v calico | grep -v cilium  |awk '{print $1}')
for container in $containers; do
    # 获取挂载点信息
    mounts=$(crictl inspect $container | jq -r '.info.config.mounts[] | "\(.container_path) -> \(.host_path)"' | grep "$LXCFS/")
    
    echo "Mounts for container $container:"
    echo "$mounts"
    
    # 检查是否有挂载到 LXCFS 路径
    echo "$mounts" | grep  "$LXCFS/"
    if [ $? -eq 0 ]; then
        echo "remount $container"
        PID=$(crictl inspect $container | jq -r '.info.pid')
        # mount /proc
        for file in meminfo cpuinfo loadavg stat diskstats swaps uptime; do
            echo nsenter --target $PID --mount -- /bin/mount -o remount,bind  -t proc "$LXCFS/proc/$file" "/proc/$file"
            nsenter --target $PID --mount --  /bin/mount -B "$LXCFS/proc/$file" "/proc/$file"
        done
        
        # mount /sys
        for file in online; do
            echo nsenter --target $PID --mount -- /bin/mount -o remount,bind "$LXCFS/sys/devices/system/cpu/$file" "/sys/devices/system/cpu/$file"
            nsenter --target $PID --mount -- /bin/mount -B "$LXCFS/sys/devices/system/cpu/$file" "/sys/devices/system/cpu/$file"
        done 
    else
        echo "容器 $container 没有挂载 /var/lib/lxcfs"
    fi
done