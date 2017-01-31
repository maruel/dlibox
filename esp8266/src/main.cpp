// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// This file uses WeMos pin numbers but they are not special, just defines to
// the actual GPIO number.

// Interfere with boot:
// - RST -> buyton
// - D3 (GPIO0) HIGH run, LOW flash via UART
// - TX (GPIO1)
// - D4 (GPIO2) HIGH; LED
// - RX (GPIO3)
// - D8 (GPIO15) LOW (or boot to SDIO)
// - D0 (GPIO16)-> pulse signal to RST to wake up from wifi

// Left:
// - TX (GPIO1)
// - RX (GPIO3)
// - D1 (GPIO5) Motor R Polarity
// - D2 (GPIO4) Motor L Polarity
// - D3 (GPIO0) Motor PWM
// - D4 (GPIO2)
// - GND
// - 5V

// Right:
// - RST button
// - A0 void
// - D0 (GPIO16)
// - D5 (GPIO14) Button
// - D6 (GPIO12) LED
// - D7 (GPIO13) Buzzer
// - D8 (GPIO15)
// - 3v3

#include <Arduino.h>
#include <Homie.h>

#include "apa102.h"
#include "painter.h"

// Web server to serve the MQTT web UI. This is NOT the web server when in
// configuration mode.
ESP8266WebServer *httpSrv;

void setup() {
  Serial.begin(115200);
  // Increase debug output to maximum level:
  //Serial.setDebugOutput(true);
  // Remove all debug output:
  // Homie.enableLogging(false);
  // Homie.disableLedFeedback(); -> use LED as button.

  // Holding this button for 10s will forcibly reset the device.
  //Homie.setResetTrigger(BUTTON, LOW, 10000);
  Homie_setFirmware("dlibox", "1.0.0");
  Homie_setBrand("dlibox");
  Serial.println();
  Homie.setup();

  if (Homie.isConfigured()) {
    httpSrv = new ESP8266WebServer(80);
    httpSrv->on("/", HTTP_GET, []() {
        String url("/index.html?device=");
        const HomieInternals::ConfigStruct& cfg = Homie.getConfiguration();
        // TODO(maruel): Escaping!!
        url += cfg.deviceId;
        url += "&host=";
        url += cfg.mqtt.server.host;
        // TODO(maruel): The websocket port number != cfg.mqtt.server.port.
        url += "&port=9001";
        //cfg.mqtt.username;
        //cfg.mqtt.password;
        httpSrv->sendHeader("Location", url, true);
        httpSrv->send(307, "text/plain", "");
    });
    httpSrv->serveStatic("/", SPIFFS, "/html/", "public; max-age=600");
    httpSrv->begin();
  }
  initAPA102();
  initPainter();
}

void loop() {
  Homie.loop();
  //buttonNode.update();
  if (httpSrv != NULL) {
    httpSrv->handleClient();
  }
  Painter.loop();
}
