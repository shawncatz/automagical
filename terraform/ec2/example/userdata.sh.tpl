#!/bin/bash -ex

exec >/root/userdata.log 2>&1

HOST="${hostname}"

# remove trailing dot
echo "$${HOST%?}" > /etc/hostname
hostname -F /etc/hostname
