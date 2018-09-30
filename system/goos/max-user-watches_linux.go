package goos

import (
	"fmt"
	"strconv"

	"github.com/ironstar-io/tokaido/utils"
)

// CheckMaxUserWatches checks if max_user_watches is sufficient for Tokaid to work properly
func CheckMaxUserWatches() {
	const maxUserWatchesRec = 524288

	m, _ := strconv.Atoi(utils.StdoutCmd("cat", "/proc/sys/fs/inotify/max_user_watches"))

	if m < maxUserWatchesRec {
		fmt.Printf(`
--- ❗ WARNING ❗ ---
Your systems's max_user_watches value [%d] is less than the recommended [%d]
This is known to cause issues problems with the Unison sync service on Linux.

Please see this article for instructions on increasing your max_user_watches
value: https://github.com/guard/listen/wiki/Increasing-the-amount-of-inotify-watchers
---
`, m, maxUserWatchesRec)
	}

}
