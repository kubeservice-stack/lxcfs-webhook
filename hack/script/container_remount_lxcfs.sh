#! /bin/bash

PATH=$PATH:/bin
LXCFS="/var/lib/lxcfs"
LXCFS_ROOT_PATH="/var/lib"

containers=$(docker ps | grep -v pause  | grep -v calico ｜ grep -v cilium | awk '{print $1}' | grep -v CONTAINE)

for container in $containers;do
	mountpoint=$(docker inspect --format '{{ range .Mounts }}{{ if eq .Destination "/var/lib/lxc" }}{{ .Source }}{{ end }}{{ end }}' $container)
	# 确保本身pod中就有mount point
	if [ "$mountpoint" = "$LXCFS_ROOT_PATH" ];then
		echo "remount $container"
		PID=$(docker inspect --format '{{.State.Pid}}' $container)
		# mount /proc
		for file in meminfo cpuinfo loadavg stat diskstats swaps uptime;do
			echo nsenter --target $PID --mount --  mount -B "$LXCFS/proc/$file" "/proc/$file"
			nsenter --target $PID --mount --  mount -B "$LXCFS/proc/$file" "/proc/$file"
		done
		# mount /sys
		for file in online;do
			echo nsenter --target $PID --mount --  mount -B "$LXCFS/sys/devices/system/cpu/$file" "/sys/devices/system/cpu/$file"
			nsenter --target $PID --mount --  mount -B "$LXCFS/sys/devices/system/cpu/$file" "/sys/devices/system/cpu/$file"
		done 
	fi 
done