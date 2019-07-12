/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2019, Carlos Neira cneirabustos@gmail.com
 */

package jail

//#cgo LDFLAGS: -lutil
//#include <libutil.h>
//#include <stdlib.h>
import "C"

import (
	"bytes"
	"crypto/rand"
	"fmt"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/plugins/drivers"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unsafe"
)

const (
	// containerMonitorIntv is the interval at which the driver checks if the
	// container is still running

	containerMonitorIntv = 2 * time.Second
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

func WaitTillStopped(jname string) (bool, error) {
	for {
		if IsJailActive(jname) == true {
			time.Sleep(containerMonitorIntv)
		} else {
			return true, nil
		}
	}
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
	cmd := exec.Command("jail", args...)
	buf := &bytes.Buffer{}
	buferr := &bytes.Buffer{}
	cmd.Stdout = buf
	cmd.Stderr = buferr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Jailcmd error args=%+v err=%s stdout=%s stderr=%s", args, err, buf.String(), buferr.String())
	}
	return nil
}

func Jailrctl(jname string, params map[string]uint64) error {
	args := make([]string, 0)
	args = append(args, "-a")
	for k, _ := range params {
		args = append(args, "jail:"+jname+k)
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
		jailparams["path"] = filepath.Join(cfg.AllocDir, cfg.Name)
	}

	if len(taskConfig.Jid) > 0 {
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

	rctlm := make(map[string]uint64)
	rctl := taskConfig.Rctl
	var amnt C.uint64_t

	if len(rctl.Cputime.Amount) > 0 {
		cs := C.CString(rctl.Cputime.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Cputime is invalid %s")
		}
		if len(rctl.Cputime.Per) > 0 {
			rctlm[":cputime:"+rctl.Cputime.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) + "/" +  rctl.Cputime.Per] = (uint64)(amnt)
		} else {
			rctlm[":cputime:"+rctl.Cputime.Action+"="] = (uint64)(amnt)
		}
	}

	if len(rctl.Stacksize.Amount) > 0 {
		cs := C.CString(rctl.Stacksize.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Stacksize is invalid %s")
		}
		if len(rctl.Stacksize.Per) > 0 {
			rctlm[":stacksize:"+rctl.Stacksize.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) + "/"+rctl.Stacksize.Per] = (uint64)(amnt)
		} else {
			rctlm[":stacksize:"+rctl.Stacksize.Action+"="] = (uint64)(amnt)
		}
	}

	if len(rctl.Coredumpsize.Amount) > 0 {
		cs := C.CString(rctl.Stacksize.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Coredumpsize is invalid %s")
		}
		if len(rctl.Coredumpsize.Per) > 0 {
			rctlm[":coredumpsize:"+rctl.Coredumpsize.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) + "/"+ rctl.Coredumpsize.Per] = (uint64)(amnt)
		} else {
			rctlm[":coredumpsize:"+rctl.Coredumpsize.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}
	}

	if len(rctl.Memoryuse.Amount) > 0 {
		cs := C.CString(rctl.Memoryuse.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Memoryuse is invalid")
		}
		if len(rctl.Memoryuse.Per) > 0 {
			rctlm[":memoryuse:"+rctl.Memoryuse.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) + "/"+rctl.Memoryuse.Per+"="] = (uint64)(amnt)
		} else {
			rctlm[":memoryuse:"+rctl.Memoryuse.Action+"="+ fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}
	}

	if len(rctl.Memorylocked.Amount) > 0 {
		cs := C.CString(rctl.Memorylocked.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Memorylocked is invalid")
		}
		if len(rctl.Memorylocked.Per) > 0 {
			rctlm[":memorylocked:"+rctl.Memorylocked.Action+ "=" + fmt.Sprintf("%d",(uint64)(amnt)) + "/"+rctl.Memorylocked.Per] = (uint64)(amnt)
		} else {
			rctlm[":memorylocked:"+rctl.Memorylocked.Action+"="+ fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}
	}

	if len(rctl.Maxproc.Amount) > 0 {
		cs := C.CString(rctl.Maxproc.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Maxproc is invalid")
		}
		if len(rctl.Maxproc.Per) > 0 {
			rctlm[":maxproc:"+rctl.Maxproc.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) + "/" +  rctl.Maxproc.Per] = (uint64)(amnt)
		} else {
			rctlm[":maxproc:"+rctl.Maxproc.Action+"="+ fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}
	}

	if len(rctl.Openfiles.Amount) > 0 {
		cs := C.CString(rctl.Openfiles.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Openfiles is invalid")
		}
		if len(rctl.Openfiles.Per) > 0 {
			rctlm[":openfiles:"+rctl.Openfiles.Action+"=" +  fmt.Sprintf("%d",(uint64)(amnt))+ "/"+rctl.Openfiles.Per] = (uint64)(amnt)
		} else {
			rctlm[":openfiles:"+rctl.Openfiles.Action+"="+ fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}
	}
	if len(rctl.Vmemoryuse.Amount) > 0 {
		cs := C.CString(rctl.Vmemoryuse.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Vmemoryuse is invalid")
		}
		if len(rctl.Vmemoryuse.Per) > 0 {
			rctlm[":vmemoryuse:"+rctl.Vmemoryuse.Action+"=" +  fmt.Sprintf("%d",(uint64)(amnt))+"/"+rctl.Vmemoryuse.Per] = (uint64)(amnt)
		} else {
			rctlm[":vmemoryuse:"+rctl.Vmemoryuse.Action+"="+  fmt.Sprintf("%d",(uint64)(amnt))] = (uint64)(amnt)
		}
	}

	if len(rctl.Pseudoterminals.Amount) > 0 {
		cs := C.CString(rctl.Pseudoterminals.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Pseudoterminals is invalid")
		}
		if len(rctl.Pseudoterminals.Per) > 0 {
			rctlm[":pseudoterminals:"+rctl.Pseudoterminals.Action+"=" +  fmt.Sprintf("%d",(uint64)(amnt))+"/"+rctl.Pseudoterminals.Per] = (uint64)(amnt)
		} else {
			rctlm[":pseudoterminals:"+rctl.Pseudoterminals.Action+"="+  fmt.Sprintf("%d",(uint64)(amnt))] = (uint64)(amnt)
		}

	}
	if len(rctl.Swapuse.Amount) > 0 {
		cs := C.CString(rctl.Swapuse.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Swapuse is invalid")
		}
		if len(rctl.Swapuse.Per) > 0 {
			rctlm[":swapuse:"+rctl.Swapuse.Action+"=" +  fmt.Sprintf("%d",(uint64)(amnt))+"/"+rctl.Swapuse.Per] = (uint64)(amnt)
		} else {
			rctlm[":swapuse:"+rctl.Swapuse.Action+"="+ fmt.Sprintf("%d",(uint64)(amnt))] = (uint64)(amnt)
		}
	}

	if len(rctl.Nthr.Amount) > 0 {
		cs := C.CString(rctl.Nthr.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Nthr is invalid")
		}
		if len(rctl.Nthr.Per) > 0 {
			rctlm[":nthr:"+rctl.Nthr.Action+"=" +  fmt.Sprintf("%d",(uint64)(amnt))+"/"+rctl.Nthr.Per] = (uint64)(amnt)
		} else {
			rctlm[":nthr:"+rctl.Nthr.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt))] = (uint64)(amnt)
		}
	}

	if len(rctl.Msgqqueued.Amount) > 0 {
		cs := C.CString(rctl.Msgqqueued.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Msgqqueued is invalid")
		}
		if len(rctl.Msgqqueued.Per) > 0 {
			rctlm[":msgqqueued:"+rctl.Msgqqueued.Action+"=" +  fmt.Sprintf("%d",(uint64)(amnt))+"/"+rctl.Msgqqueued.Per] = (uint64)(amnt)
		} else {
			rctlm[":msgqqueued:"+rctl.Msgqqueued.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt))] = (uint64)(amnt)
		}
	}

	if len(rctl.Msgqsize.Amount) > 0 {
		cs := C.CString(rctl.Msgqsize.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Msgqsize is invalid")
		}
		if len(rctl.Msgqsize.Per) > 0 {
			rctlm[":msgqsize:"+rctl.Msgqsize.Action+ "=" + fmt.Sprintf("%d",(uint64)(amnt)) + "/"+rctl.Msgqsize.Per] = (uint64)(amnt)
		} else {
			rctlm[":msgqsize:"+rctl.Msgqsize.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt))] = (uint64)(amnt)
		}
	}

	if len(rctl.Nmsgq.Amount) > 0 {
		cs := C.CString(rctl.Nmsgq.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Nmsgq is invalid")
		}
		if len(rctl.Nmsgq.Per) > 0 {
			rctlm[":nmsgq:"+rctl.Nmsgq.Action+ "=" + fmt.Sprintf("%d",(uint64)(amnt)) + "/"+rctl.Nmsgq.Per] = (uint64)(amnt)
		} else {
			rctlm[":nmsgq:"+rctl.Nmsgq.Action+"="+ fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}
	}

	if len(rctl.Nsemop.Amount) > 0 {
		cs := C.CString(rctl.Nsemop.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Nsemop is invalid")
		}
		if len(rctl.Nsemop.Per) > 0 {
			rctlm[":nsemop:"+rctl.Nsemop.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) +  "/"+rctl.Nsemop.Per] = (uint64)(amnt)
		} else {
			rctlm[":nsemop:"+rctl.Nsemop.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}

	}
	if len(rctl.Nshm.Amount) > 0 {
		cs := C.CString(rctl.Nshm.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Nshm is invalid")
		}
		if len(rctl.Nshm.Per) > 0 {
			rctlm[":nshm:"+rctl.Nshm.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) + "/"+rctl.Nshm.Per] = (uint64)(amnt)
		} else {
			rctlm[":nshm:"+rctl.Nshm.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}
	}
	if len(rctl.Shmsize.Amount) > 0 {
		cs := C.CString(rctl.Shmsize.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Shmsize is invalid")
		}
		if len(rctl.Shmsize.Per) > 0 {
			rctlm[":shmsize:"+rctl.Shmsize.Action+ "=" + fmt.Sprintf("%d",(uint64)(amnt)) +"/"+rctl.Shmsize.Per] = (uint64)(amnt)
		} else {
			rctlm[":shmsize:"+rctl.Shmsize.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}
	}

	if len(rctl.Wallclock.Amount) > 0 {
		cs := C.CString(rctl.Wallclock.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Wallclock is invalid")
		}
		if len(rctl.Wallclock.Per) > 0 {
			rctlm[":wallclock:"+rctl.Wallclock.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) + "/"+rctl.Wallclock.Per] = (uint64)(amnt)
		} else {
			rctlm[":wallclock:"+rctl.Wallclock.Action+"="+ fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}

	}

	if len(rctl.Pcpu.Amount) > 0 {
		cs := C.CString(rctl.Pcpu.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Pcpu is invalid")
		}
		if len(rctl.Pcpu.Per) > 0 {
			rctlm[":pcpu:"+rctl.Pcpu.Action+ "=" + fmt.Sprintf("%d",(uint64)(amnt)) +"/"+rctl.Pcpu.Per] = (uint64)(amnt)
		} else {
			rctlm[":pcpu:"+rctl.Pcpu.Action+"=" +  fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}
	}

	if len(rctl.Readbps.Amount) > 0 {
		cs := C.CString(rctl.Readbps.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Readbps is invalid")
		}
		if len(rctl.Readbps.Per) > 0 {
			rctlm[":readbps:"+rctl.Readbps.Action+ "="+  fmt.Sprintf("%d",(uint64)(amnt)) + "/"+rctl.Readbps.Per] = (uint64)(amnt)
		} else {
			rctlm[":readbps:"+rctl.Readbps.Action+"="+  fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}
	}

	if len(rctl.Writebps.Amount) > 0 {
		cs := C.CString(rctl.Writebps.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Writebps is invalid")
		}
		if len(rctl.Writebps.Per) > 0 {
			rctlm[":writebps:"+rctl.Writebps.Action+"="+ fmt.Sprintf("%d",(uint64)(amnt)) +"/"+rctl.Writebps.Per] = (uint64)(amnt)
		} else {
			rctlm[":writebps:"+rctl.Writebps.Action+"=" +   fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}

	}
	if len(rctl.Readiops.Amount) > 0 {
		cs := C.CString(rctl.Readiops.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Readiops is invalid")
		}
		if len(rctl.Readiops.Per) > 0 {
			rctlm[":readiops:"+rctl.Readiops.Action+"="+ fmt.Sprintf("%d",(uint64)(amnt)) +"/"+rctl.Readiops.Per] = (uint64)(amnt)
		} else {
			rctlm[":readiops:"+rctl.Readiops.Action+"=" + fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}
	}
	if len(rctl.Writeiops.Amount) > 0 {
		cs := C.CString(rctl.Writeiops.Amount)
		defer C.free(unsafe.Pointer(cs))
		err := C.expand_number(cs, &amnt)
		if err != 0 {
			return -1, fmt.Errorf("Amount for Writeiops is invalid")
		}
		if len(rctl.Writeiops.Per) > 0 {
			rctlm[":writeiops:"+rctl.Writeiops.Action+"="+ fmt.Sprintf("%d",(uint64)(amnt)) + "/"+rctl.Writeiops.Per] = (uint64)(amnt)
		} else {
			rctlm[":writeiops:"+rctl.Writeiops.Action+"="+ fmt.Sprintf("%d",(uint64)(amnt)) ] = (uint64)(amnt)
		}

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
