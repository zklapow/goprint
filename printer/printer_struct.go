package printer

import (
    "serial"
    "io"
)

type Printer struct {
    Port *serial.Config
    dev io.ReadWriteCloser
}

func (prt *Printer) Write(data []byte) (err error) {
    _, err = prt.dev.Write(data)
    if err != nil {
        return err
    }

    return nil
}
