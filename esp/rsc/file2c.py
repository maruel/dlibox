#!/usr/bin/env python
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

import sys

sys.stdout.write('const char %s[] = {' % sys.argv[1])
i = 0
while True:
  c = sys.stdin.read(1)
  if not c:
    break
  if i%16 == 0:
    sys.stdout.write('\n ')
  sys.stdout.write(' 0x%02x,' % ord(c))
  i += 1
print('\n};')
