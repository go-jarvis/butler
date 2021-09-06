package butler

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tangx/envutils"
)

type App struct {
	Name string
}

func (app *App) Conf(config interface{}) error {
	if err := envutils.SetDefaults(config); err != nil {
		return err
	}

	// write config
	data, err := envutils.Marshal(config, app.Name)
	if err != nil {
		return err
	}
	_ = os.WriteFile("default.yml", data, 0644)

	// load config from files
	for _, _conf := range []string{"default.yml", "config.yml", refConfig()} {
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
