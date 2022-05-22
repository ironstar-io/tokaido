FROM --platform=linux/amd64 nginx:stable
ENV DEBIAN_FRONTEND noninteractive
LABEL io.ironstar.platform=tokaido
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# hadolint ignore=DL3008
RUN apt-get update  \
    && apt-get install -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" --no-install-recommends \
        apt-transport-https \
        lsb-release \
        ca-certificates \
		bash \
        wget \
        python \
        python-pip \
        rsync \
        vim \
        curl \
        telnet \
	    netcat \
	&& mkdir -p /tokaido/config/nginx/conf.d/sites /tokaido/config/nginx/conf.d/redirects \
	&& mkdir -p /tokaido/logs/nginx \
	&& mkdir -p /tokaido/site/docroot \
	&& mkdir -p /var/cache/nginx \
	&& groupadd -g 1001 web  \
	&& userdel nginx \
	&& useradd -s /sbin/nologin -g web -u 1001 tok  \
    && useradd -s /sbin/nologin -d /var/cache/nginx -g web -u 1002 nginx  \
	&& chown tok:web /tokaido -R \
	&& chown nginx:web /tokaido/logs/nginx \
	&& chown nginx:web /var/cache/nginx \
	&& chmod 740 /var/cache/nginx \
	&& touch /var/run/nginx.pid \
	&& chown nginx:web /var/run/nginx.pid \
	&& curl -sSLo /usr/local/bin/yq https://github.com/mikefarah/yq/releases/download/v4.25.1/yq_linux_amd64 \
	&& chmod 777 /usr/local/bin/yq \
	&& apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY config/nginx.conf /tokaido/config/nginx/nginx.conf
COPY config/host.conf /tokaido/config/nginx/host.conf
COPY config/redirects.conf /tokaido/config/nginx/redirects.conf
COPY config/mimetypes.conf /tokaido/config/nginx/mimetypes.conf
COPY config/additional.conf /tokaido/config/nginx/additional.conf
COPY config/security.conf /tokaido/config/nginx/security.conf
COPY config/phpinfo.php /tokaido/site/docroot/phpinfo.php
COPY config/fastcgi_params /etc/nginx/fastcgi_params
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chown nginx:web /tokaido/config/nginx -R \
	&& chmod 600 /tokaido/config/nginx/*.conf \
	&& chmod 700 /tokaido/config/nginx \
	&& chown nginx:web /usr/local/bin/entrypoint.sh \
	&& chmod 750 /usr/local/bin/entrypoint.sh

EXPOSE 8082

STOPSIGNAL SIGTERM

USER nginx
VOLUME /tokaido/site

CMD ["/usr/local/bin/entrypoint.sh"]
