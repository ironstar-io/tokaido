package sshtmpl

var TokConfTmplStr = `Host {{.ProjectName}}.tok
    HostName localhost
    Port {{.DrushPort}}
    User app
    IdentityFile ~/.ssh/tok_ssh.key
    IdentitiesOnly yes
    ForwardAgent yes
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null

`
