# USER microservice
 ![Build Status](https://travis-ci.org/thomaspoignant/user-microservice.svg?branch=master) [![Coverage Status](https://coveralls.io/repos/github/thomaspoignant/user-microservice/badge.svg?branch=master)](https://coveralls.io/github/thomaspoignant/user-microservice?branch=master) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fthomaspoignant%2Fuser-microservice.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fthomaspoignant%2Fuser-microservice?ref=badge_shield)


*user-microservice* is a set of API to manage users.  
This project is written in **GO**, it store data in **dynamodb**.

## Start the project

You can build the project by using the command 
``` shell
make build
```
After that you can run the project by using _(please configure you own env variables)_ :
``` shell
export GIN_MODE: debug \
    && export APP_PORT: 8080 \
    && export RUNNING_MODE: api \
    && export DYNAMODB_ENDPOINT:http://localhost:9000 \
    && export AWS_REGION: eu-west-1 \
    && export DYNAMODB_TABLE_NAME: user \
    && ./user-microservice
```

## If you want to use a local dynamodb
Run **dynamodb-local** on docker 
```shell
docker run -d -p 9000:8000 amazon/dynamodb-local
```

Create a dynamodb table _(you need AWS cli)_
```shell
aws dynamodb create-table --table-name user \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url http://localhost:9000
```

## Use swagger
Swagger is configure in the app you can access to the APIs to the URL : ```http://localhost:8080/swagger/index.html```
