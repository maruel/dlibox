# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.


# Default: 115200
# COM_SPEED = 921600
# Running at 921.6kbps makes flashing flake out occasionally but it's
# significantly faster.
COM_SPEED_ESPTOOL := 921600

RBOOT_ENABLED := 1

SPI_SIZE := 4M

# Defaults to 40. Not sure why the default is not 80Mhz. Downgrade if there is
# stability issues.
SPI_SPEED = 80


.PHONY: push f

# Doesn't push SPIFFS rom and doesn't start the terminal.
push: all
	$(vecho) "Killing Terminal to free $(COM_PORT)"
	-$(Q) $(KILL_TERM)
	$(ESPTOOL) -p $(COM_PORT) -b $(COM_SPEED_ESPTOOL) write_flash $(flashimageoptions) 0x00000 $(RBOOT_BIN) 0x02000 $(RBOOT_ROM_0)


# Like target 'f' but doesn't push SPIFFS rom.
f: push
	$(TERMINAL)


# TODO(maruel): Rewrite the following, I'm not really good at Makefiles and it's
# not working well.


# TODO(maruel): Doesn't get triggered automatically;
app/%.cpp: rsc/%.html
	python ./rsc/file2c.py $(notdir $(basename $<)) < $< > $(addprefix app/,$(basename $(notdir $<)).cpp)

vpath %.html $(wildcard rsc/*.html)

PROTO_DIR := proto
NANOPB_DIR := proto/nanopb
NANOPB_PROTO_DIR := $(NANOPB_DIR)/generator/proto
NANOPB_CORE := $(NANOPB_DIR)/pb_encode.c $(NANOPB_DIR)/pb_decode.c $(NANOPB_DIR)/pb_common.c
PROTOC := protoc
PROTOC_OPTS := --plugin=$(NANOPB_DIR)/generator/protoc-gen-nanopb -I$(NANOPB_DIR)/generator/proto -I$(PROTO_DIR)
EXTRA_INCDIR := include $(NANOPB_DIR)
MODULES := app $(NANOPB_DIR)
NANO_LIB := $(NANOPB_PROTO_DIR)/nanopb_pb2.py $(NANOPB_PROTO_DIR)/plugin_pb2.py

# For each proto file, add the corresponding .pb.c source file as a dependency.
# TODO(maruel): It's not necessary if the files are already created. So for now
# just commit the files (ugh). This depends on a change in Makefile-rboot.mk to
# concatenate SRC and SRC_EXTRA to feed into C_OBJ.
#SRC_EXTRA := $(patsubst %.proto,%.pb.c,$(wildcard app/*.proto))

# protoc plugin.
# TODO(maruel): Commented out otherwise this becomes the default target (!?!)
#$(NANO_LIB):
#       # Define PB_BUFFER_ONLY
#	cd nanopb/generator/proto; make all

# nanopb outputs.
# TODO(maruel): Doesn't get triggered automatically;
app/%.pb.c app/%.pb.h: $(PROTO_DIR)/%.proto
	$(PROTOC) $(PROTOC_OPTS) --nanopb_out=--no-timestamp:app --proto_path=app $<

vpath %.proto $(wildcard $(PROTO_DIR)/*.proto)
