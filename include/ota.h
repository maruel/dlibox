// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __OTA_H__
#define __OTA_H__

// initSerialCommand initializes an interactive prompt over the serial port.
void initSerialCommand();

// OtaUpdate forces an OTA update.
void OtaUpdate(const char *rom_url);

// SwitchROM switches ROM bank and reboot.
void SwitchROM();

#endif
