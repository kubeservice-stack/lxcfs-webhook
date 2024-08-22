#!/bin/bash

# Cleanup
nsenter -m/proc/1/ns/mnt fusermount -u /var/lib/lxcfs 2> /dev/null || true
nsenter -m/proc/1/ns/mnt [ -L /etc/mtab ] || \
        sed -i "/^lxcfs \/var\/lib\/lxcfs fuse.lxcfs/d" /etc/mtab

# Prepare
mkdir -p /usr/local/lib/lxcfs /var/lib/lxcfs

# Update lxcfs
cp -f /lxcfs/src/lxcfs /usr/local/bin/lxcfs
cp -f /lxcfs/src/.libs/liblxcfs.so /usr/local/lib/lxcfs/liblxcfs.so
cp -f /lxcfs/src/liblxcfs.la /usr/local/lib/lxcfs/liblxcfs.la


# Mount
exec nsenter -m/proc/1/ns/mnt /usr/local/bin/lxcfs /var/lib/lxcfs/ --enable-cfs -l -o nonempty
