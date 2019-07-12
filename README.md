FreeBSD Jail Task Driver
===========================

Task driver for [FreeBSD](https://www.freebsd.org/) jails. 


- Website: https://www.nomadproject.io

Requirements
------------

- [Nomad](https://www.nomadproject.io/downloads.html) 0.9+
- [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)
- [FreeBSD 12.0-RELEASE](https://www.freebsd.org/where.html) *Should work with 11*
- [Consul](https://releases.hashicorp.com/consul/1.5.2/consul_1.5.2_freebsd_amd64.zip)

Installation
------------

Install(and compile) the jail-task-driver binary and put it in [plugin_dir](https://www.nomadproject.io/docs/configuration/index.html#plugin_dir) and then add a `plugin "jail-task-driver" {}` line in your nomad config file.


```shell
go get github.com/cneira/jail-task-driver
cp $GOPATH/bin/jail-task-driver YOURPLUGINDIR
```

Then in your nomad config file, set
```hcl
plugin "jail-task-driver" {}
```

In developer/test mode(`nomad agent -dev`) , plugin_dir is unset it seems, so you will need to mkdir plugins and then copy the jail-task-driver binary to plugins and add a `plugins_dir = "path/to/plugins"` to the above config file.
then you can run it like:

`nomad agent -dev -config nomad.config`

For more details see the nomad [docs](https://www.nomadproject.io/docs/configuration/plugin.html).

Parameters
-----------
Parameters used by the driver support most of JAIL(8) functionality, parameter names 
closely match the ones in JAIL(8).  
   

[Parameters documentation ](https://github.com/cneira/jail-task-driver/blob/master/Parameters.md)  

Examples 
---------

Basic jail 
---------

```hcl
job "test" {
  datacenters = ["dc1"]
  type        = "service"

  group "test" {
    restart {
      attempts = 0
      mode     = "fail"
    }

    task "test01" {
      driver = "jail-task-driver"

      config {
        Path    = "/zroot/iocage/jails/myjail/root"
	Persist  = true
      }
    }
  }
}
```
Non vnet jail
-------------
```hcl 
job "non-vnet" {
  datacenters = ["dc1"]
  type        = "service"

  group "test" {
    restart {
      attempts = 0
      mode     = "fail"
    }

    task "test01" {
      driver = "jail-task-driver"

      config {
        Path              = "/zroot/iocage/jails/myjail/root"
        Ip4               = "new"
        Allow_raw_sockets = true
        Allow_chflags     = true
        Ip4_addr          = "em1|192.168.1.102"
        Exec_start        = "/usr/local/bin/http-echo -listen :9999 -text hello"
      }
    }
  }
}
```

Vnet jail example 
-----------------

```hcl
job "vnet-example" {
  datacenters = ["dc1"]
  type        = "service"

  group "test" {
    restart {
      attempts = 0
      mode     = "fail"
    }

    task "test01" {
      driver = "jail-task-driver"

      config {
        Path    = "/zroot/iocage/jails/myjail/root"
	Persist  = true
 	Host_hostname = "mwl.io"
	Exec_clean = true	
	Exec_start = "sh /etc/rc"
	Exec_stop = "sh /etc/rc.shutdown"
	Mount_devfs = true
	Exec_prestart = "logger trying to start "	
	Exec_poststart = "logger jail has started"	
	Exec_prestop = "logger shutting down jail "	
	Exec_poststop = "logger has shut down jail "	
	Exec_consolelog ="/var/tmp/vnet-example"
	Vnet = true
	Vnet_nic = "e0b_loghost"
	Exec_prestart = "/usr/share/examples/jails/jib addm loghost em1"
	Exec_poststop = "/usr/share/examples/jails/jib destroy loghost "
      }
    }
  }
}
```
Setting resource limits
----------------------
```hcl

job "rctl-test" {
  datacenters = ["dc1"]
  type        = "service"

  group "test" {
    restart {
      attempts = 0
      mode     = "fail"
    }

    task "test01" {
      driver = "jail-task-driver"

      config {
        Path    = "/zroot/iocage/jails/myjail/root"
	Persist  = true
	Ip4_addr = "192.168.1.102"
	Rctl =  {
		Vmemoryuse = 1200000
	}
    }
  }
}
}
```
##  Demo
[![asciicast](https://asciinema.org/a/256519.svg)](https://asciinema.org/a/256519)

## Support

[![ko-fi](https://www.ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/J3J4YM9U)  
      
It's also possible to support the project on [Patreon](https://www.patreon.com/neirac)  
    
## References

- [Lucas, Michael W. FreeBSD Mastery: Jails (IT Mastery Book 15)](https://mwl.io/nonfiction/os#fmjail)
- [FreeBSD HandBook](https://www.freebsd.org/doc/en_US.ISO8859-1/books/handbook/)
- [RCTL(8)](https://www.freebsd.org/cgi/man.cgi?query=rctl&sektion=8)


 TODO:
-------

* ~~Implement exec interface~~
* ~~Implement RecoverTask interface~~
* Test All jail options
* Refactor to match parameters as closely as JAIL(8)
* Create jails using docker images

