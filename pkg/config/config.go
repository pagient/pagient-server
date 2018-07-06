package config

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
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
}

// General defines the general configuration.
type General struct {
	Root    string   `ini:"ROOT"`
	Secret  string   `ini:"SECRET"`
	Users   []string `ini:"USERS"`
	Clients []string `ini:"CLIENTS"`
	Pagers  []string `ini:"PAGERS"`
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

	config, err := ini.Load(path.Join(appWorkPath, "/conf/app.ini"))
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
	generalCfg.Root = path.Join(appWorkPath, generalCfg.Root)

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
	if err := checkUserClientMap(cfg); err != nil {
		return nil, errors.Wrap(err, "check user client map failed")
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
	if !allInColonFormat(cfg.General.Users) {
		return errors.New("configuration of 'users' is not formatted correctly")
	}

	if !allInColonFormat(cfg.General.Clients) {
		return errors.New("configuration of 'clients' is not formatted correctly")
	}

	if !allInColonFormat(cfg.General.Pagers) {
		return errors.New("configuration of 'pagers' is not formatted correctly")
	}

	return nil
}

func checkUserClientMap(cfg *Config) error {
	for _, user := range cfg.General.Users {
		pair := strings.SplitN(user, ":", 2)

		if _, err := getClientID(cfg, pair[0]); err != nil {
			return errors.Wrap(err, "get client id failed")
		}
	}

	for _, client := range cfg.General.Clients {
		pair := strings.SplitN(client, ":", 2)

		if _, err := getPassword(cfg, pair[0]); err != nil {
			return errors.Wrap(err, "get password failed")
		}
	}

	return nil
}

func allInColonFormat(items []string) bool {
	for _, item := range items {
		pair := strings.SplitN(item, ":", 2)

		if len(pair) != 2 {
			return false
		}
	}

	return true
}

func getClientID(cfg *Config, name string) (int, error) {
	for _, clientMapping := range cfg.General.Clients {
		clientInfo := strings.SplitN(clientMapping, ":", 2)
		if clientInfo[0] == name {
			id, err := strconv.Atoi(clientInfo[1])
			return id, errors.Wrap(err, "integer string conversion failed")
		}
	}

	return 0, errors.Errorf("No client named %s is configured", name)
}

func getPassword(cfg *Config, name string) (string, error) {
	for _, user := range cfg.General.Users {
		userInfo := strings.SplitN(user, ":", 2)
		if userInfo[0] == name {
			return userInfo[1], nil
		}
	}

	return "", errors.Errorf("No user named %s is configured", name)
}
