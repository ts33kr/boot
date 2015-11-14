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

// Function that is used to build up a provider instance. It takes a
// pointer to the provider that has been pre-allocated and preliminary
// initialized before invoking the maker function, passing it through.
// provider makers are going to be invoked during application launch.
// Please refer to the provider API for more information on usage.
type ProviderBuilder func (*Provider)

// Provider is an entity that proviedes some sort of functionality
// for the application. Good example of this is a provider that could
// provide a DB connection for application, by consuming the app config
// data, opening a connection and then making the DB object accessible
// via the application storage mechanism. Use the API to create one.
type Provider struct {

    // Description of the provider; it should be a short and succinct
    // synopsis of what this provider does, as a human readable string.
    // Keep it short yet descriptive enough to understand a basic idea
    // of what this provider is intended for. This field should be set
    // via corresponding API; please do not modify this directly.
    About string

    // Implementation of the provider. It should be the Setup typed
    // function that implements the business logic of the functionality
    // offered by the provider. It will be invoked during the app launch
    // process, before the application is actually spinned up with all
    // its services and endpoints. Please set it via special API.
    Setup UnbiasedLogic

    // Optional function that takes care of cleaning up the provider
    // related resource that might have been allocated or opened during
    // invoking the provider setup function. Cleanup function will be
    // automatically called when the application will be terminating.
    // If there is no cleanup function - nil value should be set.
    Cleanup UnbiasedLogic

    // Instant in time when the provider was invoked. The nil value
    // should indicate that current provider instance has not yet been
    // invoked. This value is used internally by the framework in the
    // multiple of ways; and may also be used by whoever is interested
    // the time of when, and if, the provider was invoked.
    Invoked time.Time
}
