package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
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
var s []string
var s2 []string
var nc2 *nats.Conn
var conectados []string
var x string

func main() {

	http.HandleFunc("/", login)
	http.HandleFunc("/cambiodivisa", cambioMoneda)
	http.ListenAndServe(":8080", nil)

}
func login(w http.ResponseWriter, r *http.Request) {
	//wait := make(chan bool)

	hostname := string("equipoB")
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		s = r.Form["username"]
		s2 = r.Form["password"]

		//wait := make(chan bool)
		//check(err)

		nc, err := nats.Connect(nats.DefaultURL, nats.UserInfo(s[0], s2[0]))

		if err != nil {
			fmt.Println(err)
			t, _ := template.ParseFiles("templates/index.html")
			t.Execute(w, nil)
		} else {

			/*tiempo := time.Now()
			fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			tiempo.Year(), tiempo.Month(), tiempo.Day(),
			tiempo.Hour(), tiempo.Minute(), tiempo.Second())*/

			nc2 = nc
			js, _ := nc.JetStream()

			kv, _ := js.KeyValue("conectados")
			kv.Put(hostname, []byte("conectado"))
			kvtareas, _ := js.KeyValue("tareas")
			entry, _ := kvtareas.Get(hostname)
			if entry != nil {
				//fmt.Printf("%s @ %d -> %q", entry.Key(), entry.Revision(), string(entry.Value()))
				t, _ := template.ParseFiles("templates/app.html")
				t.Execute(w, string("Tienes una tarea que realizar"))
			} else {
				//fmt.Printf("%s @ %d -> %q", entry.Key(), entry.Revision(), string(entry.Value()))
				t, _ := template.ParseFiles("templates/app.html")
				t.Execute(w, nil)
			}

		}
	}

}

func cambioMoneda(w http.ResponseWriter, r *http.Request) {

	//wait := make(chan bool)
	r.ParseForm()
	hostname := string("equipoB")
	name := r.Form["cambiodivisa"]
	ctx := context.Background()
	guest, err := os.ReadFile("../RUST/hello-wasm/pkg/hello_world_bg.wasm")
	if err != nil {
		panic(err)
	}

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

	result, err := instance.Invoke(ctx, "hello", []byte(name[0]))
	if err != nil {
		panic(err)
	}
	stringresult := string(result)
	split := strings.Split(stringresult, " ")
	if string(result) == "conectado" {
		js, _ := nc2.JetStream()

		kv, _ := js.KeyValue("conectados")
		listadeconectados, _ := kv.Keys()
		file, _ := template.ParseFiles("templates/cambiodivisa.html")
		file.Execute(w, listadeconectados)
	} else if string(result) == "tareas" {
		js, _ := nc2.JetStream()

		kv, _ := js.KeyValue("tareas")
		entry, _ := kv.Get(hostname)
		if entry != nil {
			file, _ := template.ParseFiles("templates/cambiodivisa.html")
			file.Execute(w, string(entry.Value()))
		} else {
			file, _ := template.ParseFiles("templates/cambiodivisa.html")
			file.Execute(w, "No hay tareas")
		}

	} else if len(split) == 2 {
		if split[1] == "conectado" {
			rep, err := nc2.Request("altavoz1", []byte("consulta"), time.Second)
			if err != nil {
				file, _ := template.ParseFiles("templates/cambiodivisa.html")
				file.Execute(w, "No")
			} else {
				fmt.Println(string(rep.Data))
				file, _ := template.ParseFiles("templates/cambiodivisa.html")
				file.Execute(w, string(rep.Data))
			}

		}

	} else if string(result) == "accion" {
		split2accion := strings.Split(name[0], " ")
		rep, err := nc2.Request("altavoz1", []byte(split2accion[0]), time.Second)

		if err != nil {
			file, _ := template.ParseFiles("templates/cambiodivisa.html")
			file.Execute(w, split2accion[1]+" esta desconectado")
		} else {
			file, _ := template.ParseFiles("templates/cambiodivisa.html")
			file.Execute(w, string(rep.Data))
		}

	} else if string(result) == "enviartarea" {
		frase := strings.Split(name[0], " ")
		destinatario := frase[0]
		js, _ := nc2.JetStream()

		kvtareas, _ := js.KeyValue("tareas")
		kvtareas.Put(destinatario, []byte(name[0]))
		fmt.Println(name[0])
		entry, _ := kvtareas.Get(hostname)
		if entry != nil {
			file, _ := template.ParseFiles("templates/app.html")
			file.Execute(w, string("Tienes una tarea que realizar"))
		} else {
			file, _ := template.ParseFiles("templates/app.html")
			file.Execute(w, nil)
		}

	} else if string(result) == "eliminartarea" {
		js, _ := nc2.JetStream()

		kvtareas, _ := js.KeyValue("tareas")
		kvtareas.Delete(hostname)
		entry, _ := kvtareas.Get(hostname)
		if entry != nil {
			file, _ := template.ParseFiles("templates/app.html")
			file.Execute(w, string("Tienes una tarea que realizar"))
		} else {
			file, _ := template.ParseFiles("templates/app.html")
			file.Execute(w, nil)
		}
	} else {
		file, _ := template.ParseFiles("templates/cambiodivisa.html")
		file.Execute(w, string(result))
	}
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
