// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/maruel/interrupt"
	"golang.org/x/exp/inotify"
)

func watchFile(fileName string) error {
	fi, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	mod0 := fi.ModTime()
	watcher, err := inotify.NewWatcher()
	if err != nil {
		return err
	}
	if err = watcher.Watch(fileName); err != nil {
		return err
	}
	for {
		select {
		case <-interrupt.Channel:
			return err
		case err = <-watcher.Error:
			return err
		case <-watcher.Event:
			if fi, err = os.Stat(fileName); err != nil || !fi.ModTime().Equal(mod0) {
				return err
			}
		}
	}
}
