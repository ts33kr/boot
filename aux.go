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
// that is stored within an aux op, in regards to supplied context
// structure that represents some sort of arbitray context. See the
// Operation interface for details. The method should be blocking; if
// asynchronous behavior intended - the caller must ensure that this
// method syncrhonizes on the asynchronous code to return onces done.
func (aux *Aux) Apply(context *Context) error {
    timer := time.After(aux.Timeout) // ticker
    value := make(chan interface {}, 1) // panic
    const einv = "undetermined endpoint panic %v"
    if e := aux.Satisfied(context); e != nil {
        elog := context.Journal.WithError(e)
        elog = elog.WithField("operation", aux)
        elog.Warn("auxiliary is not available")
        return OperationUnavailable // is N/A
    } // operation assured to be available
    go func() { // wrap as asynchronous code
        defer func() { value <- recover() }()
        aux.Business(context) // run the BL!
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
func (aux *Aux) Satisfied(*Context) error { return nil }

// Fetch prologue & epilogue code (middleware): these are required
// to be run within context prior to running the operation itself.
// Depending on the implementation of an op, middleware can either
// be stored separately in its structure, or be dynamically built
// based on the specific implementation of Operation interface.
func (aux *Aux) OnionRings() []Middleware { return aux.Middleware }

// Get a source location of where the definition of this operation
// has been made. This information may not always be available. It
// will be accordingly reflected in the return struct in this case.
// Maintenance of this information should be done within framework.
// Please refer to the SourceLocation struct for more details.
func (aux *Aux) Definition() SourceLocation { return aux.SourceLocation }

// Implementation of the Operation interface; resolve the error that
// might have occured during execution of the buiness logic implemented
// by an aux op. Depending on the application settings, this method
// would typically journal the error to an application and/or context
// journal and optionally use other mechanisms to handle the error.
func (aux *Aux) ResolveIssue(context *Context, err error) {}

// String represenation of this operation, which is used mainly
// for identification purposes when viewed by a human. The value
// is not forced to be unique, but it should unambiguously state
// the operation's identity that can be used by a developer to
// trace it down right to its implementation or definition.
func (aux *Aux) String() string { return aux.Handle }

// Auxiliary operation, not tied into HTTP stack. Aux operations are
// usually attached to services, but not necessarily. Usually, you would
// implement an aux when you need an operation that can be invoked from
// multiple endpoints or other locations that need to access to the
// same operation more than once. Uses BiasedLogic to store logic.
type Aux struct {

    // Handle is an identification tag that is both: human and machine
    // readable. It's purpose is uniquely addressing auxiliary operation
    // within the containing service. In other words, this is just the
    // name of the operation that can be used to refer to the operation.
    // Please refer to its usage for examples and better understanding.
    Handle string

    // Mark current aux operation for execution when a service is
    // getting up. Although it marks the operation to be executed when
    // up-ing the service - it is entirely up to service implementation
    // as to how or when to invoke this operation. See boot.Service
    // and its Up method for more information on the up-ing.
    WhenUp bool

    // Mark current aux operation for execution when a service is
    // going down. Although it marks the operation to be executed when
    // down-ing a service - it is entirely up to service implementation
    // as to how or when to invoke this operation. See boot.Service
    // and its Down method for more information on the down-ing.
    WhenDown bool

    // When contains a value, an aux operation is marked as peridoic
    // job and this field must contain CRON expression that defines
    // when the operation is going to be launched. A contents of the
    // field supports a reasonable subset of the CRON expression
    // specification, including most of the keywords defined.
    CronExpression string

    // Slice of middleware functions bound to this aux op. These
    // middleware shall be executed prior to actually executing the
    // business logic embedded in the auxiliary operation. For detailed
    // information on middleware, please see Middleware type signature;
    // also refer to the Operation interface definition and usage.
    Middleware []Middleware

    // Amount of time after which the operation application should be
    // considered timed out. If the operation application times out, a
    // caller will be notified of this by returning the special value to
    // it and of course unblocking the call stack. The go-routine that
    // was used to invoke the operation will continue to spin though.
    Timeout time.Duration

    // Embedded pipeline instance for this auxiliary operation. By the
    // definition, an aux operation should be ad-hoc and self-contained.
    // Therefore, this field will contain a pipeline that must be used
    // to apply this aux operation. Refer to Service implementation
    // about how the pipeline is construced. See its Up method.
    Pipeline

    // Implementation of the aux. Should be BiasedLogic typed
    // function that implements the business logic this aux op is
    // representing. It is invoked when the aux operation is being
    // requested. The context that will be passed in is determined
    // largely by the caller, so do not make any assumptions on it.
    Business BiasedLogic

    // Store source location of where the definition of this auxiliary
    // is implemented. This information may not always be available. It
    // will be accordingly reflected in the return struct in this case.
    // Maintenance of this information should be done within framework.
    // Please refer to the SourceLocation struct for more details.
    SourceLocation
}
