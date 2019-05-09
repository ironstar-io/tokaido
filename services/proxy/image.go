package proxy

// PullImages - Pull all images in compose file
func PullImages() {
	ComposeStdout("pull")
}
