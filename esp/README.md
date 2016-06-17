# dlibox-esp

The box for funny people.


## Setup

dlibox-esp uses the following tools and their dependencies:

- [Sming](https://github.com/SmingHub/Sming)
  - [esp-open-sdk](https://github.com/pfalcon/esp-open-sdk)
  - [esptool](https://github.com/themadinventor/esptool)
  - [esptool2](https://github.com/raburton/esptool2)

I wrote a [script](setup.sh) to install and/or update these on Ubuntu, otherwise
follow the wiki pages which explains all the steps for your OS.


## Design choices

- This project uses [Sming](https://github.com/SmingHub/Sming).
  - It doesn't use the [ESP8266 port for
    Arduino](https://github.com/esp8266/Arduino). Someone posted an [explanation
    why](https://primalcortex.wordpress.com/2015/10/22/esp8266-sming-how-to-start/).
    In particular, low latency is needed for this project to drive large number
    of LEDs at high refresh rate.
  - It doesn't use the raw esp-open-sdk because Sming provides a lot of higher
    level tools.
  - [Homie](https://github.com/marvinroger/homie/) looked very interesting but
    it depends on Arduino. Inspiration was taken from its
    [node](https://github.com/marvinroger/homie-esp8266) and
    [server](https://github.com/marvinroger/homie-server).
- Multiple message encoding and RPC/PubSub options were looked at.  C/C++ client
  had to be as compact as possible (<10kb ROM, <2kb RAM) and it had to work in
  Go too for interoperability with the Go version.
  - RPC/PubSub; I wanted to have a router/broker written in Go if possible in
    addition to a client:
    - [WAMP](http://wamp-proto.org/) looked
      [promising](http://wamp-proto.org/compared/). [Go server
      implemented](https://github.com/jcelliott/turnpike). They have an
      [experimental C version](https://github.com/crossbario/autobahn-c) but it
      hasn't been touched since April 2016 and their C++ version uses boost, so
      it's out. WAMP is message encoding agnostic but defaults to Message Pack
      or JSON.
    - [MQTT](http://mqtt.org/) is a complete Pub/Sub message passing protocol.
      Sming has native support. There's a [Go
      client](https://github.com/eclipse/paho.mqtt.golang) which is actively
      being developped. There's a [standard broker](http://mosquitto.org/) that
      can be installed on Raspbian. MQTT is message encoding agnostic.
    - [gRPC](https://github.com/grpc/grpc) requires HTTP/2 which is not
      available yet on embedded systems.
  - [Message
    encodings](https://en.wikipedia.org/wiki/Comparison_of_data_serialization_formats);
    zero-copy would be nice to have but not a hard requirement:
    - Unstructured:
      - JSON is generally inefficient but works everywhere so it's a good
        fallback. There's native support in Sming.
      - [MessagePack](http://msgpack.org/) looks fine. A [compact C
        library](https://github.com/camgunz/cmp),
        [Go](https://github.com/ugorji/go/) client and
        [Javascript](http://kawanet.github.io/msgpack-lite/).
      - [binn](https://github.com/liteserver/binn) looked interesting but a bit
        too experimental and no Go support.
    - Structured; the messages require the schema to be declared upfront via an
      IDL:
      - [ProtoBuf](https://developers.google.com/protocol-buffers/) is the
        Google standard encoding format. There is a surprisingly [compact C
        library](https://github.com/nanopb/nanopb). It also works in
        [Go](https://github.com/golang/protobuf) and
        [Javascript](https://github.com/dcodeIO/ProtoBuf.js/).
      - [CapnProto](https://capnproto.org) supports zero copy but defaults to 64
        bits integers which is huge on 80kb of total RAM. It has a [C
        library](https://github.com/opensourcerouting/c-capnproto) but it is
        unclear how memory efficient it is; the C++ version uses exceptions.
  - So in the end, the choice boils down to structured (ProtoBuf) vs
    unstructured (MsgPack), it's easier to work with structured data but it's
    more work upfront. For message passing, MQTT seems like the obvious choice.
- Device discovery on the local network is done via
  [mDNS](https://en.wikipedia.org/wiki/Multicast_DNS). The Espressif SDK
  has [native
  support](https://github.com/SmingHub/Sming/blob/master/samples/UdpServer_mDNS/app/application.cpp)
  and there's a [Go client+server](https://github.com/hashicorp/mdns).
