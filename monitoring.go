package main

import (
    "encoding/json"
    "fmt"
    "runtime"
)

type MonitoringStats struct {
    Alloc            uint64
    TotalAlloc       uint64
    HeapAlloc        uint64
    Sys              uint64
    NumGC            uint32
    ServerStartedAt  string
    ActiveUsersCount int
}

func (m *MonitoringStats) ToJSON() (stringified string, err error) {
    b, err := json.Marshal(m)

    if err != nil {
        fmt.Println(err)

        return "", err
    }

    fmt.Println(string(b))

    return string(b), nil
}

func PrintMemUsage() MonitoringStats {
    var m runtime.MemStats

    runtime.ReadMemStats(&m)

    // For info on each, see: https://golang.org/pkg/runtime/#MemStats
    fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
    fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
    fmt.Printf("\tHeapAlloc  = %v MiB", bToMb(m.HeapAlloc))
    fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
    fmt.Printf("\tNumGC = %v\n", m.NumGC)

    return MonitoringStats{
        Alloc:      bToMb(m.Alloc),
        TotalAlloc: bToMb(m.TotalAlloc),
        HeapAlloc:  bToMb(m.HeapAlloc),
        Sys:        bToMb(m.Sys),
        NumGC:      m.NumGC,
    }
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
