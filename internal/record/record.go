// Package record defines the metadata per file shared accros package
package record

type Record struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Mode   uint32 `json:"mode"`
	UID    uint32 `json:"uid"`
	GID    uint32 `json:"gid"`
	Size   int64  `json:"size"`
	MTime  int64  `json:"mtime"`
	Inode  uint64 `json:"inode"`
}
