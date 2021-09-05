package butler

import "testing"

type Person struct {
	Name    string `env:"name"`
	Age     int    `env:"age"`
	Address *Address
}

type Address struct {
	Street string `env:"street"`
}

func (p *Person) SetDefaults() {
	p.Name = "zhangsan"
}

func Test_App(t *testing.T) {

	app := App{
		Name: "APP",
	}

	config := &struct {
		Person *Person
	}{
		Person: &Person{
			Address: &Address{},
		},
	}

	// app.Save(config)
	app.Load(config)
}
