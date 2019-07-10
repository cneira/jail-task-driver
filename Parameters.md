Jail creation parameters
-----------------------
This is a verbatim copy for JAIL(8) parameters but modified for naming convention
used in jail-task-driver


**Jid**     The jail identifier.  This will be assigned automatically to a
            new jail (or can be explicitly set), and can be used to identify
e the jail, either by jail or from jexec(8), are run
is directory.If this parameter is omitted then it will
	     use nomad's allocation directory as default value.

              ** list of IPv4** addresses assigned to the jail.  If this is set,
             the jail is restricted to using only these addresses.  Any
              attempts to use other addresses fail, and attempts to use
              wildcard addresses silently use the jailed address instead.  For
             IPv4 the first address given will be used as the source address
             when source address selection on unbound sockets cannot find a
             better match.  It is only possible to start multiple jails with
**the same IP addre**ss if none of the jails has more than this
             single overlapping IP address assigned to itself.

**Ip4_saddrsel**
             A boolean option to change the formerly mentioned behaviour and
             disable IPv4 source address selection for the jail in favour of
             the primary IPv4 address of the jail.  Source address selection
             is enabled by default for all jails and the ip4.nosaddrsel
             setting of a parent jail is not inherited for any child jails.

 **Ip4**     Control the availability of IPv4 addresses.  Possible values are
             "inherit" to allow unrestricted access to all system addresses,
             "new" to restrict addresses via ip4.addr, and "disable" to stop
             the jail from using IPv4 entirely.  Setting the ip4.addr
             parameter implies a value of "new".

 **Ip6_addr**, **Ip6_saddrsel**, **Ip6**
             A set of IPv6 options for the jail, the counterparts to Ip4_addr,
             Ip4_saddrsel and ip4 above.

**Vnet**
	     Create the jail with its own virtual network stack, with its own
             network interfaces, addresses, routing table, etc.  The kernel
             must have been compiled with the VIMAGE option for this to be
             available.  Possible values are "inherit" to use the system
             network stack, possibly with restricted IP addresses, and "new"
             to create a new network stack.

**Host_hostname**
	     The hostname of the jail.  Other similar parameters are
             host.domainname, host.hostuuid and host.hostid.

**Host**     Set the origin of hostname and related information.  Possible
             values are "inherit" to use the system information and "new" for
             the jail to use the information from the above fields.  Setting
             any of the above fields implies a value of "new".

**Securelevel**  
             The value of the jail's kern.securelevel sysctl.  A jail never
             has a lower securelevel than its parent system, but by setting
             this parameter it may have a higher one.  If the system
             securelevel is changed, any jail securelevels will be at least as
             secure.

**Devfs_ruleset**
             The number of the devfs ruleset that is enforced for mounting
             devfs in this jail.  A value of zero (default) means no ruleset
             is enforced.  Descendant jails inherit the parent jail's devfs
             ruleset enforcement.  Mounting devfs inside a jail is possible
             only if the allow.mount and allow.mount.devfs permissions are
             effective and enforce_statfs is set to a value lower than 2.
             Devfs rules and rulesets cannot be viewed or modified from inside
             a jail.

             NOTE: It is important that only appropriate device nodes in devfs
             be exposed to a jail; access to disk devices in the jail may
             permit processes in the jail to bypass the jail sandboxing by
             modifying files outside of the jail.  See devfs(8) for
             information on how to use devfs rules to limit access to entries
             in the per-jail devfs.  A simple devfs ruleset for jails is
             available as ruleset #4 in /etc/defaults/devfs.rules.

**Children_max**
 	The number of childjails allowed to be created by this jail (or
        by other jails under this jail).  This limit is zero by default,
        indicating the jail is not allowed to create child jails.  See
        the Hierarchical Jails section for more information.

**Children_cur**
        The number of descendants of this jail, including its own child
        jails and any jails created under them.

**Enforce_statfs**
        This determines what information processes in a jail are able to
        get about mount points.  It affects the behaviour of the
        following syscalls: statfs(2), fstatfs(2), getfsstat(2), and
        fhstatfs(2)** (as well as similar compatibility syscalls).  When
        set to 0, all mount points are available without any
        restrictions.  When set to 1, only mount points below the jail's
        chroot directory are visible.  In addition to that, the path to
        the jail's chroot directory is removed from the front of their
        pathnames.  When set to 2 (default), above syscalls can operate
        only on a **mount-point where the jail's chroot directory is
        located.

