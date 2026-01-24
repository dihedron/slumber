package detect

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func HasActiveProcess(pattern string) (bool, error) {
	files, err := os.ReadDir("/proc")
	if err != nil {
		slog.Error("failed to read /proc", "error", err)
		return false, err
	}

	for _, f := range files {
		if !f.IsDir() || !isPID(f.Name()) {
			slog.Warn("ignoring non-pid directory", "name", f.Name())
			continue
		}
		slog.Debug("checking process", "pid", f.Name())
		filename := path.Clean(filepath.Join("/proc", f.Name(), "cmdline"))
		data, err := os.ReadFile(filename)
		if err != nil {
			slog.Warn("failed to read process command line", "pid", f.Name(), "error", err)
			continue
		}

		cmdline := strings.Replace(string(data), "\x00", " ", -1)
		re := regexp.MustCompile(pattern)
		if re.MatchString(cmdline) {
			slog.Debug("found process", "filename", filename, "cmdline", cmdline)
			return true, nil
		}
	}

	slog.Debug("no active process found")
	return false, nil
}

func isPID(name string) bool {
	_, err := strconv.Atoi(name)
	return err == nil
}
