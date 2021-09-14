package jarvis

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tangx/envutils"
)

type App struct {
	Name string
}

func (app *App) SetDefaults() {
	if app.Name == "" {
		app.Name = "APP"
	}
}

func (app *App) Conf(config interface{}) error {
	app.SetDefaults()

	// call SetDefaults
	if err := envutils.CallSetDefaults(config); err != nil {
		return err
	}

	// write config
	data, err := envutils.Marshal(config, app.Name)
	if err != nil {
		return err
	}
	_ = os.MkdirAll("./config", 0755)
	_ = os.WriteFile("./config/default.yml", data, 0644)

	// load config from files
	for _, _conf := range []string{"default.yml", "config.yml", refConfig()} {
		_conf := filepath.Join("./config/", _conf)
		err = envutils.UnmarshalFile(config, app.Name, _conf)
		if err != nil {
			log.Println(err)
		}
	}

	// load config from env
	err = envutils.UnmarshalEnv(config, app.Name)
	if err != nil {
		log.Print(err)
	}

	// CallInit
	if err := envutils.CallInit(config); err != nil {
		return err
	}

	return nil
}

func refConfig() string {
	// gitlab ci
	ref := os.Getenv("CI_COMMIT_REF_NAME")

	if len(ref) != 0 {
		return _refConfig(ref)
	}

	return `local.yml`
}

func _refConfig(ref string) string {
	// feat/xxxx
	parts := strings.Split(ref, "/")
	feat := parts[len(parts)-1]               // xxxx
	return fmt.Sprintf("config.%s.yml", feat) // config.xxxx.yml

}