**Persist**
             Setting this boolean parameter allows a jail to exist without any
             processes.  Normally, a command is run as part of jail creation,
             and then the jail is destroyed as its last process exits.  A new
             jail must have either the persist parameter or exec.start or
             command pseudo-parameter set.

**Osrelease**
             The string for the jail's kern.osrelease sysctl and uname -r.
**reldate****
             The number for the jail's kern.osreldate and uname -K.

**Allow_set_hostname**
                     The jail's hostname may be changed via hostname(1) or
                     sethostname(3).
**low_sysvipc****
                     A process within the jail has access to System V IPC
                     primitives.  This is deprecated in favor of the per-
                     module parameters (see below).  When this parameter is
                     set, it is equivalent to setting sysvmsg, sysvsem, and
                     sysvshm all to "inherit".
**Allow_raw_sockets**
                     The jail root is allowed to create raw sockets.  Setting
                     this parameter allows utilities like ping(8) and
                     traceroute(8) to operate inside the jail.  If this is
                     set, the source IP addresses are enforced to comply with
                     the IP address bound to the jail, regardless of whether
	             or not the IP_HDRINCL flag has been set on the socket.**
                     Since raw sockets can be used to configure and interact
                     with various network subsystems, extra caution should be
                     used where privileged access to jails is given out to
                     untrusted parties.

**Allow_chflags**
                     Normally, privileged users inside a jail are treated as
                     unprivileged by chflags(2).  When this parameter is set,
                     such users are treated as privileged, and may manipulate
                     system file flags subject to the usual constraints on
                     kern.securelevel.

**Allow_mount** 
                     privileged users inside the jail will be able to mount
                     and unmount file system types marked as jail-friendly.
                     The lsvfs(1) command can be used to find file system
                     types available for mount from within a jail.  This
                     permission is effective only if enforce_statfs is set to
                     a value lower than 2.

**Allow_mount_devfs**
**     _      **  **L pri_ileged us**e**L ins_de the jail** will be able to mount
                     and unmount the devfs file system.  This permission is
                     effective only together with allow.mount and only when
                     enforce_statfs is set to a value lower than 2.  The devfs
**       **         ruleset should be restricted from the default by using
                     the devfs_ruleset option.

**Allow_quotas**
                     The jail root may administer quotas on the jail's
                     filesystem(s).  This includes filesystems that the jail
                     may share with other jails or with non-jailed parts of
                     the system.

**Allow_read_msgbuf**
                     Jailed users may read the kernel message buffer.  If the
                     security.bsd.unprivileged_read_msgbuf MIB entry is zero,
                     this will be restricted to the root user.

**Allow_socket_af**
                Sockets within a jail are normally restricted to IPv4,
                IPv6, local (UNIX), and route.  This allows access to
                other protocol stacks that have not had jail
                functionality added to them.
**Allow_mlock**
                Locking or unlocking physical pages in memory are
                normally not available within a jail.  When this
                parameter is set, users may mlock(2) or munlock(2) memory
                subject to security.bsd.unprivileged_mlock and resource
                limits.

**Allow_reserved_ports**
                     The jail root may bind to ports lower than 1024.

Kernel modules may add their own parameters, which only exist when the
module is loaded.  These are typically headed under a parameter named
after the module, with values of "inherit" to give the jail full use of
the module, "new" to encapsulate the jail in some module-specific way,
and "disable" to make the module unavailable to the jail.  There also may
be other parameters to define jail behavior within the module.  Module-
specific parameters include:

**Allow_mount_fdescfs**
             privileged users inside the jail will be able to mount and
             unmount the fdescfs file system.  This permission is effective
             only together with allow.mount and only when enforce_statfs is
             set to a value lower than 2.

**Allow_mount_fusefs**
             privileged users inside the jail will be able to mount and
             unmount fuse-based file systems.  This permission is effective
             only together with allow.mount and only when enforce_statfs is
             set to a value lower than 2.

**Allow_mount_nullfs**
             privileged users inside the jail will be able to mount and
             unmount the nullfs file system.  This permission is effective
             only together with allow.mount and only when enforce_statfs is
             set to a value lower than 2.

**Allow_mount_procfs**
             privileged users inside the jail will be able to mount and
             unmount the procfs file system.  This permission is effective
             only together with allow.mount and only when enforce_statfs is
             set to a value lower than 2.

