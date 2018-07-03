package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	_ "github.com/kardianos/minwinsvc" // import minwinsvc for windows services
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
	"path"
)

// enumerates all database provider types
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

// GetPassword returns the password of a user
func (cfg General) GetPassword(name string) (string, error) {
	for _, user := range cfg.Users {
		userInfo := strings.SplitN(user, ":", 2)
		if userInfo[0] == name {
			return userInfo[1], nil
		}
	}

	return "", fmt.Errorf("No user named %s is configured", name)
}

// GetClientID returns the client id of a client
func (cfg General) GetClientID(name string) (int, error) {
	for _, clientMapping := range cfg.Clients {
		clientInfo := strings.SplitN(clientMapping, ":", 2)
		if clientInfo[0] == name {
			return strconv.Atoi(clientInfo[1])
		}
	}

	return 0, fmt.Errorf("No client named %s is configured", name)
}

// GetPagerName returns the name of a pager by id
func (cfg General) GetPagerName(id int) (string, error) {
	for _, pagerMapping := range cfg.Pagers {
		pagerInfo := strings.SplitN(pagerMapping, ":", 2)
		cfgID, err := strconv.Atoi(pagerInfo[0])
		if err != nil {
			return "", err
		}
		if cfgID == id {
			return pagerInfo[1], nil
		}
	}

	return "", nil
}

// Database defines the database configuration
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
	cfg, err := ini.Load(path.Join(appWorkPath, "/conf/app.ini"))
	if err != nil {
		return nil, err
	}

	serverCfg := new(Server)
	if err = cfg.Section("server").MapTo(serverCfg); err != nil {
		return nil, err
	}

	generalCfg := new(General)
	if err = cfg.Section("general").MapTo(generalCfg); err != nil {
		return nil, err
	}
	generalCfg.Root = path.Join(appWorkPath, generalCfg.Root)

	checkFormatting(generalCfg)
	checkUserClientMapping(generalCfg)

	databaseCfg := new(Database)
	if err = cfg.Section("database").MapTo(databaseCfg); err != nil {
		return nil, err
	}

	easyCallCfg := new(EasyCall)
	if err = cfg.Section("easycall").MapTo(easyCallCfg); err != nil {
		return nil, err
	}

	logCfg := new(Log)
	if err = cfg.Section("log").MapTo(logCfg); err != nil {
		return nil, err
	}

	return &Config{
		Server:   *serverCfg,
		General:  *generalCfg,
		Database: *databaseCfg,
		EasyCall: *easyCallCfg,
		Log:      *logCfg,
	}, nil
}

func init() {
	isWindows = runtime.GOOS == "windows"

	var appPath string
	var err error
	if appPath, err = getAppPath(); err != nil {
		log.Fatal().
			Err(err).
			Msg("AppPath could not be found")

		os.Exit(1)
	}

	appWorkPath = getWorkPath(appPath)
}

func getAppPath() (string, error) {
	var appPath string
	var err error

	if isWindows && filepath.IsAbs(os.Args[0]) {
		appPath = filepath.Clean(os.Args[0])
	} else {
		appPath, err = exec.LookPath(os.Args[0])
	}

	if err != nil {
		return "", err
	}
	appPath, err = filepath.Abs(appPath)
	if err != nil {
		return "", err
	}

	// Note: we don't use path.Dir here because it does not handle case
	//		 which path starts with two "/" in Windows: "//psf/Home/..."
	return strings.Replace(appPath, "\\", "/", -1), err
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

func checkFormatting(cfg *General) error {
	correctFormat := allInColonFormat(cfg.Users)
	if !correctFormat {
		return fmt.Errorf("configuration of 'users' is not formatted correctly")
	}

	correctFormat = allInColonFormat(cfg.Clients)
	if !correctFormat {
		return fmt.Errorf("configuration of 'clients' is not formatted correctly")
	}

	correctFormat = allInColonFormat(cfg.Pagers)
	if !correctFormat {
		return fmt.Errorf("configuration of 'pagers' is not formatted correctly")
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

func checkUserClientMapping(cfg *General) error {
	for _, user := range cfg.Users {
		pair := strings.SplitN(user, ":", 2)

		if _, err := cfg.GetClientID(pair[0]); err != nil {
			return err
		}
	}

	for _, client := range cfg.Clients {
		pair := strings.SplitN(client, ":", 2)

		if _, err := cfg.GetPassword(pair[0]); err != nil {
			return err
		}
	}

	return nil
}
