# Setup fzf
# ---------
if [[ ! "$PATH" == */home/app/.fzf/bin* ]]; then
  export PATH="$PATH:/home/app/.fzf/bin"
fi

# Auto-completion
# ---------------
[[ $- == *i* ]] && source "/home/app/.fzf/shell/completion.bash" 2> /dev/null

# Key bindings
# ------------
source "/home/app/.fzf/shell/key-bindings.bash"
