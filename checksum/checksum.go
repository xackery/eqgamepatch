package checksum

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	mutex sync.RWMutex
	db    Checksum
)

// Checksum represents the checksum file database
type Checksum struct {
	Base    map[string]string
	Servers []Server
}

// Server represents various servers
type Server struct {
	Name  string
	Patch map[string]string
}

func init() {
	db = Checksum{
		Base: make(map[string]string),
	}
}

func (db *Checksum) server(name string) (*Server, error) {
	for _, server := range db.Servers {
		if server.Name != name {
			continue
		}
		return &server, nil
	}

	return nil, fmt.Errorf("server not found")
}

// Get obtains a checksum from a file
func Get(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", errors.Wrap(err, "md5 check")
	}
	checksum := fmt.Sprintf("%X", h.Sum(nil))
	return checksum, nil
}

// IsDirty returns true if file is known to be changed since last checksum run
func IsDirty(serverName string, path string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	path = strings.ToLower(path)
	server, err := db.server(serverName)
	if err != nil {
		log.Warn().Err(err).Str("serverName", serverName).Msg("isDirty")
		return false
	}

	checksum, ok := server.Patch[path]
	if !ok {
		checksum, ok = db.Base[path]
		if !ok {
			return true
		}
	}
	sum, err := Get(path)
	if err != nil {
		log.Warn().Err(err).Msg("checksum")
	}
	return checksum != sum
}

// Add adds a new checksum
func Add(category string, serverName string, path string, checksum string) error {
	mutex.Lock()
	defer mutex.Unlock()

	path = strings.ToLower(path)
	switch strings.ToLower(category) {
	case "base":
		val := db.Base[path]
		if val == checksum {
			return nil
		}
		db.Base[path] = checksum
	case "patch":
		server, err := db.server(serverName)
		if err != nil {
			server = &Server{
				Name:  serverName,
				Patch: make(map[string]string),
			}
			db.Servers = append(db.Servers, *server)
		}

		val := server.Patch[path]
		if val == checksum {
			return nil
		}
		server.Patch[path] = checksum
	default:
		return fmt.Errorf("unknown category: %s", category)
	}

	if err := save(); err != nil {
		return errors.Wrap(err, "save")
	}
	return nil
}

// save saves the checksum database
func save() error {
	path := "checksum.dat"
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := toml.NewEncoder(f)

	if err = enc.Encode(db); err != nil {
		return errors.Wrapf(err, "encode %s", path)
	}
	return nil

}

// Load loads the initial checksum database
func Load() error {
	var f *os.File
	newDB := Checksum{}
	isNewDatabase := false
	path := "checksum.dat"
	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrapf(err, "open %s", path)
		}
		f, err = os.Create(path)
		if err != nil {
			return errors.Wrapf(err, "create %s", path)
		}
		fi, err = os.Stat(path)
		if err != nil {
			return errors.Wrap(err, "new checksum info")
		}
		isNewDatabase = true
	}
	if !isNewDatabase {
		f, err = os.Open(path)
		if err != nil {
			return errors.Wrap(err, "open checksum")
		}
	}

	defer f.Close()
	if fi.IsDir() {
		return fmt.Errorf("talkeq.conf is a directory, should be a file")
	}

	if isNewDatabase {
		return nil
	}

	_, err = toml.DecodeReader(f, &newDB)
	if err != nil {
		return errors.Wrapf(err, "decode %s", path)
	}

	mutex.Lock()
	db = newDB
	mutex.Unlock()
	return nil
}
