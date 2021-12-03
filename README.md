## Scratch - kit of projects templates

`Scratch` allows you to single time write a project template(scratch) and distribute all across your company or your needs.
Key feature of `scratch` is that all projects templates are stored in binary file, so to share `scratch` with someone you need just send a single binary file. 

## Adding New template

Scratch supports adding multiple project templates.
All scratches stored in `data/scratch/templates`

To new add template you should describe it in `manifest.yaml`, like an example

```yaml
name: go.cli
description: Golang CLI application with BuildKit support
rename:
    #   Ignored files '.' excluded from embedding, will be fixed with
    #   https://github.com/golang/go/issues/43854
    gitignore: .gitignore
    golangci.yml: .golangci.yml
    #    When using go:embed go tries to import go.mod as part of another module
    #    So rename to go_mod and go_sum to avoid this
    go_mod: go.mod
    go_sum: go.sum

```

Then add project template files in folder with name of project template

```text
go.cli
├── bin
│   └── gitignore
├── cmd
│   └── {{ProjectID}}
│       └── main.go
├── Dockerfile
├── golangci.yml
├── go_mod
├── go_sum
├── LICENSE
├── Makefile
├── README.md
└── rules
    └── builder.mk
```

Then build your `scratch`

## Build

Tools to build scratch

- `make` - To aggregate docker build commands
- [Docker BuildKit](https://github.com/moby/buildkit) - Improved build system
- [Docker BuildX](https://github.com/docker/buildx) - docker build plugin

After you install all of it just run
```shell
make
```

It will run all tests linter check and build binary to `/bin/scratch`


To run build **without any checks run**

```shell
make build-dry # may be useful for fast project building  
```