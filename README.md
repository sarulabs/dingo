![DINGO](https://raw.githubusercontent.com/sarulabs/assets/master/dingo/logo.png)

Generation of dependency injection containers for go programs (golang).

**Dingo** is a code generator. It generates dependency injection containers based on [sarulabs/di](https://github.com/sarulabs/di).

It is better than [sarulabs/di](https://github.com/sarulabs/di) alone because:

- Generated containers have **typed** methods to retrieve each object. You do not need to cast them before they can be used. That implies less runtime errors.
- Definitions are easy to write. Some dependencies can be guessed, allowing **shorter definitions**.

The disadvantage is that the code must be generated. But this can be compensated by the use of a file watcher.

# Table of contents

[![Build Status](https://travis-ci.org/sarulabs/dingo.svg?branch=master)](https://travis-ci.org/sarulabs/dingo)
[![GoDoc](https://godoc.org/github.com/sarulabs/dingo?status.svg)](http://godoc.org/github.com/sarulabs/dingo)
[![codebeat badge](https://codebeat.co/badges/833d6016-e4dd-482f-bcfe-210a1be48b94)](https://codebeat.co/projects/github-com-sarulabs-dingo-master)
[![goreport](https://goreportcard.com/badge/github.com/sarulabs/dingo)](https://goreportcard.com/report/github.com/sarulabs/dingo)

- [Dependencies](#dependencies)
- [Similarities with di](#similarities-with-di)
- [Setup](#setup)
    * [Code structure](#code-structure)
    * [Generating the code](#generating-the-code)
- [Definitions](#definitions)
    * [Name and scope](#name-and-scope)
    * [Build based on a structure](#build-based-on-a-structure)
    * [Build based on a function](#build-based-on-a-function)
    * [Parameters](#parameters)
    * [Close function](#close-function)
    * [Avoid automatic filling](#avoid-automatic-filling)
- [Generated container](#generated-container)
    * [Basic container](#basic-container)
    * [Additional methods](#additional-methods)
    * [Logging errors](#logging-errors)
    * [C function](#c-function)
    * [Retrieval functions](#retrieval-functions)
- [Upgrade from v3](#upgrade-from-v3)

# Dependencies

This module depends on `github.com/sarulabs/di/v2`. You will need it to generate and use the dependency injection container.

# Similarities with di

Dingo is very similar to [sarulabs/di](https://github.com/sarulabs/di) as it mainly a wrapper around it. This documentation mostly covers the differences between the two libraries. You probably should read the di documentation before going further.

# Setup

## Code structure

You will have to write the service definitions and register them in a `Provider`. The recommended structure to organize the code is the following:

```txt
- services/
    - provider/
        - provider.go
    - servicesA.go
    - servicesB.go
```

In the service files, you can write the service definitions:

```go
// servicesA.go
package services

import (
    "github.com/sarulabs/dingo/v4"
)

var ServicesADefs = []dingo.Def{
    {
        Name: "definition-1",
        // ...
    },
    {
        Name: "definition-2",
        // ...
    },
}
```

In the provider file, the definitions are registered with the `Load` method.

```go
// provider.go
package provider

import (
    "github.com/sarulabs/dingo/v4"
    services "YOUR_OWN_SERVICES_PACKAGE"
)

// Redefine your own provider by overriding the Load method of the dingo.BaseProvider.
type Provider struct {
    dingo.BaseProvider
}

func (p *Provider) Load() error {
    if err := p.AddDefSlice(services.ServicesADefs); err != nil {
        return err
    }
    if err := p.AddDefSlice(services.ServicesBDefs); err != nil {
        return err
    }
    return nil
}
```

An example of this can be found in the `tests` directory.

## Generating the code

You will need to create your own command to generate the container. You can adapt the following code:

```go
package main

import (
    "fmt"
    "os"

    "github.com/sarulabs/dingo/v4"
    provider "YOUR_OWN_PROVIDER_PACKAGE"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run main.go path/to/output/directory")
        os.Exit(1)
    }

    err := dingo.GenerateContainer((*provider.Provider)(nil), os.Args[1])
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
}
```

Running the following command will generate the code in the `path/to/generated/code` directory:

```sh
go run main.go path/to/generated/code
```

### Custom package name

If you want to customize the package name for the generated code you can use `dingo.GenerateContainerWithCustomPkgName` instead of `dingo.GenerateContainer`.

```go
// The default package name is "dic", but you can replace it with the name of your choosing.
// Go files will be generated in os.Args[1]+"/dic/".
err := dingo.GenerateContainerWithCustomPkgName((*provider.Provider)(nil), os.Args[1], "dic")
```

# Definitions

## Name and scope

Dingo definitions are not that different from [sarulabs/di](https://github.com/sarulabs/di) definitions.

They have a `name` and a `scope`. For more information about scopes, refer to the documentation of [sarulabs/di](https://github.com/sarulabs/di).

The default scopes are `di.App`, `di.Request` and `di.SubRequest`.

 The `Unshared` field is also available (see [sarulabs/di unshared objects](https://github.com/sarulabs/di#unshared-objects)).

## Build based on a structure

`Def.Build` can be a pointer to a structure. It defines the type of the registered object.

```go
type MyObject struct{}

dingo.Def{
    Name: "my-object",
    Build: (*MyObject)(nil),
}
```

You can use a nil pointer, like `(*MyObject)(nil)`, but it is not mandatory. `&MyObject{}` is also valid.

If your object has fields that must be initialized when the object is created, you can configure them with `Def.Params`.

```go
type OtherObject struct {
    FieldA *MyObject
    FieldB string
    FieldC int
}

dingo.Def{
    Name: "other-object",
    Build: (*OtherObject)(nil),
    Params: dingo.Params{
        "FieldA": dingo.Service("my-object"),
        "FieldB": "value",
    },
}
```

`dingo.Params` is a `map[string]interface{}`. The key is the name of the field. The value is the one that the associated field should take.

You can use `dingo.Service` to use another object registered in the container.

Some fields can be omitted like `FieldC`. In this case, the field will have the default value `0`. But it may have a different behaviour. See the [parameters](#parameters) section to understand why.

## Build based on a function

`Def.Build` can also be a function. Using a pointer to a structure is a simple way to declare an object, but it lacks flexibility.

To declare `my-object` you could have written:

```go
dingo.Def{
    Name: "my-object",
    Build: func() (*MyObject, error) {
        return &MyObject{}, nil
    },
}
```

It is similar to the `Build` function of [sarulabs/di](https://github.com/sarulabs/di), but without the container as an input parameter, and with `*MyObject` instead of `interface{}` in the output.

To build the `other-object` definition, you need to use the `my-object` definition. This can be achieved with `Def.Params`:

```go
dingo.Def{
    Name: "other-object",
    Build: func(myObject *MyObject) (*OtherObject, error) {
        return &OtherObject{
            FieldA: myObject,
            FieldB: "value",
        }, nil
    },
    Params: dingo.Params{
        "0": dingo.Service("my-object"),
    },
}
```

The build function can actually take as many input parameters as needed. In `Def.Params` you can define their values.

The key is the index of the input parameter.

## Parameters

As explained before, the key of `Def.Params` is either the field name (for build structures) or the index of the input parameter (for build functions).

When an item is not defined in `Def.Params`, there are different situations:

- If there is exactly one service of this type also defined in the container, its value is used.
- If there is none, the default value for this type is used.
- If there are more than one service with this type, the container can not be compiled. You have to specify the value for this parameter.

With these properties, it is possible to avoid writing some parameters. They will be automatically filled. This way you can have shorter definitions:

```go
dingo.Def{
    Name: "other-object",
    Build: func(myObject *MyObject) (*OtherObject, error) {
        return &OtherObject{
            FieldA: myObject, // no need to write the associated parameter
            FieldB: "value",
        }, nil
    },
}
```

It works well for specific structures. But for basic types it can become a little bit risky. Thus it is better to only store pointers to structures in the container and avoid types like `string` or `int`.

It is possible to force the default value for a parameter, instead of using the associated object. You have to set the parameter with `dingo.AutoFill(false)`:

```go
dingo.Def{
    Name: "other-object",
    Build: (*OtherObject)(nil),
    Params: dingo.Params{
        // *MyObject is already defined in the container,
        // so you have to use Autofill(false) to avoid
        // using this instance.
        "FieldA": dingo.AutoFill(false),
    },
}
```

## Close function

Close functions are identical to those of [sarulabs/di](https://github.com/sarulabs/di). But they are typed. No need to cast the object anymore.

```go
dingo.Def{
    Name: "my-object",
    Build: func() (*MyObject, error) {
        return &MyObject{}, nil
    },
    Close: func(obj *MyObject) error {
        // Close object.
        return nil
    }
}
```

## Avoid automatic filling

Each definition in the container is a candidate to automatically fill another (if its parameters are not specified).

You can avoid that with `Def.NotForAutoFill`:

```go
dingo.Def{
    Name: "my-object",
    Build: (*MyObject)(nil),
    NotForAutoFill: true,
}
```

This can be useful if you have more than one object of a given type, but one should be used by default to automatically fill the other definitions. Use `Def.NotForAutoFill` on the definition you do not want to use automatically.

# Generated container

## Basic container

The container is generated in the `dic` package inside the destination directory. The container is more or less similar to the one from [sarulabs/di](https://github.com/sarulabs/di).

It implements this interface:

```go
interface {
    Scope() string
    Scopes() []string
    ParentScopes() []string
    SubScopes() []string
    Parent() *Container
    SubContainer() (*Container, error)
    SafeGet(name string) (interface{}, error)
    Get(name string) interface{}
    Fill(name string, dst interface{}) error
    UnscopedSafeGet(name string) (interface{}, error)
    UnscopedGet(name string) interface{}
    UnscopedFill(name string, dst interface{}) error
    Clean() error
    DeleteWithSubContainers() error
    Delete() error
    IsClosed() bool
}
```

To create the container, there is the `NewContainer` function:

```go
func NewContainer(scopes ...string) (*Container, error)
```

You need to specify the scopes. By default `di.App`, `di.Request` and `di.SubRequest` are used.

A `NewBuilder` function is also available. It allows you to redefine some services (`Add` and `Set` methods) before generating the container with its `Build` method. It is not recommended but can be useful for testing.

## Additional methods

For each object, four other methods are generated. These methods are typed so it is probably the one you will want to use.

They match the `SafeGet`, `Get`, `UnscopedSafeGet` and `UnscopedGet` methods. They have the name of the definition as suffix.

For example for the `my-object` definition:

```go
interface {
    SafeGetMyObject() (*MyObject, error)
    GetMyObject() *MyObject
    UnscopedSafeGetMyObject() (*MyObject, error)
    UnscopedGetMyObject() *MyObject
}
```

`my-object` has been converted to `MyObject`.

The name conversion follow these rules:

- only letters and digits are kept
- it starts with an uppercase character
- after a character that is not a letter or a digit, there is another uppercase character

For example `--apple--orange--2--PERRY--` would become `AppleOrange2PERRY`.

Note that you can not have a name beginning by a digit.

## C function

There is also a `C` function in the dic package. Its role is to turn an interface into a `*Container`.

By default, the `C` function can:

- cast a container into a `*Container` if it is possible
- retrieve a `*Container` from the context of an `*http.Request` (the key being `dingo.ContainerKey("dingo")`)

This function can be redefined to fit your use case:

```go
dic.C = func(i interface{}) *Container {
    // Find and return the container.
}
```

The `C` function is used in retrieval functions.

## Retrieval functions

For each definition, another function is generated.

Its name is the formatted definition name. It takes an interface as input parameter, and returns the object.

For `my-object`, the generated function would be:

```go
func MyObject(i interface{}) *MyObject {
    // ...
}
```

The generated function uses the `C` function to retrieve the container from the given interface. Then it builds the object.

It can be useful if you do not have the `*Container` but only an interface wrapping the container:

```go
type MyContainer interface {
    // Only basic untyped methods.
    Get(string) interface{}
}

func (c MyContainer) {
    obj := c.Get("my-object").(*MyObject)

    // or

    obj := c.(*dic.Container).GetMyObject()

    // can be replaced by

    obj := dic.MyObject(c)
}
```

It can also be useful in an http handler. If you add a middleware to store the container in the request context with:

```go
// Create a new request req, which is like request r but with the container in its context.
ctx := context.WithValue(r.Context(), dingo.ContainerKey("dingo"), container)
req := r.WithContext(ctx)
```

Then you can use it in the handler:

```go
func (w http.ResponseWriter, r *http.Request) {
    // The function can find the container in r thanks to dic.C.
    // That is why it can create the object.
    obj := dic.MyObject(r)
}
```

# Upgrade from v3

- You need to register the definitions in a `Provider`. The `dingo` binary is also not available anymore as it would depends on the `Provider` now. So you have to write it yourself. See the [Setup](#setup) section.
- `dingo.App`, `dingo.Request` and `dingo.SubRequest` have been removed. Use `di.App`, `di.Request` and `di.SubRequest` instead.
