package serial

import "io"

type Config struct {
    Name string
    Baud int
}

func OpenTTY(c *Config) (io.ReadWriteCloser, error) {
    return openTTY(c.Name, c.Baud)
}
