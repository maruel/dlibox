// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// package fastbezier implements a fast cubic bezier curve evaluator for curves
// of type (0, 0), (x0, y0), (x1, y1), (1, 1), in the uint16 domain.
//
// The implementation trades off precision for performance.
package fastbezier
