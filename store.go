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

import "sync"

// A general purpose storage that is intended as a frequently and
// easily embeddable piece of functionality used by many structures
// within the framework. This intended as a simple, in-memory storage
// exposing the key/value data model to its user. Values are all typed
// as empty interface. Storage intended to hold few (not many) records.
type Storage struct {

    // The underlying data holder. A simple map of strings to untyped
    // values. In order to be useful, values have to be casted (or type
    // asserted) after they have been retrieved from the container. All
    // access to the underlying container has to be regulated using by
    // using the appropriate techniques to manage concurrent access.
    Container map[string] interface {}

    // A read-write mutually exclusive lock that is embedded in the
    // storage structure for synchronizing concurrent access to the
    // underlying data holder of the storage instance. This mutex must
    // be used around all the data access operations to ensure that
    // contained data does not get corrupted by concurrent access.
    sync.RWMutex
}
