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

import "os"
import "time"
import "net/http"
import "path/filepath"
import "strings"
import "regexp"
import "sync"

import "github.com/pelletier/go-toml"
import "github.com/Sirupsen/logrus"
import "github.com/satori/go.uuid"
import "github.com/blang/semver"

// Create and initialize a new application. This is a front gate for
// the framework, since you should start by creating a new app struct.
// Every application should have a valid slug (name) and a version. So
// this function makes sure they have been passed and are all valid.
// Generally, you should not be creating more than one application.
func MakeApplication (slug, version string) *App {
    const url = "https://github.com/ts33kr/boot"
    const eslug = "slug is not of correct format"
    const eversion = "version is not valid semver"
    pattern := regexp.MustCompile("^[a-zA-Z0-9-_]+$")
    var parsed semver.Version = semver.MustParse(version)
    if !pattern.MatchString(slug) { panic(eslug) }
    if parsed.Validate() != nil { panic(eversion) }
    reference := uuid.NewV5(uuid.NamespaceURL, url)
    application := &App { Slug: slug, Version: parsed }
    application.Reference = reference // set UUID
    application.Done = &sync.WaitGroup {}
    return application // prepared app
}

// Launch the application. This means that the application should have
// all the services installed and all the necessary configurations done
// before invoking the boot sequence. This method will then carry out
// the boot sequence, which involves all the necessary setups. It will
// block until the application will be gracefully terminated.
func (app *App) Boot(env, level, root string) {
    const eenv = "environment name must be 1 word"
    const estat = "could not open the specified root"
    pattern := regexp.MustCompile("^[a-zA-Z0-9]+$")
    parsedLevel, err := logrus.ParseLevel(level)
    if err != nil { panic("wrong logging level") }
    if !pattern.MatchString(env) { panic(eenv) }
    if _, e := os.Stat(root); e != nil { panic(estat) }
    app.Journal = app.PrepareJournal(parsedLevel)
    app.Env = strings.ToLower(strings.TrimSpace(env))
    app.RootDirectory = filepath.Clean(root)
    app.Storage = make(map[string] interface {})
    app.Config = app.LoadConfig(app.Env, app.RootDirectory)
    for _, srv := range app.Services { srv.Loading(app) }
    app.Booted = time.Now() // mark as booted
    app.InstallServers() // listen to ports
    app.Done.Wait() // waiting to stop
}

func (app *App) ServeHTTP(rw *http.ResponseWriter, r *http.Response) {}
func (app *App) LoadConfig(name, base string) *toml.TomlTree { return nil }
func (app *App) PrepareJournal(level logrus.Level) *logrus.Logger { return nil }
func (app *App) InstallServers() {}

// Core data structure of the framework; represents a web application
// built with the framework. Contains all the necessary API to create
// and launch the application, as well as to maintain its lifecyle and
// the operational business logic. Please refer to the fields of the
// structure as well as the methods for a detailed information.
type App struct {

    // Slug is a short name (or tag) that identifies the application.
    // It is advised to keep it machine & human readable: in a form of
    // of a slug - no spaces, all lower case, et cetera. The framework
    // itself, as well as any other code could use this variable to
    // unique identify an instance of the specific application.
    Slug string

    // Complement the application slug; represents a version of the
    // running application instance. The version format should conform
    // to the semver (Semantical Versioning) format. Typically, version
    // looks like 0.0.1, according to the semver formatting. Refer to
    // the Semver package for more info on how to work with versions.
    Version semver.Version

    // A path within the local file system where an instance of the
    // running application should be residing. The framework will use
    // this path to lookup configuration directories, optional static
    // assets and a number of other things it may need. By default, it
    // will be set to the CWD directory that the app was launched in.
    RootDirectory string

    // Short identifier of the logical environment that this instance
    // of the application is running in, such as: production, staging,
    // development and a number of any other possible environments that
    // could be defined and used by the application creators. It should
    // be kept as short, prererrably a 1-word ID, for convenience.
    Env string

    // Unique identifier of the application instance, conforming to a
    // version 4 of the commonly known UUID standards. Every time an
    // application is launched - it gets a new UUID identifier that
    // uniquely represents the specific instance of the application.
    // So every time you start your application, it gets a new ID.
    Reference uuid.UUID

    // Root level logger, as configured by the framework, according to
    // the application and environment settings. Since the framework
    // makes extensive use of a structured logger, this field contains
    // a pre-configured root logging structure, with no fields set yet.
    // Please refer to the Logrus package for more information on it.
    Journal *logrus.Logger

    // General purpose storage for keeping key/value records per the
    // application instance. The storage may be used by the framework
    // as well as the application code, to store and retrieve any sort
    // of values that may be required by the application logic or the
    // framework logic. Beware, values are empty-interface typed.
    Storage map[string] interface {}

    // Configuration data for the application instance. This will be
    // populated by the framework, when the app is being launched. It
    // will locate the necessary TOML configuration file, based on the
    // environment configured, load it and make it availale to the app.
    // Please refer to the corresponding method for more details.
    Config *toml.TomlTree

    // Instant in time when the application was booted. A nil value
    // should indicate that the application instance has not yet been
    // booted up. This value is used internally by the framework in a
    // multiple of ways; and may also be used by whoever is interested
    // the time of when exactly the application was launched.
    Booted time.Time

    // Application wide stop signal, implement as a wait group. After
    // the app is being booted the caller should wait on this group to
    // be resumed once the application has been gracefully stopped. Do
    // prefer this construct instead of abruptly terminating the app
    // using other, likely more destructive, ways of terminating it.
    Done *sync.WaitGroup

    // Slice of HTTP servers that will be used to server application
    // instance. Servers are automatically created by the framework
    // for every corresponding section in the config file. This is
    // needed for applications that must be served on multiple ports
    // or network interfaces at the same time, within one process.
    Servers []*http.Server

    // Slice of providers installed within this application. Provider
    // is an entity, with a piece of code attached, that provides some
    // kind of functionality for the application, such as: a database
    // connection, etc. Providers will be invoked when the application
    // is being launched. Refer to Provider for more information.
    Providers []*Provider

    // Slice of services mounted in the application instance. Service
    // is a collection of endpoints (HTTP request handlers), amongst
    // other things. This slice should not be manipulated directly;
    // but rather through the provided API to manage services within
    // an application instance; please refer to it for details.
    Services []*Service
}
