// mgo - MongoDB driver for Go
//
// Copyright (c) 2010-2012 - Gustavo Niemeyer <gustavo@niemeyer.net>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package mgo

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// ---------------------------------------------------------------------------
// Logging integration.

// Avoid importing the log type information unnecessarily.  There's a small cost
// associated with using an interface rather than the type.  Depending on how
// often the logger is plugged in, it would be worth using the type instead.
type log_Logger interface {
	Output(calldepth int, s string) error
}

var (
	// To ensure goroutine safety, use the getter
	// and setter methods below to access these.
	globalLogger unsafe.Pointer
	globalDebug  uint32
)

// Specify the *log.Logger object where log messages should be sent to.
func SetLogger(logger log_Logger) {
	if logger == nil {
		atomic.StorePointer(&globalLogger, nil)
		return
	}
	atomic.StorePointer(&globalLogger, unsafe.Pointer(&logger))
}

func getLogger() log_Logger {
	if loggerPtr := (*log_Logger)(atomic.LoadPointer(&globalLogger)); loggerPtr != nil {
		return *loggerPtr
	}
	return nil
}

// Enable the delivery of debug messages to the logger.  Only meaningful
// if a logger is also set.
func SetDebug(debug bool) {
	value := uint32(0)
	if debug {
		value = 1
	}
	atomic.StoreUint32(&globalDebug, value)
}

func getDebug() bool {
	return atomic.LoadUint32(&globalDebug) != 0
}

func log(v ...interface{}) {
	if logger := getLogger(); logger != nil {
		_ = logger.Output(2, fmt.Sprint(v...))
	}
}

func logln(v ...interface{}) {
	if logger := getLogger(); logger != nil {
		_ = logger.Output(2, fmt.Sprintln(v...))
	}
}

func logf(format string, v ...interface{}) {
	if logger := getLogger(); logger != nil {
		_ = logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func debug(v ...interface{}) {
	if logger := getLogger(); logger != nil && getDebug() {
		_ = logger.Output(2, fmt.Sprint(v...))
	}
}

func debugln(v ...interface{}) {
	if logger := getLogger(); logger != nil && getDebug() {
		_ = logger.Output(2, fmt.Sprintln(v...))
	}
}

func debugf(format string, v ...interface{}) {
	if logger := getLogger(); logger != nil && getDebug() {
		_ = logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func debugFunc(debugFunc func() string) {
	if logger := getLogger(); logger != nil && getDebug() {
		out := debugFunc()
		if out != "" {
			_ = logger.Output(2, out)
		}
	}
}
