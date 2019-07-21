package testcafe

import "github.com/ironstar-io/tokaido/system/fs"

func generateEnvFile(path, url string) {
	fs.TouchByteArray(path, buildEnvFileContents(url))
}

func buildEnvFileContents(url string) []byte {
	return []byte(`{
	"system": {
		"domain": "` + url + `"
	},
	"node": {
		"create": {
			"path": "/node/add"
		}
	},
	"users": {
		"admin": {
			"username": "testcafe_admin",
			"password": "testcafe_admin",
			"role": "administrator"
		},
		"editor": {
			"username": "testcafe_editor",
			"password": "testcafe_editor",
			"role": "editor"
		},
		"authenticated_user": {
			"username": "testcafe_user",
			"password": "testcafe_user"
		}
	},
	"user": {
		"login": {
			"path": "/user/login",
			"selectors": {
				"username": "#edit-name",
				"password": "#edit-pass",
				"login_button": "form.user-login-form #edit-submit"
			}
		},
		"add": {
			"path": "/admin/people/create"
		}
	}
}
`)
}
