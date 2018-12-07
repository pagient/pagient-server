package config

import (
	"github.com/rs/zerolog"
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

	// General config
	General = &general{}
	// Server config
	Server = &server{}
	// DB config
	DB = &db{}
	// Log config
	Log = &log{}

	// Bridge to internal system config
	Bridge = &bridge{}
	// EasyCall config
	EasyCall = &easyCall{}

	// AppWorkPath of binary
	AppWorkPath string
	isWindows   bool
)

// General defines the general configuration.
type general struct {
	Root   string `ini:"ROOT"`
	Secret string `ini:"SECRET"`
	DB     db     `ini:"general"`
}

// Server defines the server configuration.
type server struct {
	Address       string `ini:"ADDRESS"`
	Host          string `ini:"HOST"`
	Root          string `ini:"ROOT"`
	Cert          string `ini:"CERT"`
	Key           string `ini:"KEY"`
	StrictCurves  bool   `ini:"STRICT_CURVES"`
	StrictCiphers bool   `ini:"STRICT_CIPHERS"`
}

// Database defines the repository configuration
type db struct {
	Driver   string `ini:"DB_DRIVER"`
	Host     string `ini:"DB_HOST,omitempty"`
	Port     string `ini:"DB_PORT,omitempty"`
	Name     string `ini:"DB_NAME,omitempty"`
	User     string `ini:"DB_USER,omitempty"`
	Password string `ini:"DB_PASSWORD,omitempty"`
	// in case of sqlite
	Path string `ini:"DB_PATH,omitempty"`
}

// Log defines the logging configuration.
type log struct {
	Level   string `ini:"LEVEL"`
	Colored bool   `ini:"COLORED"`
	Pretty  bool   `ini:"PRETTY"`
}

// Bridge defines the surgery software bridge configuration
type bridge struct {
	DB                      db     `ini:"bridge"`
	PollingInterval         int    `ini:"POLLING_INTERVAL"`
	CallActionWZ            string `ini:"CALL_ACTION_WZ"`
	CallActionQueuePosition uint   `ini:"CALL_ACTION_QUEUE_POSITION"`
	RemoveActionWZ          string `ini:"REMOVE_ACTION_WZ"`
}

// EasyCall defines the easycall pager backend configuration
type easyCall struct {
	URL      string `ini:"URL"`
	User     string `ini:"USER"`
	Password string `ini:"PASSWORD"`
	Port     string `ini:"PORT"`
}

// Load loads the configuration from `Path`
func Load() error {
	isWindows = runtime.GOOS == "windows"

	var appPath string
	var err error
	if appPath, err = getAppPath(); err != nil {
		return errors.Wrap(err, "could not get application path")
	}
	AppWorkPath = getWorkPath(appPath)

	if !filepath.IsAbs(Path) {
		Path = path.Join(AppWorkPath, Path)
	}

	config, err := ini.Load(Path)
	if err != nil {
		return errors.Wrap(err, "could not load ini config")
	}

	if err = config.Section("general").MapTo(General); err != nil {
		return errors.Wrap(err, "could not map general section")
	}

	if !filepath.IsAbs(General.Root) {
		General.Root = path.Join(AppWorkPath, General.Root)
	}
	if err := os.MkdirAll(General.Root, os.ModePerm); err != nil {
		return errors.Wrap(err, "could not create folders of root path")
	}

	if err = config.Section("server").MapTo(Server); err != nil {
		return errors.Wrap(err, "read config server section failed")
	}

	if err = config.Section("db").MapTo(DB); err != nil {
		return errors.Wrap(err, "could not map db section")
	}
	DB.Path = sanitizePath(DB.Path)
	// TODO: move to db package
	dbFolder, err := getDBDirectory()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dbFolder, os.ModePerm); err != nil {
		return errors.Wrap(err, "could not create folder of database path")
	}

	if err = config.Section("log").MapTo(Log); err != nil {
		return errors.Wrap(err, "could not map log section")
	}

	if _, err := zerolog.ParseLevel(Log.Level); err != nil {
		return errors.Wrap(err, "could not parse log level")
	}

	if err = config.Section("bridge").MapTo(Bridge); err != nil {
		return errors.Wrap(err, "read config bridge section failed")
	}

	if err = config.Section("easycall").MapTo(EasyCall); err != nil {
		return errors.Wrap(err, "read config easycall section failed")
	}
	EasyCall.URL = strings.TrimSuffix(EasyCall.URL, "/")

	return nil
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

func sanitizePath(dirtyPath string) string {
	if !filepath.IsAbs(DB.Path) {
		return path.Join(General.Root, dirtyPath)
	}

	return dirtyPath
}

func getDBDirectory() (string, error) {
	// Note: we don't use path.Dir here because it does not handle case
	//		 which path starts with two "/" in Windows: "//psf/Home/..."
	dbPath := strings.Replace(DB.Path, "\\", "/", -1)

	i := strings.LastIndex(dbPath, "/")
	if i == -1 {
		return "", errors.New("database path has to specify sqlite database file")
	}

	return dbPath[:i], nil
}
