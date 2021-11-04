
# docker-compose based Storj environment

`storj-up` is a swiss-army tool to create / customize Storj clusters with the help of `docker-compose` (not just storagenode but all satellite and edge services).

This is useful for Storj development and not for Storage Node operator.

## Getting started

Install the tool:

```
go install storj.io/storj-up@latest
```

Go an empty directory an initialize docker-compose:

```
storj-up init
```

Start the cluster:

```
docker-compose up -d
docker-compose ps
```

You can check the generated credentials with:

```
storj-up credentials
```

You can set the required environment variables with `eval $(storj-up credentials -e)` (at least on Linux/OSX)

Or you can update the credentials of local `rclone` setup with `storj-up credentials -w`

## More features

There are dedicated subcommands to modify the `docker-compose` easily. The generic form of these commands:

```
storj-up <subcommand> <selector> <argument>
```

Here `selector` can be either a service (like `storagenode`) or a name of a service group. (like `edge`). To find out all the groups, please use `storj-up services` 

