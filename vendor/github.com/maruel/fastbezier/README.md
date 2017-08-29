# fastbezier

Fast cubic bezier curve evaluation lookup table for curves `(0, 0), (x0, y0), (x1, y1), (1,
1)` in uint16 space. 

- Trades off precision for performance.
- Particularly optimized for ARM cores.
- Includes a C code generator for embedded devices without a FPU (e.g. ESP8266).

[![GoDoc](https://godoc.org/github.com/maruel/fastbezier?status.svg)](https://godoc.org/github.com/maruel/fastbezier)
