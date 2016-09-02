
anim1d_msg.proto  config_msg.proto
PROTO_DIR := ../proto
NANOPB_DIR := ../proto/nanopb
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
#       cd nanopb/generator/proto; make all

# nanopb outputs.
# TODO(maruel): Doesn't get triggered automatically;
app/%.pb.c app/%.pb.h: $(PROTO_DIR)/%.proto
        $(PROTOC) $(PROTOC_OPTS) --nanopb_out=--no-timestamp:app --proto_path=app $<

vpath %.proto $(wildcard $(PROTO_DIR)/*.proto)
