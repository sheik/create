# Create #
Create is a language agnostic build tool that allows you to define steps and dependencies.

## Using Create in Your Project ##
In order to use create, first you need to install create:

     go install github.com/sheik/create/cmd/create@latest

Next, you need to make a "createfile". In order to do this, you need to create a package called "createfile" under your "cmd" directory in your project. For example:

	.
	├── cmd
	│   ├── createfile
	│   │   └── createfile.go
	│   └── myproject
	│       └── main.go
	├── go.mod
	├── go.sum
	├── internal
	│   └── utils
	│       └── utils.go
	└── pkg
	    └── shell
		└── shell.go


## Building ##
Run "make". This will install the create tool and run it to build. The resulting output is an rpm with the "create" tool in it.

## Example Createfile ##
See an [example here](https://github.com/sheik/create/blob/main/cmd/createfile/createfile.go)

## License ##
This project is licensed under the terms of the MIT license. See LICENSE.md for details.
