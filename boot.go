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

import "errors"

// Function that encapsulates a unit of application's business logic.
// It is a function of a context struct instance; function is used for
// resolving (handling) HTTP requests that come into the app. Although
// this type sigunature could be as well used to represent other kinds
// of application logic that is derived from the context object.
type BiasedLogic func (*Context)

// Function that encapsulates a unit of application's business logic.
// It is a function of an application struct instance. This function is
// typically used to implement logic that is not bound to any other data
// except the data that is encapsulated in the application instance. A
// common usage is a setup function that only needs an app to work.
type UnbiasedLogic func (*App)

// Pipeline is a structure that wraps an operation with all required
// pieces of data and implementation to properly run it. It Basically
// is a way of providing a permanent context for the operation that
// is always constant, within one instance of the application. Please
// see the structure implementation and usage for more information.
type Pipeline struct { Operation Operation; Service *Service }

// Run the embedded business logic with the supplied context struct.
// This method is responsible for running all pre-requisites prior to
// the operation itself, such as - middleware and/or other utilities.
// See the implementation code for more information. Also, please take
// a look at the Apply method of the Operation interface definition.
func (pipe *Pipeline) Invoke(*Context) {}

// Error value to represent a situation when operation application
// has timed out. This error value will be used by the framework to
// indicate when some operation has failed to execute in the allocated
// amount of time (supposedly configurable). Please see the usage of
// this value by the framework or app code for more information.
var OperationTimeout = errors.New("operation timed out")

// Error value to represent a situation when a requested operation is
// not available within the configured environment. The framework will
// use this error value to indicate when some sort of operation invoked
// but not available according to the app configuration. See usage of
// this value by the framework or app code for more information.
var OperationUnavailable = errors.New("operation is not available")

// Something that contains a piece of application's business logic and
// knows how to invoke it. Any operation within the framework can only
// be invoked in with regards to an instance of the context structure.
// This interface abstracts away of knowledge of what logic type is
// used, it only cares about the ability to apply it to a context.
type Operation interface {

    // Apply whatever business logic is stored in this operation to
    // an instance of the context structure, effectively - executing
    // the business logic. Panic handling must be encapsulated within
    // this method's implementation and in case if there was a panic,
    // its error value should be returned as the method's result.
    Apply(*Context) error

    // Request to make a report of an error that might have occured
    // while applying (executing) the operation. The way how an error
    // is reported entirely depends on the interface implementation.
    // This method should be invoked with error that might have been
    // handed off by the Apply method, upon method's completion.
    ReportIssue(*Context, error)
}
