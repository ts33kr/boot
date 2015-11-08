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
// It is a function of a context and an application that is used for
// resolving (handling) HTTP requests that come into the app. Although
// this type sigunature could be as well used to represent other kinds
// of application logic that is derived from an app and a context.
type Logic func (*Context, *App)

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

    // Slice of HTTP methods (also known as verbs) that could be used
    // to invoke this endpoint through an HTTP request. Same endpoint
    // can respond to multiple HTTP methods, with possibly different
    // behavior that is encoded in the endpoint implementation logic.
    // This field should not be, as a general, manipulated directly.
    Methods []string

    // Pattern that is used to match an HTTP request against this
    // endpoint. Usually it is a mask of a partial URL (a path) that
    // contains parameter placeholders and other pettern expressions.
    // The exact details on the pattern format should be obtained from
    // the router documentation; please refer to it for more info.
    Pattern string

    // Implementation of the endpoint. It should be the Logic typed
    // function that implements the business logic this endpoint is
    // representing. It is invoked to handle an HTTP request matched
    // to this endpoint. A unique per-request context is going to be
    // passed to the function. See Logic type info for details.
    Business Logic
}
