# plugsys

Simple plugin system for go applications.

# example

```go
package main

// every plugin needs to provide ID method to identify plugin
type Plugin struct {
}

func(Plugin) ID() string { return "my-plugin" }

// instantiate plugin storage
storage := plugsys.New()

// add plugin instance
storage.Add(&Plugin{})

```

# Author
Peter Vrba