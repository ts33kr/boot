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

// Create and mount a new endpoint into the current service. Method
// takes the origin function that will take the endpoint instance and
// properly set it up. An endpoint instance itself will be allocated by
// this method and then automatically mounted within the service. Any
// modifications could be made to the endpoint instance afterwards.
func (srv *Service) Endpoint(origin func(*Endpoint)) *Endpoint {
    if !srv.Erected.IsZero() { // service is up?
        panic("refusing to modify erected service")
    } // service is not yet up; we are good to go
    if origin == nil { // origin points to nowhere?
        panic("missing the endpoint origin function")
    } // origin is intact, we shall invoke it later
    var endpoint *Endpoint = &Endpoint {} // allocate
    endpoint.Methods = make(map[string] bool) // HTTP
    endpoint.Timeout = time.Second * 3 // default!
    origin(endpoint) // endpoint is made right here
    if len(endpoint.Methods) == 0 { // no methods?
        endpoint.Methods["GET"] = true
    } // ensure at least one env is in the map
    if len(endpoint.Pattern) == 0 { // empty URL
        panic("missing URL pattern for endpoint")
    } // looks like endpoint was properly assembled
    srv.Lock() // accquire mutex lock on the app
    srv.Endpoints = append(srv.Endpoints, endpoint)
    srv.Unlock() // release the accquired mutex
    return endpoint // is ready for usage
}

// Create and install a new service into the current app. Method
// takes the origin function that will take the service instance and
// properly set it up. The service instance itself will be allocated by
// this method and automatically installed within the application. Any
// modifications could be made to the service instance afterwards.
func (app *App) Service(origin func(*Service)) *Service {
    if !app.Booted.IsZero() { // app is booted?
        panic("refusing to modify the booted app")
    } // app is not yet booted; we are good to go
    if origin == nil { // origin points to nowhere?
        panic("missing the service origin function")
    } // origin is intact, we shall invoke it later
    var service *Service = &Service {} // allocate
    var room = make(map[string] interface {})
    service.Available = make(map[string] bool)
    service.Storage = Storage { Container: room }
    service.Auxes = make(map[string] *Aux)
    origin(service) // service is made right here
    if len(service.Available) == 0 { // no envs?
        service.Available[app.Env] = true
    } // ensure at least one env is in the map
    if len(service.Prefix) == 0 { // empty prefix
        panic("missing mandatory service prefix")
    } // looks like service was properly assembled
    app.Lock() // accquire mutex lock on the app
    app.Services = append(app.Services, service)
    app.Unlock() // release the accquired mutex
    return service // is ready for usage
}
