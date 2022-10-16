package settings

import (
	"os"
	"log"
	"gopkg.in/yaml.v3"
	"errors"
)

type MqttSettings struct {
	Hostname	string
	Port		string
	Username	string
	Password	string
	ClientId	string
}

type BackendSettings struct {
	SmartCtlEnabled	bool
	NetDataEnabled bool
}

type AllSettings struct {
	Mqtt MqttSettings
	Backend BackendSettings
}

var All AllSettings

func init() {

	log.Println("Opening configuration file...")
	Save()
	Load()
}

// Save Saves the current settings
func Save() {

	if _, err := os.Stat("settings.yaml"); errors.Is(err, os.ErrNotExist) {

		log.Println("Config file not found, writing default config...")

		All.Mqtt = MqttSettings{"192.168.1.4", "1883", "", "", "sensible_mqtt_client"}
		All.Backend = BackendSettings{false, false}

		yaml, err := yaml.Marshal(&All)
		if err != nil {
			log.Fatal(err)
		}

		f, err2 := os.Create("settings.yaml")
		if err2 != nil {
			log.Fatal(err)
		}
		defer f.Close()	

		_, err2 = f.Write(yaml)
		if err2 != nil {
			log.Fatal(err)
		}
	}
}

// Load Loads the current settings
func Load() {

	f, err := os.Open("settings.yaml")
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

// EnsureOk Checks if the loaded configuration is intact
func EnsureOk() {

}
