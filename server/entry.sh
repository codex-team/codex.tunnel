read HOST
read PORT
sed "s/{HOST}/$HOST/" /home/sample.config | sed "s/{PORT}/$PORT/" > /etc/nginx/sites-enabled/$HOST.dev.conf
service nginx restart
while true; do sleep 5; done