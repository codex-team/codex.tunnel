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

COPY entry.py /home/entry.py
RUN chmod +x /home/entry.py

RUN rm /etc/nginx/sites-enabled/default
RUN service nginx restart

RUN sed -i 's/#MaxSessions 10/MaxSessions 1000/' /etc/ssh/sshd_config

EXPOSE 22

CMD ["/usr/sbin/sshd", "-D"]