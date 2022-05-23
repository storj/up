
# docker compose based Storj environment

`storj-up` is a swiss-army tool to create / customize Storj clusters with the help of `docker compose` (not just storagenode but all satellite and edge services).

This is useful for Storj development and not for Storage Node operator.

## Getting started

It is recommended that you use [Docker Compose V2](https://docs.docker.com/compose/cli-command/) as well as [enable Docker BuildKit builds](https://docs.docker.com/develop/develop-images/build_enhancements/) for all features to work correctly.

Install the tool:

```
go install storj.io/storj-up@latest
```

Go to an empty directory and initialize docker compose:

```
storj-up init
```

Start the cluster:

```
docker compose up -d
docker compose ps
```

You can check the generated credentials with:

```
storj-up credentials
```

You can set the required environment variables with `eval $(storj-up credentials -e)` (at least on Linux/OSX)

Or you can update the credentials of local `rclone` setup with `storj-up credentials -w`

## More features

There are dedicated subcommands to modify the `docker-compose.yaml` file easily. The generic form of these commands:

```
storj-up <subcommand> <selector> <argument>
```

Here `selector` can be either a service (like `storagenode`) or a name of a service group. (like `edge`). To find out all the groups, please use `storj-up services` 

### Example: Building specific binaries based on a Gerrit change

After running `storj-up init`, you can use the following command to replace binaries based on a specific Gerrit changeset:

```
storj-up build remote gerrit -f refs/changes/65/6365/1 satellite-api satellite-core satellite-admin uplink versioncontrol
```

You will need to change `refs/changes/65/6365/1` to the Gerrit patchset you want to use, and change the binaries that follow it based on what you are trying to replace.

Then, run `docker compose build` followed by `docker compose up` in order to spin everything up.

### Example: Modify the configuration variable of a service

You can modify configuration variable by setting the environment variables:

```
storj-up env setenv satellite-api STORJ_CONSOLE_CREDENTIALS_REQUEST_URL=http://myservername
```

Available configuration (for selected services) can be listed by `storj-up configs <service>`, for example:

```
storj-up configs storagenode

STORJ_IDENTITY_CERT_PATH                                               path to the certificate chain for this identity (default: $IDENTITYDIR/identity.cert)
STORJ_IDENTITY_KEY_PATH                                                path to the private key for this identity (default: $IDENTITYDIR/identity.key)
STORJ_SERVER_CONFIG_REVOCATION_DBURL                                   url for revocation database (e.g. bolt://some.db OR redis://127.0.0.1:6378?db=2&password=abc123) *(default: bolt://$CONFDIR/revocations.db)
STORJ_SERVER_CONFIG_PEER_CAWHITELIST_PATH                              path to the CA cert whitelist (peer identities must be signed by one these to be verified). this will override the default peer whitelist
...
```

### Example: Using your local satellite installation rather than a remote change

#### Backend

After running `storj-up init`, you can use the following command to replace binaries from your local machine:

**On Linux:**

This will mount the correct binaries from your `$GOPATH/bin`

```
storj-up local-bin satellite-core satellite-admin satellite-api
```

**Mac and Windows:**

This will mount the correct binaries from your `$GOPATH/bin/linux_amd64`

```
storj-up local-bin -s linux_amd64 satellite-core satellite-admin satellite-api
```

You will also need to cross-compile to Linux when you update your local satellite, e.g.

```
GOOS=linux GOARCH=amd64 go install ./cmd/satellite
```

Then if you are not currently running the containers, run

```
docker compose up -d
```

to start the containers.

Or run 

```
docker restart up-satellite-core-1 up-satellite-api-1 up-satellite-admin-1
```

(the "up" prefix may be different depending on the location of your docker-compose.yaml file)

to restart already-running containers.

#### Frontend

Here, you will need to attach your local web/satellite directory as a volume. Do this with

```
storj-up local-ws /path/to/storj/web/satellite/
```

When you run `npm run build` from your local web/satellite directory, the webapp should be automatically updated, no need to restart any docker containers.

The exception is if you are making a frontend change in web/satellite that requires a corresponding backend change. In this case, you will need to also run `go install ./cmd/satellite` followed by a restart of the relevant containers (see command at the end of the "Backend" section above).

### Interacting with and resetting your Satellite database

`docker compose ps` will list your running containers. Find the one that looks like `<prefix>-cockroach-1`

To run sql on this container,

```
docker exec -it <prefix>-cockroach-1 ./cockroach sql --insecure 
```

`show databases;` will list all the databases you can query from. `master` will contain most satellite tables, and `metainfo` contains... metainfo tables.

`use <database>;` will switch to one of those.

`show tables;` will give a list of tables accessible from the selected database.

Then you can run queries like `update users set project_limit=3 where ...;`

#### Resetting

There is a chance that due to going back and forth between database versions will result in errors that look like this in your logs:

```
up-satellite-api-1    | 2022-05-16T17:34:53.916Z        DEBUG   process/exec_conf.go:403        Unrecoverable error     {"error": "Error checking version for satellitedb: validate db version mismatch: expected 196 != 195\n\tstorj.io/storj/private/migrate.(*Migration).ValidateVersions:138\n\tstorj.io/storj/satellite/satellitedb.(*satelliteDB).CheckVersion:138
```

If you are okay with starting with a fresh satellite database, this can be accomplished by running

```
docker compose down -v
```
