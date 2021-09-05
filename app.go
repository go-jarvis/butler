package butler

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tangx/envutils"
)

type App struct {
	Name string
}

func (app *App) Save(config interface{}) error {
	// envutils.Marshal(config)

	if err := envutils.SetDefaults(config); err != nil {
		return err
	}
	data, err := envutils.Marshal(config, app.Name)
	if err != nil {
		return err
	}

	return os.WriteFile("config.yml", data, 0644)
}

func (app *App) Load(config interface{}) error {

	confname := `config.yml`

	file, err := os.Open(confname)
	if err != nil {
		log.Printf("no file %s \n", confname)
	}

	setEnv(file)

	return envutils.LoadEnv(config, app.Name)
}

func setEnv(r io.Reader) {
	fmt.Scanln()
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("read line=>", line)
		parts := strings.Split(line, ":")
		os.Setenv(parts[0], parts[1])

		e2 := os.Getenv(parts[0])
		fmt.Println("e2=>", e2)
	}
}
