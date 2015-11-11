#Overview
This project is a yet another micro framework developed in Go. What
stands out in this case, is it's probably the only web framework for
Go that is both - lightweight and ready for writing production code,
right out of the box. It provides zero-configuration & hassle free
essentials that are absolutely necessary for any production app, but
which are missing from most of the web frameworks for Go. The goal
of the framework is to let you write your business logic since 1-st
line of code, without having to write all the usual boilerplate code
that is necessary to power a production application and its deployment.
Here is a list of some of the features supported out of the box.
  
* **Strong REST & JSON architecture with SOA approach**
* **Out of the box statistics and measurements essentials**
* Switching app environments: dev, production, staging, etc
* Automatically load a config (TOML) based on the environment
* Flexible URL routing with placeholders, hosts and patterns
* Serving inventory of available APIs defined in the app
* Optional built-in UI for browsing and testing the APIs
* Structured logger is automatically available in the app
* Middleware support for the endpoints and the services
* Shipped with essentials for authentication & access control
* Auto connecting to multiple DBs by just adding config sections
* Support for MongoDB & major SQL databases out of the box
* Support for Redis as session, cache & general purpose storage
* High performance and efficient memory consumption
* Efficiently serving static resources (asset files)

![arch](https://raw.github.com/ts33kr/boot/master/design/boot-arch.png)

The framework is designed to enable maximum code composition! This is
achieved by allowing any piece of code that does something useful to be
declared in a way that it can be reused by any other piece of code in an
application, whether it is an HTTP request handler or an auxiliary method.
Important point is that it keeps architecture of an application solid and
and predictable and allows for execution of any code within the app to be
done in a very controllable and segregated fashion.

The code will ultimately reside in different **boot.Service** instances;
its immediate implementation will usually be typed under **boot.BiasedLogic**
type, meaning a function that takes in the context. Any application code
will eventually be exposed via **boot.Operation** interface that knows how
to execute that code within a context of **boot.Context** structure. Under
the hood, you would use one of the two follwing structures to store the code.

* **boot.Endpoint** - is a piece of application's business logic that
is exposed through the HTTP interface directly. It has an HTTP verb and
an URL routing pattern, among other parameters. It also has access to
fetching the HTTP request with all its data, as well as an ability to
respond to it with a JSON encoded HTTP response.

* **boot.Aux** - is a piece of application's business logic that
is exposed as an auxilirary operation. It could be invoked by endpoints
or by another auxiliary operations. Operations in one service can easily
call auxiliary operations of other services. All the communication between
the caller and an aux operation should be implemented via the context.

![proc](https://raw.github.com/ts33kr/boot/master/design/boot-proc.png)
