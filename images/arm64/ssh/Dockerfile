ARG TOK_VERSION
ARG PHP_VERSION

USER root
FROM --platform=linux/arm64 tokaido/php${PHP_VERSION}:${TOK_VERSION}-arm64
ENV DEBIAN_FRONTEND noninteractive
ENV SHELL /bin/bash
COPY configs/.bash_profile /etc/skel/.bash_profile
COPY configs/.bashrc /etc/skel/.bashrc
COPY configs/sudoers /etc/sudoers
COPY configs/sshd_config /etc/ssh/sshd_config
COPY configs/motd.sh /etc/motd.sh
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY configs/fzf.bash /home/tok/.fzf.bash
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
 && cp /etc/skel/.bash_profile /home/tok/.bash_profile \
 && cp /etc/skel/.bashrc /home/tok/.bashrc \
 && chown tok:web /home/tok/.bash* \
 && chmod 600 /home/tok/.bash* \
 && mkdir /run/sshd  \
 && rm -rf /tmp/*

RUN git clone --depth 1 https://github.com/junegunn/fzf.git /home/tok/.fzf \
 && cd /home/tok && ./.fzf/install \
 && chown tok:web /home/tok/.fzf -R \
 && chown tok:web /home/tok/.fzf.bash \
 && echo "" >> /home/tok/.bash_profile \
 && echo "source /home/tok/.fzf.bash" >> /home/tok/.bash_profile \
 && find /home/tok -type d -a -print0 | xargs -0 chmod 2700 \
 && find /home/tok -type d -iname ".*" -a -print0 | xargs -0 chmod 2700 \
 && find /home/tok -type f -a -print0 | xargs -0 chmod 600 \
 && find /home/tok -type f -iname ".*" -a -print0 | xargs -0 chmod 600 \
 && ln -s /usr/share/bash-completion/completions/git /etc/bash_completion.d/git-completion.bash \
 && chmod 700 /home/tok/.fzf/bin/fzf

RUN chmod 700 /usr/local/bin/entrypoint.sh \
    && chmod 755 /etc/motd.sh  \
    && rm /etc/motd

USER app
WORKDIR /app/site
EXPOSE 22
CMD ["/usr/local/bin/entrypoint.sh"]
