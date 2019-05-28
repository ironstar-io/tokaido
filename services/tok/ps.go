package tok

import (
	"fmt"
	"strconv"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/services/docker"
	. "github.com/logrusorgru/aurora"
	"github.com/ryanuber/columnize"
)

// Ps will return a list of all containers in the environment and their state
func Ps() {
	c := conf.GetConfig()
	pn := c.Tokaido.Project.Name
	failure := false

	fmt.Println()
	fmt.Println(Cyan(Sprintf("Your main Tokaido HTTPS entrypoint is:  %s", Bold("https://"+pn+".local.tokaido.io:5154/"))))
	fmt.Println(Cyan(Sprintf("You can open the entrypoint by running: %s", Bold("tok open"))))
	fmt.Println(Cyan(Sprintf("You can SSH in by running the command:  %s", Bold("ssh "+pn+".tok"))))

	o := []string{}
	o = append(o, "Container|Local Endpoint|Shortcut|Status")
	o = append(o, "--------|--------|--------|--------")
	// Output Unison status if relevant
	if c.Global.Syncservice == "unison" {
		unison := docker.GetContainer("unison", pn)
		if unison.State == "running" {
			o = append(o, Sprintf("unison(sync)|-|-|%s", Green(unison.State)))
		} else {
			o = append(o, Sprintf("unison(sync)|-|-|%s", Yellow("offline")))
			failure = true
		}
	}

	if conf.GetConfig().Services.Chromedriver.Enabled {
		chromedriver := docker.GetContainer("chromedriver", pn)
		if chromedriver.State == "running" {
			o = append(o, Sprintf("chromedriver|-|-|%s", Green(chromedriver.State)))
		} else {
			o = append(o, Sprintf("chromedriver|-|-|%s", Yellow("offline")))
			failure = true
		}
	}

	// Output Drush status
	admin := docker.GetContainer("drush", pn)
	if admin.State == "running" {
		port := strconv.Itoa(int(admin.Ports[0].PublicPort))
		o = append(o, Sprintf("admin(drush/ssh)|ssh://localhost:%s|ssh %s.tok|%s", port, pn, Green(admin.State)))
	} else {
		o = append(o, Sprintf("admin(drush/ssh)|-|-|%s", Yellow("offline")))
		failure = true
	}

	// Output Haproxy status
	haproxy := docker.GetContainer("haproxy", pn)
	if haproxy.State == "running" {
		port := 0
		for _, v := range haproxy.Ports {
			if v.PrivatePort == uint16(constants.HaproxyInternalPort) {
				port = int(v.PublicPort)
			}
		}
		o = append(o, Sprintf("haproxy|https://localhost:%d|tok open haproxy|%s", port, Green(haproxy.State)))
	} else {
		o = append(o, Sprintf("haproxy|-|-|%s", Yellow("offline")))
		failure = true

	}

	// Output Varnish status
	varnish := docker.GetContainer("varnish", pn)
	if varnish.State == "running" {
		port := 0
		for _, v := range varnish.Ports {
			if v.PrivatePort == uint16(constants.VarnishInternalPort) {
				port = int(v.PublicPort)
			}
		}
		o = append(o, Sprintf("varnish|http://localhost:%d|tok open varnish|%s", port, Green(varnish.State)))
	} else {
		o = append(o, Sprintf("varnish|-|-|%s", Yellow("offline")))
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
		o = append(o, Sprintf("nginx|http://localhost:%d|tok open nginx|%s", port, Green(nginx.State)))
	} else {
		o = append(o, Sprintf("nginx|-|-|%s", Yellow("offline")))
		failure = true
	}

	// Output FPM status
	fpm := docker.GetContainer("fpm", pn)
	if fpm.State == "running" {
		o = append(o, Sprintf("fpm|-|-|%s", Green(fpm.State)))
	} else {
		o = append(o, Sprintf("fpm|-|-|%s", Yellow("offline")))
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
		o = append(o, Sprintf("mysql|mysql://localhost:%d|%s|%s", port, adminerMsg, Green(mysql.State)))
	} else {
		o = append(o, Sprintf("mysql|-|-|%s", Yellow("offline")))
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
			o = append(o, Sprintf("mailhog|mailhog://localhost:%d|tok open mailhog|%s", port, Green(mailhog.State)))
		} else {
			o = append(o, Sprintf("mailhog|-|-|%s", Yellow("offline")))
			failure = true
		}
	}

	// Output Adminer Status
	if c.Services.Adminer.Enabled {
		adminer := docker.GetContainer("adminer", pn)
		if adminer.State == "running" {
			port := strconv.Itoa(int(adminer.Ports[0].PublicPort))
			o = append(o, Sprintf("adminer|adminer://localhost:%s|tok open adminer|%s", port, Green(adminer.State)))
		} else {
			o = append(o, Sprintf("adminer|-|-|%s", Yellow("offline")))
			failure = true
		}
	}

	// Output Redis Status
	if c.Services.Redis.Enabled {
		redis := docker.GetContainer("redis", pn)
		if redis.State == "running" {
			port := strconv.Itoa(int(redis.Ports[0].PublicPort))
			o = append(o, Sprintf("redis|redis://localhost:%s|-|%s", port, Green(redis.State)))
		} else {
			o = append(o, Sprintf("redis|-|-|%s", Yellow("offline")))
			failure = true
		}
	}

	// Output Solr Status
	if c.Services.Solr.Enabled {
		solr := docker.GetContainer("solr", pn)
		if solr.State == "running" {
			port := strconv.Itoa(int(solr.Ports[0].PublicPort))
			o = append(o, Sprintf("solr|http://localhost:%s|tok open solr|%s", port, Green(solr.State)))
		} else {
			o = append(o, Sprintf("solr|-|-|%s", Yellow("offline")))
			failure = true
		}
	}

	// Output Memcache Status
	if c.Services.Memcache.Enabled {
		memcache := docker.GetContainer("memcache", pn)
		if memcache.State == "running" {
			o = append(o, Sprintf("memcache|-|-|%s", Green(memcache.State)))
		} else {
			o = append(o, Sprintf("memcache|-|-|%s", Yellow("offline")))
			failure = true
		}
	}

	// Output Syslog Status
	syslog := docker.GetContainer("syslog", pn)
	if syslog.State == "running" {
		o = append(o, Sprintf("logging|-|-|%s", Green(syslog.State)))
	} else {
		o = append(o, Sprintf("logging|-|-|%s", Yellow("offline")))
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
		fmt.Println(Red("It looks like one of your Tokaido containers is offline"))
		fmt.Println(Sprintf("You can try to fix this by running '%s' again, or you can use", Blue("tok up")))
		fmt.Println(Sprintf("'%s %s' to see the docker logs for that container", Blue("tok logs"), BrightBlue("{container name}")))
	}

	fmt.Println()

}
