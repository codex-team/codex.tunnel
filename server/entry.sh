echo $SSH_ORIGINAL_COMMAND > /tmp/log.txt
echo $XDG_SESSION_ID >> /tmp/log.txt
echo $SSH_TTY >> /tmp/log.txt

ssh_tty=$SSH_TTY | cut -d'/' -f4

read HOST
read PORT
echo "$HOST:$PORT" > /tmp/params.txt

#ssh_pid1=$(cat /var/log/auth.log | grep $XDG_SESSION_ID -B1 | head -n 1)
#ssh_pid1=$(echo $ssh_pid1 | cut -d' ' -f5 | cut -d'[' -f2 | cut -d']' -f1)
#ssh_pid2=$(pstree -p $ssh_pid1 | awk -F"-+" '{ print $2 }' | cut -d'(' -f2 | cut -d')' -f1)
#ssh_pid2=$(pstree -p $ssh_pid1)
#ssh_pid2=$(echo $ssh_pid2 | awk -F"-+" '{ print $2 }' | cut -d'(' -f2 | cut -d')' -f1)
#ssh_port=$(netstat -tulpn | grep $ssh_pid2 | grep -E "^tcp\s+.*?[0-9]+/sshd: ladmin@" | awk -F"[ ]+" '{print $4}' | cut -d: -f2)

#echo $ssh_pid1 >> /tmp/nginx.txt
#echo $ssh_pid2 >> /tmp/nginx.txt
#echo $ssh_port >> /tmp/nginx.txt

#pstree -p $pid1 | awk -F"-+" '{ print $2 }' >> /tmp/nginx.txt

#echo $port >> /tmp/nginx.txt

sed "s/{HOST}/$HOST/" /home/sample.config | sed "s/{PORT}/$PORT/" > /etc/nginx/sites-enabled/$HOST.dev.conf
#rm /etc/nginx/sites-enabled/$SSH_ORIGINAL_COMMAND.dev.conf 2> /dev/null
#ln -s /etc/nginx/sites-available/$SSH_ORIGINAL_COMMAND.dev.conf /etc/nginx/sites-enabled
service nginx reload

#echo "$SSH_ORIGINAL_COMMAND:$ssh_port" >> /home/active.txt

#tty=`env | grep SSH_TTY | cut -d/ -f4`
#netstat -tulpn | grep -E "^tcp\s+.*?[0-9]+/[0-9]+" | awk -F"[ ]+" '{print $7}' > /tmp/netstat.txt
#netstat -tulpn | grep -E "^tcp\s+.*?[0-9]+/[0-9]+" | awk -F"[ ]+" '{print $7}' | grep $tty > /tmp/result.txt
#id
while true; do sleep 1; done