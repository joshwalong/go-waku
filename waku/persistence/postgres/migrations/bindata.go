// Code generated by go-bindata. DO NOT EDIT.
// sources:
// 1_messages.down.sql (124B)
// 1_messages.up.sql (452B)
// 2_messages_index.down.sql (60B)
// 2_messages_index.up.sql (226B)
// 3_rendezvous.down.sql (65B)
// 3_rendezvous.up.sql (181B)
// 4_signed_peer_record.down.sql (166B)
// 4_signed_peer_record.up.sql (178B)
// doc.go (74B)

package migrations

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __1_messagesDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\xf0\xf4\x73\x71\x8d\x50\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\xc8\x4d\x2d\x2e\x4e\x4c\x4f\x8d\x2f\x4e\xcd\x4b\x49\x2d\x0a\xc9\xcc\x4d\x2d\x2e\x49\xcc\x2d\xb0\xe6\xc2\xab\xba\x28\x35\x39\x35\xb3\x0c\x53\x7d\x88\xa3\x93\x8f\x2b\xa6\x7a\x6b\x2e\x40\x00\x00\x00\xff\xff\xc2\x48\x8c\x05\x7c\x00\x00\x00")

func _1_messagesDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__1_messagesDownSql,
		"1_messages.down.sql",
	)
}

func _1_messagesDownSql() (*asset, error) {
	bytes, err := _1_messagesDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1_messages.down.sql", size: 124, mode: os.FileMode(0664), modTime: time.Unix(1673044753, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xff, 0x4a, 0x8e, 0xa9, 0xd9, 0xa8, 0xa4, 0x73, 0x3a, 0x54, 0xe4, 0x35, 0xfd, 0xea, 0x87, 0x4c, 0xa, 0x5c, 0xc0, 0xc9, 0xe7, 0x8, 0x8c, 0x6f, 0x60, 0x9e, 0x54, 0x77, 0x59, 0xd0, 0x2b, 0xfe}}
	return a, nil
}

var __1_messagesUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x90\x41\x4f\x83\x40\x10\x85\xcf\xec\xaf\x98\x23\x24\x1c\xbc\x73\x5a\xda\x69\x33\x11\x17\xb3\x4c\x93\x72\x32\x14\x26\x66\x13\x59\x08\x4b\x1b\xfd\xf7\x46\xad\x4a\x5a\x35\x9e\xbf\x37\x6f\xde\x7b\x2b\x8b\x9a\x11\x58\xe7\x05\x02\x6d\xc0\x94\x0c\xb8\xa7\x8a\x2b\xe8\x25\x84\xe6\x51\x20\x56\x91\xeb\x20\xaf\x19\x75\xaa\xa2\x49\x5a\x71\x27\x99\xd8\xf5\x12\xe6\xa6\x1f\x21\xa7\x2d\x19\x7e\xbf\x34\xbb\xa2\x48\x55\x14\xc4\x77\x7f\x2b\xda\xc1\xcf\xe2\x67\x1e\x46\xd7\x7e\x58\x2f\xe9\x78\x3c\x84\xe3\xe1\x37\xd8\xbc\x3c\x0d\xcd\x77\xa0\x93\x4c\xc1\x0d\x1e\xc8\x30\x6e\xd1\x7e\x49\x61\x8d\x1b\xbd\x2b\x18\x6e\x52\x15\xad\x4a\x53\xb1\xd5\x6f\x29\xce\xb5\xc8\x77\xf2\x0c\xf7\x96\xee\xb4\xad\xe1\x16\x6b\x88\x5d\x97\xc2\xe2\x75\xa2\x92\x4c\xa9\xf3\x40\x64\xd6\xb8\xff\x79\xa0\x87\xcb\xba\xa5\xf9\x44\xf1\x05\x4a\xb2\xff\xf8\x5d\x4f\xbc\x70\xbc\x82\x49\xf6\x1a\x00\x00\xff\xff\xa0\x46\xcd\x13\xc4\x01\x00\x00")

func _1_messagesUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1_messagesUpSql,
		"1_messages.up.sql",
	)
}

func _1_messagesUpSql() (*asset, error) {
	bytes, err := _1_messagesUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1_messages.up.sql", size: 452, mode: os.FileMode(0664), modTime: time.Unix(1673044753, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xe4, 0x17, 0xde, 0xd4, 0x55, 0x47, 0x7f, 0x61, 0xe6, 0xbd, 0x2e, 0x89, 0xb5, 0x7, 0xe1, 0x31, 0x1b, 0xd3, 0x20, 0x3d, 0x3e, 0x68, 0x54, 0xfe, 0xd3, 0x62, 0x51, 0x87, 0x5f, 0xbf, 0x57, 0x64}}
	return a, nil
}

