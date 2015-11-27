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
import "os/signal"
import "net/http"
import "path/filepath"
import "strings"
import "regexp"
import "sync"
import "fmt"

import "github.com/naoina/denco"
import "github.com/pelletier/go-toml"
import "github.com/renstrom/shortuuid"
import "github.com/Sirupsen/logrus"
import "github.com/blang/semver"

// Create and initialize a new application. This is a front gate for
// the framework, since you should start by creating a new app struct.
// Every application should have a valid slug (name) and a version. So
// this function makes sure they have been passed and are all valid.
// Generally, you should not be creating more than one application.
func New (slug, version string) *App {
    const url = "https://github.com/ts33kr/boot"
    const eslug = "slug is not of correct format"
    const eversion = "version is not valid semver"
    pattern := regexp.MustCompile("^[a-zA-Z0-9-_]+$")
    var parsed semver.Version = semver.MustParse(version)
    if !pattern.MatchString(slug) { panic(eslug) }
    if parsed.Validate() != nil { panic(eversion) }
    application := &App { Slug: slug, Version: parsed }
    application.Servers = make(map[string]*http.Server)
    application.Reference = shortuuid.New().UUID(url)
    application.Providers = make([]*Provider, 0)
    application.Services = make([]*Service, 0)
    application.TimeLayout = time.RFC850
    application.Namespace = url // set
    return application // prepared app
}

// Erect the application. Once completed, the application should have
// all the services installed and all the necessary configurations done
// before invoking the deploy sequence. Basically, this method will do
// everything to get the application configured and be ready to launch.
// Method itself however will not launch the app; see Deploy for that.
func (app *App) Boot(env, level, root string) {
    const eenv = "environment name must be 1 word"
    const estat = "could not open the specified root"
    pattern := regexp.MustCompile("^[a-zA-Z0-9]+$")
    parsedLevel, err := logrus.ParseLevel(level)
    if err != nil { panic("wrong logging level") }
    if !pattern.MatchString(env) { panic(eenv) }
    if _, e := os.Stat(root); e != nil { panic(estat) }
    app.RootDirectory = filepath.Clean(root)
    app.Journal = app.makeJournal(parsedLevel)
    app.Env = strings.ToLower(strings.TrimSpace(env))
    app.Storage = make(map[string] interface {})
    app.Config = app.loadConfig(app.Env, "config")
    app.Booted = time.Now() // mark app as booted
    for _, p := range app.Providers { p.Setup(app) }
    for _, s := range app.Services { s.Up(app) }
    log := app.Journal.WithField("env", app.Env)
    log = log.WithField("root", app.RootDirectory)
    log = log.WithField("level", parsedLevel)
    log.Info("application has been booted")
    app.router = app.assembleRouter()
}

// Deploy the application. Spawn one or more of HTTP(s) servers, as
// defined in the loaded config, and make them listen on respective
// addresses and ports. Every server will have this application set as
// the HTTP requests handler. Method will block until all servers are
// stopped. See boot.App and this method implementation for details.
func (app *App) Deploy(s *Supervisor) {
    var volume int = len(app.Services)
    log := app.Journal.WithField("slug", app.Slug)
    log = log.WithField("version", app.Version)
    log = log.WithField("ref", app.Reference) // UID
    log.Infof("deploying app with %v services", volume)
    cancelled := make(chan os.Signal, 1) // killed
    signal.Notify(cancelled, os.Interrupt, os.Kill)
    app.unfoldHttpsServers() // spawn HTTPS and listen
    app.unfoldHttpServers() // spawn HTTP and listen
    go func() { // this runs in the background
        _ = <- cancelled // waiting for signal
        signal.Stop(cancelled) // stop monitoring
        fmt.Fprintln(app.Journal.Out) // write ^C\n
        moment := time.Now().Format(app.TimeLayout)
        uptime := time.Now().Sub(app.Booted) // calc
        for _, s := range app.Services { s.Down(app) }
        for _, p := range app.Providers { p.Cleanup(app) }
        log := app.Journal.WithField("time", moment)
        log = log.WithField("uptime", uptime.String())
        log.Warn("shutting the application down")
        os.Exit(2) // emulate Ctrl-C exit code
    }() // run go-routine & wait to finish
    app.finish.Wait()
}

