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

  // http://www.daviddurman.com/flexi-color-picker/
  ColorPicker.fixIndicators(
      document.getElementById('slider-indicator'),
      document.getElementById('picker-indicator'));
  ColorPicker(
      document.getElementById('slider'),
      document.getElementById('picker'),
      //document.getElementById('fancy'),
      function(hex, hsv, rgb, pickerCoordinate, sliderCoordinate) {
        ColorPicker.positionIndicators(
            document.getElementById('slider-indicator'),
            document.getElementById('picker-indicator'),
            sliderCoordinate, pickerCoordinate);
        document.body.style.backgroundColor = hex;
        document.getElementById('rgb_r').value = rgb.r;
        document.getElementById('rgb_g').value = rgb.g;
        document.getElementById('rgb_b').value = rgb.b;
        updatePattern('"' + hex + '"');
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
    node.id = 'button-' + k;
    node.attributes['data-mode'] = v;
    node.innerHTML = '<img src="/thumbnail/' + encodeURI(btoa(v)) + '" /> ' + k;
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
      patterns = oReq.response;
      loadButtons();
    }
    // TODO(maruel): Handle failure.
  };
  oReq.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
  oReq.send('pattern=' + btoa(JSON.stringify(JSON.parse(document.getElementById('patternBox').value))));
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
    // TODO(maruel): Handle failure.
  };
  oReq.send();
}

