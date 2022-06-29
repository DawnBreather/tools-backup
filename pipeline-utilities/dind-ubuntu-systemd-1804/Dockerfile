FROM ubuntu:18.04

# ENVIRONMENT VARIABLES
ENV container=docker \
    LANGUAGE=en_US.UTF-8 \
    LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8 \
    TERM=xterm \
    DEBIAN_FRONTEND="noninteractive"

# APT INITIALIZATION
RUN apt update

# DIND IMPLEMENTATION
USER root
RUN apt-get -y install vim git wget openssh-server net-tools iptables

RUN wget https://download.docker.com/linux/static/edge/x86_64/docker-18.06.0-ce.tgz
RUN tar xvfz docker-18.06.0-ce.tgz
RUN  cp ./docker/* /usr/bin/

# SYSTEMD IMPLEMENTATION
RUN find /etc/systemd/system \
    /lib/systemd/system \
    -path '*.wants/*' \
    -not -name '*journald*' \
    -not -name '*systemd-tmpfiles*' \
    -not -name '*systemd-user-sessions*' \
    -print0 | xargs -0 rm -vf

RUN apt update && \
    INSTALL_PKGS="python sudo bash apt-utils locales iproute2 locales ca-certificates dbus gnupg systemd iptables net-tools iputils-ping" && \
    apt install -y $INSTALL_PKGS && \
    #apt-get clean && \
    #localedef -f UTF-8 -i en_US en_US.UTF-8 \
    apt install -y python sudo bash ca-certificates gcc python-dev build-essential python-pip libssl-dev libffi-dev git \
    && rm -rf /var/lib/apt/lists/*

RUN cp /bin/true /sbin/agetty

STOPSIGNAL SIGRTMIN+3


# MOLECULE IMPLEMENTATION


RUN pip install --upgrade pip
RUN pip install docker-py pyOpenSSL pytest PyYAML

# DOCKERFILE
COPY Dockerfile /Dockerfile

# WRAPDOCKER
COPY wrapdocker /bin/wrapdocker
COPY wrapdocker.service /etc/systemd/system/wrapdocker.service
RUN systemctl enable wrapdocker.service

# VOLUMES
VOLUME [ "/var/lib/docker" ]
VOLUME [ "/sys/fs/cgroup" ]

RUN localedef -f UTF-8 -i en_US en_US.UTF-8

# COMMAND
CMD ["/bin/bash", "-c", "echo root:$ROOT_PASSWORD | chpasswd root && /sbin/init"]