// Load config file that contains the configuration data for the app
// instance. Config file should be a valid TOML file that has a bare
// minimum data to make it a valid config. Method will panic in case if
// there is an error loading the config or interpreting data inside.
// Must have the app.slug and app.version fields defined correctly.
// Refer to implementation code for more details on the loading.
func (app *App) loadConfig(name, base string) *toml.TomlTree {
    const eload = "failed to load TOML config\n %v"
    const estat = "could not open config file at %v"
    const eold = "config version is older than app"
    const eforeign = "config is from different app"
    var root string = app.RootDirectory // root dir
    var fileName string = fmt.Sprintf("%s.toml", name)
    resolved := filepath.Join(root, base, fileName)
    var clean string = filepath.Clean(resolved)
    log := app.Journal.WithField("file", clean)
    log.Info("loading application config file")
    _, err := os.Stat(clean) // check if file exists
    if err != nil { panic(fmt.Errorf(estat, clean)) }
    configTree, err := toml.LoadFile(clean) // load up!
    if err != nil { panic(fmt.Errorf(eload, err.Error())) }
    verStr := configTree.GetDefault("app.version", "")
    slug := configTree.GetDefault("app.slug", "n/a")
    version := semver.MustParse(verStr.(string))
    if version.LT(app.Version) { panic(eold) }
    if slug != app.Slug { panic(eforeign) }
    return configTree // config is ready
}

// Build an adequate instance of the structured logger for this
// application instance. The journal builder may draw data from the
// app instance to configure the journal correctly. This method only
// instantiates a very basic journal; anything more complicated than
// that should be implementing using a boot.Provider to do it.
func (app *App) makeJournal(level logrus.Level) *logrus.Logger {
    const m = "begin writing application journal"
    var journal *logrus.Logger = &logrus.Logger {}
    formatter := new(logrus.TextFormatter) // std
    journal.Level = level // use requested level
    journal.Out = os.Stdout // all goes to stdout
    journal.Hooks = make(logrus.LevelHooks) // empty
    journal.Formatter = formatter // set formatter
    formatter.ForceColors = false // act smart
    formatter.DisableColors = false // make pretty
    formatter.DisableTimestamp = false // is useful
    formatter.FullTimestamp = false // numbers
    formatter.TimestampFormat = time.StampMilli
    formatter.DisableSorting = false // order!
    moment := time.Now().Format(app.TimeLayout)
    journal.WithField("time", moment).Info(m)
    return journal // is ready to use
}

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

    // Default time layout (formatting template) to be used by the
    // framework and application code, when it needs to output or send
    // some dates & times that have no specific formatting requirement.
    // This is prefferable way of doing it, rather than have different,
    // inconsistent formatting in different portions of the code.
    TimeLayout string

    // Short identifier of the logical environment that this instance
    // of the application is running in, such as: production, staging,
    // development and a number of any other possible environments that
    // could be defined and used by the application creators. It should
    // be kept as short, prererrably a 1-word ID, for convenience.
    Env string

    // Default application namespace that will be used to generate
    // UUID identificators of v5. This may potentially used by other
    // consumers that need namespace declaration. Typically, value is
    // going to be a stringified URL that designates some namespace.
    // This field will be set by the framework automatically mostly.
    Namespace string

    // Unique identifier of the application instance, conforming to a
    // version 5 of the commonly known UUID standards. Every time an
    // application is launched - it gets a new UUID identifier that
    // uniquely represents the specific instance of the application.
    // So every time you start your application, it gets a new ID.
    Reference string

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

    // Supervisor instance to use with this application instance. A
    // supervisor is responsible for handling issues that might occur
    // during the normal operation mode. Please refer to Supervisor
    // interface definition & documentation for more information.
    // Normally, a default supervisor should be used, as it is.
    Supervisor Supervisor

    // An HTTP request router that the app will use to match incoming
    // requests against the registered endpoints that should handle the
    // requests. The framework will build and maintain this router
    // automatically; normally you should not be refering to this
    // field directly. See Denco library docs for more details.
    router *denco.Router

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

    // Map of HTTP servers that will be used to server application
    // instance. Servers are automatically created by the framework
    // for every corresponding section in the config file. This is
    // needed for applications that must be served on multiple ports
    // or network interfaces at the same time, within one process.
    Servers map[string]*http.Server

    // Application wide stop signal, implement as a wait group. After
    // the app is being booted the caller should wait on this group to
    // be resumed once the application has been gracefully stopped. Do
    // prefer this construct instead of abruptly terminating the app
    // using other, likely more destructive, ways of terminating it.
    finish sync.WaitGroup

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
