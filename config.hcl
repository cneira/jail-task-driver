plugin "nomad-zone-driver" {
  config {}
}

client {
  options = {
    "driver.whitelist" = "jail-task-driver"
  }
#     "cpu_total_compute" = 2813
}

#acl {
#enabled    = true
#token_ttl  = "30s"
#policy_ttl = "60s"
#}

#vault {
#enabled          = true
#address          = "http://127.0.0.1:8200"
#task_token_ttl   = "1h"
#create_from_role = "nomad-cluster"
#token            = "s.lhGVgrtjcmTKc7XttpL5MZkp"
#}

consul {
  address = "0.0.0.0:8500"
}
