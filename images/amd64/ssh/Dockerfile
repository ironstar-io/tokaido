ARG TOK_VERSION
ARG PHP_VERSION

FROM --platform=linux/amd64 tokaido/php${PHP_VERSION}:${TOK_VERSION}-amd64
ENV DEBIAN_FRONTEND noninteractive
ENV SHELL /bin/bash
USER root
COPY configs/.bash_profile /etc/skel/.bash_profile
COPY configs/.bashrc /etc/skel/.bashrc
COPY configs/sudoers /etc/sudoers
COPY configs/sshd_config /etc/ssh/sshd_config
COPY configs/motd.sh /etc/motd.sh
COPY start.sh /app/bin/start.sh
COPY configs/fzf.bash /home/app/.fzf.bash
COPY ZscalerRootCertificate-2048-SHA256.crt /usr/local/share/ca-certificates/ZscalerRootCertificate-2048-SHA256.crt
RUN apt-get update  \
 && apt-get install -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" \
      openssh-server \
      multitail \
      sendmail \
      locales-all \
      bash-completion \
      sudo \
      jq \
      nano \
 && /usr/bin/ssh-keygen -A  \
 && chmod 600 /etc/skel/.bash* \
 && cp /etc/skel/.bash_profile /home/app/.bash_profile \
 && cp /etc/skel/.bashrc /home/app/.bashrc \
 && chown app:app /home/app/.bash* \
 && chmod 600 /home/app/.bash* \
 && mkdir /run/sshd  \
 && rm -rf /tmp/*

RUN git clone --depth 1 https://github.com/junegunn/fzf.git /home/app/.fzf \
 && cd /home/app && ./.fzf/install \
 && chown app:app /home/app/.fzf -R \
 && chown app:app /home/app/.fzf.bash \
 && echo "" >> /home/app/.bash_profile \
 && echo "source /home/app/.fzf.bash" >> /home/app/.bash_profile \
 && find /home/app -type d -a -print0 | xargs -0 chmod 2700 \
 && find /home/app -type d -iname ".*" -a -print0 | xargs -0 chmod 2700 \
 && find /home/app -type f -a -print0 | xargs -0 chmod 600 \
 && find /home/app -type f -iname ".*" -a -print0 | xargs -0 chmod 600 \
 && usermod -s /bin/bash app \
 && ln -s /usr/share/bash-completion/completions/git /etc/bash_completion.d/git-completion.bash \
 && chmod 700 /home/app/.fzf/bin/fzf

RUN chmod 700 /app/bin/start.sh \
    && chmod 755 /etc/motd.sh  \
    && rm /etc/motd

USER root
WORKDIR /app/site
EXPOSE 22
ENTRYPOINT ["/app/bin/start.sh"]
