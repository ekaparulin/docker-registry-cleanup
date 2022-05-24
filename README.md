# Algorithm

- Get all tags
- Get manifests for each tags
- Get create date from blobs
- Compare create date against expiry date
- Delete expired images

# Build 

```
cd src
go mod vendor
go build .
```

# Convert manifest schema version V1 to V2

From https://docs.docker.com/registry/spec/deprecated-schema-v1/:

One way to upgrade an image from image manifest version 2, schema 1 to schema 2 is to docker pull the image and then docker push the image with a current version of Docker. Doing so will automatically convert the image to use the latest image manifest specification.

Use "-list-schema-v1" flag to get a list of the images with schema 1 in order to perform pull&push

E.g.:

```
docker-registry-cleanup --registry https://REGISTRY -repo REPO -list-schema-v1 > list

for IMG in $(cat list); do docker pull $IMG; docker push $IMG; done
for IMG in $(cat list); do docker rmi $IMG; done
```