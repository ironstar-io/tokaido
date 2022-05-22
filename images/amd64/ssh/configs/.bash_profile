# Vars we want to see in both interactive and non-interactive shells
drupal_root=${DRUPAL_ROOT:-docroot}
export PATH=$PATH:$HOME/.composer/vendor/bin
export DRUSH_LAUNCHER_FALLBACK=/usr/local/drush/global/bin/drush
source /tokaido/config/.env

# If not running interactively, don't do anything
case $- in
    *i*) ;;
      *) return;;
esac

# Set umask so drupal file uploads are accessible by the web group
umask 002

# don't put duplicate lines or lines starting with space in the history.
# See bash(1) for more options
HISTCONTROL=ignoreboth

# append to the history file, don't overwrite it
shopt -s histappend

# for setting history length see HISTSIZE and HISTFILESIZE in bash(1)
HISTSIZE=10000
HISTFILESIZE=20000
# check the window size after each command and, if necessary,
# update the values of LINES and COLUMNS.
shopt -s checkwinsize

# If set, the pattern "**" used in a pathname expansion context will
# match all files and zero or more directories and subdirectories.
#shopt -s globstar

# make less more friendly for non-text input files, see lesspipe(1)
[ -x /usr/bin/lesspipe ] && eval "$(SHELL=/bin/sh lesspipe)"

# set variable identifying the chroot you work in (used in the prompt below)
if [ -z "${debian_chroot:-}" ] && [ -r /etc/debian_chroot ]; then
    debian_chroot=$(cat /etc/debian_chroot)
fi

# set a fancy prompt (non-color, unless we know we "want" color)
case "$TERM" in
    xterm-color|*-256color) color_prompt=yes;;
esac

# uncomment for a colored prompt, if the terminal has the capability; turned
# off by default to not distract the user: the focus in a terminal window
# should be on the output of commands, not on the prompt
#force_color_prompt=yes

if [ -n "$force_color_prompt" ]; then
    if [ -x /usr/bin/tput ] && tput setaf 1 >&/dev/null; then
        # We have color support; assume it's compliant with Ecma-48
        # (ISO/IEC-6429). (Lack of such support is extremely rare, and such
        # a case would tend to support setf rather than setaf.)
        color_prompt=yes
    else
        color_prompt=
    fi
fi

if [ "$color_prompt" = yes ]; then
    if [[ $APP_ENV == "local" ]]; then
        PS1='🚅 \[\e[1;32m\]\[\e[1mLOCAL\[\e[33m\]\[\e[m\] \[\e[36m\]${PROJECT_NAME}\[\e[m\] \[\e[33m\]\t\[\e[m\] \[\e[34m\]\w\[\e[m\]\n\[\e[32m\]\$\[\e[m\] '
    elif [[ $APP_ENV == "dev" ]]; then
        PS1='\[\e[1;33m\]\[\e[1mDEV\[\e[33m\]\[\e[m\] \[\e[36m\]${PROJECT_NAME}\[\e[m\] \[\e[33m\]\t\[\e[m\] \[\e[34m\]\w\[\e[m\]\n\[\e[32m\]\$\[\e[m\] '
    elif [[ $APP_ENV == "stg" ]]; then
        PS1='\[\e[1;33m\]\[\e[1mSTG\[\e[33m\]\[\e[m\] \[\e[36m\]${PROJECT_NAME}\[\e[m\] \[\e[33m\]\t\[\e[m\] \[\e[34m\]\w\[\e[m\]\n\[\e[32m\]\$\[\e[m\] '
    elif [[ $APP_ENV == "prd" ]]; then
        PS1='☢️ \[\e[1;31m\]\[\e[1mPRD\[\e[33m\]\[\e[m\] \[\e[36m\]${PROJECT_NAME}\[\e[m\] \[\e[33m\]\t\[\e[m\] \[\e[34m\]\w\[\e[m\]\n\[\e[32m\]\$\[\e[m\] '
    else
        PS1='\[\e[36m\]${PROJECT_NAME}\[\e[m\] \[\e[33m\]\t\[\e[m\] \[\e[34m\]\w\[\e[m\]\n\[\e[32m\]\$\[\e[m\] '
    fi
else
    PS1='${debian_chroot:+($debian_chroot)}\u@\h:\w\$ '
fi
unset color_prompt force_color_prompt

# If this is an xterm set the title to user@host:dir
case "$TERM" in
xterm*|rxvt*)
    PS1="\[\e]0;${debian_chroot:+($debian_chroot)}\u@\h: \w\a\]$PS1"
    ;;
*)
    ;;
esac

# enable color support of ls and also add handy aliases
if [ -x /usr/bin/dircolors ]; then
    test -r ~/.dircolors && eval "$(dircolors -b ~/.dircolors)" || eval "$(dircolors -b)"
    alias ls='ls --color=auto'
    #alias dir='dir --color=auto'
    #alias vdir='vdir --color=auto'

    alias grep='grep --color=auto'
    alias fgrep='fgrep --color=auto'
    alias egrep='egrep --color=auto'
fi

# colored GCC warnings and errors
#export GCC_COLORS='error=01;31:warning=01;35:note=01;36:caret=01;32:locus=01:quote=01'

# some more ls aliases
alias ll='ls -alF'
alias la='ls -A'
alias l='ls -CF'

# Add an "alert" alias for long running commands.  Use like so:
#   sleep 10; alert
alias alert='notify-send --urgency=low -i "$([ $? = 0 ] && echo terminal || echo error)" "$(history|tail -n1|sed -e '\''s/^\s*[0-9]\+\s*//;s/[;&|]\s*alert$//'\'')"'

# Alias definitions.
# You may want to put all your additions into a separate file like
# ~/.bash_aliases, instead of adding them here directly.
# See /usr/share/doc/bash-doc/examples in the bash-doc package.

if [ -f ~/.bash_aliases ]; then
    . ~/.bash_aliases
fi

# enable programmable completion features (you don't need to enable
# this, if it's already enabled in /etc/bash.bashrc and /etc/profile
# sources /etc/bash.bashrc).
if ! shopt -oq posix; then
  if [ -f /usr/share/bash-completion/bash_completion ]; then
    . /usr/share/bash-completion/bash_completion
  elif [ -f /etc/bash_completion ]; then
    . /etc/bash_completion
  fi
fi

cd /tokaido/site/${drupal_root}

echo "Checking your site's database connection. This might take a moment..."
echo "You can press CTRL+C to continue at any time"

bash /etc/motd.sh
