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

import "net/http"
import "fmt"

import "github.com/pelletier/go-toml"

// Implementation of http.Handler interface for boot.App struct. It
// will be used to mount the application as HTTP request handler for
// http.Server instances that will be created by the app deployment.
// The boot.App application can, in fact, be mounted into any servers
// that support the standard http.Handler interface and its methods.
func (app *App) ServeHTTP(rw http.ResponseWriter, r *http.Request) {}

// Find all HTTP application server declarations in the app config
// and use the configuration data to create and deploy every declared
// app server. Deploying a server means configuring it with correct
// parameters and bind it to the declared address to listen and accept
// incoming HTTP requests. See boot.App.Deploy method for details.
func (app *App) deployHttpServers() {
    const expression = "$.app.servers.http[:]"
    const ecomp = "can't to compile config query: %s"
    const eempty = "no HTTP app servers in a config"
    log := app.Journal.WithField("server", "HTTP")
    query, err := toml.CompileQuery(expression)
    if err != nil { panic(fmt.Errorf(ecomp, err)) }
    servers := query.Execute(app.Config) // find all!
    if len(servers.Values()) == 0 { panic(eempty) }
    for _, subtree := range servers.Values() {
        config := subtree.(*toml.TomlTree) // cast
        intent := config.Get("intent").(string)
        host := config.Get("hostname").(string)
        port := config.Get("port-number").(uint)
        addr := fmt.Sprintf("%s:%d", host, port)
        server := &http.Server { Addr: addr }
        server.Handler = app // set HTTP handler
        app.Servers[intent] = server // store server
        app.finish.Add(1) // wait for one server
        go func() { // do not block on listening
            log = log.WithField("binding", addr)
            log = log.WithField("intent", intent)
            log.Info("deploy application server")
            defer app.finish.Done() // finished
            server.ListenAndServe() // listen!
        }()
    }
}
