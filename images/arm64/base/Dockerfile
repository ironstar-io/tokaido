ARG TOK_VERSION
FROM debian:bullseye-slim
LABEL io.ironstar.platform=tokaido
LABEL maintainers="tokaido@ironstar.io"
LABEL security="security@ironstar.io"
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ARG DEBIAN_FRONTEND=noninteractive

# Install baseline packages
# hadolint ignore=DL3008
ENV TOK_VERSION=${TOK_VERSION}
WORKDIR /tmp
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        apt-transport-https \
        build-essential \
        ca-certificates \
        curl \
        dirmngr \
        gawk\
        git \
        gnupg2 \
        iproute2 \
        iputils-ping \
        less \
        locales-all \
        lsb-release \
        moreutils \
        net-tools \
        parallel \
        pigz \
        python3 \
        python3-pip \
        python3-setuptools \
        rsync \
        sed \
        software-properties-common \
        traceroute \
        unzip \
        vim \
        wget

RUN ln -sf /usr/bin/python3 /usr/bin/python \
    && curl -sSLo ep.tgz https://github.com/kreuzwerker/envplate/releases/download/v1.0.2/envplate_1.0.2_Linux_arm64.tar.gz \
    && tar xzf ep.tgz \
    && mv /tmp/envplate /bin/ep \
    && chmod +x /bin/ep \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Configure user accounts and default permissions
RUN mkdir /etc/skel/.ssh \
    && find /etc/skel -type d -a -print0 | xargs -0 chmod 700 \
    && find /etc/skel -type d -name ".*" -a -print0 | xargs -0 chmod 700 \
    && find /etc/skel -type f -a -print0 | xargs -0 chmod 600 \
    && find /etc/skel -type f -name ".*" -a -print0 | xargs -0 chmod 600 \
    && echo "umask 002" > /etc/profile \
    && groupadd -g 1001 app \
    && useradd -m -s /bin/false -g app  -K UMASK=002 -u 1001 app \
    && usermod -p '*' app \
    && chown app:root /home/app -R \
    && find /home/ -type d -a -print0 | xargs -0 chmod 700 \
    && find /home/ -type f -a -print0 | xargs -0 chmod 600 \
    && chmod 711 /home

# Create directory structure to be used in child containers
RUN mkdir /app/ \
    && mkdir -p /app/logs \
    && mkdir -p /app/config \
    && mkdir -p /app/site \
    && mkdir -p /app/bin \
    && mkdir -p /app/cache/ \
    && chown app:app /app -R \
    && chmod 2771 /app/ -R \
    && ln -s /app /tokaido \
    && ln -s /app /joetsu

# Lock down the environment by ensuring that non-essential users and paths are sterilised
RUN userdel backup \
    && userdel daemon \
    && userdel games \
    && userdel gnats \
    && userdel irc \
    && userdel list \
    && userdel lp \
    && userdel mail \
    && userdel man \
    && userdel news \
    && userdel proxy \
    && userdel sync \
    && userdel uucp \
    && userdel www-data \
    && groupdel audio \
    && groupdel cdrom \
    && groupdel dialout \
    && groupdel dip \
    && groupdel disk \
    && groupdel fax \
    && groupdel floppy \
    && groupdel kmem \
    && groupdel operator \
    && groupdel plugdev \
    && groupdel sasl \
    && groupdel shadow \
    && groupdel src \
    && groupdel sudo \
    && groupdel tape \
    && groupdel users \
    && groupdel utmp \
    && groupdel video \
    && groupdel voice

# Add tini to provide a simple init service and signal forwarding
ADD https://github.com/krallin/tini/releases/download/v0.19.0/tini-amd64 /usr/bin/tini
ADD https://github.com/krallin/tini/releases/download/v0.19.0/tini-amd64.asc /tmp/tini.asc
RUN gpg --batch --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 595E85A6B1B4779EA4DAAEC70B588DFF0527A9B7 \
    && gpg --batch --verify /tmp/tini.asc /usr/bin/tini \
    && chmod 555 /usr/bin/tini \
    && rm /tmp/tini.asc

# Quality of life improvements
COPY --chmod=644 configs/inputrc /etc/inputrc
RUN echo "UTC" > /etc/timezone \
    && echo "stty werase undef" >> /etc/bash.bashrc \
    && echo "bind '\C-w:unix-filename-rubout'" >> /etc/bash.bashrc

# System cleanup and security hardening
RUN find /sbin -type f -print0 | xargs -0 chmod o-x \
    && find /usr/sbin -type f -print0 | xargs -0 chmod o-x \
    && find /var/lib/dpkg/ -type f -print0 | xargs -0 chmod o-x \
    && rm -rf /var/local /var/mail /var/tmp
