ARG TOK_VERSION

FROM --platform=linux/amd64 php:7.4-fpm AS php
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        libbz2-dev \
        zlib1g-dev \
        libpng-dev \
        libgmp-dev \
        libc-client-dev \
        libkrb5-dev \
        libicu-dev \
        libldap2-dev \
        libpq-dev \
        libpspell-dev \
		libxml2-dev \
        libtidy-dev \
        libxslt-dev \
        libmagickwand-dev \
		libcurl4-openssl-dev \
        libgpgme11-dev \
        libyaml-dev \
        libmemcached-dev \
        libzip-dev \
        libwebp-dev \
 && docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp \
 && docker-php-ext-configure imap --with-kerberos --with-imap-ssl \
 && pecl update-channels \
 && printf "\n" | pecl install raphf \
 && docker-php-ext-enable raphf \
 && printf "\n" | pecl install imagick \
 && docker-php-ext-enable imagick \
 && printf "\n" | pecl install propro \
 && docker-php-ext-enable propro \
 && printf "\n" | pecl install pecl_http-3.2.4 \
 && docker-php-ext-enable http \
 && printf "\n" | pecl install gnupg \
 && docker-php-ext-enable gnupg \
 && printf "\n" | pecl install yaml \
 && docker-php-ext-enable yaml \
 && printf "\n" | pecl install memcached \
 && docker-php-ext-enable memcached \
 && printf "\n" | pecl install redis \
 && docker-php-ext-enable redis \
 && printf "\n" | pecl install oauth \
 && docker-php-ext-enable oauth \
 && printf "\n" | pecl install igbinary \
 && docker-php-ext-enable igbinary \
 && printf "\n" | pecl install apcu \
 && docker-php-ext-enable apcu \
 && printf "\n" | pecl install uploadprogress \
 && docker-php-ext-enable uploadprogress \
 && printf "\n" | pecl install zip \
 && docker-php-ext-enable zip \
 && printf "\n" | pecl install xdebug \
 && rm -rf /tmp/pear ~/.pearrc \
 && docker-php-ext-install \
        bcmath \
        bz2 \
        calendar \
        dba \
        exif \
        gd \
        gettext \
        gmp \
        imap \
        intl \
        ldap \
        mysqli \
        pcntl \
        pdo_mysql \
        pdo_pgsql \
        pgsql \
        pspell \
        shmop \
        soap \
        sysvmsg \
        sysvsem \
        sysvshm \
        tidy \
        xsl \
 && docker-php-source delete \
 && echo $PHP_VERSION > /tmp/version.txt

FROM --platform=linux/amd64 tokaido/base:${TOK_VERSION}-amd64
ARG TOK_VERSION
ARG LIBRARY_PATH
LABEL io.ironstar.platform=tokaido
LABEL maintainers="tokaido@ironstar.io"
LABEL security="security@ironstar.io"
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ARG DEBIAN_FRONTEND=noninteractive
WORKDIR /tmp
COPY config/php.ini-template /app/config/php/php.ini-template
COPY config/php.ini-template /app/config/php/php.ini
COPY config/www.conf-template /app/config/php/www.conf-template
COPY config/www.conf-template /app/config/php/www.conf
COPY --from=php /usr/local/bin/phar /usr/local/bin/php
COPY --from=php /usr/local/bin/php /usr/local/bin/php
COPY --from=php /usr/local/sbin/php /usr/local/sbin/php
COPY --from=php /usr/local/php/php /usr/local/php/php
COPY --from=php /usr/local/lib/php /usr/local/lib/php
COPY --from=php /usr/lib /usr/lib
COPY --from=php /tmp/version.txt /app/config/php/version.txt
COPY --from=php /usr/local/etc/php/conf.d /app/config/php/conf.d

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        libjpeg-dev \
        libxrender1 \
        libfreetype6 \
        libxext6 \
        libx11-6 \
        libpng-dev \
        fontconfig \
        xfonts-base \
        xfonts-75dpi \
        postfix \
        mariadb-client-core-10.5 \
        msmtp \
        qpdf \
        poppler-utils \
        swaks \
        libnss3 \
        libxss1 \
        libasound2 \
        libatk-bridge2.0-0 \
        libgtk-3-0 \
        libgbm-dev \
        libatk1.0-0 

RUN mkdir -p /usr/local/etc/php/ \
 && ln -s /app/config/php/conf.d /usr/local/etc/php/conf.d

RUN curl -s https://getcomposer.org/installer > composer-setup.php && /usr/local/bin/php composer-setup.php && mv composer.phar /usr/local/bin/composer && rm composer-setup.php  \
    && su - app -s /bin/bash -c "/usr/local/bin/composer global require \"drush/drush\""  \
    && curl -sSLo /usr/local/bin/drush https://github.com/drush-ops/drush-launcher/releases/download/0.10.1/drush.phar \
    && chmod ug+x /usr/local/bin/drush \
    && chown app:app /usr/local/bin/drush \
    && ln -sf /usr/bin/msmtp /usr/sbin/sendmail \
    && ln -sf /app/config/msmtp/msmtprc /etc/msmtprc

RUN ep /app/config/php/php.ini \
 && ep /app/config/php/www.conf

RUN curl -sSLo /tmp/fnm.zip https://github.com/Schniz/fnm/releases/download/v1.31.0/fnm-linux.zip \
 && unzip /tmp/fnm.zip -d /usr/local/bin/ \
 && chmod +x /usr/local/bin/fnm \
 && chown app:app /usr/local/bin/fnm \
 && chmod 770 /usr/local/bin/fnm \
 && mkdir -p /home/app/.fnm \
 && chmod 2770 /home/app/.fnm 

COPY config/.bashrc /home/app/.bashrc
COPY config/.bash_profile /home/app/.bash_profile

RUN apt-get clean \
 && rm -rf /tmp/* \
 && chown app:app /home/app -R \
 && chmod o-rwx /home/app -R 

USER app
WORKDIR /app/site
RUN eval $(fnm env) \
 && fnm install v12 \
 && fnm install v14 \
 && fnm install v16 \
 && fnm default system #disables node
