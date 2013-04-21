scp forum.linux james@carbocation.com:/data/bin/upload.forum.linux

ssh -n -f james@carbocation.com "sh -c 'cd /data/bin/; killall forum.linux; mv upload.forum.linux forum.linux;  nohup ./forum.linux > /dev/null 2>&1 &'"
