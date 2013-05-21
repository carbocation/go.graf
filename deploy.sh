scp forum.linux james@carbocation.com:/data/bin/upload.forum.linux
ssh james@carbocation.com "sh -c 'cd /data/bin/; sudo service go-askbitcoin-upstart stop; mv upload.forum.linux forum.linux; sudo service go-askbitcoin-upstart start'"


#scp forum.linux james@carbocation.com:/data/bin/upload.forum.linux
#ssh -n -f james@carbocation.com "sh -c 'cd /data/bin/; killall forum.linux; mv upload.forum.linux forum.linux;GOMAXPROCS=8  nohup ./forum.linux > /dev/null 2>&1 &'"

