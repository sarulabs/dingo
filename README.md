![DINGO](https://raw.githubusercontent.com/sarulabs/assets/master/dingo/logo.png)

Generation of dependency injection containers for go programs (golang).

**Dingo** is a code generator. It generates dependency injection containers based on [sarulabs/di](https://github.com/sarulabs/di).

It is better than [sarulabs/di](https://github.com/sarulabs/di) alone because:
- Generated containers have **typed** methods to retrieve each object. You do not need to cast them before they can be used. That implies less runtime errors.
- Definitions are easy to write. The generated code is also **safer** compared to what you may write yourself. That also means that the code is **easier to debug**.
- Some dependencies can be guessed, allowing **shorter definitions**.

The disadvantage is that the code must be generated. But this can be compensated by the use of a file watcher.



# Table of contents

[![Build Status](https://travis-ci.org/sarulabs/dingo.svg?branch=master)](https://travis-ci.org/sarulabs/dingo)
[![GoDoc](https://godoc.org/github.com/sarulabs/dingo?status.svg)](http://godoc.org/github.com/sarulabs/dingo)
[![codebeat badge](https://codebeat.co/badges/833d6016-e4dd-482f-bcfe-210a1be48b94)](https://codebeat.co/projects/github-com-sarulabs-dingo-master)
[![goreport](https://goreportcard.com/badge/github.com/sarulabs/dingo)](https://goreportcard.com/report/github.com/sarulabs/dingo)

- [Installation](#installation)
- [Usage](#usage)
- [Similarities with di](#similarities-with-di)
- [Definitions](#definitions)
    * [Code structure](#code-structure)
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



# Installation

You can use `go get` to install the latest dingo binary:

```sh
go get -u github.com/sarulabs/dingo/dingo
```

In your project, you need to import `github.com/sarulabs/dingo` to write your objects definitions. You can use `go get` or a dependency manager like `vgo`:

```sh
go get github.com/sarulabs/dingo
```



# Usage

There is only one way to use `dingo`:

```sh
dingo -src="$projectDir/services" -dest="$projectDir/generatedServices"
````

- **src**: the directory containing your definitions
- **dest**: the directory where the code will be generated

The only requirement is that these two directories must be in your `$GOPATH`.

There is no watch mode included in dingo. You need to bring your own file watcher.



# Similarities with di

Dingo is very similar to [sarulabs/di](https://github.com/sarulabs/di). This documentation mostly covers the differences between the two libraries. You probably should read the di documentation before going further.



# Definitions

## Code structure

The definitions must be in a single directory. They must be contained in go files, and the package name must be the same as the directory name.

A definition is a `dingo.Def` from the `github.com/sarulabs/dingo` package.

The definitions must be contained in public variables. The allowed types for these variables are:

- `dingo.Def`
- `*dingo.Def`
- `[]dingo.Def`
- `[]*dingo.Def`
- `func() dingo.Def`
- `func() *dingo.Def`
- `func() []dingo.Def`
- `func() []*dingo.Def`

A definition file should look like this:

```go
package services

import "github.com/sarulabs/dingo"

var MyDefinitions = []dingo.Def{
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

## Name and scope

Dingo definitions are not that different from [sarulabs/di](https://github.com/sarulabs/di) definitions.

They have a `name` and a `scope`. For more information about scopes, refer to the documentation of [sarulabs/di](https://github.com/sarulabs/di).

The default scopes are `dingo.App`, `dingo.Request` and `dingo.SubRequest`.


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
````

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
````

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
````

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
````


## Close function

Close functions are identical to those of [sarulabs/di](https://github.com/sarulabs/di). But they are typed. No need to cast the object anymore.

```go
dingo.Def{
    Name: "my-object",
    Build: func() (*MyObject, error) {
        return &MyObject{}, nil
    },
    Close: func(obj *MyObject) {
        // Close object.
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
    func Scope() string {
    func Scopes() []string {
    func ParentScopes() []string {
    func SubScopes() []string {
    func Parent() *Container {
    func SubContainer() (*Container, error) {
    func SafeGet(name string) (interface{}, error) {
    func Get(name string) interface{} {
    func UnscopedSafeGet(name string) (interface{}, error) {
    func UnscopedGet(name string) interface{} {
    func Clean() {
    func DeleteWithSubContainers() {
    func Delete() {
    func IsClosed() bool {
}
````

To create the container, there is the `NewContainer` function:

```go
func NewContainer(scopes ...string) (*Container, error)
```

You need to specify the scopes. By default `dingo.App`, `dingo.Request` and `dingo.SubRequest` are used.


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
````

`my-object` has been converted to `MyObject`.

The name conversion follow these rules:

- only letters and digits are kept
- it starts with an uppercase character
- after a character that is not a letter or a digit, there is another uppercase character

For example `--apple--orange--2--PERRY--` would become `AppleOrange2PERRY`.

Note that you can not have a name beginning by a digit.


## Logging errors

In the generated `dic` package, there is a function used to log the container errors.

```go
var ErrorCallback = func(err error) {
	log.Println(err.Error())
}
````

You can change its behavior and use your own logger:

```go
dic.ErrorCallback = func(err error) {
	// Use your own logger.
}
````


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