**Allow_mount_linprocfs**
             privileged users inside the jail will be able to mount and
             unmount the linprocfs file system.  This permission is effective
             only together with allow.mount and only when enforce_statfs is
             set to a value lower than 2.

**Allow_mount_linsysfs**
             privileged users inside the jail will be able to mount and
             unmount the linsysfs file system.  This permission is effective
             only together with allow.mount and only when enforce_statfs is
             set to a value lower than 2.

**Allow_mount_tmpfs**
             privileged users inside the jail will be able to mount and
             unmount the tmpfs file system.  This permission is effective only
             together with allow.mount and only when enforce_statfs is set to
             a value lower than 2.

**Allow_mount_zfs**
             privileged users inside the jail will be able to mount and
             unmount the ZFS file system.  This permission is effective only
             together with allow.mount and only when enforce_statfs is set to
             a value lower than 2.  See zfs(8) for information on how to
             configure the ZFS filesystem to operate from within a jail.

**Allow_vmm**
             The jail may access vmm(4).  This flag is only available when the
             vmm(4) kernel module is loaded.

**Linux**    Determine how a jail's Linux emulation environment appears.  A
             value of "inherit" will keep the same environment, and "new" will
             give the jail it's own environment (still originally inherited
             when the jail is created).

**Linux_osname**, **Linux_osrelease**,**Linux_oss_version**
             The Linux OS name, OS release, and OSS version associated with
             this jail.

**Sysvmsg**
             Allow access to SYSV IPC message primitives.  If set to
             "inherit", all IPC objects on the system are visible to this
             jail, whether they were created by the jail itself, the base
             system, or other jails.  If set to "new", the jail will have its
             own key namespace, and can only see the objects that it has
             created; the system (or parent jail) has access to the jail's
             objects, but not to its keys.  If set to "disable", the jail
             cannot perform any sysvmsg-related system calls.

**Sysvsem**, **sysvshm**
             Allow access to SYSV IPC semaphore and shared memory primitives,
             in the same manner as sysvmsg.

There are pseudo-parameters that are not passed to the kernel, but are
used by jail to set up the jail environment, often by running specified
commands when jails are created or removed.  The exec.* command
parameters are sh(1) command lines that are run in either the system or
jail environment.  They may be given multiple values, which would run the
specified commands in sequence.  All commands must succeed (return a zero
exit status), or the jail will not be created or removed, as appropriate.

The pseudo-parameters are:

**Exec_prestart**
             Command(s) to run in the system environment before a jail is
             created.

**Exec_created**
             Command(s) to run in the system environment right after a jail
             has been created, but before commands (or services) get executed
             in the jail.

**Exec_start**
             Command(s) to run in the jail environment when a jail is created.
             A typical command to run is "sh /etc/rc".

**Command**
             A synonym for exec.start for use when specifying a jail directly
             on the command line.  Unlike other parameters whose value is a
             single string, command uses the remainder of the jail command
             line as its own arguments.

**Exec_poststart**
             Command(s) to run in the system environment after a jail is
             created, and after any exec.start commands have completed.

**Exec_prestop**
             Command(s) to run in the system environment before a jail is
             removed.

**Exec_stop**
             Command(s) to run in the jail environment before a jail is
             removed, and after any exec.prestop commands have completed.  A
             typical command to run is "sh /etc/rc.shutdown".

**Exec_poststop**
             Command(s) to run in the system environment after a jail is
             removed.

**Exec_clean**
             Run commands in a clean environment.  The environment is
             discarded except for HOME, SHELL, TERM and USER.  HOME and SHELL
             are set to the target login's default values.  USER is set to the
             target login.  TERM is imported from the current environment.
             The environment variables from the login class capability
             database for the target login are also set.

**Exec_jail_user**
             The user to run commands as, when running in the jail
             environment.  The default is to run the commands as the current
             user.

**Exec_system_jail_user**
             This boolean option looks for the exec.jail_user in the system
             passwd(5) file, instead of in the jail's file.

**Exec_system_user**
             The user to run commands as, when running in the system
             environment.  The default is to run the commands as the current
             user.

**Exec_timeout**
             The maximum amount of time to wait for a command to complete, in
             seconds.  If a command is still running after this timeout has
             passed, the jail will not be created or removed, as appropriate.

**Exec_consolelog**
             A file to direct command output (stdout and stderr) to.

**Exec_fib**
             The FIB (routing table) to set when running commands inside the
             jail.

