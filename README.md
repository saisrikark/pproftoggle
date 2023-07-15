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
Resource consumption scenarios are hard to replicate due to this restart. <br>
We would prefer not to host an extra server if not required. <br><br>
This package helps us to toggle the pprof server without restarting. <br>
Some preexisting rules are provided which are implemented behind interface. <br>

## Features

Provide your own '''*http.Server''' we will use it to host a http server while overwritting the handler <br>

Toggle using
- Environment Variable
    - Based on the key value pair match of an environment variable.
- Simple yaml file
    - Yaml file key value pair match only for elements present at the root.
- Implement your own rule
    - Rules are behind an interface so you can try your own implementation.

## Usage

To use <br>
```
go get -d github.com/saisrikark/pproftoggle
```

Below is an example using a simple yaml rule. <br>
It will read a yaml file at the specified path and check if the key value pair match. <br>
On a match pprof will be server via http. <br>
When unmatched pprof will no longer be served. <br>

```go
package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/saisrikark/pproftoggle"
	"github.com/saisrikark/pproftoggle/rules"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	toggler, err := pproftoggle.NewToggler(
		pproftoggle.Config{
			PollInterval: time.Second * 1,
			HttpServer: &http.Server{
				Addr: ":8080",
			},
			Rules: []pproftoggle.Rule{
				rules.SimpleYamlRule{
					Key:   "enablepprof",
					Value: "true",
					Path:  "example.yaml",
				},
			},
		})
	if err != nil {
		log.Println("received error while trying to create new toggler", err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := toggler.Serve(ctx); err != nil {
			log.Println("received error while trying to serve using toggler", err)
		}
	}()
	time.Sleep(time.Minute * 10)

	cancel()
	wg.Wait()
}
```

Attempt to call the server. <br>
```
curl http://localhost:8080/debug/pprof/heap
```

Refer to the **examples** folder for better examples. <br>
**NOTE**: to run examples main.go must be run from the directory it is present in. <br>

## Contributing

Just me for now 🙂

## License

The code in this repository is licensed under the terms of the Apache License 2.0.
