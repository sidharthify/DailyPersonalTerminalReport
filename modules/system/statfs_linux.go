//go:build linux

package system

import "syscall"

func statfs(path string) (total, used uint64, err error) {
	var s syscall.Statfs_t
	if err = syscall.Statfs(path, &s); err != nil {
		return
	}
	total = s.Blocks * uint64(s.Bsize)
	avail := s.Bavail * uint64(s.Bsize)
	used = total - avail
	return
}
