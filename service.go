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

import "github.com/Sirupsen/logrus"
import "github.com/renstrom/shortuuid"

// Get the service up and running. This method is typically called
// by the framework, during the application deployment sequence. As
// a rule, you would not need to call this method yourself. It will
// initialize the service, run all the relevant aux operation that
// might have been marked for execution during service up-ing.
func (srv *Service) Up(app *App) {
    context := &Context { App: app, Service: srv }
    log := app.Journal.WithFields(logrus.Fields {
        "service": srv.Slug, // short service name
        "prefix": srv.Prefix, // URL prefix used
    }) // logger descriptively identifies a service
    if !srv.Available[app.Env] { // check the env
        log.Warn("is not available in this env")
        return // stop booting, is not available
    } // assume that service is available to run
    context.Created = srv.Erected // creation stamp
    context.Journal = log // setup derived logger
    context.Reference = shortuuid.New() // V4
    log.Info("booting application service up")
    srv.Erected = time.Now() // mark service up
    for _, aux := range srv.Auxes { // walk auxes
        aux.Pipeline = Pipeline { Operation: aux }
        aux.Pipeline.Service = srv // bound the op
        aux.Pipeline.Compile(app) // compile pipe
        oplog := log.WithField("aux", aux) // OP log
        if aux.Satisfied(context) != nil { continue }
        if ce := aux.CronExpression; len(ce) > 0 {
            oplog.Infof("schedule CRON at %v", ce)
            app.CronEngine.AddFunc(ce, func() {
                aux.Run(context) // CRON-called
            }) // schedule as a new CRON task
        } // see if it needs to be invoked on up
        if aux.WhenUp { // invoke when service up
            oplog.Info("running aux on service up")
            aux.Run(context) // invoke on up-ing
        }
    }
}

// Strip the service down and stop. This method is typically called
// by the framework, during the application termination sequence. As
// a rule, you would not need to call this method yourself. It will
// clean-up the service, run all the relevant aux operation that
// might have been marked for execution during service down-ing.
func (srv *Service) Down(app *App) {
    if srv.Erected.IsZero() { return } // down
    srv.Erected = time.Time {} // set service down
    context := &Context { App: app, Service: srv }
    log := app.Journal.WithFields(logrus.Fields {
        "service": srv.Slug, // short service name
        "prefix": srv.Prefix, // URL prefix used
    }) // logger descriptively identifies a service
    context.Created = time.Now() // creation stamp
    context.Journal = log // setup derived logger
    context.Reference = shortuuid.New() // V4
    log.Info("taking application service down")
    for _, aux := range srv.Auxes { // walk auxes
        oplog := log.WithField("aux", aux) // OP log
        if aux.Satisfied(context) != nil { continue }
        if aux.WhenDown { // invoke when service down
            oplog.Info("running aux on service down")
            aux.Run(context) // invoke on down-ing
        }
    }
}

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

    // Mounting point of the service. All the endpoints in the current
    // service will share the same URL prefix, as it is specified when
    // building up a service structure. Therefore, an endpoint that is
    // installed in this service will only be matched by its pattern if
    // the HTTP request URL contains the prefix set in the service.
    Prefix string

    // Map of environment names that designates where this service
    // should be made available. If an application is being booted with
    // the configured environment that is not in this slice - service
    // will not be available in that instance of the application. Refer
    // to the App structure and its Env field for more information.
    Available map[string] bool

    // Map of aux operations belonging to a service. Normally, field
    // should not be manipulated directly, but rather using framework
    // API for that. All aux ops within a group should usually share
    // the same purpose or intention. Please refer to the Aux type
    // for detailed information on the aux operations themselves.
    Auxes map[string] *Aux

    // Slice of middleware functions bound to this service. These
    // middleware shall be executed prior to actually executing the
    // business logic embedded in any aux or endpoint. For detailed
    // information on middleware, please see Middleware type signature;
    // also refer to the Operation interface definition and usage.
    Middleware []Middleware

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

    // Instant in time when the service was brought up. A nil value
    // should indicate that current service instance has not yet been
    // loaded up. This value is used internally by the framework in a
    // multiple of ways; and may also be used by whoever is interested
    // the time of when the service was loaded, if it was at all.
    Erected time.Time
}
