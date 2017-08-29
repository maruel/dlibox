# Stripped MQTT Go client

This is a fork of https://github.com/eclipse/paho.mqtt.golang with two
modifications to make it smaller when vendoring:

- Removed a lot of unnecessary files
- Removed dependency on golang.org/x/net by removing websocket support
