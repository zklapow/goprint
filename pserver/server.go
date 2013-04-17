package main

import (
    "flag"
    "net/http"
    "html/template"
    "log"
    "printer"
)

var addr = flag.String("addr", ":8000", "http service address")
var tty = flag.String("tty", "/dev/tty.usbserial-A900fN7x", "tty device file")
var baud = flag.Int("baud", 19200, "baud rate for tty device")

var p *printer.Printer = nil

func init() {
    flag.Parse()

    // Print the time using microseconds
    // and the filename/line before each log
    log.SetFlags(log.Lmicroseconds | log.Lshortfile)
    log.Printf("Running on %v", *addr)

    tmp, err := printer.NewPrinter(*tty, *baud)

    p = tmp

    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    // Register the hanlder function
    http.Handle("/", http.HandlerFunc(index))
    http.Handle("/text", http.HandlerFunc(printText))


    // Listen on localhost
    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func index(w http.ResponseWriter, req *http.Request) {
    log.Printf("%v Request: %v", req.Method, req.URL)

    t, err := template.ParseFiles("tmpl/index.html")
    if err != nil {
        log.Fatal(err)
    }

    p := "Welcome"
    t.Execute(w, p)
}

func printText(w http.ResponseWriter, req *http.Request) {
    log.Printf("%v Request: %v", req.Method, req.URL)

    t, err := template.ParseFiles("tmpl/text.html")
    if err != nil {
        log.Fatal(err)
    }

    if req.Method == "GET" {
    } else if req.Method == "POST" {
        data := []byte(req.FormValue("input"))
        log.Printf("Printing: %v", data)

        p.Write(data)
        p.Write([]byte("\n"))

        w.WriteHeader(http.StatusAccepted)

    }
    t.Execute(w, nil)
}
