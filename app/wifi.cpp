// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>

#include "conf.h"
#include "wifi.h"

namespace {

// Note that debugf() in SmingCore/Platform/WifiEvents.cpp already prints these
// so we may instead keep debugf() enabled.

void onStationConnect(String ssid, unsigned char ssid_len, unsigned char* bssid, unsigned char channel) {
  Serial.printf("stationOnConnect(ssid:%s, channel:%d)\n", ssid.c_str(), channel);
}

void onStationDisconnect(String ssid, unsigned char ssid_len, unsigned char* bssid, unsigned char channel) {
  Serial.printf("onStationDisconnect(ssid:%s, channel:%d)\n", ssid.c_str(), channel);
}

void onStationAuthModeChange(unsigned char old_mode, unsigned char new_mode) {
  Serial.printf("onStationAuthModeChange(old:%d, new:%d)\n", old_mode, new_mode);
}

void onStationGotIP(IPAddress ip, IPAddress mask, IPAddress gateway) {
  Serial.printf("onStationGotIP(ip:%s, mask:%s, gateway:%s)\n", ip.toString().c_str(), mask.toString().c_str(), gateway.toString().c_str());
}

void onAccessPointConnect(unsigned char* mac, unsigned char aid) {
  Serial.printf("onAccessPointConnect(mac:%s, aid:%d)\n", mac, aid);
}

void onAccessPointDisconnect(unsigned char* mac, unsigned char aid) {
  Serial.printf("onAccessPointDisconnect(mac:%s, aid:%d)\n", mac, aid);
}

void onAccessPointProbeReqRecved(short int rssi, unsigned char* mac) {
  Serial.printf("onAccessPointProbeReqRecved(rssi:%d, mac:%s)\n", rssi, mac);
}

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

  WifiEvents.onStationConnect(onStationConnect);
  WifiEvents.onStationDisconnect(onStationDisconnect);
  WifiEvents.onStationAuthModeChange(onStationAuthModeChange);
  WifiEvents.onStationGotIP(onStationGotIP);
  WifiEvents.onAccessPointConnect(onAccessPointConnect);
  WifiEvents.onAccessPointDisconnect(onAccessPointDisconnect);
  WifiEvents.onAccessPointProbeReqRecved(onAccessPointProbeReqRecved);

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
  } else {
    Serial.printf("wifi default AP: \"%s\"\n", config.host.name);
    // TODO(maruel): Scan networks.
    // TODO(maruel): Use WifiStation.smartConfigStart()
    //if (!WifiAccessPoint.config(config.host.name, chipID, AUTH_WPA2_PSK)) {
    if (!WifiAccessPoint.config(config.host.name, "", AUTH_OPEN)) {
      Serial.println("failure");
    }
  }
  //WifiEvents.onStationDisconnect(STADisconnect);
  //WifiEvents.onStationGotIP(STAGotIP);
}
