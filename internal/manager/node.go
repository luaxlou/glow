package manager

import (
	"github.com/luaxlou/glow/pkg/api"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetNodeStatus() (*api.Node, error) {
	// Host Info
	hInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	// CPU
	cpuPercents, err := cpu.Percent(0, false)
	var cpuUsage float64
	if err == nil && len(cpuPercents) > 0 {
		cpuUsage = cpuPercents[0]
	}

	// Mem
	vMem, err := mem.VirtualMemory()
	var memUsage float64
	if err == nil {
		memUsage = vMem.UsedPercent
	}

	// Disk (Root)
	dUsage, err := disk.Usage("/")
	var diskUsage float64
	if err == nil {
		diskUsage = dUsage.UsedPercent
	}

	// Resources
	resources, _ := ListResources()

	return &api.Node{
		TypeMeta: api.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name: hInfo.Hostname,
		},
		Status: api.NodeStatus{
			Hostname:  hInfo.Hostname,
			OS:        hInfo.OS,
			Arch:      hInfo.KernelArch,
			Kernel:    hInfo.KernelVersion,
			CPUUsage:  cpuUsage,
			MemUsage:  memUsage,
			DiskUsage: diskUsage,
			Resources: resources,
		},
	}, nil
}
