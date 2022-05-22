# Setup fzf
# ---------
if [[ ! "$PATH" == */home/tok/.fzf/bin* ]]; then
  export PATH="$PATH:/home/tok/.fzf/bin"
fi

# Auto-completion
# ---------------
[[ $- == *i* ]] && source "/home/tok/.fzf/shell/completion.bash" 2> /dev/null

# Key bindings
# ------------
source "/home/tok/.fzf/shell/key-bindings.bash"

