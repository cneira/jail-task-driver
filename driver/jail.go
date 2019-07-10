/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2019, Carlos Neira cneirabustos@gmail.com
 */

package jail

import (
	"crypto/rand"
	"fmt"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/plugins/drivers"
	"os/exec"
	"strings"
)

func simple_uuid() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error calling rand.Read")
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}

func IsJailActive(jailname string) bool {
	args := []string{"-n", "name"}

	out, err := exec.Command("jls", args...).Output()
	if err != nil {
		return false
	}
	jails := strings.Fields(string(out))
	jname := "name=" + jailname
	for _, name := range jails {
		if jname == name {
			return true
		}
	}
	return false
}

// These params don't take a value
func isparamboolean(category string) bool {
	switch category {
	case
		"ip4.saddrsel",
		"nopersist",
		"exec.system_jail_user",
		"exec.clean",
		"vnet",
		"mount.devfs",
		"persist":
		return true
	}
	return false
}

func Jailcmd(params map[string]string) error {
	args := make([]string, 0)
	args = append(args, "-cmr")

	for k, v := range params {
		if isparamboolean(k) {
			args = append(args, k)
		} else {
			param := fmt.Sprintf(k + "=" + v)
			args = append(args, param)
		}
	}
	out, err := exec.Command("jail", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("Jailcmd error args=%+v err=%s out=%s", args, err, string(out))

	}
	return nil
}

func Jailrctl(jname string, params map[string]uint) error {
	args := make([]string, 0)
	args = append(args, "-a")
	for k, v := range params {
		args = append(args, "jail:"+jname+k+ fmt.Sprintf("%d",v))
		out, err := exec.Command("rctl", args...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("applying rctl error args=%+v err=%s out=%s", args, err, string(out))
		}
	}
	return nil 
}

