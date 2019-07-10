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
