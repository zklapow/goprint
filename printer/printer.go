package printer

import (
    "serial"
//    "io"
)

func NewPrinter(port string, baud int) (prt *Printer, err error) {
    prt = new(Printer)

    prt.Port = &serial.Config{Name: port, Baud: baud}
    s, err := serial.OpenTTY(prt.Port)
    if err != nil {
        return nil, err
    }

    prt.dev = s

    return prt, nil
}
