module github.com/ironstar-io/tokaido-installer

go 1.13

replace github.com/ironstar-io/tokaido v0.0.0 => ../

require (
	github.com/docker/docker v0.0.0-20190509004234-994007dd89fd
	github.com/ironstar-io/tokaido v0.0.0
	github.com/logrusorgru/aurora v0.0.0-20191017060258-dc85c304c434
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.5.0
)
