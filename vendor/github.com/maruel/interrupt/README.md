interrupt
=========

Package interrupt is a single global way to handle process interruption.

It is useful for both long lived process to implement controlled shutdown and
for CLI tool to handle early termination.

The interrupt signal can be set exactly once in the process lifetime and cannot
be unset. The signal can optionally be set automatically on Ctrl-C/os.Interrupt.
When set, it is expected the process to abort any on-going execution early.

The signal can be read via two ways:

    select {
    case <- interrupt.Channel:
      // Handle abort.
    case ...
      ...
    default:
    }

or

    if interrupt.IsSet() {
      // Handle abort.
    }

[![GoDoc](https://godoc.org/github.com/maruel/interrupt?status.svg)](https://godoc.org/github.com/maruel/interrupt)
[![Build Status](https://travis-ci.org/maruel/interrupt.svg?branch=master)](https://travis-ci.org/maruel/interrupt)
[![Coverage Status](https://img.shields.io/coveralls/maruel/interrupt.svg)](https://coveralls.io/r/maruel/interrupt?branch=master)
