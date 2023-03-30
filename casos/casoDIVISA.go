package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/wapc/wapc-go"
	"github.com/wapc/wapc-go/engines/wazero"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var subject = "my_subject"

func main() {
	/*
		if len(os.Args) < 2 {
			fmt.Println("usage: hello <name>")
			return
		}
	*/
	//name, err := os.Hostname()
	nc, err := nats.Connect(nats.DefaultURL, nats.UserInfo("root", "xxxx"))
	check(err)
	js, _ := nc.JetStream()
	hostname, err := os.Hostname()
	kv, _ := js.KeyValue("conectados")
	tiempo := time.Now()
	fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		tiempo.Year(), tiempo.Month(), tiempo.Day(),
		tiempo.Hour(), tiempo.Minute(), tiempo.Second())

	kv.Put(hostname, []byte(fecha))
	fmt.Println("conectado")
	if err := nc.Publish(subject, []byte("All is Well")); err != nil {
		log.Fatal(err)
	}
	name := "convertir 10 euros a dolares"
	ctx := context.Background()
	guest, err := os.ReadFile("../../RUST/hello-wasm/pkg/hello_world_bg.wasm")
	if err != nil {
		panic(err)
	}

	engine := wazero.Engine()

	module, err := engine.New(ctx, host, guest, &wapc.ModuleConfig{
		Logger: wapc.PrintlnLogger,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		panic(err)
	}
	defer module.Close(ctx)

	instance, err := module.Instantiate(ctx)
	if err != nil {
		panic(err)
	}
	defer instance.Close(ctx)

	result, err := instance.Invoke(ctx, "hello", []byte(name))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(result))
}

func host(ctx context.Context, binding, namespace, operation string, payload []byte) ([]byte, error) {
	// Route the payload to any custom functionality accordingly.
	// You can even route to other waPC modules!!!
	switch namespace {
	case "example":
		switch operation {
		case "capitalize":
			name := string(payload)
			name = strings.Title(name)
			return []byte(name), nil
		}
	}
	return []byte("default"), nil
}
