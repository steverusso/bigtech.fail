#!/bin/sh

# PROVIDE: bigtechfaild
# REQUIRE: DAEMON
# KEYWORD: shutdown

. /etc/rc.subr

name="bigtechfaild"

rcvar="${name}_enable"
load_rc_config ${name}
: ${bigtechfaild_enable:=NO}

pidfile="/var/run/${name}.pid"
logfile="/var/log/${name}.log"
procname="/usr/local/bin/${name}"

command="/usr/sbin/daemon"
command_args="-p ${pidfile} -o ${logfile} ${procname}"

run_rc_command "$1"
