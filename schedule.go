package main

import (
    "fmt"
    "time"
)

type Schedule struct {
    Pool *Pool
}

func (schedule *Schedule) Start(d time.Duration, name string, callback func(s *Schedule)) {
    go func() {
        for {
            fmt.Println("Running " + name)
            time.Sleep(d)
            callback(schedule)
        }
    }()
}
