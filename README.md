#Overview
This project is a yet another micro framework developed in Go. What
stands out in this case, is it's probably the only web framework for
Go that is both - lightweight and ready for writing production code
right out of the box. It provides zero-configuration & hassle free
essentials that are absolutely necessary for any production app, but
which are missing from most of the web frameworks for Go. The goal
of the framework to let you write your application logic since 1-st
line of code, without having to code all the usual boilerplate code
that is necessary to power and a production app and its deployment.
Here is a list of some features support out of the box.
  
* Switching app environments: dev, production, staging, etc
* Automatically load a config file based on the environment
* Declarative API for endpoint definition, with metadata
* Serving inventory of available APIs defined in the app
* Structured logger is automatically available in the app
* Middleware support for the endpoints and endpoint groups
* High performance and as efficient memory consumption
* Efficiently serving static resources (asset files)
