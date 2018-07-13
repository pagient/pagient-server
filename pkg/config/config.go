package config

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	_ "github.com/kardianos/minwinsvc" // import minwinsvc for windows services
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
	"runtime"
)

// enumerates all repository provider types
const (
	DatabaseProviderFile string = "file"
)

var (
	Path string

	isWindows   bool
	appWorkPath string
)

// Server defines the server configuration.
type Server struct {
	Address       string `ini:"ADDRESS"`
	Host          string `ini:"HOST"`
	Root          string `ini:"ROOT"`
	Cert          string `ini:"CERT"`
	Key           string `ini:"KEY"`
	StrictCurves  bool   `ini:"STRICT_CURVES"`
	StrictCiphers bool   `ini:"STRICT_CIPHERS"`
	SecretKey     string `ini:"SECRET_KEY"`
}

// General defines the general configuration.
type General struct {
	Root       string   `ini:"ROOT"`
	Secret     string   `ini:"SECRET"`
	Users      []string `ini:"USERS"`
	Clients    []string `ini:"CLIENTS"`
	UserClient []string `ini:"USER_CLIENT"`
	Pagers     []string `ini:"PAGERS"`
}

// Database defines the repository configuration
type Database struct {
	Provider string `ini:"PROVIDER"`
}

// EasyCall defines the easycall pager backend configuration
type EasyCall struct {
	URL      string `ini:"URL"`
	User     string `ini:"USER"`
	Password string `ini:"PASSWORD"`
}

// Log defines the logging configuration.
type Log struct {
	Level   string `ini:"LEVEL"`
	Colored bool   `ini:"COLORED"`
	Pretty  bool   `ini:"PRETTY"`
}

// Config defines the general configuration.
type Config struct {
	Server   Server
	General  General
	Database Database
	EasyCall EasyCall
	Log      Log
}

// New prepares a new default configuration.
func New() (*Config, error) {
	isWindows = runtime.GOOS == "windows"

	var appPath string
	var err error
	if appPath, err = getAppPath(); err != nil {
		return nil, errors.Wrap(err, "get app path failed")
	}

	appWorkPath = getWorkPath(appPath)

	if !filepath.IsAbs(Path) {
		Path = path.Join(appWorkPath, Path)
	}

	config, err := ini.Load(Path)
	if err != nil {
		return nil, errors.Wrap(err, "load config ini file failed")
	}

	serverCfg := new(Server)
	if err = config.Section("server").MapTo(serverCfg); err != nil {
		return nil, errors.Wrap(err, "read config server section failed")
	}

	generalCfg := new(General)
	if err = config.Section("general").MapTo(generalCfg); err != nil {
		return nil, errors.Wrap(err, "read config general section failed")
	}

	if !filepath.IsAbs(generalCfg.Root) {
		generalCfg.Root = path.Join(appWorkPath, generalCfg.Root)
	}
	if err := os.MkdirAll(generalCfg.Root, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "create folders for storage root failed")
	}

	databaseCfg := new(Database)
	if err = config.Section("repository").MapTo(databaseCfg); err != nil {
		return nil, errors.Wrap(err, "read config repository section failed")
	}

	easyCallCfg := new(EasyCall)
	if err = config.Section("easycall").MapTo(easyCallCfg); err != nil {
		return nil, errors.Wrap(err, "read config easycall section failed")
	}

	logCfg := new(Log)
	if err = config.Section("log").MapTo(logCfg); err != nil {
		return nil, errors.Wrap(err, "read config log section failed")
	}

	cfg := &Config{
		Server:   *serverCfg,
		General:  *generalCfg,
		Database: *databaseCfg,
		EasyCall: *easyCallCfg,
		Log:      *logCfg,
	}

	if err := checkFormat(cfg); err != nil {
		return nil, errors.Wrap(err, "check format failed")
	}

	return cfg, nil
}