func (d *Driver) initializeContainer(cfg *drivers.TaskConfig, taskConfig TaskConfig) (int32, error) {

	jailparams := make(map[string]string)

	jailparams["name"] = fmt.Sprintf("%s-%s", cfg.Name, cfg.AllocID)
	jailparams["host.hostname"] = fmt.Sprintf("%s-%s", cfg.Name, cfg.AllocID)


  if len(taskConfig.Path) > 0 {
     jailparams["path"] = taskConfig.Path
  } else {
     jailparams["path"] = cfg.AllocDir
  }

	if len(taskConfig.Jid) > 1 {
		jailparams["jid"] = taskConfig.Jid
	}

	if len(taskConfig.Ip4_addr) > 1 {
		jailparams["ip4.addr"] = taskConfig.Ip4_addr
	}

	if taskConfig.Ip4_saddrsel {
		jailparams["ip4.saddrsel"] = "true"
	}

	if len(taskConfig.Ip4) > 1 {
		jailparams["ip4"] = taskConfig.Ip4
	}

	if taskConfig.Ip6_saddrsel {
		jailparams["ip6.saddrsel"] = "true"
	}
	if len(taskConfig.Ip6_addr) > 1 {
		jailparams["ip6.addr"] = taskConfig.Ip6_addr
	}

	if len(taskConfig.Vnet) > 1 {
		jailparams["vnet"] = taskConfig.Vnet
	}

	if len(taskConfig.Host_hostname) > 1 {
		jailparams["host.hostname"] = taskConfig.Host_hostname
	}

	if len(taskConfig.Host) > 1 {
		jailparams["host"] = taskConfig.Host
	}

	if len(taskConfig.Securelevel) > 1 {
		jailparams["securelevel"] = taskConfig.Securelevel
	}

	if len(taskConfig.Devfs_ruleset) > 1 {
		jailparams["devfs_ruleset"] = taskConfig.Devfs_ruleset
	}

	if taskConfig.Children_max > 0 {
		jailparams["children.max"] = fmt.Sprintf("%d", taskConfig.Children_max)
	}

	if taskConfig.Children_cur > 0 {
		jailparams["children.cur"] = fmt.Sprintf("%d", taskConfig.Children_cur)
	}

	if taskConfig.Enforce_statfs > 0 {
		jailparams["enforce_statfs"] = fmt.Sprintf("%d", taskConfig.Enforce_statfs)
	}
	//  A new jail must have either the persist parameter or exec.start or
	//  command pseudo-parameter set.

	if len(taskConfig.Exec_start) > 1 {
		jailparams["exec.start"] = taskConfig.Exec_start
	} else if len(taskConfig.Command) > 1 {
		jailparams["command"] = taskConfig.Command
	} else if taskConfig.Persist == true {
		jailparams["persist"] = "true"
	}

	if len(taskConfig.Osreldate) > 1 {
		jailparams["osreldate"] = taskConfig.Osreldate
	}

	if len(taskConfig.Osrelease) > 1 {
		jailparams["osrelease"] = taskConfig.Osrelease
	}

	if taskConfig.Allow_set_hostname {
		jailparams["allow.set_hostname"] = "true"
	}

	if taskConfig.Allow_sysvipc {
		jailparams["allow.sysvipc"] = "true"
	}

	if taskConfig.Allow_raw_sockets {
		jailparams["allow.raw_sockets"] = "true"
	}

	if taskConfig.Allow_chflags {
		jailparams["allow.chflags"] = "true"
	}

	if taskConfig.Allow_mount {
		jailparams["allow.mount"] = "true"
	}

	if taskConfig.Allow_mount_fdescfs {
		jailparams["allow.mount_fdescfs"] = "true"
	}

	if taskConfig.Allow_mount_fusefs {
		jailparams["allow.mount_fusefs"] = "true"
	}

	if taskConfig.Allow_mount_nullfs {
		jailparams["allow.mount_nullfs"] = "true"
	}

	if taskConfig.Allow_mount_procfs {
		jailparams["allow.mount_procfs"] = "true"
	}

	if taskConfig.Allow_mount_linprocfs {
		jailparams["allow.mount_linprocfs"] = "true"
	}

	if taskConfig.Allow_mount_linsysfs {
		jailparams["allow.mount_linsysfs"] = "true"
	}
	if taskConfig.Allow_mount_tmpfs {
		jailparams["allow.mount_tmpfs"] = "true"
	}

	if taskConfig.Allow_mount_zfs {
		jailparams["allow.mount_zfs"] = "true"
	}

	if taskConfig.Allow_vmm {
		jailparams["allow.vmm"] = "true"
	}

	if len(taskConfig.Linux) > 1 {
		jailparams["linux"] = taskConfig.Linux
	}

	if len(taskConfig.Linux_osname) > 1 {
		jailparams["linux.osname"] = taskConfig.Linux_osname
	}

	if len(taskConfig.Linux_osrelease) > 1 {
		jailparams["linux.osrelease"] = taskConfig.Linux_osrelease
	}

	if len(taskConfig.Sysvsem) > 1 {
		jailparams["sysvsem"] = taskConfig.Sysvsem
	}

	if len(taskConfig.Sysvmsg) > 1 {
		jailparams["sysvmsg"] = taskConfig.Sysvmsg
	}

	if len(taskConfig.Sysvshm) > 1 {
		jailparams["sysvshm"] = taskConfig.Sysvshm
	}

	if len(taskConfig.Exec_prestart) > 1 {
		jailparams["exec.prestart"] = taskConfig.Exec_prestart
	}

	if len(taskConfig.Exec_prestop) > 1 {
		jailparams["exec.prestop"] = taskConfig.Exec_prestop
	}

	if len(taskConfig.Exec_created) > 1 {
		jailparams["exec.created"] = taskConfig.Exec_created
	}

	if len(taskConfig.Exec_poststart) > 1 {
		jailparams["exec.poststart"] = taskConfig.Exec_poststart
	}

	if len(taskConfig.Exec_stop) > 1 {
		jailparams["exec.stop"] = taskConfig.Exec_stop
	}

	if len(taskConfig.Exec_poststop) > 1 {
		jailparams["exec.poststop"] = taskConfig.Exec_poststop
	}

	if taskConfig.Exec_clean {
		jailparams["exec.clean"] = "true"
	}

	if len(taskConfig.Exec_jail_user) > 1 {
		jailparams["exec.jail_user"] = taskConfig.Exec_jail_user
	}

	if len(taskConfig.Exec_system_user) > 1 {
		jailparams["exec.system_user"] = taskConfig.Exec_system_user
	}

	if taskConfig.Exec_timeout > 0 {
		jailparams["exec.timeout"] = fmt.Sprintf("%d", taskConfig.Exec_timeout)
	}

	if len(taskConfig.Exec_consolelog) > 1 {
		jailparams["exec.consolelog"] = taskConfig.Exec_consolelog
	}

	if taskConfig.Stop_timeout > 0 {
		jailparams["stop.timeout"] = fmt.Sprintf("%d", taskConfig.Stop_timeout)
	}

	if len(taskConfig.Nic) > 1 {
		jailparams["interface"] = taskConfig.Nic
	}

	if len(taskConfig.Vnet_nic) > 1 {
		jailparams["vnet.interface"] = taskConfig.Vnet_nic
	}

	if len(taskConfig.Ip_hostname) > 1 {
		jailparams["ip_hostname"] = taskConfig.Ip_hostname
	}

	if taskConfig.Mount {
		jailparams["mount"] = "true"
	}

	if len(taskConfig.Mount_fstab) > 1 {
		jailparams["mount.fstab"] = taskConfig.Mount_fstab
	}

	if taskConfig.Mount_devfs {
		jailparams["mount.devfs"] = "true"
	}

	if taskConfig.Mount_devfs {
		jailparams["mount.devfs"] = "true"
	}

	if len(taskConfig.Depend) > 1 {
		jailparams["depend"] = taskConfig.Depend
	}

	//RCTL options

	rctlm := make(map[string]uint)
	rctl := taskConfig.Rctl

	if rctl.Cputime > 0 {
		rctlm[":cputime:deny="] = rctl.Cputime
	}

	if rctl.Stacksize > 0 {
		rctlm[":stacksize:deny="] = rctl.Stacksize
	}

	if rctl.Coredumpsize > 0 {
		rctlm[":coredumpsize:deny="] = rctl.Coredumpsize
	}
	if rctl.Memoryuse > 0 {
		rctlm[":memoryuse:deny="] = rctl.Memoryuse
	}
	if rctl.Memorylocked > 0 {
		rctlm[":memorylocked:deny="] = rctl.Memorylocked
	}
	if rctl.Maxproc > 0 {
		rctlm[":maxproc:deny="] = rctl.Maxproc
	}
	if rctl.Openfiles > 0 {
		rctlm[":openfile:deny="] = rctl.Openfiles
	}
	if rctl.Vmemoryuse > 0 {
		rctlm[":vmemoryuse:deny="] = rctl.Vmemoryuse
	}
	if rctl.Pseudoterminals > 0 {
		rctlm[":pseudote:deny="] = rctl.Pseudoterminals
	}
	if rctl.Swapuse > 0 {
		rctlm[":swapuse:deny="] = rctl.Swapuse
	}
	if rctl.Nthr > 0 {
		rctlm[":nthr:deny="] = rctl.Nthr
	}
	if rctl.Msgqqueued > 0 {
		rctlm[":msgqqueued:deny="] = rctl.Msgqqueued
	}

	if rctl.Msgqsize > 0 {
		rctlm[":msgqsize:deny="] = rctl.Msgqsize
	}

	if rctl.Nmsgq > 0 {
		rctlm[":nmsg:deny="] = rctl.Nmsgq
	}
	if rctl.Nsemop > 0 {
		rctlm[":nsem:deny="] = rctl.Nsemop
	}
	if rctl.Nshm > 0 {
		rctlm[":nshm:deny="] = rctl.Nshm
	}
	if rctl.Shmsize > 0 {
		rctlm[":shmsize:deny="] = rctl.Shmsize
	}
	if rctl.Wallclock > 0 {
		rctlm[":wallclock:deny="] = rctl.Wallclock
	}
	if rctl.Pcpu > 0 {
		rctlm[":pcpu:deny="] = rctl.Pcpu
	}
	if rctl.Readbps > 0 {
		rctlm[":readbps:deny="] = rctl.Readbps
	}
	if rctl.Writebps > 0 {
		rctlm[":writebps:deny="] = rctl.Writebps
	}
	if rctl.Readiops > 0 {
		rctlm[":readiops:deny="] = rctl.Readiops
	}
	if rctl.Writeiops > 0 {
		rctlm[":writeiops:deny="] = rctl.Writeiops
	}

	err := Jailcmd(jailparams)

	if err != nil {
		d.logger.Info("Error Creating Jail", "driver_initialize_container", hclog.Fmt("Params %+v", jailparams))
		d.logger.Info("Error Creating Jail", "driver_initialize_container", hclog.Fmt("%s", err))
		return -1, fmt.Errorf("Calling jail failed %s", err)
	}

	err = Jailrctl(jailparams["name"], rctlm)

	if err != nil {
		d.logger.Info("Error setting resource control ", "driver_initialize_container", hclog.Fmt("%s", err))
		return -1, fmt.Errorf("Calling rctl failed %s", err)
	}
	return 0, nil
}
