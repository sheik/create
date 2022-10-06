# Create #
Create is a language agnostic build tool that allows you to define steps and dependencies.

## Using Create in Your Project ##
In order to use create, first go into your project directory and execute the following:

    go install github.com/sheik/create/cmd/create@latest
    create update 

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


## Example Createfile ##
Click to see a [createfile example](https://github.com/sheik/create/blob/main/cmd/createfile/createfile.go)

## Building Your Project ##
Once you have installed "create" into your project and created a createfile, building your is as simple as running:

    user@host:/home/user/myProject$ create <target>

You can run it without specifying a target, in which case it will run whichever target has "Default" set to "true".

## License ##
This project is licensed under the terms of the MIT license. See [LICENSE.md](https://github.com/sheik/create/blob/main/LICENSE.md) for details.
