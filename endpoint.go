// Copyright (c) 2015, Alexander Cherniuk <ts33kr@gmail.com>
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

package boot

import "time"
import "fmt"

// Implementation of the Operation interface; execute business logic
// that is stored within an endpoint, in regards to supplied context
// structure that should normally represent an HTTP request. See the
// Operation interface for details. The method should be blocking; if
// asynchronous behavior intended - the caller must ensure that this
// method syncrhonizes on the asynchronous code to return onces done.
func (ep *Endpoint) Apply(context *Context) error {
    timer := time.After(ep.Timeout) // ticker
    value := make(chan interface {}, 1) // panic
    const einv = "undetermined endpoint panic %v"
    if e := ep.Satisfied(context); e != nil {
        elog := context.Journal.WithError(e)
        elog = elog.WithField("operation", ep)
        elog.Warn("epiliary is not available")
        return OperationUnavailable // is N/A
    } // operation assured to be available
    go func() { // wrap as asynchronous code
        defer func() { value <- recover() }()
        ep.Business(context) // run the BL!
    }() // spin off go-routine to execute it
    select { // wait for either of 2 channels
        case <- timer: return OperationTimeout
        case x := <- value: switch e := x.(type) {
            case error: return e // regular panic
            case nil: return nil // executed OK
            // operation paniced with non-error
            default: return fmt.Errorf(einv, e)
        }
    }
}

// Check whether the operation is satisfied with supplied context.
// If not - then it is safe to assume that the operation will not
// be available, and its application with yield the corresponding
// error. The exact logic behind this check is determined by the
// implementation. Must return some error value is not satisfied.
func (ep *Endpoint) Satisfied(*Context) error { return nil }

// Fetch prologue & epilogue code (middleware): these are required
// to be run within context prior to running the operation itself.
// Depending on the implementation of an op, middleware can either
// be stored separately in its structure, or be dynamically built
// based on the specific implementation of Operation interface.
func (ep *Endpoint) OnionRings() []Middleware { return ep.Middleware }

// Implementation of the Operation interface; report the error that
// might have occured during execution of the buiness logic implemented
// by an endpoint. Depending on the application settings, this method
// would typically let an HTTP client know about the error, by writing
// to the Context.Responder with the appropriate code and message.
func (ep *Endpoint) ReportIssue(context *Context, err error) {}

// String represenation of this operation, which is used mainly
// for identification purposes when viewed by a human. The value
// is not forced to be unique, but it should unambiguously state
// the operation's identity that can be used by a developer to
// trace it down right to its implementation or definition.
func (ep *Endpoint) String() string { return ep.Pattern }

// Final destination of where an HTTP request lands when it comes via
// the web application. This data structure holds the implementation
// function as well as a number of additional fields that accompany
// the actualy business logic. This data structure should not be
// created or manipulated directly; use framework API for that.
type Endpoint struct {

    // Description of the endpoint; it should be a short and succinct
    // synopsis of what this endpoint does, as a human readable string.
    // Keep it short yet descriptive enough to understand a basic idea
    // of what this endpoint is intended for. This field should be set
    // via corresponding API; please do not modify this directly.
    About string

    // Map of HTTP methods (also known as verbs) that could be used
    // to invoke this endpoint through an HTTP request. Same endpoint
    // can respond to multiple HTTP methods, with possibly different
    // behavior that is encoded in the endpoint implementation logic.
    // This field should not be, as a general, manipulated directly.
    Methods map[string] bool

    // Map of environment names that designates where this endpoint
    // should be made available. If an application is being booted with
    // the configured environment that is not in this slice - endpoint
    // will not be available in that instance of the application. Refer
    // to the App structure and its Env field for more information.
    Available map[string] bool

    // Slice of middleware functions bound to this endpoint. These
    // middleware shall be executed prior to actually executing the
    // business logic embedded in the endpoint structure. For detailed
    // information on middleware, please see Middleware type signature;
    // also refer to the Operation interface definition and usage.
    Middleware []Middleware

    // Amount of time after which the operation application should be
    // considered timed out. If the operation application times out, a
    // caller will be notified of this by returning the special value to
    // it and of course unblocking the call stack. The go-routine that
    // was used to invoke the operation will continue to spin though.
    Timeout time.Duration

    // Pattern that is used to match an HTTP request against this
    // endpoint. Usually it is a mask of a partial URL (a path) that
    // contains parameter placeholders and other pettern expressions.
    // The exact details on the pattern format should be obtained from
    // the router documentation; please refer to it for more info.
    Pattern string

    // Implementation of the endpoint. Should be BiasedLogic typed
    // function that implements the business logic this endpoint is
    // representing. It is invoked to handle an HTTP request matched
    // to this endpoint. A unique per-request context is going to be
    // passed to the function. See BiasedLogic type info for info.
    Business BiasedLogic
}
