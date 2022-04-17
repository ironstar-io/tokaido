package tok

import (
	"fmt"
	"strconv"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/logrusorgru/aurora"
	"github.com/ryanuber/columnize"
)

// Ps will return a list of all containers in the environment and their state
func Ps() {
	c := conf.GetConfig()
	pn := c.Tokaido.Project.Name
	failure := false

	fmt.Println()
	fmt.Println(aurora.Cyan(aurora.Sprintf("Your main Tokaido HTTPS entrypoint is:  %s", aurora.Bold("https://"+pn+".local.tokaido.io:5154/"))))
	fmt.Println(aurora.Cyan(aurora.Sprintf("You can open the entrypoint by running: %s", aurora.Bold("tok open"))))
	fmt.Println(aurora.Cyan(aurora.Sprintf("You can SSH in by running the command:  %s", aurora.Bold("ssh "+pn+".tok"))))

	o := []string{}
	o = append(o, "Container|Local Endpoint|Shortcut|Status")
	o = append(o, "--------|--------|--------|--------")

	// Output Drush status
	admin := docker.GetContainer("drush", pn)
	if admin.State == "running" {
		port := strconv.Itoa(int(admin.Ports[0].PublicPort))
		o = append(o, aurora.Sprintf("admin(drush/ssh)|ssh://localhost:%s|ssh %s.tok|%s", port, pn, aurora.Green(admin.State)))
	} else {
		o = append(o, aurora.Sprintf("admin(drush/ssh)|-|-|%s", aurora.Yellow("offline")))
		failure = true
	}

	// Output Nginx status
	nginx := docker.GetContainer("nginx", pn)
	if nginx.State == "running" {
		port := 0
		for _, v := range nginx.Ports {
			if v.PrivatePort == uint16(constants.NginxInternalPort) {
				port = int(v.PublicPort)
			}
		}
		o = append(o, aurora.Sprintf("nginx|http://localhost:%d|tok open nginx|%s", port, aurora.Green(nginx.State)))
	} else {
		o = append(o, aurora.Sprintf("nginx|-|-|%s", aurora.Yellow("offline")))
		failure = true
	}

	// Output FPM status
	fpm := docker.GetContainer("fpm", pn)
	if fpm.State == "running" {
		o = append(o, aurora.Sprintf("fpm|-|-|%s", aurora.Green(fpm.State)))
	} else {
		o = append(o, aurora.Sprintf("fpm|-|-|%s", aurora.Yellow("offline")))
		failure = true
	}

	// Output MySQL status
	mysql := docker.GetContainer("mysql", pn)
	if mysql.State == "running" {
		adminerMsg := "-"
		if c.Services.Adminer.Enabled {
			adminerMsg = "tok open adminer"
		}

		port := 0
		for _, v := range mysql.Ports {
			if v.PrivatePort == uint16(constants.MysqlInternalPort) {
				port = int(v.PublicPort)
			}
		}
		o = append(o, aurora.Sprintf("mysql|mysql://localhost:%d|%s|%s", port, adminerMsg, aurora.Green(mysql.State)))
	} else {
		o = append(o, aurora.Sprintf("mysql|-|-|%s", aurora.Yellow("offline")))
		failure = true
	}

	// Output Mailhog Status
	if c.Services.Mailhog.Enabled {
		mailhog := docker.GetContainer("mailhog", pn)
		if mailhog.State == "running" {
			port := 0
			for _, v := range mailhog.Ports {
				if v.PrivatePort == uint16(constants.MailhogInternalHTTPPort) {
					port = int(v.PublicPort)
				}
			}
			o = append(o, aurora.Sprintf("mailhog|mailhog://localhost:%d|tok open mailhog|%s", port, aurora.Green(mailhog.State)))
		} else {
			o = append(o, aurora.Sprintf("mailhog|-|-|%s", aurora.Yellow("offline")))
			failure = true
		}
	}

	// Output Adminer Status
	if c.Services.Adminer.Enabled {
		adminer := docker.GetContainer("adminer", pn)
		if adminer.State == "running" {
			port := strconv.Itoa(int(adminer.Ports[0].PublicPort))
			o = append(o, aurora.Sprintf("adminer|adminer://localhost:%s|tok open adminer|%s", port, aurora.Green(adminer.State)))
		} else {
			o = append(o, aurora.Sprintf("adminer|-|-|%s", aurora.Yellow("offline")))
			failure = true
		}
	}

	// Output Redis Status
	if c.Services.Redis.Enabled {
		redis := docker.GetContainer("redis", pn)
		if redis.State == "running" {
			port := strconv.Itoa(int(redis.Ports[0].PublicPort))
			o = append(o, aurora.Sprintf("redis|redis://localhost:%s|-|%s", port, aurora.Green(redis.State)))
		} else {
			o = append(o, aurora.Sprintf("redis|-|-|%s", aurora.Yellow("offline")))
			failure = true
		}
	}

	// Output Solr Status
	if c.Services.Solr.Enabled {
		solr := docker.GetContainer("solr", pn)
		if solr.State == "running" {
			port := strconv.Itoa(int(solr.Ports[0].PublicPort))
			o = append(o, aurora.Sprintf("solr|http://localhost:%s|tok open solr|%s", port, aurora.Green(solr.State)))
		} else {
			o = append(o, aurora.Sprintf("solr|-|-|%s", aurora.Yellow("offline")))
			failure = true
		}
	}

	// Output Memcache Status
	if c.Services.Memcache.Enabled {
		memcache := docker.GetContainer("memcache", pn)
		if memcache.State == "running" {
			o = append(o, aurora.Sprintf("memcache|-|-|%s", aurora.Green(memcache.State)))
		} else {
			o = append(o, aurora.Sprintf("memcache|-|-|%s", aurora.Yellow("offline")))
			failure = true
		}
	}

	// Output Syslog Status
	syslog := docker.GetContainer("syslog", pn)
	if syslog.State == "running" {
		o = append(o, aurora.Sprintf("logging|-|-|%s", aurora.Green(syslog.State)))
	} else {
		o = append(o, aurora.Sprintf("logging|-|-|%s", aurora.Yellow("offline")))
		failure = true
	}

	fmt.Println()

	cc := columnize.DefaultConfig()
	cc.Delim = "|"
	cc.Glue = "  "
	cc.Prefix = ""
	cc.Empty = ""
	cc.NoTrim = false

	result := columnize.Format(o, cc)
	fmt.Println(result)

	if failure {
		fmt.Println()
		fmt.Println(aurora.Red("It looks like one of your Tokaido containers is offline"))
		fmt.Println(aurora.Sprintf("You can try to fix this by running '%s' again, or you can use", aurora.Blue("tok up")))
		fmt.Println(aurora.Sprintf("'%s %s' to see the docker logs for that container", aurora.Blue("tok logs"), aurora.BrightBlue("{container name}")))
	}

	fmt.Println()

}
