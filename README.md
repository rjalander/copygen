# Copygen

[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge&logo=appveyor&logo=appveyor)](https://pkg.go.dev/github.com/switchupcb/copygen)
[![Go Report Card](https://goreportcard.com/badge/github.com/switchupcb/copygen?style=for-the-badge)](https://goreportcard.com/report/github.com/switchupcb/copygen)
[![MIT License](https://img.shields.io/github/license/switchupcb/copygen.svg?style=for-the-badge)](https://github.com/switchupcb/copygen/blob/main/LICENSE)

Copygen is a command-line [code generator](https://github.com/gophersgang/go-codegen) that generates type-to-type and field-to-field struct code without adding any reflection or dependencies to your project. Manual-copy code generated by copygen is [**391x faster**](https://github.com/gotidy/copy#benchmark) than [jinzhu/copier](https://github.com/jinzhu/copier), and adds no allocation to your program. Copygen is the most customizable type-copy generator to-date and features a rich yet simple setup inspired by [goverter](https://github.com/jmattheis/goverter).

| Topic                           | Categories                                                                                    |
| :------------------------------ | :-------------------------------------------------------------------------------------------- |
| [Usage](#Usage)                 | [Types](#types), [Setup](#setup), [Command Line](#command-line), [Output](#output)            |
| [Customization](#customization) | [Custom Types](#custom-types), [Templates](#templates)                                        |
| [Matcher](#matcher)             | [Automatch](#automatch), [Depth](#depth)                                                      |
| [Optimization](#optimization)   | [Shallow Copy vs. Deep Copy](#shallow-copy-vs-deep-copy), [When to Use](#when-to-use-copygen) |

## Usage

Each example has a **README**.

| Example                                                                         | Description                                                       |
| :------------------------------------------------------------------------------ | :---------------------------------------------------------------- |
| main                                                                            | The default example.                                              |
| [manual](https://github.com/switchupcb/copygen/tree/main/examples/manual)       | Uses the manual map feature.                                      |
| [automatch](https://github.com/switchupcb/copygen/tree/main/examples/automatch) | Uses the automatch feature with depth _(doesn't require fields)_. |
| [new](https://github.com/switchupcb/copygen/tree/main/examples/new)             | Uses a new type to assist with type-conversion.                   |
| deepcopy _(Roadmap Feature)_                                                    | Uses the deepcopy option.                                         |
| [error](https://github.com/switchupcb/copygen/tree/main/examples/error)         | Uses templates to return an error (temporarily unsupported).      |

**NOTE: The following guide is set for v0.2 ([view v0.1](https://github.com/switchupcb/copygen/tree/v0.1.0))**

This [example](https://github.com/switchupcb/copygen/blob/main/examples/main) uses three type-structs to generate the `ModelsToDomain()` function.

### Types

`./domain/domain.go`

```go
// Package domain contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID     int
	UserID int
	Name   string
	Other  string // The other field is not used.
}
```

`./models/model.go`

```go
// Package models contains data storage models (i.e database).
package models

// Account represents the data model for account.
type Account struct {
	ID       int
	Name     string
	Password string
	Email    string
}

// A User represents the data model for a user.
type User struct {
	UserID   int
	Name     int
	UserData string
}
```

### Setup

Setting up copygen is a 2-step process involving a `YML` and `GO` file.

**setup.yml**

```yml
# Define where the code will be generated.
generated:
  setup: ./setup.go
  output: ./copygen.go
  package: copygen

# Define the optional custom templates used to generate the file.
templates:
  header: ./templates/header.go
  function: ./templates/function.go

# Define custom options for customization.
# Templates are passed to the generator options.
custom:
  option: The possibilities are endless.
```

The main example ignores the template fields.

**setup.go**

Create an interface in the specified setup file with a `type Copygen interface`. In each function, specify _the types you want to copy from_ as parameters, and _the type you want to copy to_ as return values.

```go
/* Copygen defines the functions that will be generated. */
type Copygen interface {
  // custom: see table below for options
  ModelsToDomain(models.Account, models.User) *domain.Account
}
```

Copygen uses no allocation with pointers which means fields are assigned to _objects passed as parameters_. In contrast, using a type with no pointer will return a copy of the new type.

**options**

You can specify options for your functions using comments. Do **NOT** put empty lines between comments that pertain to one function.

| Option              | Use                         | Description                                                                                                                                                                        | Example                                                                      |
| :------------------ | :-------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------- |
| `map from to`       | Manual Field Mapping        | Copygen uses its [automatcher](#automatch) default. <br /> Override this using `map` with _regex_ to identify <br /> fields that will be mapped to and from eachother.             | `map .* package.Type.Field` <br /> `map models.Account.ID domain.Account.ID` |
| `depth field level` | Use a specific field depth. | Copygen uses the full-field [depth](#depth) by default. <br /> Override this using `depth` with _regex_ <br /> and a [depth-level](#depth) integer.                                | `depth .* 2` <br /> `depth models.Account.* 1`                               |
| `deepcopy field`    | Deepcopy from-fields.       | Copygen shallow copies fields by default. <br /> Override this using `deepcopy` with _regex_. <br /> For more info, view [Shallow Copy vs. Deep Copy](#shallow-copy-vs-deep-copy). | `deepcopy package.Type.Field` <br /> `deepcopy .*` _(all fields)_            |
| `custom option`     | Specify custom options.     | You may want to use custom [templates](#templates). <br /> `custom` options are passed to a function's options. <br /> Returns `map[string][]string` _(trim-spaced)_.              | `ignore true` <br /> `swap false`                                            |

_[View a reference on Regex.](https://cheatography.com/davechild/cheat-sheets/regular-expressions/)_

#### Convert

In certain cases, you may want to specify a how a specific type or field is copied with a function. This can be done by defining a function with a `convert` option.
```go
/* Define the fields this converter is applied to using regex. If unspecified, converters are applied to all valid fields. */
// convert: models.User.ID
// comment: Itoa converts an integer to an ascii value.
func Itoa(i int) string {
  return strconv.Itoa(i)
}
```

### Command Line

Install the command line utility. Copygen is an executable and not a dependency, so use `go install`.

```
go install github.com/switchupcb/copygen@latest
```

Install a specific version by specifying a tag version.
```
go install github.com/switchupcb/copygen@v0.0.0
```

Run the executable with given options.

```bash
# Specify the .yml configuration file.
copygen -yml path/to/yml
```

_The path to the YML file is specified in reference to the current working directory._

### Output

This example outputs a `copygen.go` file with the specified imports and functions.

```go
// Code generated by github.com/switchupcb/copygen
// DO NOT EDIT.

package copygen

import (
	"github.com/switchupcb/copygen/examples/main/converter"
	"github.com/switchupcb/copygen/examples/main/domain"
	"github.com/switchupcb/copygen/examples/main/models"
)

// ModelsToDomain copies a User, Account to a Account.
func ModelsToDomain(tA *domain.Account, fU models.User, fA models.Account) {
	// Account fields
	tA.UserID = c.Itoa(fU.ID)
	tA.ID = fA.ID
	tA.Name = fA.Name

}
```

## Customization

Copygen's method of input and output allows you to generate code not limited to copying fields.

#### Custom Types

Custom types external to your application can be created for use in the `setup.go` file. When a file is generated, all types _(structs, interfaces, funcs)_ are copied **EXCEPT** the `type Copygen interface`.

```go
type DataTransferObject struct {
  // ...
}

type DataTransferObject interface {
  // ...
}

func ExternalFunc() {
  // ...
}
```

#### Templates

Templates can be created using **Go** to customize the generated code algorithm. The `copygen` generator uses the `package tenplates` `Header(*models.Generator)` to generate header code and `Function(*models.Function)` to generate code for each function. As a result, these _(package templates with functions)_ are **required** for your templates to work. View [models.Generator](https://github.com/switchupcb/copygen/blob/main/cli/models/function.go) and [models.Function](https://github.com/switchupcb/copygen/blob/main/cli/models/function.go) for context on the parameters passed to each function. Templates are interpreted by [yaegi](https://github.com/traefik/yaegi) which has limitations on module imports _(that are being fixed)_: As a result, **templates are temporarily unsupported.** The [error example](https://github.com/switchupcb/copygen/blob/main/examples/main) modifies the .yml to use **custom functions** which `return error`. This is done by modifying the .yml and creating **custom template files**.

## Matcher

Copygen provides two ways to configure fields: **Manually** and the **Automatcher**. Matching is specified in a `.go` file _(which functions as a schema in relation to other generators)_. Tags are complicated to use with other generators which is why they aren't used.

### Automatch

When fields aren't specified using options, copygen will attempt to automatch type-fields by name. Automatch **supports field-depth** (where types are located within fields) **and recursive types** (where the same type is in another type). Automatch loads types from Go modules _(in GOPATH)_. Ensure your modules are up to date by using `go get -u <insert/module/import/path>`.

#### Depth

The automatcher uses a field-based depth system. A field with a depth-level of 0 will only match itself. Increasing the depth-level allows its sub-fields to be matched. This system allows you to specify the depth-level for whole types **and** specific fields.

```go
// depth-level in relation to the first-level fields.
type Account
  // 0
  ID      int
  Name    string
  Email   string
  Basic   domain.T // int
  User    domain.DomainUser
              // 1
              UserID   string
              Name     string
              UserData map[string]interface{}
  // 0
  Log     log.Logger
              // 1
              mu      sync.Mutex
                          // 2
                          state   int32
                          sema    uint32
              // 1
              prefix  string
              flag    int
              out     io.Writer
                          // 2
                          Write   func(p []byte) (n int, err error)
              buf     []byte
```

## Optimization 

### Shallow Copy vs. Deep Copy
The library generates a [shallow copy](https://en.m.wikipedia.org/wiki/Object_copying#Shallow_copy) by default. An easy way to deep-copy fields with the same return type is by using `new()` as/in a converter function or by using a custom template.

### When to Use Copygen

Copygen's customizability gives it many potential usecases. However, copygen's main purpose is save you time by generating boilerplate code to map objects together.

#### Why would I do that?

In order to keep a program adaptable _(to new features)_, a program may contain two types of models. The first type of model is the **domain model** which is **used throughout your application** to model its business logic. For example, the [domain models of Copygen](https://github.com/switchupcb/copygen/tree/main/cli/models) focus on field relations and manipulation. In contrast, the ideal way to store your data _(such as in a database)_ may not match your domain model. In order to amend this problem, you create a **data model**. The [data models of Copygen](https://github.com/switchupcb/copygen/blob/main/cli/loader/models.go) are located in its loader(s). In many cases, you will need a way to map these models together to exchange information from your data-model to your domain _(and vice-versa)_. It's tedious to repeateadly do this in the application _(through assignment or function definitions)_. Copygen solves this problem.

## Contributing

You can contribute to this repository by viewing the [Project Structure, Code Specifications, and Roadmap](CONTRIBUTING.md).