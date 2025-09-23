#! /bin/bash +x

PATH=$PATH:/bin
LXCFS="/var/lib/lxc/lxcfs"
LXCFS_ROOT_PATH="/var/lib/lxc"

containers=$(docker ps | grep -v pause  | grep -v calico | grep -v cilium | awk '{print $1}' | grep -v CONTAINE)

for container in $containers;do
	mountpoint=$(docker inspect --format '{{ range .Mounts }}{{ if eq .Destination "/var/lib/lxc" }}{{ .Source }}{{ end }}{{ end }}' $container)
	# 确保本身pod中就有mount point
	if [ "$mountpoint" = "$LXCFS_ROOT_PATH" ];then
		echo "remount $container"
		PID=$(docker inspect --format '{{.State.Pid}}' $container)
		if nsenter --target $PID --mount -- -- which mount >/dev/null 2>&1; then
		    # mount /proc
		    for file in meminfo cpuinfo loadavg stat diskstats swaps uptime;do
			    echo nsenter --target $PID --mount -- mount -o bind "$LXCFS/proc/$file" "/proc/$file"
			    nsenter --target $PID --mount -- mount -o bind "$LXCFS/proc/$file" "/proc/$file"
		    done
		    # mount /sys
		    for file in online;do
		        echo nsenter --target $PID --mount -- mount -o bind "$LXCFS/sys/devices/system/cpu/$file" "/sys/devices/system/cpu/$file"
			    nsenter --target $PID --mount -- mount -o bind "$LXCFS/sys/devices/system/cpu/$file" "/sys/devices/system/cpu/$file"
		    done 
        else
            echo "容器 $container 没有 mount 软件"
        fi
    else
        echo "容器 $container 没有挂载 /var/lib/lxc/lxcfs"
	fi
done
exit 0