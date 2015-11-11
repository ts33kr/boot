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
* Switching app environments: dev, production, staging, etc
* Automatically load a config (TOML) based on the environment
* Flexible URL routing with placeholders, hosts and patterns
* Declarative API defining the endpoints, with metadata
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

![arch](https://raw.github.com/ts33kr/boot/master/arch.png)
