package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func conf(c *cli.Context) {
	if c.Bool("dump") {
		dumpFormat := c.Args().First()
		if dumpFormat == "" {
			dumpFormat = "all"
		}
		Config.Dump(defaultConfigFile, dumpFormat)
		return
	}

	url := c.Args().First()
	if url == "" && c.Bool("rm") {
		fmt.Println("assign target name")
		return
	}

	if !c.Bool("add") && !c.Bool("dump") && !c.Bool("rm") {
		fmt.Println("please --add or --rm option")
		return
	}

	for i, _ := range Config.HostInfo {
		if Config.HostInfo[i].Name == c.String("name") {
			if c.String("rm") != "" {
				// FIXME
				Config.HostInfo = append(Config.HostInfo[:i+1], Config.HostInfo[i+1:]...)
				fmt.Println(len(Config.HostInfo))
			} else {
				Config.HostInfo[i] = JcConfigHostInfo{Name: c.String("name"), Hostname: c.Args().First()}
				fmt.Println("url:", c.Args().First())
			}
			Config.Save(defaultConfigFile)
			return
		}
	}

	Config.HostInfo = append(Config.HostInfo, JcConfigHostInfo{Name: c.String("name"), Hostname: c.Args().First()})
	Config.Save(defaultConfigFile)
}

var ConfCommand = cli.Command{
	Name:   "conf",
	Usage:  "config jc command setting param",
	Action: conf,
	Flags: []cli.Flag{
		cli.StringFlag{
			"name, n",
			"default",
			"host key name(default is 'default')",
			"",
		},
		cli.BoolFlag{
			"add, a",
			"add new host key name",
			"",
		},
		cli.BoolFlag{
			"rm",
			"remove host key name",
			"",
		},
		cli.BoolFlag{
			"dump, d",
			"print configuration (all, list). default is 'all'",
			"",
		},
	},
}
