package main

import (
    "github.com/codegangsta/cli"
    "github.com/hhatto/nanairo"
    "bytes"
    "fmt"
    "io"
    "encoding/json"
)

type healthReport struct {
    Score int `json:"score"`
}

type Job struct {
    Name string `json:"name"`
    Url string `json:"url"`
    Color string `json:"color"`
    HealthReport []healthReport `json:"healthReport"`
}

func getJobs(hostname string) ([]Job, error) {
    client := NewClient(hostname)
    res, err := client.get("api/json?depth=1")
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    var r struct {
        Jobs []Job `json:"jobs"`
    }
    err = json.NewDecoder(io.Reader(res.Body)).Decode(&r)
    if err != nil {
        return nil, err
    }

    return r.Jobs, nil
}

func jobs(c *cli.Context) {
    if len(c.Args()) == 0 {
        fmt.Println("set target hostname")
        return
    }

    jobs, _ := getJobs(c.Args()[0])
    for _, job := range jobs {
        // S
        var j = bytes.NewBufferString("")
        if job.Color == "blue" {
            j.WriteString(nanairo.FgColor("#0c0", "✔"))
        } else {
            j.WriteString(nanairo.FgColor("#c00", "✔"))
        }
        j.WriteString("  ")

        // W
        if job.HealthReport[0].Score >= 80 {
            j.WriteString(nanairo.FgColor("#f90", "☀"))
        } else if job.HealthReport[0].Score >= 20 {
            j.WriteString(nanairo.FgColor("#999", "☁"))
        } else {
            j.WriteString(nanairo.FgColor("#03c", "☂"))
        }

        // Name
        j.WriteString("  " + job.Name)

        fmt.Println(j.String())
    }
}

var Jobs = cli.Command {
    Name: "jobs",
    Usage: "print status for all jobs",
    Action: jobs,
}