var __2_messages_indexDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\xf0\xf4\x73\x71\x8d\x50\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\xc8\x8c\xcf\x2d\x4e\x8f\x37\xb4\xe6\xc2\x23\x6b\x64\xcd\x05\x08\x00\x00\xff\xff\x53\x77\x9e\x4d\x3c\x00\x00\x00")

func _2_messages_indexDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__2_messages_indexDownSql,
		"2_messages_index.down.sql",
	)
}

func _2_messages_indexDownSql() (*asset, error) {
	bytes, err := _2_messages_indexDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "2_messages_index.down.sql", size: 60, mode: os.FileMode(0664), modTime: time.Unix(1673044753, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x6e, 0xcb, 0x70, 0x82, 0x33, 0x13, 0x70, 0xd5, 0xbd, 0x3e, 0x68, 0x9, 0x4f, 0x78, 0xa9, 0xc, 0xd6, 0xf4, 0x64, 0xa0, 0x8c, 0xe4, 0x0, 0x15, 0x71, 0xf0, 0x5, 0xdb, 0xa6, 0xf2, 0x12, 0x60}}
	return a, nil
}

var __2_messages_indexUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x0e\x72\x75\x0c\x71\x55\xf0\xf4\x73\x71\x8d\x50\xf0\x74\x53\xf0\xf3\x0f\x51\x70\x8d\xf0\x0c\x0e\x09\x56\xc8\x8c\xcf\x2d\x4e\x8f\x37\x54\xf0\xf7\x53\xc8\x4d\x2d\x2e\x4e\x4c\x4f\xd5\x48\xce\xcf\x2b\x49\xcd\x2b\x09\xc9\x2f\xc8\x4c\x56\x70\x0c\x76\xd6\x51\x28\x28\x4d\x2a\x2e\x4d\x42\x12\x28\x4e\xcd\x4b\x49\x2d\x0a\xc9\xcc\x4d\x2d\x2e\x49\xcc\x2d\x80\x08\x66\xa6\x80\x68\x4d\x6b\x2e\x82\xd6\x19\xe1\xb4\xce\xc5\x15\xdd\x3e\x88\x08\xba\x85\x10\xd1\xcc\x14\x30\x43\xd3\x9a\x0b\x10\x00\x00\xff\xff\x2a\x3b\xab\xf4\xe2\x00\x00\x00")

func _2_messages_indexUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__2_messages_indexUpSql,
		"2_messages_index.up.sql",
	)
}

func _2_messages_indexUpSql() (*asset, error) {
	bytes, err := _2_messages_indexUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "2_messages_index.up.sql", size: 226, mode: os.FileMode(0664), modTime: time.Unix(1673044753, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xce, 0xb1, 0xc8, 0x2d, 0xa8, 0x6f, 0x83, 0xfb, 0xf2, 0x40, 0x30, 0xe9, 0xd, 0x18, 0x54, 0xe8, 0xf5, 0xf5, 0xc4, 0x5b, 0xf5, 0xa4, 0x94, 0x50, 0x56, 0x4a, 0xc8, 0x73, 0x3f, 0xf1, 0x56, 0xce}}
	return a, nil
}

var __3_rendezvousDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x08\x4a\x4d\xcf\x2c\x2e\x29\x4a\x2c\xc9\xcc\xcf\x2b\xb6\xe6\xe2\xc2\xaa\xc8\x2f\x3f\x2f\x39\xd5\x9a\x0b\x10\x00\x00\xff\xff\x58\x3f\x49\x49\x41\x00\x00\x00")

func _3_rendezvousDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__3_rendezvousDownSql,
		"3_rendezvous.down.sql",
	)
}

func _3_rendezvousDownSql() (*asset, error) {
	bytes, err := _3_rendezvousDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "3_rendezvous.down.sql", size: 65, mode: os.FileMode(0664), modTime: time.Unix(1678392960, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x1f, 0x4b, 0xc0, 0x7d, 0x4f, 0xac, 0xc4, 0x75, 0x59, 0xcc, 0xfc, 0x1a, 0x6c, 0x18, 0x81, 0x29, 0x24, 0x33, 0x3, 0x10, 0x39, 0xd0, 0x67, 0x28, 0xa0, 0xe0, 0xfd, 0x36, 0x91, 0x25, 0x37, 0x83}}
	return a, nil
}

