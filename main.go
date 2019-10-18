package main

import (
	"github.com/cneira/jail-task-driver/driver"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/plugins"
)

func main() {
	// Serve the plugin
	plugins.Serve(factory)
}

func factory(log log.Logger) interface{} {
	return jail.NewJailDriver(log)
}
