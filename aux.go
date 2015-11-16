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

// Function that is used to build up an aux op instance. It takes a
// pointer to the aux op that has been pre-allocated and preliminary
// initialized before invoking the maker function, passing it through.
// aux op makers are going to be invoked during application launch.
// Please refer to the aux op API for more information on usage.
type AuxBuilder func (*Aux)

// Implementation of the Operation interface; execute business logic
// that is stored within an aux op, in regards to supplied context
// structure that represents some sort of arbitray context. See the
// Operation interface for details. The method should be blocking; if
// asynchronous behavior needed - must be implemented by the caller.
func (aux *Aux) Apply(context *Context) error { return nil }

// Implementation of the Operation interface; report the error that
// might have occured during execution of the buiness logic implemented
// by an aux op. Depending on the application settings, this method
// would typically journal the error to an application and/or context
// journal and optionally use other mechanisms to expose the error.
func (aux *Aux) ReportIssue(context *Context, err error) {}

// Auxiliary operation, not tied into HTTP stack. Aux operations are
// usually attached to services, but not necessarily. Usually, you would
// implement an aux when you need an operation that can be invoked from
// multiple endpoints or other locations that need to access to the
// same operation more than once. Uses BiasedLogic to store logic.
type Aux struct {

    // Slug is a short name (or tag) that identifies specific aux op.
    // It is advised to keep it machine & human readable: in a form of
    // of a slug - no spaces, all lower case, et cetera. The framework
    // itself, as well as any other code could use this variable to
    // unique identify and label some aux op for referencing it.
    Slug string

    // Description of the aux; it should be a short and succinct
    // synopsis of what this aux does, as a human readable string.
    // Keep it short yet descriptive enough to understand a basic idea
    // of what this aux is intended for. This field should be set
    // via corresponding API; please do not modify this directly.
    About string

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

    // Map of environment names that designates where this aux op
    // should be made available. If an application is being booted with
    // the configured environment that is not in this slice - aux op
    // will not be available in that instance of the application. Refer
    // to the App structure and its Env field for more information.
    Available map[string] bool

    // Implementation of the aux. Should be BiasedLogic typed
    // function that implements the business logic this aux op is
    // representing. It is invoked when the aux operation is being
    // requested. The context that will be passed in is determined
    // largely by the caller, so do not make any assumptions on it.
    Business BiasedLogic
}
