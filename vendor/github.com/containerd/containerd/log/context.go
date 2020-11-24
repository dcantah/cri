/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package log

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

var (
	// G is an alias for GetLogger.
	//
	// We may want to define this locally to a package to get package tagged log
	// messages.
	G = GetLogger

	// L is an alias for the standard logger.
	L = logrus.NewEntry(logrus.StandardLogger())
)

type (
	loggerKey struct{}
)

// RFC3339NanoFixed is time.RFC3339Nano with nanoseconds padded using zeros to
// ensure the formatted time is always the same number of characters.
const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"

// WithLogger returns a new context with the provided logger. Use in
// combination with logger.WithField(s) for great effect.
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// GetLogger retrieves the current logger from the context. If no logger is
// available, the default logger is returned. If the context has a span
// associated with it, add correlating information to the returned logger.
func GetLogger(ctx context.Context) *logrus.Entry {
	e, _ := ctx.Value(loggerKey{}).(*logrus.Entry)
	if e == nil {
		e = L
	}
	if span := trace.FromContext(ctx); span != nil {
		spanCtx := span.SpanContext()
		// The field names used are are specified in the OpenCensus log
		// correlation document, tweaked to match with general Golang naming
		// conventions.
		// https://github.com/census-instrumentation/opencensus-specs/blob/master/trace/LogCorrelation.md
		e = e.WithFields(logrus.Fields{
			"traceID": spanCtx.TraceID.String(),
			"spanID":  spanCtx.SpanID.String(),
		})
	}
	return e
}