var __3_rendezvousUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\xce\xc1\x0a\x82\x40\x10\xc6\xf1\xfb\x3e\xc5\x77\x54\xf0\x18\x5d\x3a\x8d\x32\x94\x64\x11\xd3\x12\x78\x14\x1d\xc2\xcb\x2a\xb3\x1b\xf4\xf8\x81\x66\xd0\x6d\xf8\xcd\x77\xf8\x57\xc2\xe4\x19\x9e\xca\x86\x61\xfa\x1c\x63\xb2\x2e\x8d\x53\x88\xc8\x1c\x00\xf4\xd3\x2b\x24\x35\xdc\x59\x6a\x6a\x70\x93\xfa\x42\xd2\xe2\xcc\x6d\xb1\xfc\x67\x55\xc3\x83\xa4\x3a\x91\x64\xfb\x5d\xbe\x6a\x88\x9b\x15\x58\x40\xdf\xf3\x68\x8a\xfa\xea\xf9\xc8\x1b\x76\xc3\x60\x11\x65\xeb\x99\x5c\x7e\x70\xee\x2f\x26\x4c\xa1\xd7\x6f\xc4\x7a\xff\x86\x9f\x00\x00\x00\xff\xff\x00\x70\x80\x83\xb5\x00\x00\x00")

func _3_rendezvousUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__3_rendezvousUpSql,
		"3_rendezvous.up.sql",
	)
}

func _3_rendezvousUpSql() (*asset, error) {
	bytes, err := _3_rendezvousUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "3_rendezvous.up.sql", size: 181, mode: os.FileMode(0664), modTime: time.Unix(1678392960, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x5c, 0xb8, 0x4f, 0x88, 0xe9, 0xc6, 0xc, 0xbb, 0x2e, 0x56, 0xa2, 0xcd, 0x9, 0xfa, 0x33, 0x94, 0xd7, 0x73, 0xc1, 0xa, 0xc5, 0x69, 0xfb, 0x9f, 0x75, 0xdb, 0x75, 0x58, 0x20, 0x5e, 0xf, 0x14}}
	return a, nil
}

var __4_signed_peer_recordDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\xc9\xc1\x0a\x82\x40\x10\x06\xe0\xfb\x3e\xc5\x7f\x54\xf0\x18\x5d\x3c\x8d\x36\x94\x64\x25\xd3\x12\x78\x94\x1c\xc2\xcb\x2a\xb3\x1b\xf4\xf8\x81\xe1\xad\xeb\xf7\x1d\xe4\xd6\xc1\x53\xd5\x32\x4c\x5f\x53\x4c\x36\xa4\x69\x0e\xb1\x74\xae\x16\x26\xcf\xff\x12\x99\x03\x80\xe7\xfc\x0e\x49\x0d\x77\x96\x86\x5a\x74\xd2\x5c\x48\x7a\x9c\xb9\x2f\xd6\x5f\x54\x0d\x0f\x92\xfa\x44\x92\xed\x77\xf9\x4f\x43\xdc\xac\xc0\x0a\xfa\x59\x26\x53\x34\x57\xcf\x47\xde\x70\x18\x47\x8b\xa8\x7a\xcf\xe4\xf2\xf2\x1b\x00\x00\xff\xff\x2b\xea\xec\xa2\xa6\x00\x00\x00")

func _4_signed_peer_recordDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__4_signed_peer_recordDownSql,
		"4_signed_peer_record.down.sql",
	)
}

func _4_signed_peer_recordDownSql() (*asset, error) {
	bytes, err := _4_signed_peer_recordDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "4_signed_peer_record.down.sql", size: 166, mode: os.FileMode(0664), modTime: time.Unix(1685640630, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x9c, 0x4b, 0x6c, 0x4d, 0xe5, 0x90, 0x9e, 0xd5, 0xf9, 0x86, 0xe5, 0x23, 0xe2, 0xd, 0xbe, 0x3c, 0x28, 0x8f, 0x68, 0x8e, 0x8e, 0xe9, 0xc5, 0xf4, 0x60, 0x9b, 0xc7, 0x40, 0xca, 0x30, 0xb6, 0x9e}}
	return a, nil
}

