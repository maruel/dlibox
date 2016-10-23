// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

"use strict";

var patterns = [];
var settings = {};

function onload() {
  fetchPatterns();
  fetchSettings();

  // Set background.
  var text = '';
  for (var i=0; i<50; i++) {
    text += '🐉🐢🐇🌴';
  }
  document.getElementById('background').innerText = text;

  var slideInd = document.getElementById('slide-indicator');
  var pickerInd = document.getElementById('picker-indicator');
  ColorPicker(
      document.getElementById('slide'),
      document.getElementById('picker'),
      function(hex, hsv, rgb, pickerCoordinate, slideCoordinate) {
        ColorPicker.positionIndicators(slideInd, pickerInd, slideCoordinate, pickerCoordinate);
        updateColor(rgb.r, rgb.g, rgb.b);
      });
}

// View.

// Reconstructs the pattern buttons.
function loadButtons() {
  var dst = document.getElementById('boutons');
  dst.innerHTML = '';
  for (var k in patterns) {
    var node = document.createElement('button');
    var v = patterns[k];
    var i = parseInt(k);
    node.id = 'button-' + i + 1;
    node.attributes['data-mode'] = v;
    node.innerHTML = '<img src="/thumbnail/' + encodeURI(btoa(v)) + '" /> ' + (i + 1);
    node.addEventListener('click', function (event) {
      updatePattern(this.attributes['data-mode']);
    });
    dst.appendChild(node);
    dst.appendChild(document.createElement('br'));
  }
}

// Updates the textarea and set the new pattern.
function updatePattern(data) {
  document.getElementById('patternBox').value = data;
  setPattern();
}

function componentToHex(c) {
  var hex = c.toString(16);
  return hex.length == 1 ? "0" + hex : hex;
}

function updateColor(r, g, b) {
  var hex = "#" + componentToHex(r) + componentToHex(g) + componentToHex(b);
  document.body.style.backgroundColor = hex;
  document.getElementById('rgb_r').value = r;
  document.getElementById('rgb_g').value = g;
  document.getElementById('rgb_b').value = b;
  document.getElementById('rgb').value = hex;
  updatePattern('"' + hex + '"');
}

function updateFromHEX() {
  var hex = document.getElementById('rgb').value;
  var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  if (result) {
    updateColor(parseInt(result[1], 16), parseInt(result[2], 16), parseInt(result[3], 16));
  }
}

function updateFromRGB() {
  updateColor(
      parseInt(document.getElementById('rgb_r').value, 10),
      parseInt(document.getElementById('rgb_g').value, 10),
      parseInt(document.getElementById('rgb_b').value, 10));
}

function patternKeyDown() {
  if (event.keyCode == 13) {
    setPattern();
  }
  return false;
}

// API calls.

function fetchPatterns() {
  getJSON('/api/patterns', function(data) {
    patterns = data;
    loadButtons();
  })
}

function fetchSettings() {
  getJSON('/api/settings', function(data) {
    settings = data;
    document.getElementById('settingsBox').value = JSON.stringify(data, null, 2);
  })
}

function setPattern() {
  document.getElementById('patternBox').value = JSON.stringify(
      JSON.parse(document.getElementById('patternBox').value), null, 2);
  var oReq = new XMLHttpRequest();
  oReq.open('post', '/switch', true);
  oReq.responseType = 'json';
  oReq.onreadystatechange = function () {
    if (oReq.readyState === XMLHttpRequest.DONE && oReq.status === 200) {
      document.getElementById('patternBox').value = JSON.stringify(oReq.response, null, 2);
      fetchPatterns();
    }
    // TODO(maruel): Handle failure.
  };
  oReq.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
  oReq.send('pattern=' + btoa(JSON.stringify(JSON.parse(document.getElementById('patternBox').value))));
  return false;
}

function setSettings() {
  document.getElementById('settingsBox').value = JSON.stringify(
      JSON.parse(document.getElementById('settingsBox').value), null, 2);
  /* TODO
  var oReq = new XMLHttpRequest();
  oReq.open('post', '/api/settings', true);
  oReq.responseType = 'json';
  oReq.setRequestHeader('Content-type', 'application/json');
  oReq.send(document.getElementById('settingsBox').value);
  */
  return false;
}

// Misc.

function getJSON(url, onGET) {
  var oReq = new XMLHttpRequest();
  oReq.open('get', url, true);
  oReq.responseType = 'json';
  oReq.onreadystatechange = function () {
    if (oReq.readyState === XMLHttpRequest.DONE && oReq.status === 200) {
      onGET(oReq.response);
    }
    // TODO(maruel): Handle failure by adding a red X at top right.
  };
  oReq.send();
}

