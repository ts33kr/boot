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

// Function that is used to build up a service instance. It takes a
// pointer to the service that has been pre-allocated and preliminary
// initialized before invoking the maker function, passing it through.
// Service makers are going to be invoked during application launch.
// Please refer to the service API for more information on usage.
type ServiceBuilder func (*Service)

// Get the service up and running. This method is typically called
// by the framework, during the application deployment sequence. As
// a rule, you would not need to call this method yourself. It will
// initialize the service, run all the relevant aux operation that
// might have been marked for execution during service up-ing.
func (srv *Service) Up(app *App) {}

// Strip the service down and stop. This method is typically called
// by the framework, during the application termination sequence. As
// a rule, you would not need to call this method yourself. It will
// clean-up the service, run all the relevant aux operation that
// might have been marked for execution during service down-ing.
func (srv *Service) Down(app *App) {}

// Service is a group of endpoints that are functionally related. It
// also serves as a common data exchange bus between the endpoints that
// belong to the same service. Endpoints may store data in the service,
// as well as use it for coordination. Besides this, the data structure
// also contains fields related to the internals of the framework.
type Service struct {

    // Slug is a short name (or tag) that identifies specific service.
    // It is advised to keep it machine & human readable: in a form of
    // of a slug - no spaces, all lower case, et cetera. The framework
    // itself, as well as any other code could use this variable to
    // unique identify and label some service for referencing it.
    Slug string

    // Description of the service; it should be a short and succinct
    // synopsis of what this service does, as a human readable string.
    // Keep it short yet descriptive enough to understand a basic idea
    // of what this service is intended for. This field should be set
    // via corresponding API; please do not modify this directly.
    About string

    // Mounting point of the service. All the endpoints in the current
    // service will share the same URL prefix, as it is specified when
    // building up a service structure. Therefore, an endpoint that is
    // installed in this service will only be matched by its pattern if
    // the HTTP request URL contains the prefix set in the service.
    Prefix string

    // Slice of environment names that designates where this service
    // should be made available. If an application is being booted with
    // the configured environment that is not in this slice - service
    // will not be available in that instance of the application. Refer
    // to the App structure and its Env field for more information.
    Available []string

    // Map of aux operations belonging to a service. Normally, field
    // should not be manipulated directly, but rather using framework
    // API for that. All aux ops within a group should usually share
    // the same purpose or intention. Please refer to the Aux type
    // for detailed information on the aux operations themselves.
    Auxes map[string]*Aux

    // Slice of endpoints that make up this service. Normally, field
    // should not be manipulated directly, but rather using framework
    // API for that. All endpoints within a group should usually share
    // the same purpose or intention. Please refer to the Endpoint type
    // for detailed information on the endpoints themselves.
    Endpoints []*Endpoint

    // General purpose storage for keeping key/value records per the
    // service instance. This storage may be used by the framework
    // as well as application code, to store and retrieve any sort
    // of values that may be required by the service logic or the
    // framework logic. Beware, values are empty-interface typed.
    Storage map[string] interface {}

    // Instant in time when the service was loaded up. A nil value
    // should indicate that current service instance has not yet been
    // loaded up. This value is used internally by the framework in a
    // multiple of ways; and may also be used by whoever is interested
    // the time of when the service was loaded, if it was at all.
    Loaded time.Time
}
