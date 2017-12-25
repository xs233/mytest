package log

import (
	"os"
	//"path/filepath"
)

func logFileExploreToDelete(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	os.Remove(path)
	return nil
}

func init() {
	//filepath.Walk(cpath, logFileExploreToDelete)
}
