// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>

#include "conf.h"
#include "wifi.h"

namespace {

void stationOnConnect() {
  Serial.println("stationOnConnect()");
  //Serial.printf("IP: %s\n", WifiStation.getIP().toString().c_str());
  // getNetworkMask
  // getNetworkGateway
  // getNetworkBroadcast
  // getSSID
  // getRssi
  // getChannel
}

void stationOnConnectFail() {
  // TODO(maruel): Enable Soft AP and continuously retry connecting in case the
  // router is rebooting.
}

}  // namespace

void resetWifi() {
  if (config.host.highSpeed) {
    wifi_set_sleep_type(NONE_SLEEP_T);
  } else {
    wifi_set_sleep_type(LIGHT_SLEEP_T);
  }
  if (config.has_wifiClient && config.wifiClient.has_ssid && config.wifiClient.has_password) {
    Serial.printf("wifi client \"%s\"\n", config.wifiClient.ssid);
    WifiAccessPoint.enable(false);
    WifiStation.waitConnection(stationOnConnect);
    if (!WifiStation.config(config.wifiClient.ssid, config.wifiClient.password)) {
      Serial.println("failure");
    }
  } else if (config.has_wifiAP && config.wifiAP.has_ssid && config.wifiAP.has_password) {
    Serial.printf("wifi AP \"%s\"\n", config.wifiAP.ssid);
    WifiStation.enable(false);
    if (!WifiAccessPoint.config(config.wifiAP.ssid, config.wifiAP.password, AUTH_WPA2_PSK)) {
      Serial.println("failure");
    }
  } else {
    Serial.printf("wifi default AP: \"%s\"\n", config.host.name);
    // TODO(maruel): Scan networks.
    // TODO(maruel): Use WifiStation.smartConfigStart()
    WifiStation.enable(false);
    //if (!WifiAccessPoint.config(config.host.name, chipID, AUTH_WPA2_PSK)) {
    if (!WifiAccessPoint.config(config.host.name, "", AUTH_OPEN)) {
      Serial.println("failure");
    }
  }
  //WifiEvents.onStationDisconnect(STADisconnect);
  //WifiEvents.onStationGotIP(STAGotIP);
}
