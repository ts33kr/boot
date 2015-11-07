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
import "github.com/pelletier/go-toml"
import "github.com/satori/go.uuid"

type Service struct {}
type Endpoint struct {}
type Provider struct {}
type Context struct {}

type BuildService func (*Service)
type BuildEndpoint func (*Endpoint)
type BuildProvider func (*Provider)

type App struct {

    // Slug is a short name that identifies the application instance.
    // It is advised to keep it machine & human readable: in a form of
    // of a slug - no spaces, all lower case, et cetera. The framework
    // itself, as well as any other code could use this variable to
    // unique identify an instance of the running application.
    Slug string

    // Complement the application slug; represents a version of the
    // running application instance. The version format should conform
    // to the semver format, like 0.0.1. The version should be strict,
    // as in it should not containt undefined masks like 1.1.x. It is
    // up to the user to assign meaningul semantics to app versions.
    Version string

    // A path within the local file system where an instance of the
    // running application should be residing. The framework will use
    // this path to lookup configuration directories, optional static
    // assets and a number of other things it may need. By default, it
    // will be set to the CWD directory that the app was launched in.
    Root string

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
    Instance uuid.UUID

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

    // Slice of HTTP servers that will be used to server application
    // instance. Servers are automatically created by the framework
    // for every corresponding section in the config file. This is
    // needed for applications that must be served on multiple ports
    // or network interfaces at the same time, within one process.
    Servers []*http.Server
}
