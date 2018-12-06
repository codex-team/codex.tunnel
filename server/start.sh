#!/bin/bash

/usr/sbin/service nginx start
/usr/sbin/service ssh start
/usr/bin/python3 /home/server.py
tail -f /dev/null