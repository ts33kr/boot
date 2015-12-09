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

// Seal up the pipeline and prepare for execution cycles. Current
// implementation is responsible for building up the middleware chain.
// This chain is an onion-like structure of closures that allow for
// middlewares to be invoked in a fashion that allows for a middleware
// to control ongoing flow of execution of the rest of the chain.
func (pipe *Pipeline) Compile(app *App) {
    pipe.App = app // remember application
    pipe.onion = func (c *Context) { // prepare
        err := pipe.Operation.Apply(c) // run op
        if err != nil { // operation ended with error
            var op Operation = pipe.Operation // shortcut
            var sv Supervisor = app.Supervisor // shortcut
            switch err { // switch on the application error value
                case OperationUnavailable: sv.OperationUnavailable(c, op)
                case OperationTimeout: sv.OperationTimeout(c, op)
                default: sv.OperationPaniced(c, op, err)
            } // we have dispatched the error value
            pipe.Operation.ReportIssue(c, err)
        } // operation application has finished
    } // innermost function actually executes op
    var middleware = make([]Middleware, 0) // alloc
    var inherited = pipe.Service.Middleware // inherit
    items := pipe.Operation.Intermediate() // obtained
    middleware = append(middleware, inherited...) // add
    middleware = append(middleware, items...) // add
    for i := len(middleware) - 1; i >= 0; i-- {
        // reversed for natural order of chaining
        peek := pipe.onion // remember peek layer
        current := middleware[i] // a middleware
        pipe.onion = func (c *Context) {
            current(c, peek) // run it
        }
    }
}

// Run the embedded business logic with the supplied context struct.
// This method is responsible for running all pre-requisites prior to
// the operation itself, such as - middleware and/or other utilities.
// See the implementation code for more information. Also, please take
// a look at the Apply method of the Operation interface definition.
func (pipe *Pipeline) Run(context *Context) { pipe.onion(context) }

// Pipeline is a structure that wraps an operation with all required
// pieces of data and implementation to properly run it. It Basically
// is a way of providing a permanent context for the operation that
// is always constant, within one instance of the application. Please
// see the structure implementation and usage for more information.
type Pipeline struct {

    // Currently compiled onion of closures that represents operation
    // application wrapped in 0 to N layers of middleware. This is an
    // internal field and it should only be used by the Pipeline struct
    // implementation. Take a look at the Compile and Cycle methods
    // for more information on initializing and using the field
    onion BiasedLogic

    // Holds a pointer to an operation that this pipeline is intended
    // to execute. An operation is something that contains a piece of
    // application's business logic and knows how to invoke it. A pipe
    // is agnostic of what that operation exactly is, as long as it is
    // properly implementing the Operation interface abstraction.
    Operation Operation

    // Pointer to an Application structure that represent currently
    // running application. Normally, there can be only one app struct
    // within a process; but that's not a strict requirement. Pointer
    // will always point to a valid App structure and can never be nil.
    // The framework will take care of setting this pointer up.
    App *App

    // Holds a pointer to a service that represents kind of permanent
    // context to be used to run the operation. Service is a group of
    // endpoints that are functionally related. It also serves as a
    // common data exchange bus between the endpoints that belong to
    // the same service. Refer to Service struct for more info.
    Service *Service
}
