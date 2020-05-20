#!/bin/sh

# PROVIDE: bigtechfaild
# REQUIRE: DAEMON
# KEYWORD: shutdown

. /etc/rc.subr

name="bigtechfaild"
rcvar="${name}_enable"

load_rc_config ${name}

: ${bigtechfaild_enable:=NO}
: ${bigtechfaild_rootdir:="/var/www/${name}"}
: ${bigtechfaild_logfile="/var/log/${name}.log"}
: ${bigtechfaild_user="root"}
: ${bigtechfaild_group="wheel"}

pidfile="/var/run/${name}.pid"
procname="/usr/local/bin/${name}"

command="/usr/sbin/daemon"
command_args="-p ${pidfile} -o ${bigtechfaild_logfile} ${procname} -prod -dir=${bigtechfaild_rootdir}"

start_precmd()
{
	if [ ! -e "${pidfile}" ]; then
		install -o "${bigtechfaild_user}" -g "${bigtechfaild_group}" "/dev/null" "${pidfile}"
	fi

	if [ ! -e "${bigtechfaild_logfile}" ]; then
		install -o "${bigtechfaild_user}" -g "${bigtechfaild_group}" "/dev/null" "${bigtechfaild_logfile}"
	fi
}

run_rc_command "$1"
