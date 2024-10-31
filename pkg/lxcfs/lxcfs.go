/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lxcfs

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

// -v /var/lib/lxc/lxcfs/proc/cpuinfo:/proc/cpuinfo:ro
// -v /var/lib/lxc/lxcfs/proc/diskstats:/proc/diskstats:ro
// -v /var/lib/lxc/lxcfs/proc/meminfo:/proc/meminfo:ro
// -v /var/lib/lxc/lxcfs/proc/stat:/proc/stat:ro
// -v /var/lib/lxc/lxcfs/proc/swaps:/proc/swaps:ro
// -v /var/lib/lxc/lxcfs/proc/uptime:/proc/uptime:ro
// -v /var/lib/lxc/lxcfs/proc/loadavg:/proc/loadavg:ro
// -v /var/lib/lxc/lxcfs/sys/devices/system/cpu/online:/sys/devices/system/cpu/online:ro
var VolumeMountsTemplate = []corev1.VolumeMount{

	{
		Name:      "lxcfs-proc-cpuinfo",
		MountPath: "/proc/cpuinfo",
		ReadOnly:  true,
	},
	{
		Name:      "lxcfs-proc-meminfo",
		MountPath: "/proc/meminfo",
		ReadOnly:  true,
	},
	{
		Name:      "lxcfs-proc-diskstats",
		MountPath: "/proc/diskstats",
		ReadOnly:  true,
	},
	{
		Name:      "lxcfs-proc-stat",
		MountPath: "/proc/stat",
		ReadOnly:  true,
	},
	{
		Name:      "lxcfs-proc-swaps",
		MountPath: "/proc/swaps",
		ReadOnly:  true,
	},
	{
		Name:      "lxcfs-proc-uptime",
		MountPath: "/proc/uptime",
		ReadOnly:  true,
	},
	{
		Name:      "lxcfs-proc-loadavg",
		MountPath: "/proc/loadavg",
		ReadOnly:  true,
	},
	{
		Name:      "lxcfs-sys-devices-system-cpu-online",
		MountPath: "/sys/devices/system/cpu/online",
		ReadOnly:  true,
	},
	{
		Name:             "var-lib-lxc",
		MountPath:        "/var/lib/lxc",
		ReadOnly:         true,
		MountPropagation: ptr.To(corev1.MountPropagationHostToContainer),
	},
}
var VolumesTemplate = []corev1.Volume{
	{
		Name: "lxcfs-proc-cpuinfo",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxc/lxcfs/proc/cpuinfo",
				Type: ptr.To(corev1.HostPathFile),
			},
		},
	},
	{
		Name: "lxcfs-proc-diskstats",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxc/lxcfs/proc/diskstats",
				Type: ptr.To(corev1.HostPathFile),
			},
		},
	},
	{
		Name: "lxcfs-proc-meminfo",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxc/lxcfs/proc/meminfo",
				Type: ptr.To(corev1.HostPathFile),
			},
		},
	},
	{
		Name: "lxcfs-proc-stat",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxc/lxcfs/proc/stat",
				Type: ptr.To(corev1.HostPathFile),
			},
		},
	},
	{
		Name: "lxcfs-proc-swaps",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxc/lxcfs/proc/swaps",
				Type: ptr.To(corev1.HostPathFile),
			},
		},
	},
	{
		Name: "lxcfs-proc-uptime",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxc/lxcfs/proc/uptime",
				Type: ptr.To(corev1.HostPathFile),
			},
		},
	},
	{
		Name: "lxcfs-proc-loadavg",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxc/lxcfs/proc/loadavg",
				Type: ptr.To(corev1.HostPathFile),
			},
		},
	},
	{
		Name: "lxcfs-sys-devices-system-cpu-online",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxc/lxcfs/sys/devices/system/cpu/online",
				Type: ptr.To(corev1.HostPathFile),
			},
		},
	},
	{
		Name: "var-lib-lxc",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/lib/lxc/",
				Type: ptr.To(corev1.HostPathDirectory),
			},
		},
	},
}
