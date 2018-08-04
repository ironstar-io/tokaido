package ssltmpl

// ProxyCaConfig - proxy_ca-config.json
func ProxyCaConfig() []byte {
	return []byte(`{
	"signing": {
		"default": {
			"expiry": "87600h"
		},
		"profiles": {
			"tokaido": {
				"usages": [
					"signing",
					"key encipherment",
					"server auth",
					"client auth"
				],
				"expiry": "87600h"
			}
		}
	}
}
	`)
}

// ProxyCaCsr - proxy_ca-csr.json
func ProxyCaCsr() []byte {
	return []byte(`{
	"CN": "Tokaido Local Certificate Authority",
	"key": {
		"algo": "rsa",
		"size": 2048
	},
	"names": [
		{
			"C": "JP",
			"L": "Osaka",
			"O": "Tokaido Local Development Environment",
			"OU": "tokaido.io",
			"ST": "Kansai"
		}
	]
}
	`)
}

// ProxyCsr - tokaido-csr.json
func ProxyCsr() []byte {
	return []byte(`{
	"CN": "tokaido.local",
	"key": {
		"algo": "rsa",
		"size": 2048
	},
	"names": [{
		"C": "JP",
		"L": "Osaka",
		"O": "Tokaido Local Development Environment",
		"OU": "tokaido.io",
		"ST": "Kansai"
	}],
	"hosts": [
		"*.tokaido.local",
		"tokaido.local"
	]
}
	`)
}
