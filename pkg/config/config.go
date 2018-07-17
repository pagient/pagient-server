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
	// Path of config file
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
}

// General defines the general configuration.
type General struct {
	Root   string `ini:"ROOT"`
	Secret string `ini:"SECRET"`
}

// Database defines the repository configuration
type Database struct {
	Path string `ini:"PATH"`
}

// EasyCall defines the easycall pager backend configuration
type EasyCall struct {
	URL      string `ini:"URL"`
	User     string `ini:"USER"`
	Password string `ini:"PASSWORD"`
	Port     string `ini:"PORT"`
}

// Bridge defines the surgery software bridge configuration
type Bridge struct {
	DbURL                   string `ini:"DB_URL"`
	DbUser                  string `ini:"DB_USER"`
	DbPassword              string `ini:"DB_PASSWORD"`
	DbName                  string `ini:"DB_NAME"`
	PollingInterval         int    `ini:"POLLING_INTERVAL"`
	CallActionWZ            string `ini:"CALL_ACTION_WZ"`
	CallActionQueuePosition int    `ini:"CALL_ACTION_QUEUE_POSITION"`
	RemoveActionWZ          string `ini:"REMOVE_ACTION_WZ"`
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
	Bridge   Bridge
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
	if err = config.Section("database").MapTo(databaseCfg); err != nil {
		return nil, errors.Wrap(err, "read config repository section failed")
	}

	if !filepath.IsAbs(databaseCfg.Path) {
		databaseCfg.Path = path.Join(generalCfg.Root, databaseCfg.Path)
	}

	// Note: we don't use path.Dir here because it does not handle case
	//		 which path starts with two "/" in Windows: "//psf/Home/..."
	databasePath := strings.Replace(databaseCfg.Path, "\\", "/", -1)
	databaseFolder := ""

	i := strings.LastIndex(databasePath, "/")
	if i == -1 {
		return nil, errors.New("database path has to specify sqlite database file")
	}
	databaseFolder = databasePath[:i]

	if err := os.MkdirAll(databaseFolder, os.ModePerm); err != nil {
		return nil, err
	}

	easyCallCfg := new(EasyCall)
	if err = config.Section("easycall").MapTo(easyCallCfg); err != nil {
		return nil, errors.Wrap(err, "read config easycall section failed")
	}

	easyCallCfg.URL = strings.TrimSuffix(easyCallCfg.URL, "/")

	bridgeCfg := new(Bridge)
	if err = config.Section("bridge").MapTo(bridgeCfg); err != nil {
		return nil, errors.Wrap(err, "read config bridge section failed")
	}

	logCfg := new(Log)
	if err = config.Section("log").MapTo(logCfg); err != nil {
		return nil, errors.Wrap(err, "read config log section failed")
	}

	return &Config{
		Server:   *serverCfg,
		General:  *generalCfg,
		Database: *databaseCfg,
		EasyCall: *easyCallCfg,
		Bridge:   *bridgeCfg,
		Log:      *logCfg,
	}, nil
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
