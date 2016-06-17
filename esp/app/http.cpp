// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>

#include "conf.h"
#include "http.h"

HttpServer server;

// config_page.cpp
extern const char config_page[];

namespace {

void onIndex(HttpRequest &request, HttpResponse &response) {
  response.setCache(60, true);
  response.setContentType(ContentType::HTML);
  response.sendString(config_page);
}

void onConfig(HttpRequest &request, HttpResponse &response) {
  /*
  strcpy(config.wifiClient.password, "");
  strcpy(config.wifiClient.ssid, "");
  config.wifiClient.has_ssid = true;
  config.has_wifiClient = true;
  config.wifiClient.has_password = true;
  config.wifiClient.has_ssid = true;
  saveConfig();
  */
}

void on404(HttpRequest &request, HttpResponse &response) {
  response.forbidden();
}

}  // namespace

void startWebServer() {
  server.listen(80);
  server.addPath("/", onIndex);
  server.addPath("/config", onConfig);
  server.setDefaultHandler(on404);
}
