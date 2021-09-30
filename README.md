
#

This is a quick prototype of running Storj cluster in docker. (Not just the storagenodes but full cluster).

Today it depends on my own local containers (see `rebuild.sh`) but it may work with any official image...


## Quick start

Just start the cluster:

```
#remove previous state
docker-compose down
docker-compose up -d
```

Check the log files:

```
docker-compose logs satellite-api
```

Use the cluster:

```
docker-compose exec satellite-api bash
devrun credentials satellite-api test@mailinator.com
export STORJ_ACCESS=...

uplink mb sj://bucket1

uplink share --auth-service http://authservice:8000 --url --not-after=none sj://bucket1/file1
```

## Modifications

Use local version from any of the services:

```
volumes:
    - /home/elek/go/bin/storagenode:/var/lib/storj/go/bin/storagenode

```


Local development if entrypoint can be done with:

```
volumes:
    - ./runner/entrypoint.sh:/var/lib/storj/entrypoint.sh
```


Remote debug is can be turned on with

```
environment:
   GO_DLV: "true"
```
