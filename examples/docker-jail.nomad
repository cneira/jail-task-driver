job "test3" {
  datacenters = ["dc1"]
  type        = "service"

  group "test" {

    task "test01" {
      driver = "jail-task-driver"
      config {
	Docker = "gitea/gitea latest"
        Path   = "/home/cneira/dockerjail"
      }
    }
  }
}