func getAppPath() (string, error) {
	var appPath string

	if isWindows && filepath.IsAbs(os.Args[0]) {
		appPath = filepath.Clean(os.Args[0])
	} else {
		var err error
		appPath, err = exec.LookPath(os.Args[0])
		if err != nil {
			return "", errors.Wrap(err, "app path lookup failed")
		}
	}

	appPath, err := filepath.Abs(appPath)
	if err != nil {
		return "", errors.Wrap(err, "get absolute app path failed")
	}

	// Note: we don't use path.Dir here because it does not handle case
	//		 which path starts with two "/" in Windows: "//psf/Home/..."
	return strings.Replace(appPath, "\\", "/", -1), nil
}

func getWorkPath(appPath string) string {
	workPath := ""

	i := strings.LastIndex(appPath, "/")
	if i == -1 {
		workPath = appPath
	} else {
		workPath = appPath[:i]
	}

	// Note: we don't use path.Dir here because it does not handle case
	//		 which path starts with two "/" in Windows: "//psf/Home/..."
	return strings.Replace(workPath, "\\", "/", -1)
}

func checkFormat(cfg *Config) error {
	if !itemsInColonNotation(cfg.General.Users) || !itemsUniqueByColonNotationSide(cfg.General.Users, 0) {
		return errors.New("configuration of 'users' is not formatted correctly")
	}

	if !itemsInColonNotation(cfg.General.Clients) || !itemsUniqueByColonNotationSide(cfg.General.Clients, 0) {
		return errors.New("configuration of 'clients' is not formatted correctly")
	}

	if !itemsInColonNotation(cfg.General.UserClient) || !itemsUniqueByColonNotationSide(cfg.General.UserClient, 0) || !itemsUniqueByColonNotationSide(cfg.General.UserClient, 1) {
		return errors.New("configuration of 'user_client' is not formatted correctly")
	}

	if !itemsInColonNotation(cfg.General.Pagers) || !itemsUniqueByColonNotationSide(cfg.General.Pagers, 0) {
		return errors.New("configuration of 'pagers' is not formatted correctly")
	}

	if err := checkUserClientMap(cfg); err != nil {
		return errors.Wrap(err, "configuration of 'user_client' is incorrect")
	}

	return nil
}

func checkUserClientMap(cfg *Config) error {
	// check that every mapping has a valid user
UserClientLoop:
	for _, userClientInfo := range cfg.General.UserClient {
		userClientPair := strings.SplitN(userClientInfo, ":", 2)

		for _, userInfo := range cfg.General.Users {
			userPair := strings.SplitN(userInfo, ":", 2)

			if userClientPair[0] == userPair[0] {
				continue UserClientLoop
			}
		}
		return errors.Errorf("No user is configured for user client mapping %s", userClientPair[0])
	}

	// check that every client is mentioned in the mapping
	if len(cfg.General.Clients) != len(cfg.General.UserClient) {
		return errors.New("client and user_client configuration is not valid")
	}

ClientLoop:
	for _, clientInfo := range cfg.General.Clients {
		clientPair := strings.SplitN(clientInfo, ":", 2)

		for _, userClientInfo := range cfg.General.UserClient {
			userClientPair := strings.SplitN(userClientInfo, ":", 2)

			if clientPair[0] == userClientPair[1] {
				continue ClientLoop
			}
		}
		return errors.Errorf("No client user mapping is configured for client %s", clientPair[0])
	}

	return nil
}

// checks whether items are in colon notation e.g. {id}:{name}
func itemsInColonNotation(items []string) bool {
	for _, item := range items {
		pair := strings.SplitN(item, ":", 2)

		if len(pair) != 2 {
			return false
		}
	}

	return true
}

// checks whether items aren't unique on the given side of the split
func itemsUniqueByColonNotationSide(items []string, side int) bool {
	for i, item := range items {
		pair := strings.SplitN(item, ":", 2)

		for y, otherItem := range items {
			otherPair := strings.SplitN(otherItem, ":", 2)

			if i == y {
				continue
			}

			if pair[side] == otherPair[side] {
				return false
			}
		}
	}

	return true
}
