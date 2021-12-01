
# docker-compose based Storj environment

`storj-up` is a swiss-army tool to create / customize Storj clusters with the help of `docker-compose` (not just storagenode but all satellite and edge services).

This is useful for Storj development and not for Storage Node operator.

## Getting started

You may need to [enable Docker BuildKit builds](https://docs.docker.com/develop/develop-images/build_enhancements/) for certain features to work correctly (e.g. creating an environment based partially on a Gerrit change).

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

### Example: Building specific binaries based on a Gerrit change

After running `storj-up init`, you can use the following command to replace binaries based on a specific Gerrit changeset:

```
storj-up build remote gerrit -f refs/changes/65/6365/1 satellite-api satellite-core satellite-admin uplink versioncontrol
```

You will need to change `refs/changes/65/6365/1` to the Gerrit patchset you want to use, and change the binaries that follow it based on what you are trying to replace.

Then, run `docker-compose build` followed by `docker-compose up` in order to spin everything up.

### Modify the configuration variable of a service

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