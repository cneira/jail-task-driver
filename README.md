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


Parameters
-----------
Parameters used by the driver support most of JAIL(8) functionality, parameter names closely match the ones in 
JAIL(8)  
.
[Parameters documentation ](https://github.com/cneira/jail-task-driver/blob/master/Parameters.md)  

Examples 
---------

Basic jail 

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
	Ip4_addr = "192.168.1.102"
      }
    }
  }
}
```
Vnet jail example 

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
	Exec_prestart = "/usr/share/examples/jails/jib addm loghost jailether"
	Exec_poststop = "/usr/share/examples/jails/jib destroy loghost "
      }
    }
  }
}
```

## Support

It's also possible to support the project on [Patreon](https://www.patreon.com/neirac)


## References

- Lucas, Michael W. FreeBSD Mastery: Jails (IT Mastery Book 15). 
- [FreeBSD HandBook](https://www.freebsd.org/doc/en_US.ISO8859-1/books/handbook/)

 TODO:
-------

* Implement exec interface
* Test All jail options
* Refactor to match parameters as closely as JAIL(8)
* Create jails using docker images
