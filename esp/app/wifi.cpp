// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>
//extern "C" {
//#include <lwip/mdns.h>
//}

#include "conf.h"
#include "wifi.h"

namespace {

Timer mqttTimer;
MqttClient* mqtt;

void startMqttClient();

void startmDNS() {
  mdns_info info = {0};
  info.host_name = config.host.name;
  info.ipAddr = WifiStation.getIP();
  info.server_name = (char*)"dlibox";
  info.server_port = 80;
  info.txt_data[0] = (char*)"esp8266=1";
  //mdns_init(&info);
  espconn_mdns_init(&info);
  espconn_mdns_server_register();
  espconn_mdns_enable();
}

void onMessageReceived(String topic, String message) {
  Serial.printf("MQTT: \"%s\" : \"%s\"\n", topic.c_str(), message.c_str());
}

void publishUpdate() {
  if (mqtt->getConnectionState() != eTCS_Connected) {
    startMqttClient(); // Auto reconnect
  }
  if (!mqtt->publish("dlibox/%s/stats", "TODO(maruel): Add message")) {
    Serial.println("mqtt publish failed");
  }
  // mqtt->publishWithQoS().
}

void checkMQTTDisconnect(TcpClient& client, bool flag) {
  if (flag == true) {
    Serial.println("MQTT Broker Disconnected!");
  } else {
    Serial.println("MQTT Broker Unreachable!");
  }
  mqttTimer.initializeMs(2000, startMqttClient).start();
}


void startMqttClient() {
  mqttTimer.stop();
  if (!mqtt) {
    mqtt = new MqttClient(config.mqtt.host, config.mqtt.port, onMessageReceived);
    mqtt->setCompleteDelegate(checkMQTTDisconnect);
  }
  if (mqtt->connect(config.host.name, config.mqtt.username, config.mqtt.password)) {
    mqttTimer.initializeMs(2000, startMqttClient).start();
    return;
  }
  if (!mqtt->setWill("last/will", "Dying", 1, true)) {
    Serial.println("Unable to die, device is probably saturated.");
  }
  if (!mqtt->subscribe("dlibox/ota/#")) {
    Serial.println("Unable to subscribe.");
  }
  mqttTimer.initializeMs(1000, publishUpdate).start();
}

/*
void onStationConnect(String ssid, unsigned char ssid_len, unsigned char* bssid, unsigned char channel) {
  Serial.printf("stationOnConnect(ssid:%s, bssid:" MACSTR ", channel:%d)\n", ssid.c_str(), MAC2STR(bssid), channel);
}

void onStationDisconnect(String ssid, unsigned char ssid_len, unsigned char* bssid, unsigned char channel) {
  Serial.printf("onStationDisconnect(ssid:%s, bssid:" MACSTR ", channel:%d)\n", ssid.c_str(), MAC2STR(bssid), channel);
}

void onStationAuthModeChange(unsigned char old_mode, unsigned char new_mode) {
  Serial.printf("onStationAuthModeChange(old:%d, new:%d)\n", old_mode, new_mode);
}
*/

void onStationGotIP(IPAddress ip, IPAddress mask, IPAddress gateway) {
  Serial.printf("onStationGotIP(ip:%s, mask:%s, gateway:%s)\n", ip.toString().c_str(), mask.toString().c_str(), gateway.toString().c_str());
  // - Start mDNS
  startmDNS();
  // TODO(maruel): Query for network local MQTT server.
  // - Start MQTT client.
  if (*config.mqtt.host) {
    startMqttClient();
  }
}

/*
void onAccessPointConnect(unsigned char* mac, unsigned char aid) {
  Serial.printf("onAccessPointConnect(mac:" MACSTR ", aid:%d)\n", MAC2STR(mac), aid);
}

void onAccessPointDisconnect(unsigned char* mac, unsigned char aid) {
  Serial.printf("onAccessPointDisconnect(mac:" MACSTR ", aid:%d)\n", MAC2STR(mac), aid);
}

void onAccessPointProbeReqRecved(short int rssi, unsigned char* mac) {
  Serial.printf("onAccessPointProbeReqRecved(rssi:%d, mac:" MACSTR ")\n",
      rssi, MAC2STR(mac));
}
*/

// Resets the settings stored by the firmware in flash.
// This may happen if a previous firmware hard stored stuff in there by
// accident. Take no chance.
void hardReset() {
  WifiStation.enable(false, true);
  WifiAccessPoint.enable(false, true);
  //wifi_station_set_auto_connect(true);
}

}  // namespace

void initWifi() {
  hardReset();
  if (config.host.highSpeed) {
    wifi_set_sleep_type(NONE_SLEEP_T);
  } else {
    wifi_set_sleep_type(LIGHT_SLEEP_T);
  }

  // Connecting to an access point currently causes a hang. :(
  // TODO(maruel): Figure it out.
  //WifiEvents.onStationConnect(onStationConnect);
  //WifiEvents.onStationDisconnect(onStationDisconnect);
  //WifiEvents.onStationAuthModeChange(onStationAuthModeChange);
  WifiEvents.onStationGotIP(onStationGotIP);
  //WifiEvents.onAccessPointConnect(onAccessPointConnect);
  //WifiEvents.onAccessPointDisconnect(onAccessPointDisconnect);
  //WifiEvents.onAccessPointProbeReqRecved(onAccessPointProbeReqRecved);
  if (config.has_wifiClient && config.wifiClient.has_ssid && config.wifiClient.has_password) {
    Serial.printf("wifi client \"%s\"\n", config.wifiClient.ssid);
    if (!WifiStation.config(config.wifiClient.ssid, config.wifiClient.password, true)) {
      Serial.println("failure");
    }
    WifiStation.enable(true, false);
    WifiStation.connect();
  } else if (config.has_wifiAP && config.wifiAP.has_ssid && config.wifiAP.has_password) {
    Serial.printf("wifi AP \"%s\"\n", config.wifiAP.ssid);
    // TODO(maruel): Channel is hardcoded to 7, beacon at 200ms.
    if (!WifiAccessPoint.config(config.wifiAP.ssid, config.wifiAP.password, AUTH_WPA2_PSK)) {
      Serial.println("failure");
    }
    WifiAccessPoint.enable(true, false);
  } else {
    Serial.printf("wifi default AP: \"%s\"\n", config.host.name);
    // TODO(maruel): Scan networks.
    // TODO(maruel): Use WifiStation.smartConfigStart()
    //if (!WifiAccessPoint.config(config.host.name, chipID, AUTH_WPA2_PSK)) {
    if (!WifiAccessPoint.config(config.host.name, "", AUTH_OPEN)) {
      Serial.println("failure");
    }
    WifiAccessPoint.enable(true, false);
  }
}
