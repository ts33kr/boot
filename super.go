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

// Watchdog is a default implementation of the app supervisor to
// be used out of the box, without having to write your own one. It
// should satifsy the very basic needs and functionality for handling
// the conditions as defined by the Supervisor interface. If something
// more custom is necessary, you should implement your own Supervisor.
type Watchdog struct {}

// Invoked when an incoming HTTP request cannot be routed to any
// endpoint, because it does not match any of the endpoints within
// the application. This method should respond to the client with
// the corresponding message and optionally perform other, internal
// routines, such as writing to the application journal.
func (wd *Watchdog) EndpointNotFound(*Context) {}

// Invoked when an incoming HTTP request has been routed to one
// endpoint while that endpoint does not allow for an HTTP method
// (also known as verb) that have been requested. This method should
// respond to the client with the corresponding message and maybe
// perform other, internal routines, such as write app journal.
func (wd *Watchdog) MethodNotAllowed(*Context) {}

// Invoked when an operation application has timed out. This could
// have happened due to different reasons. This can happen for aux
// operation as well as for endpoint. There is no strict algorithm
// as to when this method will be called, as the issues could be
// entirely handled within Operation and Pipeline coding.
func (wd *Watchdog) OperationTimeout(*Context, Operation) {}

// Invoked when an operation is not availe in a current env. Could
// have happened due to different reasons. This can happen for aux
// operation as well as for endpoint. There is no strict algorithm
// as to when this method will be called, as the issues could be
// entirely handled within Operation and Pipeline coding.
func (wd *Watchdog) OperationUnavailable(*Context, Operation) {}

// Invoked when an operation application has paniced. This could
// have happened due to different reasons. This can happen for aux
// operation as well as for endpoint. There is no strict algorithm
// as to when this method will be called, as the issues could be
// entirely handled within Operation and Pipeline coding.
func (wd *Watchdog) OperationPaniced(*Context, Operation, error) {}

// Invoked when the framework detects that the process has been
// running out of the memory limits as configured for application.
// It is then a responsibility of a supervisor to take (or not)
// action, such as reboot or stop the application process and/or
// notify the staff about a problem through available methods.
func (wd *Watchdog) HittingMemLimits(*App) {}

// Supervisor is responsible for handling issues that might occur
// during the normal operation mode. These issues are typically needed
// to be handled in a uniformed fashion, despite their origin. Once an
// issue of that kind arises somewhere, the framework will delegate its
// handling to a supervisor instance assigned for the app instance.
type Supervisor interface {

    // Invoked when an incoming HTTP request cannot be routed to any
    // endpoint, because it does not match any of the endpoints within
    // the application. This method should respond to the client with
    // the corresponding message and optionally perform other, internal
    // routines, such as writing to the application journal.
    EndpointNotFound(*Context)

    // Invoked when an incoming HTTP request has been routed to one
    // endpoint while that endpoint does not allow for an HTTP method
    // (also known as verb) that have been requested. This method should
    // respond to the client with the corresponding message and maybe
    // perform other, internal routines, such as write app journal.
    MethodNotAllowed(*Context)

    // Invoked when an operation application has timed out. This could
    // have happened due to different reasons. This can happen for aux
    // operation as well as for endpoint. There is no strict algorithm
    // as to when this method will be called, as the issues could be
    // entirely handled within Operation and Pipeline coding.
    OperationTimeout(*Context, Operation)

    // Invoked when an operation is not availe in a current env. Could
    // have happened due to different reasons. This can happen for aux
    // operation as well as for endpoint. There is no strict algorithm
    // as to when this method will be called, as the issues could be
    // entirely handled within Operation and Pipeline coding.
    OperationUnavailable(*Context, Operation)

    // Invoked when an operation application has paniced. This could
    // have happened due to different reasons. This can happen for aux
    // operation as well as for endpoint. There is no strict algorithm
    // as to when this method will be called, as the issues could be
    // entirely handled within Operation and Pipeline coding.
    OperationPaniced(*Context, Operation, error)

    // Invoked when the framework detects that the process has been
    // running out of the memory limits as configured for application.
    // It is then a responsibility of a supervisor to take (or not)
    // action, such as reboot or stop the application process and/or
    // notify the staff about a problem through available methods.
    HittingMemLimits(*App)
}
