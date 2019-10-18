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
      }
    }
  }
}