**Stop_timeout**
             The maximum amount of time to wait for a jail's processes to exit
             after sending them a SIGTERM signal (which happens after the
             exec.stop commands have completed).  After this many seconds have
             passed, the jail will be removed, which will kill any remaining
             processes.  If this is set to zero, no SIGTERM is sent and the
             jail is immediately removed.  The default is 10 seconds.

**Nic** 
             A network interface to add the jail's IP addresses (ip4_addr and
             ip6_addr) to.  An alias for each address will be added to the
             interface before the jail is created, and will be removed from
             the interface after the jail is removed.

**Ip4_addr**
             In addition to the IP addresses that are passed to the kernel, an
             interface, netmask and additional parameters (as supported by
             ifconfig(8)) may also be specified, in the form
             "interface|ip-address/netmask param ...".  If an interface is
             given before the IP address, an alias for the address will be
             added to that interface, as it is with the interface parameter.
             If a netmask in either dotted-quad or CIDR form is given after an
             IP address, it will be used when adding the IP alias.  If
             additional parameters are specified then they will also be used
             when adding the IP alias.

**Ip6_addr**
             In addition to the IP addresses that are passed to the kernel, an
             interface, prefix and additional parameters (as supported by
             ifconfig(8)) may also be specified, in the form
             "interface|ip-address/prefix param ...".

**Vnet_nic**
             A network interface to give to a vnet-enabled jail after is it
             created.  The interface will automatically be released when the
             jail is removed.

**Ip_hostname**
             Resolve the host.hostname parameter and add all IP addresses
             returned by the resolver to the list of addresses (ip4.addr or
             ip6.addr) for this jail.  This may affect default address
             selection for outgoing IPv4 connections from jails.  The address
             first returned by the resolver for each address family will be
             used as the primary address.

**Mount**    A filesystem to mount before creating the jail (and to unmount
             after removing it), given as a single fstab(5) line.

**Mount_fstab**
             An fstab(5) format file containing filesystems to mount before
             creating a jail.

**Mount_devfs**
             Mount a devfs(5) filesystem on the chrooted /dev directory, and
             apply the ruleset in the devfs_ruleset parameter (or a default of
             ruleset 4: devfsrules_jail) to restrict the devices visible
             inside the jail.

**Mount_fdescfs**
             Mount a fdescfs(5) filesystem on the chrooted /dev/fd directory.

**Mount_procfs**
             Mount a procfs(5) filesystem on the chrooted /proc directory.


**Depend**   Specify a jail (or jails) that this jail depends on.  When this
             jail is to be created, any jail(s) it depends on must already
             exist.  If not, they will be created automatically, up to the
             completion of the last exec.poststart command, before any action
             will taken to create this jail.  When jails are removed the
             opposite is true: this jail will be removed, up to the last
             exec.poststop command, before any jail(s) it depends on are
             stopped.
   
  

Jail resource control parameters
--------------------------------
This is a verbatim copy for RCTL(8) parameters but modified for naming convention
used in jail-task-driver, the rctl action is always deny for these parameters.

**Cputime**		   CPU time, in	seconds
**Datasize**		   data	size, in bytes
**Stacksize**	   	   stack size, in bytes
**Coredumpsize**	   core	dump size, in bytes
**Memoryuse**	   	   resident set	size, in bytes
**Memorylocked**	   locked memory, in bytes
**Maxproc**		   number of processes
**Openfiles**	   	   file	descriptor table size
**Vmemoryuse** 	  	   address space limit,in bytes
**Pseudoterminals**	   number of PTYs
**Swapuse**		   swap	space that may be reserved or used, in bytes
**Nthr**		   number of threads
**Msgqqueued** 	   	   number of queued SysV messages
**Msgqsize**	 	   SysV	message	queue size, in bytes
**Nmsgq**		   number of SysV message queues
**Nsem**		   number of SysV semaphores
**Nsemop**		   number of SysV semaphores modified in a single
			   semop(2) call
**Nshm**		   number of SysV shared memory	segments
**Shmsize**		   SysV	shared memory size, in bytes
**Wallclock**		   wallclock time, in seconds
**Pcpu**		   %CPU, in percents of	a single CPU core
**Readbps**		   filesystem reads, in	bytes per second
**Writebps**		   filesystem writes, in bytes per second
**Readiops**	           filesystem reads, in	operations per second
**Writeiops**	   	   filesystem writes, in operations per	second
