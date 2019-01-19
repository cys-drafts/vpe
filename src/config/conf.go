package config

import (
	"fmt"
	"os"
	"encoding/json"
)

const (
	ROLE_CLIENT = 0
	ROLE_SERVER = 1
)

type RawConf struct {
	Proto string `json:"proto"`
	Mtu int `json:"mtu"`
	Ifname string `json:"ifname"`
	Ifscript string `json:"ifup"`
	Role string `json:"role"`
	Server string `json:"server"`
}

type Conf struct {
	Proto string
	Mtu int
	Ifname string
	Ifscript string
	Role int
	Server string
}

func (c *Conf) String() string {
	return fmt.Sprintf("role(%d): %s : %d : %s : %s : %s", c.Role, c.Proto, c.Mtu, c.Ifname, c.Ifscript, c.Server)
}

func parsefile(confname string) (RawConf, error) {
	var rc RawConf
	confFile, e := os.Open(confname)
	if e != nil {
		return rc, e
	}
	jsonParser := json.NewDecoder(confFile)
	e = jsonParser.Decode(&rc)
	if e != nil {
		return rc, e
	}
	return rc, nil
}

//~/.config/vpe/ /etc/vpe/
func Parse() (*Conf, error) {
	var rc RawConf
	var e error

	confdir := fmt.Sprintf("%s/.config/vpe", os.Getenv("HOME"))

	if _, e := os.Stat(confdir); os.IsNotExist(e) {
		fmt.Printf("no config dir\n")
		return nil, e
	}

	confname := fmt.Sprintf("%s/vpe.json", confdir)
	if _, e := os.Stat(confdir); os.IsNotExist(e) {
		fmt.Printf("no config file\n");
		return nil, e
	}

	if rc, e = parsefile(confname); e != nil {
		return nil, e
	}

	c := Conf {
		Proto: rc.Proto,
		Mtu: rc.Mtu,
		Ifname: rc.Ifname,
		Ifscript: rc.Ifscript,
		Server: rc.Server,
	}

	if rc.Role == "client" {
		c.Role = ROLE_CLIENT
	} else {
		c.Role = ROLE_SERVER
	}
	if c.Ifscript == "" {
		c.Ifscript = fmt.Sprintf("%s/if-up", confdir)
	}
	if _, e := os.Stat(c.Ifscript); os.IsNotExist(e) {
		fmt.Printf("no ifup script\n");
		return nil, e
	}

	return &c, nil
}
