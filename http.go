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
import "fmt"

import "github.com/pelletier/go-toml"
import "github.com/naoina/denco"

// Implementation of http.Handler interface for boot.App struct. It
// will be used to mount the application as HTTP request handler for
// http.Server instances that will be created by the app deployment.
// The boot.App application can, in fact, be mounted into any servers
// that support the standard http.Handler interface and its methods.
func (app *App) ServeHTTP(rw http.ResponseWriter, r *http.Request) {}

// Create and configure an implementation of a HTTP request router.
// It will be used by the application to match incoming requests against
// the endpoints that are meant to handle those requests. Current way
// of implementation uses Denco library for routing. Please see the
// Application.Router field, as well as the library documentation.
func (app *App) assembleRouter() *denco.Router {
    var router *denco.Router = denco.New()
    const mloaded = "registered %v URL patterns"
    app.Journal.Info("assembling request router")
    records := make([]denco.Record, 0) // alloc
    for _, srv := range app.Services {
        for _, ep := range srv.Endpoints {
            epp := strings.TrimPrefix(ep.Pattern, "/")
            mask := fmt.Sprintf("%v/%v", srv.Prefix, epp)
            pipe := &Pipeline {Operation: ep, Service: srv}
            log := app.Journal.WithField("url", mask)
            log = log.WithField("service", srv.Slug)
            log.Debug("mounting endpoint into router")
            record := denco.NewRecord(mask, pipe)
            records = append(records, record)
        } // inner loop actually builds records
    } // build the router and check for errors
    if err := router.Build(records); err != nil {
        app.Journal.Fatal("failed to build router")
        panic(err) // inability to build is fatal
    } // if built successfully, return a router
    app.Journal.Infof(mloaded, len(records))
    return router // is ready for use
}

// Find all HTTPS application server declarations in the app config
// and use the configuration data to create and run every declared
// app server. Running an app server means configuring it with correct
// parameters and bind it to the declared address to listen and accept
// incoming HTTP requests. See boot.App.Deploy method for details.
func (app *App) unfoldHttpsServers() {
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
        app.Servers[intent] = server // store server
        app.finish.Add(1) // wait for one server
        go func() { // do not block on listening
            log = log.WithField("bind", server.Addr)
            log = log.WithField("intent", intent)
            log.Info("spawn application server")
            defer app.finish.Done() // clean up
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
        app.Servers[intent] = server // store server
        app.finish.Add(1) // wait for one server
        go func() { // do not block on listening
            log = log.WithField("bind", server.Addr)
            log = log.WithField("intent", intent)
            log.Info("spawn application server")
            defer app.finish.Done() // clean up
            panic(server.ListenAndServe())
        }()
    }
}
