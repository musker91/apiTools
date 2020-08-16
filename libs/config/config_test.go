package config

import (
	"fmt"
	"testing"
)

func TestInitConfig(t *testing.T) {
	err := InitConfig()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", AppConfig.Service)
}

func TestGet(t *testing.T) {
	err := InitConfig()
	if err != nil {
		t.Fatal(err)
	}
	val := Get("web::startTime")
	t.Log(val)
}