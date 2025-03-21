package settings

import (
	"TheTinkerDad/sensible/utility"
	"errors"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type GeneralSettings struct {
	Logfile        string
	LogLevel       string
	ScriptLocation string
}

type MqttSettings struct {
	Hostname string
	Port     string
	Username string
	Password string
	ClientId string
}

type DiscoverySettings struct {
	DeviceName string
	Prefix     string
}

type ApiSettings struct {
	Enabled bool
	Port    int
	Token   string
}

type AllSettings struct {
	General   GeneralSettings
	Mqtt      MqttSettings
	Discovery DiscoverySettings
	Api       ApiSettings
	Plugins   []Plugin
}

type Plugin struct {
	Name              string
	Kind              string
	SensorId          string
	Script            string
	UnitOfMeasurement string
	Icon              string
	DeviceClass       string
	Period            float32
    LastExecuted      time.Time
}

var All AllSettings

var settingsFile string = "/etc/sensible/settings.yaml"

// Backs up the existing settings file - if there's any
func BackupSettingsFile() {

	if _, err := os.Stat(settingsFile); errors.Is(err, os.ErrNotExist) {
		return
	} else {
		utility.Copy(settingsFile, settingsFile+".bkp")
	}
}

// Generates the default configuration file
func GenerateDefaults() {

	All.General = GeneralSettings{"/var/log/sensible/sensible.log", "info", "/etc/sensible/scripts/"}
	All.Mqtt = MqttSettings{"127.0.0.1", "1883", "", "", "sensible_mqtt_client"}
	All.Discovery = DiscoverySettings{"sensible-demo", "homeassistant"}
	All.Api = ApiSettings{Port: 8090, Enabled: false, Token: utility.NewRandomUUID()}
	All.Plugins = make([]Plugin, 1)
	
	All.Plugins[0] = Plugin{"Heartbeat", "internal", "heartbeat", "", "", "mdi:wrench-check", "", 1, time.Now()}
	/*
	All.Plugins[1] = Plugin{"Boot Time", "internal", "boot_time", "", "", "mdi:clock", ""}
	All.Plugins[2] = Plugin{"System Time", "internal", "system_time", "", "", "mdi:clock", ""}
	All.Plugins[3] = Plugin{"Root Disk Free", "script", "root_free", "root_free.sh", "GB", "mdi:harddisk", ""}
	All.Plugins[4] = Plugin{"Host IP Address", "script", "ip_address", "ip_address.sh", "", "mdi:network", ""}
	All.Plugins[5] = Plugin{"Hostname", "script", "hostname", "hostname.sh", "", "mdi:network", ""}
	All.Plugins[6] = Plugin{"Platform", "script", "platform", "platform.sh", "", "mdi:wrench-check", ""}
	*/
	yaml, err := yaml.Marshal(&All)
	if err != nil {
		log.Fatal(err)
	}

	f, err2 := os.Create(settingsFile)
	if err2 != nil {
		log.Fatal(err)
	}
	_, err2 = f.Write(yaml)
	if err2 != nil {
		log.Fatal(err)
	}
	f.Close()
}

// CreateFolders Creates the default folders used by Sensible
func CreateFolders() {

	log.Info("Creating default folders...")
	utility.CreateFolder("/etc/sensible/scripts/")
	utility.CreateFolder("/var/log/sensible")
}

// GenerateDefaultIfNotExists Generates the default configuration file
func GenerateDefaultIfNotExists() {

	if _, err := os.Stat(settingsFile); errors.Is(err, os.ErrNotExist) {

		log.Warn("Config file not found, writing default config...")
		GenerateDefaults()
	}
}

// Load Loads the current settings
func Load() {

	f, err := os.Open(settingsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fi, _ := f.Stat()
	raw := make([]byte, fi.Size())
	f.Read(raw)

	err = yaml.Unmarshal(raw, &All)
	if err != nil {
		log.Fatal(err)
	}
}

// Initialize Tries to load the current settings - initializes a base settings file if there's none available
func Initialize() {

	log.Debug("Opening configuration file...")
	GenerateDefaultIfNotExists()
	Load()
}
