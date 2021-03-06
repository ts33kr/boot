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

import "strings"
import "net/http"
import "time"
import "fmt"

import stdlog "log"

import "github.com/renstrom/shortuuid"
import "github.com/pelletier/go-toml"
import "github.com/Sirupsen/logrus"
import "github.com/naoina/denco"

// Implementation of http.Handler interface for boot.App struct. It
// will be used to mount the application as HTTP request handler for
// http.Server instances that will be created by the app deployment.
// The boot.App application can, in fact, be mounted into any servers
// that support the standard http.Handler interface and its methods.
// Note, it will be invoked in a new go-routine by std HTTP stack.
func (app *App) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
    context := &Context { App: app, Request: r }
    context.Created = time.Now() // mark an instant
    context.ResponseWriter = rw // embed responder
    context.Reference = shortuuid.New() // V4
    log := app.Journal.WithFields(logrus.Fields {
        "ref": context.Reference, // a short UUID
        "url": r.RequestURI, // the URL requested
        "method": r.Method, // an HTTP method (verb)
        "ip": r.RemoteAddr, // remote host & port
    }) // the logger is compiled and ready for use
    log.Info("accepted an incoming HTTP request")
    if app.routers[r.Method] == nil { // allowed?
        log.Warn("request method is not allowed")
        app.Supervisor.MethodNotAllowed(context)
        return // we are done with this request
    } // ok, looks like request method fits in
    router := app.routers[r.Method] // accquire
    rec, ps, hit := router.Lookup(r.RequestURI)
    if !hit { // request did not match any endpoint
        log.Warn("request did not match any route")
        app.Supervisor.EndpointNotFound(context)
        return // we are done with this request
    } // ok, looks like request match an endpoint
    var pipe *Pipeline = rec.(*Pipeline) // cast
    context.Data = make(map[string] string)
    context.Service = pipe.Service // restore
    context.Journal = log // structured logger
    d := context.Data // for convenient access
    for _,p := range ps { d[p.Name] = p.Value }
    pipe.Run(context) // fire up the pipeline
    log.Info("finish accepted HTTP request")
}

// Given the map of HTTP methods to a vector of routables that may
// respond to the specific verb, fill it with the relevant records.
// These records shall be built out of the endpoints registered with
// the application, inside of its services. Please refer to methods
// that build up the HTTP routers for more information and details.
func (app *App) collectRecords(records map[string] []denco.Record) {
    for _, srv := range app.Services {
        for _, ep := range srv.Endpoints {
            epp := strings.TrimPrefix(ep.Pattern, "/")
            mask := fmt.Sprintf("%v/%v", srv.Prefix, epp)
            pipe := &Pipeline {Operation: ep, Service: srv}
            pipe.Compile(app) // seal up pipeline instance
            log := app.Journal.WithField("url", mask)
            log = log.WithField("service", srv)
            log.Debug("mounting endpoint into router")
            record := denco.NewRecord(mask, pipe)
            for m, _ := range ep.Methods { // HTTP verbs
                records[m] = append(records[m], record)
            } // vector of records per each method
        } // inner loop actually builds records
    } // finish up with collecting the records
}

// Create and configure an implementation of the HTTP request routers.
// Will be used by the application to match incoming requests against
// the endpoints that are meant to handle those requests. Current way
// of implementation uses Denco library for routing. Please see the
// Application.routers field, as well as the library documentation.
func (app *App) assembleRouters() map[string] *denco.Router {
    var volume int = 0 // how many records?
    routers := make(map[string] *denco.Router)
    records := make(map[string] []denco.Record)
    const mloaded = "registered %v URL patterns"
    app.Journal.Info("assembling request routers")
    app.collectRecords(records) // build records
    for method, vector := range records { // walk
        routers[method] = denco.New() // allocating
        err := routers[method].Build(vector) // build
        if sz := len(vector); sz > volume { volume = sz }
        if err != nil { // check if built was successful
            app.Journal.Fatal("failed to build router")
            panic(err) // inability to build is fatal
        } // if built successfully, move to the next
    } // done with building routers per HTTP method
    app.Journal.Infof(mloaded, volume) // verbose
    return routers // routers are ready for use
}

// Find all HTTPS application server declarations in the app config
// and use the configuration data to create and run every declared
// app server. Running an app server means configuring it with correct
// parameters and bind it to the declared address to listen and accept
// incoming HTTP requests. See boot.App.Deploy method for details.
func (app *App) unfoldHttpsServers() {
    writer := app.Journal.Writer() // log writer
    log := app.Journal.WithField("proto", "HTTPS")
    const eempty = "no HTTPS app servers in a config"
    sections := app.Config.Get("app.servers.https")
    servers, ok := sections.([]*toml.TomlTree)
    if !ok { panic("invalid app.servers.https") }
    if len(servers) == 0 { panic(eempty) }
    for _, config := range servers {
        key := config.Get("key").(string)
        cert := config.Get("cert").(string)
        intent := config.Get("intent").(string)
        host := config.Get("hostname").(string)
        port := config.Get("port-number").(int64)
        server := &http.Server { Handler: app }
        server.Addr = fmt.Sprintf("%v:%d", host, port)
        server.ErrorLog = stdlog.New(writer, "", 0)
        app.Servers[intent] = server // store server
        app.finish.Add(1) // wait for one server
        go func() { // do not block on listening
            log = log.WithField("bind", server.Addr)
            log = log.WithField("intent", intent)
            log.Info("spawn application server")
            defer app.finish.Done() // clean up
            defer writer.Close() // close writer
            panic(server.ListenAndServeTLS(cert, key))
        }()
    }
}

// Find all HTTP application server declarations in the app config
// and use the configuration data to create and run every declared
// app server. Running an app server means configuring it with correct
// parameters and bind it to the declared address to listen and accept
// incoming HTTP requests. See boot.App.Deploy method for details.
func (app *App) unfoldHttpServers() {
    writer := app.Journal.Writer() // log writer
    log := app.Journal.WithField("proto", "HTTP")
    const eempty = "no HTTP app servers in a config"
    sections := app.Config.Get("app.servers.http")
    servers, ok := sections.([]*toml.TomlTree)
    if !ok { panic("invalid app.servers.http") }
    if len(servers) == 0 { panic(eempty) }
    for _, config := range servers {
        intent := config.Get("intent").(string)
        host := config.Get("hostname").(string)
        port := config.Get("port-number").(int64)
        server := &http.Server { Handler: app }
        server.Addr = fmt.Sprintf("%v:%d", host, port)
        server.ErrorLog = stdlog.New(writer, "", 0)
        app.Servers[intent] = server // store server
        app.finish.Add(1) // wait for one server
        go func() { // do not block on listening
            log = log.WithField("bind", server.Addr)
            log = log.WithField("intent", intent)
            log.Info("spawn application server")
            defer app.finish.Done() // clean up
            defer writer.Close() // close writer
            panic(server.ListenAndServe())
        }()
    }
}
