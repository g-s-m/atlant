# Description
It is supposed to use k8s cluster for this service. Envoy uses as a load balancer for the service with round robin algorithm. Service uses free cluster of MongoDB provided by https://cloud.mongodb.com/

# How do build docker container
```
$ cd ./atlant
$ make container
```
To build service on your machine without docker install protoc and go 1.13+ before and run the following command
```
$ cd atlant
$ make
```

# How to deploy
Deploy in k8s cluster.
Make sure you have generic secret in your k8s cluster named mongo-pswd contained password to mongo db.
```
$ cd ./deploy
$ ./deploy.sh
```
This service has been deployed into gce k8s cluster. It is available through ip address 35.197.54.52:443

# Try with grpcurl
```
$ grpcurl -d '{"url":"http://my/csv/resource"}' -proto ./atlant/atlant/interface/service.proto -insecure -v 35.197.54.52:443 interface.ProductService/Fetch
```