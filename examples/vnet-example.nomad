job "vnet-example2" {
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
