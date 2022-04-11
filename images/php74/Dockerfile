FROM php:7.4-fpm
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
 && docker-php-ext-configure imap --with-kerberos --with-imap-ssl \
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
 && docker-php-source delete
