package serial

// #include <termios.h>
// #include <unistd.h>
import "C"

import (
    "errors"
    "os"
    "io"
    "fmt"
    "syscall"
)

func openTTY(name string, baud int) (rwc io.ReadWriteCloser, err error) {
    f, err := os.OpenFile(name, syscall.O_RDWR | syscall.O_NOCTTY | syscall.O_NONBLOCK, 0666)

    if err != nil {
        return
    }

    fd := C.int(f.Fd())
    if C.isatty(fd) != 1 {
        return nil, errors.New("File is not a tty")
    }

    var st C.struct_termios
    _, err = C.tcgetattr(fd, &st)
    if err != nil {
        f.Close()
        return nil, err
    }

    var speed C.speed_t

    switch baud {
    case 115200:
        speed = C.B115200
    case 57600:
        speed = C.B57600
    case 38400:
        speed = C.B38400
    case 19200:
        speed = C.B19200
    case 9600:
        speed = C.B9600
    default:
        f.Close()
        return nil, fmt.Errorf("Unknown baud rate %v", baud)
    }

    _, err = C.cfsetispeed(&st, speed);
    if err != nil {
        f.Close()
        return nil, err
    }

    _, err = C.cfsetospeed(&st, speed)
    if err != nil {
        f.Close()
        return nil, err
    }

    st.c_cflag |= (C.CLOCAL | C.CREAD)

    st.c_lflag &= ^C.tcflag_t(C.ICANON | C.ECHO | C.ECHOE | C.ISIG)
    st.c_oflag &= ^C.tcflag_t(C.OPOST)

    _, err = C.tcsetattr(fd, C.TCSANOW, &st)
    if err != nil {
        f.Close()
        return nil, err
    }

    r1, _, e := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), uintptr(syscall.F_SETFL), uintptr(0))
    if e != 0 || r1 != 0 {
        s := fmt.Sprint("Clearing NONBLOCK syscall error: ", e, r1)
        f.Close()

        return nil, errors.New(s)
    }

    return f, nil
}
