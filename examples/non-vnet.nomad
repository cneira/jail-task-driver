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
