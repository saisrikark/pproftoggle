# pproftoggle - in development
run pprof without restarting your application

## Table of Contents

- [Description](#description)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Description

[pprof](https://github.com/google/pprof) is a tool to view resources used by go applications. <br>
To use it we must host a http server. <br><br>
Often to switch it on, users are forced to restart their application.
Resource consumption scenarioes are hard to replicate due to this restart. <br>
In general we would prefer not to host an extra server if not required. <br><br>
This package helps us to toggle on and off the pprof server without restarting. <br>
Some preexisting rules are provided which are implemented as part of an interface. <br>

## Features

Toggle pprof <br>

Toggle using
- Environment Variable
    - Based on the existence of an environment variable
- more in development

## Usage

TODO

## Contributing

Just me for now 🙂

## License

The code in this repository is licensed under the terms of the Apache License 2.0.
