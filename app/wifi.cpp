// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>

#include "conf.h"
#include "wifi.h"

namespace {

void stationOnConnect() {
  Serial.printf("IP: %s\r\n", WifiStation.getIP().toString().c_str());
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
  if (config.highSpeed || true) {
    // ??
    wifi_set_sleep_type(NONE_SLEEP_T);
  } else {
    wifi_set_sleep_type(LIGHT_SLEEP_T);
  }
  if (config.has_wifiClient && config.wifiClient.has_ssid && config.wifiClient.has_password) {
    Serial.println("wifi client");
    WifiAccessPoint.enable(false);
    WifiStation.waitConnection(stationOnConnect);
    if (!WifiStation.config(config.wifiClient.ssid, config.wifiClient.password)) {
      Serial.println("failure");
    }
  } else if (config.has_wifiAP && config.wifiAP.has_ssid && config.wifiAP.has_password) {
    Serial.println("wifi AP");
    WifiStation.enable(false);
    if (!WifiAccessPoint.config(config.wifiClient.ssid, config.wifiClient.password, AUTH_WPA2_PSK)) {
      Serial.println("failure");
    }
  } else {
    Serial.printf("wifi soft AP: %s / %s\r\n", hostName, chipID);
    // TODO(maruel): Scan networks.
    // TODO(maruel): Use WifiStation.smartConfigStart()
    WifiStation.enable(false);
    //if (!WifiAccessPoint.config(hostName, chipID, AUTH_WPA2_PSK)) {
    if (!WifiAccessPoint.config(hostName, "", AUTH_OPEN)) {
      Serial.println("failure");
    }
  }
  //WifiEvents.onStationDisconnect(STADisconnect);
  //WifiEvents.onStationGotIP(STAGotIP);

}
