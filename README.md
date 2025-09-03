# reviewService -  Ingestion service to ingest customer reviews in our system
  
## Setup Instructions

Install the following for running the setup in your local system:

- [Install Redis](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis-on-linux/) 

- [AMQP in local (RabbitMQ) 4.1.3](https://www.rabbitmq.com/docs/install-debian)  

- [MariaDB in local](https://mariadb.com/docs/server/server-installation/mariadb-package-repositories/)  
  
-  Add virtual host (review_dev) to rabbitMQ: 
   ```bash
   sudo rabbitmqctl add_vhost review_dev
    ```
  
-  Create RabbitMq service user and password  :
    ```bash
    sudo rabbitmqctl add_user reviewService <PWD>    (username , pwd)
    ```


-  Grant read and write permission to the user reviewService on the vhost review_dev  
   ```bash
    sudo rabbitmqctl set_permissions -p review_dev  reviewService  ".*" ".*" ".*"
    ```

- To run in local , create an .env file (will share over email) with the following variables,pls change the  mariabdb ,amqp  , redis port , credentials as per your setup

- Database and  tables creation in mariadb

  ```bash
            run the file /reviewService/deploy/db/tablecreation.sql
    ```


- Clone the project

```bash
  git clone https://github.com/Arpitmovers/reviewService
```

- Go to the project root ,install dependencies

```bash
  go mod tidy
```

- Start the server

```bash
  go run cmd/reviewService/main.go
```

-  reviewService/deploy/Dockerfile - file to build container of our service - working
- reviewService/.github/workflows/docker-image.yml  - pipeline to build and push a Docker image to GitHub Container Registry (GHCR). - working

## Design Decisions:
 - Tech Stack: Golang1.22, RabbitMQ 4.1.3, Redis, MariaDB 15.1
-  Asynchronous Processing: The API responds immediately after validating the request and discovering files in S3/
-  Actual file ingestion and publishing of each json record  happens in the background , the consumer is responsible for persisiting it to Datastore.
-  Worker Pool : to allow processing of multiple files concurrently , we are using  n go routines where n=min (no of cpus, no of files in s3)
-  Idempotency handling: using redis for locking using SETNX , so processing of same file is not triggered multiple times.
-  Streaming s3 file reads : to avoid loading file in memory.
-  Invalid json objects ( not having hotelreviewId) are skipped and are not sent to rabbit mq.
- In case no files are present in the s3 path , api will respond with message :"No files found in s3 path"
- Each file will be processed only if it has failed previosuly or it has not been processed before.
- Incase of  failure from Broker / Exchange in publishing msg , the sender will be retying using exponential backoff strategy at max 5 times. 
-  The consumer will send ACK once the record is successfully saved in DB - post this the broker will remove the record .
-  We are using Manual ACK approach from consumer , to confirm that the record is processed, in case of consumer is down , the message will be
redelivered from broker , and ensures atleast once semantic.
- We are only allowing 5 records to be consumed by a consumer at a time , post the ack from consumer , further records will be sent by AMQP broker.
- The enties defined are Hotel , Reviews, Providers, Reviewers,ProviderSummary
- There is 1:many relation defined between Hotel to Reviews, Provider to Reviews, Reviewer to Reviews, Hotel. 
- Hotel , Provider have a composite relation  to ProviderSummary
- All db operations run in a single transaction (per record) , this gaurantees atomocity .
- In order for data integrity and consistency, model relation : have defined foreign keys




##
## Flow:
- The api checks if files are present in the s3 bucket path , and we launch n  go routines where n is  minm(no of cpus, no of files in s3)
- The api will return HTTP 201 in case all validations passed and processing is
- Each go routine reads and parses the records in .jl file , and does validation of json object 
- If the records is Valid it is published to "reviews"  AMQP exchange. The exchange declared is of direct type.
- In case of errors in publishing to AMQP exchange , there is exponential backoff strategy implemented.
- 


- Consumer side:
   In case DB operation fails (db timeout/ db crash)- the mesage will be retried from broker as we are using Nack(false, true)
   In case of bad message , we are informing broker to discard the message = Nack(false, false)
   In case of succesful  db insert , we are sending ACK to broker , the broker will discard the msg ,and will not resend the message

   

![Sender side flow](docs/sender.png)


- DB Schema with relatiionships 
![ER Diagram](docs/er-diagram.png "Title")
_________________________
Assumptions:

1. The input files in s3 are of unique name , as we will not process data from file for 15 days with the same file name.This is done in order to  handle idempotency at file level.


____________________
 ## Instructions to run the ingestion flow::

 ```bash
curl --location --request POST 'localhost:8080/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"admin",
    "password":"3#%sdf"
}'
```

Sample output:

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzU2NjU2NjgyLCJpYXQiOjE3NTY2NTMwODJ9.ifE80l3K_4bRC9S1Dn0Cq2Iy4O7W6cx1m10dELKDT2w"}


Trigger Review Ingestion :

```bash
curl --location --request POST 'http://localhost:8080/v1/reviews/injest' \
--header 'Authorization: Bearer {token}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "PathPrefix":"reviews"
}'
``` 

Sample output in case processing is started:

{
    "success": true
}




API Parameter:
PathPrefix - path within Bucket where 1 or  multiple *.jl files exists.

| Validations                                   | API Behaviour                         |
|-----------------------------------------------|---------------------------------------|
| Invalid JSON, `PathPrefix` is missing in payload | HTTP **400 Bad Request** |
| Passed `PathPrefix` is wrong / S3 connection issue | HTTP **500 Internal Server Error**      |
| No files present in the s3 path |HTTP **400 Bad Request**     |


## CI Setup:
- DockerFile is defined in  deploy/Dockerfile  (working) 
- .github/workflows/docker-image.yml :: Github job that that runs on main branch, it build the container image and pushes the image  to github container registry(GHCR)

## Further Steps : 
- For deployment in aws kubernetes , i have defined a sample deployment file reviewService/deploy/k8s/reviewservice.yaml that picks image from GHCR.
- For local we are using .env , but for EKS deployment we can connect to aws secret manger at boottime , and add these in our enviroment.