var __4_signed_peer_recordUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\xc9\xb1\xaa\x83\x40\x10\x05\xd0\x7e\xbf\xe2\x96\x0a\x96\x8f\xd7\x58\x8d\x66\x48\x24\x26\x91\xc9\x12\xb0\x0c\x3a\xc8\x36\xab\xcc\x1a\xc8\xe7\x07\x0c\x76\x69\xcf\x39\xc8\xad\x83\xa7\xaa\x65\x98\x4e\x21\xad\xf6\x5c\xc3\x1c\x53\xe9\x5c\x2d\x4c\x9e\x7f\x25\x32\x07\x00\xc3\xfc\x8a\xab\x1a\xee\x2c\x0d\xb5\xe8\xa4\xb9\x90\xf4\x38\x73\x5f\x6c\xbf\xa8\x1a\x1e\x24\xf5\x89\x24\xfb\xff\xcb\xbf\x1a\xd3\x6e\x05\x36\xd0\xf7\x12\x4c\xd1\x5c\x3d\x1f\x79\xc7\x14\xa6\xa8\x63\xa7\x6a\xa2\xc3\x6c\x23\xaa\xde\x33\xb9\xbc\x74\x9f\x00\x00\x00\xff\xff\xa8\x37\x10\xd4\xb2\x00\x00\x00")

func _4_signed_peer_recordUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__4_signed_peer_recordUpSql,
		"4_signed_peer_record.up.sql",
	)
}

func _4_signed_peer_recordUpSql() (*asset, error) {
	bytes, err := _4_signed_peer_recordUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "4_signed_peer_record.up.sql", size: 178, mode: os.FileMode(0664), modTime: time.Unix(1685640560, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xb4, 0xa7, 0xfe, 0xfc, 0xdb, 0xd4, 0xad, 0x54, 0xb7, 0x67, 0xb, 0x97, 0x2a, 0xc5, 0x88, 0xc7, 0xe3, 0x49, 0x33, 0xa3, 0x13, 0x7e, 0xb4, 0x6a, 0x15, 0xf3, 0x6b, 0x7e, 0xc5, 0xb6, 0x38, 0x5f}}
	return a, nil
}

var _docGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2c\xc9\xb1\x0d\xc4\x20\x0c\x05\xd0\x9e\x29\xfe\x02\xd8\xfd\x6d\xe3\x4b\xac\x2f\x44\x82\x09\x78\x7f\xa5\x49\xfd\xa6\x1d\xdd\xe8\xd8\xcf\x55\x8a\x2a\xe3\x47\x1f\xbe\x2c\x1d\x8c\xfa\x6f\xe3\xb4\x34\xd4\xd9\x89\xbb\x71\x59\xb6\x18\x1b\x35\x20\xa2\x9f\x0a\x03\xa2\xe5\x0d\x00\x00\xff\xff\x60\xcd\x06\xbe\x4a\x00\x00\x00")

func docGoBytes() ([]byte, error) {
	return bindataRead(
		_docGo,
		"doc.go",
	)
}

func docGo() (*asset, error) {
	bytes, err := docGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "doc.go", size: 74, mode: os.FileMode(0664), modTime: time.Unix(1673044753, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xde, 0x7c, 0x28, 0xcd, 0x47, 0xf2, 0xfa, 0x7c, 0x51, 0x2d, 0xd8, 0x38, 0xb, 0xb0, 0x34, 0x9d, 0x4c, 0x62, 0xa, 0x9e, 0x28, 0xc3, 0x31, 0x23, 0xd9, 0xbb, 0x89, 0x9f, 0xa0, 0x89, 0x1f, 0xe8}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"1_messages.down.sql": _1_messagesDownSql,

	"1_messages.up.sql": _1_messagesUpSql,

	"2_messages_index.down.sql": _2_messages_indexDownSql,

	"2_messages_index.up.sql": _2_messages_indexUpSql,

	"3_rendezvous.down.sql": _3_rendezvousDownSql,

	"3_rendezvous.up.sql": _3_rendezvousUpSql,

	"4_signed_peer_record.down.sql": _4_signed_peer_recordDownSql,

	"4_signed_peer_record.up.sql": _4_signed_peer_recordUpSql,

	"doc.go": docGo,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"1_messages.down.sql":           &bintree{_1_messagesDownSql, map[string]*bintree{}},
	"1_messages.up.sql":             &bintree{_1_messagesUpSql, map[string]*bintree{}},
	"2_messages_index.down.sql":     &bintree{_2_messages_indexDownSql, map[string]*bintree{}},
	"2_messages_index.up.sql":       &bintree{_2_messages_indexUpSql, map[string]*bintree{}},
	"3_rendezvous.down.sql":         &bintree{_3_rendezvousDownSql, map[string]*bintree{}},
	"3_rendezvous.up.sql":           &bintree{_3_rendezvousUpSql, map[string]*bintree{}},
	"4_signed_peer_record.down.sql": &bintree{_4_signed_peer_recordDownSql, map[string]*bintree{}},
	"4_signed_peer_record.up.sql":   &bintree{_4_signed_peer_recordUpSql, map[string]*bintree{}},
	"doc.go":                        &bintree{docGo, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
