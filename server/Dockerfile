FROM ubuntu:18.04
LABEL Maintainer="Alexander Menshchikov @nostr <ntpcp@yandex.ru>"

RUN apt-get update
RUN apt-get install -y nginx g++ ssh python3

RUN mkdir /var/run/sshd
RUN useradd -ms /bin/bash codex
RUN mkdir /home/codex/.ssh

COPY authorized_keys /home/codex/.ssh/authorized_keys
COPY entry.cpp /home/codex/entry.cpp

RUN g++ /home/codex/entry.cpp -o /home/codex/entry
RUN rm /home/codex/entry.cpp
RUN chmod 4755 /home/codex/entry
RUN apt-get remove -y g++ --purge
RUN apt-get -y autoremove
RUN chmod 777 /home/codex/.ssh/authorized_keys

COPY entry.py /home/entry.py
#COPY default.conf /etc/nginx/sites-available/default
RUN chmod +x /home/entry.py

RUN apt-get install -y python3-pip
RUN pip3 install flask
COPY config /home/config
COPY server.py /home/server.py

RUN sed -i 's/#MaxSessions 10/MaxSessions 1000/' /etc/ssh/sshd_config
RUN sed -i 's/#StrictModes yes/StrictModes no/' /etc/ssh/sshd_config

COPY default.conf /etc/nginx/sites-available/default

EXPOSE 22

COPY start.sh /root/start.sh
RUN chmod +x /root/start.sh
WORKDIR /root
CMD ["./start.sh", "-D"]