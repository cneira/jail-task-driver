/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2019, Carlos Neira cneirabustos@gmail.com
 */

package jail

import (
	"context"
	"fmt"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/drivers/shared/eventer"
	"github.com/hashicorp/nomad/plugins/base"
	"github.com/hashicorp/nomad/plugins/drivers"
	"github.com/hashicorp/nomad/plugins/shared/hclspec"
	pstructs "github.com/hashicorp/nomad/plugins/shared/structs"
)

const (
	// pluginName is the name of the plugin
	pluginName = "jail-task-driver"

	// fingerprintPeriod is the interval at which the driver will send fingerprint responses
	fingerprintPeriod = 30 * time.Second

	// taskHandleVersion is the version of task handle which this driver sets
	// and understands how to decode driver state
	taskHandleVersion = 1
)

var (
	// pluginInfo is the response returned for the PluginInfo RPC
	pluginInfo = &base.PluginInfoResponse{
		Type:              base.PluginTypeDriver,
		PluginApiVersions: []string{drivers.ApiVersion010},
		PluginVersion:     "0.1.1-dev",
		Name:              pluginName,
	}

	// taskConfigSpec is the hcl specification for the driver config section of
	// a task within a job. It is returned in the TaskConfigSchema RPC
	taskConfigSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		"Path":                  hclspec.NewAttr("Path", "string", false),
		"Docker":                hclspec.NewAttr("Docker", "string", false),
		"Jid":                   hclspec.NewAttr("Jid", "string", false),
		"Ip4_addr":              hclspec.NewAttr("Ip4_addr", "string", false),
		"Ip4_saddrsel":          hclspec.NewAttr("Ip4_saddrsel", "bool", false),
		"Ip4":                   hclspec.NewAttr("Ip4", "string", false),
		"Ip6_addr":              hclspec.NewAttr("Ip6_addr", "string", false),
		"Ip6_saddrsel":          hclspec.NewAttr("Ip6_saddrsel", "bool", false),
		"Ip6":                   hclspec.NewAttr("Ip6", "string", false),
		"Vnet":                  hclspec.NewAttr("Vnet", "string", false),
		"Host_hostname":         hclspec.NewAttr("Host_hostname", "string", false),
		"Host":                  hclspec.NewAttr("Host", "string", false),
		"Securelevel":           hclspec.NewAttr("Securelevel", "string", false),
		"Devfs_ruleset":         hclspec.NewAttr("Devfs_ruleset", "string", false),
		"Children_max":          hclspec.NewAttr("Children_max", "number", false),
		"Children_cur":          hclspec.NewAttr("Children_cur", "number", false),
		"Enforce_statfs":        hclspec.NewAttr("Enforce_statfs", "number", false),
		"Persist":               hclspec.NewAttr("Persist", "bool", false),
		"Osrelease":             hclspec.NewAttr("Osrelease", "string", false),
		"Osreldate":             hclspec.NewAttr("Osreldate", "string", false),
		"Allow_set_hostname":    hclspec.NewAttr("Allow_set_hostname", "bool", false),
		"Allow_sysvipc":         hclspec.NewAttr("Allow_sysvipc", "bool", false),
		"Allow_raw_sockets":     hclspec.NewAttr("Allow_raw_sockets", "bool", false),
		"Allow_chflags":         hclspec.NewAttr("Allow_chflags", "bool", false),
		"Allow_mount":           hclspec.NewAttr("Allow_mount", "bool", false),
		"Allow_mount_devfs":     hclspec.NewAttr("Allow_mount_devfs", "bool", false),
		"Allow_quotas":          hclspec.NewAttr("Allow_quotas", "bool", false),
		"Allow_read_msgbuf":     hclspec.NewAttr("Allow_read_msgbug", "bool", false),
		"Allow_socket_af":       hclspec.NewAttr("Allow_socket_af", "bool", false),
		"Allow_reserved_ports":  hclspec.NewAttr("Allow_reserved_ports", "bool", false),
		"Allow_mlock":           hclspec.NewAttr("Allow_mlock", "bool", false),
		"Allow_mount_fdescfs":   hclspec.NewAttr("Allow_mount_fdescds", "bool", false),
		"Allow_mount_fusefs":    hclspec.NewAttr("Allow_mount_fusefs", "bool", false),
		"Allow_mount_nullfs":    hclspec.NewAttr("Allow_mount_nullfs", "bool", false),
		"Allow_mount_procfs":    hclspec.NewAttr("Allow_mount_procfs", "bool", false),
		"Allow_mount_linprocfs": hclspec.NewAttr("Allow_mount_linprocfs", "bool", false),
		"Allow_mount_linsysfs":  hclspec.NewAttr("Allow_mount_linsysfs", "bool", false),
		"Allow_mount_tmpfs":     hclspec.NewAttr("Allow_mount_tmpfs", "bool", false),
		"Allow_mount_zfs":       hclspec.NewAttr("Allow_mount_zfs", "bool", false),
		"Allow_vmm":             hclspec.NewAttr("Allow_vmm", "bool", false),
		"Linux":                 hclspec.NewAttr("Linux", "string", false),
		"Linux_osname":          hclspec.NewAttr("Linux_osname", "string", false),
		"Linux_osrelease":       hclspec.NewAttr("Linux_osrelease", "string", false),
		"Linux_oss_version":     hclspec.NewAttr("Linux_oss_version", "string", false),
		"Sysvmsg":               hclspec.NewAttr("Sysvmsg", "string", false),
		"Sysvsem":               hclspec.NewAttr("Sysvsem", "string", false),
		"Sysvshm":               hclspec.NewAttr("Sysvshm", "string", false),
		"Exec_prestart":         hclspec.NewAttr("Exec_prestart", "string", false),
		"Exec_prestop":          hclspec.NewAttr("Exec_prestop", "string", false),
		"Exec_created":          hclspec.NewAttr("Exec_created", "string", false),
		"Exec_start":            hclspec.NewAttr("Exec_start", "string", false),
		"Exec_stop":             hclspec.NewAttr("Exec_stop", "string", false),
		"Exec_poststart":        hclspec.NewAttr("Exec_poststart", "string", false),
		"Exec_poststop":         hclspec.NewAttr("Exec_poststop", "string", false),
		"Exec_clean":            hclspec.NewAttr("Exec_clean", "bool", false),
		"Exec_jail_user":        hclspec.NewAttr("Exec_jail_user", "string", false),
		"Exec_system_jail_user": hclspec.NewAttr("Exec_system_jail_user", "string", false),
		"Exec_system_user":      hclspec.NewAttr("Exec_system_user", "string", false),
		"Exec_timeout":          hclspec.NewAttr("Exec_timeout", "number", false),
		"Exec_consolelog":       hclspec.NewAttr("Exec_consolelog", "string", false),
		"Exec_fib":              hclspec.NewAttr("Exec_fib", "string", false),
		"Stop_timeout":          hclspec.NewAttr("Stop_timeout", "number", false),
		"Nic":                   hclspec.NewAttr("Nic", "string", false),
		"Vnet_nic":              hclspec.NewAttr("Vnet_nic", "string", false),
		"Ip_hostname":           hclspec.NewAttr("Ip_hostname", "string", false),
		"Mount":                 hclspec.NewAttr("Mount", "string", false),
		"Mount_fstab":           hclspec.NewAttr("Mount_fstab", "string", false),
		"Mount_devfs":           hclspec.NewAttr("Mount_devfs", "bool", false),
		"Mount_fdescfs":         hclspec.NewAttr("Mount_fdescfs", "bool", false),
		"Depend":                hclspec.NewAttr("Depend", "string", false),
		"Rctl": hclspec.NewBlock("Rctl", false, hclspec.NewObject(map[string]*hclspec.Spec{
			"Cputime": hclspec.NewBlock("Cputime", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),
			"Datasize": hclspec.NewBlock("Datasize", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Stacksize": hclspec.NewBlock("Stacksize", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Coredumpsize": hclspec.NewBlock("Coredumpsize", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Memoryuse": hclspec.NewBlock("Memoryuse", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Memorylocked": hclspec.NewBlock("Memorylocked", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Maxproc": hclspec.NewBlock("Maxproc", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Openfiles": hclspec.NewBlock("Openfiles", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Vmemoryuse": hclspec.NewBlock("Vmemoryuse", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Pseudoterminals": hclspec.NewBlock("Pseudoterminals", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Swapuse": hclspec.NewBlock("Swapuse", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Nthr": hclspec.NewBlock("Nthr", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Msgqqueued": hclspec.NewBlock("Msgqqueued", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Msgqsize": hclspec.NewBlock("Msgqsize", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Nmsgq": hclspec.NewBlock("Nmsgq", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Nsem": hclspec.NewBlock("Nsem", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Nsemop": hclspec.NewBlock("Nsemop", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Nshm": hclspec.NewBlock("Nshm", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Shmsize": hclspec.NewBlock("Shmsize", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Wallclock": hclspec.NewBlock("Wallclock", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Pcpu": hclspec.NewBlock("Pcpu", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Readbps": hclspec.NewBlock("Readbps", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Writebps": hclspec.NewBlock("Writebps", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Readiops": hclspec.NewBlock("Readiops", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),

			"Writeiops": hclspec.NewBlock("Writeiops", false, hclspec.NewObject(map[string]*hclspec.Spec{
				"Action": hclspec.NewAttr("Action", "string", true),
				"Amount": hclspec.NewAttr("Amount", "string", true),
				"Per":    hclspec.NewAttr("Per", "string", false),
			})),
		})),
	})

	// capabilities is returned by the Capabilities RPC and indicates what
	// optional features this driver supports
	capabilities = &drivers.Capabilities{
		SendSignals: false,
		Exec:        true,
		FSIsolation: drivers.FSIsolationImage,
	}
)

type Driver struct {
	// eventer is used to handle multiplexing of TaskEvents calls such that an
	// event can be broadcast to all callers
	eventer *eventer.Eventer

	// config is the driver configuration set by the SetConfig RPC
	config *Config

	// nomadConfig is the client config from nomad
	nomadConfig *base.ClientDriverConfig

	// tasks is the in memory datastore mapping taskIDs to rawExecDriverHandles
	tasks *taskStore

	// ctx is the context for the driver. It is passed to other subsystems to
	// coordinate shutdown
	ctx context.Context

	// signalShutdown is called when the driver is shutting down and cancels the
	// ctx passed to any subsystems
	signalShutdown context.CancelFunc

	// logger will log to the Nomad agent
	logger hclog.Logger
}

// Config is the driver configuration set by the SetConfig RPC call
type Config struct {
}

type RctlOpts struct {
	Action string `codec:"Action"`
	Amount string `codec:"Amount"`
	Per    string `codec:"Per"`
}

type Rctl struct {
	Cputime         RctlOpts `codec:"Cputime"`
	Datasize        RctlOpts `codec:"Datasize"`
	Coredumpsize    RctlOpts `codec:"Coredumpsize"`
	Stacksize       RctlOpts `codec:"Stacksize"`
	Memoryuse       RctlOpts `codec:"Memoryuse"`
	Memorylocked    RctlOpts `codec:"Memorylocked"`
	Maxproc         RctlOpts `codec:"Maxproc"`
	Openfiles       RctlOpts `codec:"Openfiles"`
	Vmemoryuse      RctlOpts `codec:"Vmemoryuse"`
	Pseudoterminals RctlOpts `codec:"Pseudoterminals"`
	Swapuse         RctlOpts `codec:"Swapuse"`
	Nthr            RctlOpts `codec:"Nthr"`
	Msgqqueued      RctlOpts `codec:"Msgqqueued"`
	Msgqsize        RctlOpts `codec:"Msgqsize"`
	Nmsgq           RctlOpts `codec:"Nmsgq"`
	Nsemop          RctlOpts `codec:"Nsemop"`
	Nshm            RctlOpts `codec:"Nshm"`
	Shmsize         RctlOpts `codec:"Shmsize"`
	Wallclock       RctlOpts `codec:"Wallclock"`
	Pcpu            RctlOpts `codec:"Pcpu"`
	Readbps         RctlOpts `codec:"Readbps"`
	Writebps        RctlOpts `codec:"Writebps"`
	Readiops        RctlOpts `codec:"Readiops"`
	Writeiops       RctlOpts `codec:"Writeiops"`
}

// TaskConfig is the driver configuration of a task within a job
type TaskConfig struct {
	Path                  string `codec:"Path"`
	Docker                string `codec:"Docker"`
	Jid                   string `codec:"Jid"`
	Ip4_addr              string `codec:"Ip4_addr"`
	Ip4_saddrsel          bool   `codec:"Ip4_saddrsel"`
	Ip4                   string `codec:"Ip4"`
	Ip6_addr              string `codec:"Ip6_addr"`
	Ip6_saddrsel          bool   `codec:"Ip6_saddrsel"`
	Ip6                   string `codec:"Ip6"`
	Vnet                  string `codec:"Vnet"`
	Host_hostname         string `codec:"Host_hostname"`
	Host                  string `codec:"Host"`
	Securelevel           string `codec:"Securelevel"`
	Devfs_ruleset         string `codec:"Devfs_ruleset"`
	Children_max          uint   `codec:"Children_max"`
	Children_cur          uint   `codec:"Children_cur"`
	Enforce_statfs        uint   `codec:"Enforce_statfs"`
	Persist               bool   `codec:"Persist"`
	Osrelease             string `codec:"Osrelease"`
	Osreldate             string `codec:"Osreldate"`
	Allow_set_hostname    bool   `codec:"Allow_set_hostname"`
	Allow_sysvipc         bool   `codec:"Allow_sysvipc"`
	Allow_raw_sockets     bool   `codec:"Allow_raw_sockets"`
	Allow_chflags         bool   `codec:"Allow_chflags"`
	Allow_mount           bool   `codec:"Allow_mount"`
	Allow_mount_devfs     bool   `codec:"Allow_mount.devfs"`
	Allow_quotas          bool   `codec:"Allow_quotas"`
	Allow_read_msgbuf     bool   `codec:"Allow_read_msgbuf"`
	Allow_socket_af       bool   `codec:"Allow_socket_af"`
	Allow_reserved_ports  bool   `codec:"Allow_reserved_ports"`
	Allow_mlock           bool   `codec:"Allow_mlock"`
	Allow_mount_fdescfs   bool   `codec:"Allow_mount_fdescfs"`
	Allow_mount_fusefs    bool   `codec:"Allow_mount_fusefs"`
	Allow_mount_nullfs    bool   `codec:"Allow_mount_nullfs"`
	Allow_mount_procfs    bool   `codec:"Allow_mount_procfs"`
	Allow_mount_linprocfs bool   `codec:"Allow_mount_linprocfs"`
	Allow_mount_linsysfs  bool   `codec:"Allow_mount_linsysfs"`
	Allow_mount_tmpfs     bool   `codec:"Allow_mount_tmpfs"`
	Allow_mount_zfs       bool   `codec:"Allow_mount_zfs"`
	Allow_vmm             bool   `codec:"Allow_vmm"`
	Linux                 string `codec:"Linux"`
	Linux_osname          string `codec:"Linux_osname"`
	Linux_osrelease       string `codec:"Linux_osrelease"`
	Linux_oss_version     string `codec:"Linux_oss_version"`
	Sysvmsg               string `codec:"Sysvmsg"`
	Sysvsem               string `codec:"Sysvsem"`
	Sysvshm               string `codec:"Sysvshm"`
	Exec_prestart         string `codec:"Exec_prestart"`
	Exec_prestop          string `codec:"Exec_prestop"`
	Exec_created          string `codec:"Exec_created"`
	Exec_start            string `codec:"Exec_start"`
	Exec_stop             string `codec:"Exec_stop"`
	Exec_poststart        string `codec:"Exec_postart"`
	Exec_poststop         string `codec:"Exec_poststop"`
	Exec_clean            bool   `codec:"Exec_clean"`
	Exec_jail_user        string `codec:"Exec_jail_user"`
	Exec_system_jail_user string `codec:"Exec_system_jail_user"`
	Exec_system_user      string `codec:"Exec_system_user"`
	Exec_timeout          uint   `codec:"Exec_timeout"`
	Exec_consolelog       string `codec:"Exec_consolelog"`
	Exec_fib              string `codec:"Exec_fib"`
	Stop_timeout          uint   `codec:"Stop_timeout"`
	Nic                   string `codec:"Nic"`
	Vnet_nic              string `codec:"Vnet_nic"`
	Ip_hostname           string `codec:"Ip_hostname"`
	Mount                 bool   `codec:"Mount"`
	Mount_fstab           string `codec:"Mount_fstab"`
	Mount_devfs           bool   `codec:"Mount_devfs"`
	Mount_fdescfs         bool   `codec:"Mount_fdescfs"`
	Depend                string `codec:"Depend"`
	Rctl                  Rctl   `codec:"Rctl"`
}

// TaskState is the state which is encoded in the handle returned in
// StartTask. This information is needed to rebuild the task state and handler
// during recovery.
type TaskState struct {
	TaskConfig    *drivers.TaskConfig
	ContainerName string
	StartedAt     time.Time
}

func NewJailDriver(logger hclog.Logger) drivers.DriverPlugin {
	ctx, cancel := context.WithCancel(context.Background())
	logger = logger.Named(pluginName)
	return &Driver{
		eventer:        eventer.NewEventer(ctx, logger),
		config:         &Config{},
		tasks:          newTaskStore(),
		ctx:            ctx,
		signalShutdown: cancel,
		logger:         logger,
	}
}

func (d *Driver) PluginInfo() (*base.PluginInfoResponse, error) {
	return pluginInfo, nil
}

func (d *Driver) ConfigSchema() (*hclspec.Spec, error) {
	return nil, nil
}

func (d *Driver) SetConfig(cfg *base.Config) error {
	var config Config
	if len(cfg.PluginConfig) != 0 {
		if err := base.MsgPackDecode(cfg.PluginConfig, &config); err != nil {
			return err
		}
	}

	d.config = &config
	if cfg.AgentConfig != nil {
		d.nomadConfig = cfg.AgentConfig.Driver
	}

	return nil
}

func (d *Driver) Shutdown(ctx context.Context) error {
	d.signalShutdown()
	return nil
}

func (d *Driver) TaskConfigSchema() (*hclspec.Spec, error) {
	return taskConfigSpec, nil
}

func (d *Driver) Capabilities() (*drivers.Capabilities, error) {
	return capabilities, nil
}

func (d *Driver) Fingerprint(ctx context.Context) (<-chan *drivers.Fingerprint, error) {
	ch := make(chan *drivers.Fingerprint)
	go d.handleFingerprint(ctx, ch)
	return ch, nil
}

func (d *Driver) handleFingerprint(ctx context.Context, ch chan<- *drivers.Fingerprint) {
	defer close(ch)
	ticker := time.NewTimer(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			ticker.Reset(fingerprintPeriod)
			ch <- d.buildFingerprint()
		}
	}
}

func (d *Driver) buildFingerprint() *drivers.Fingerprint {
	var health drivers.HealthState
	var desc string
	attrs := map[string]*pstructs.Attribute{"driver.jail": pstructs.NewStringAttribute("1")}
	health = drivers.HealthStateHealthy
	desc = "ready"
	d.logger.Info("buildFingerprint()", "driver.FingerPrint", hclog.Fmt("%+v", health))
	return &drivers.Fingerprint{
		Attributes:        attrs,
		Health:            health,
		HealthDescription: desc,
	}
}

func (d *Driver) RecoverTask(handle *drivers.TaskHandle) error {
	if handle == nil {
		return fmt.Errorf("error: handle cannot be nil")
	}

	if _, ok := d.tasks.Get(handle.Config.ID); ok {
		return nil
	}

	var driverConfig TaskConfig
	if err := handle.Config.DecodeDriverConfig(&driverConfig); err != nil {
		return fmt.Errorf("failed to decode driver config: %v", err)
	}

	var taskState TaskState
	if err := handle.GetDriverState(&taskState); err != nil {
		return fmt.Errorf("failed to decode task state from handle: %v", err)
	}

	_, err := d.initializeContainer(handle.Config, driverConfig)
	if err != nil {
		d.logger.Info("Error RecoverTask k", "driver_cfg", hclog.Fmt("%+v", err))
		return fmt.Errorf("task with ID %q failed", handle.Config.ID)

	}

	h := &taskHandle{
		taskConfig: taskState.TaskConfig,
		State:      drivers.TaskStateRunning,
		startedAt:  taskState.StartedAt,
		exitResult: &drivers.ExitResult{},
		logger:     d.logger,
	}

	d.tasks.Set(taskState.TaskConfig.ID, h)
	go h.run()
	return nil
}

func (d *Driver) StartTask(cfg *drivers.TaskConfig) (*drivers.TaskHandle, *drivers.DriverNetwork, error) {

	if _, ok := d.tasks.Get(cfg.ID); ok {
		return nil, nil, fmt.Errorf("task with ID %q already started", cfg.ID)
	}

	var driverConfig TaskConfig
	if err := cfg.DecodeDriverConfig(&driverConfig); err != nil {
		return nil, nil, fmt.Errorf("failed to decode driver config: %v", err)
	}

	d.logger.Info("starting jail task", "driver_cfg", hclog.Fmt("%+v", driverConfig))
	handle := drivers.NewTaskHandle(taskHandleVersion)
	handle.Config = cfg

	_, err := d.initializeContainer(cfg, driverConfig)
	if err != nil {
		d.logger.Info("Error starting jail task", "driver_cfg", hclog.Fmt("%+v", err))
		return nil, nil, fmt.Errorf("task with ID %q failed", cfg.ID)

	}

	h := &taskHandle{
		taskConfig: cfg,
		State:      drivers.TaskStateRunning,
		startedAt:  time.Now().Round(time.Millisecond),
		logger:     d.logger,
	}

	driverState := TaskState{
		ContainerName: fmt.Sprintf("%s-%s", cfg.Name, cfg.AllocID),
		TaskConfig:    cfg,
		StartedAt:     h.startedAt,
	}

	if err := handle.SetDriverState(&driverState); err != nil {
		d.logger.Error("failed to start task, error setting driver state", "error", err)
		return nil, nil, fmt.Errorf("failed to set driver state: %v", err)
	}

	d.tasks.Set(cfg.ID, h)

	go h.run()

	return handle, nil, nil
}

func (d *Driver) WaitTask(ctx context.Context, taskID string) (<-chan *drivers.ExitResult, error) {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}

	ch := make(chan *drivers.ExitResult)
	go d.handleWait(ctx, handle, ch)

	return ch, nil
}

func (d *Driver) handleWait(ctx context.Context, handle *taskHandle, ch chan *drivers.ExitResult) {
	defer close(ch)

	// Going with simplest approach of polling for handler to mark exit.
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			s := handle.TaskStatus()
			if s.State == drivers.TaskStateExited {
				ch <- handle.exitResult
			}
		}
	}
}

func (d *Driver) StopTask(taskID string, timeout time.Duration, signal string) error {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return drivers.ErrTaskNotFound
	}

	if err := handle.shutdown(timeout); err != nil {
		return fmt.Errorf("executor Shutdown failed: %v", err)
	}

	return nil
}

func (d *Driver) DestroyTask(taskID string, force bool) error {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return drivers.ErrTaskNotFound
	}

	if handle.IsRunning() && !force {
		return fmt.Errorf("cannot destroy running task")
	}

	if handle.IsRunning() {
		// grace period is chosen arbitrary here
		if err := handle.shutdown(1 * time.Minute); err != nil {
			handle.logger.Error("failed to destroy executor", "err", err)
		}
	}

	d.tasks.Delete(taskID)
	return nil
}

func (d *Driver) InspectTask(taskID string) (*drivers.TaskStatus, error) {
	handle, ok := d.tasks.Get(taskID)

	if !ok {
		return nil, drivers.ErrTaskNotFound
	}

	return handle.TaskStatus(), nil
}

func (d *Driver) TaskStats(ctx context.Context, taskID string, interval time.Duration) (<-chan *drivers.TaskResourceUsage, error) {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}

	return handle.stats(ctx, interval)
}

func (d *Driver) TaskEvents(ctx context.Context) (<-chan *drivers.TaskEvent, error) {
	return d.eventer.TaskEvents(ctx)
}

func (d *Driver) SignalTask(taskID string, signal string) error {
	return fmt.Errorf("Jail driver does not support signals")
}

func (d *Driver) ExecTask(taskID string, cmd []string, timeout time.Duration) (*drivers.ExecTaskResult, error) {

	if len(cmd) == 0 {
		return nil, fmt.Errorf("cmd is required, but was empty")
	}
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return handle.Exec(ctx, cmd[0], cmd[1:])

}
