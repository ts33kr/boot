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

// Something that contains a piece of application's business logic and
// knows how to invoke it. Any operation within the framework can only
// be invoked in with regards to an instance of the context structure.
// This interface abstracts away of knowledge of what logic type is
// used, it only cares about the ability to apply it to a context.
type Operation interface {

    // Apply whatever business logic is stored in this operation to
    // an instance of the context structure, effectively - executing
    // the business logic. Panic handling must be encapsulated within
    // this method's implementation and may use the context to obtain
    // or provide whatever might be needed to handle the errors.
    // The code must write to a chan to indicate completion!
    Apply(*Context, chan<-error)

    // Request to make a report of an error that might have occured
    // while applying (executing) the operation. The way how an error
    // is reported entirely depends on the interface implementation.
    // This method should be invoked with error that might have been
    // handed off by the Apply method, through its completion chan.
    ReportIssue(*Context, error)
}
