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
import "net/http"
import "github.com/satori/go.uuid"

// Unique object that captures the details needed to invoke the Logic
// typed function. Usually, context will include an HTTP request object,
// a means of responding to the HTTP request, as well as references to
// the App and Service instances that have been used to process this
// HTTP request up to the point when it has reached the Logic func.
type Context struct {

    // Instant in time when this context object was created. This value
    // is used internally by the framework in a multiple of ways; and
    // may also be used by whoever is interested the time of when the
    // context object has been instantiated. The value will be set by
    // the framework, so please do not modify this value directly.
    Created time.Time

    // Unique identifier of the context instance, conforming to a
    // version 4 of the commonly known UUID standards. Every time a
    // new context is created - it gets a new UUID identifier that
    // uniquely represents the specific instance of the contex, which
    // effectively represents every HTTP request that comes in.
    Reference uuid.UUID

    // General purpose storage for keeping key/value records per the
    // context instance. This storage may be used by the framework
    // as well as application code, to store and retrieve any sort
    // of values that may be required by the application logic or the
    // framework logic. Beware, values are empty-interface typed.
    Storage map[string] interface {}

    // Pointer to the HTTP requested that triggered the creation of
    // a context instance. This field will be automatically set by the
    // framework; please do not manipulate it directly. In some very
    // rare occasions, it is possible that the pointer will have nil
    // value, indicating that there was no HTTP context to attach.
    Request *http.Request

    // Pointer to the HTTP response writer instance that can be used
    // to write out a response to the incoming HTTP request that has
    // been wrapped by this context instance. Note that in some very
    // rare occasions, it is possible that the pointer will have nil
    // value, indicating that there is no valid response writer.
    Responder *http.ResponseWriter

    // Pointer to a Service struct instance that a context could be
    // bound to. Usually, when an HTTP request comes in - it is being
    // handled by an Endpoint that resides within a Service. In some
    // rare occasions, it is possible that the pointer will have nil
    // value, indicating that there was no Service to attach.
    Service *Service

    // Pointer to an Application structure that represent currently
    // running application. Normally, there can be only one app struct
    // within a process; but that's not a strict requirement. Pointer
    // will always point to a valid App structure and can never be nil.
    // The framework will take care of setting this pointer up.
    Application *App
}
