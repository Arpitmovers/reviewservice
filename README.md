# reviewService
microservice  





## Setup Instructions

Install the following for running the setup in your local system:

- [Install Redis](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis-on-linux/) 
-  
- [AMQP in local (RabbitMQ) 4.1.3](https://www.rabbitmq.com/docs/install-debian)  
- 
- [MariaDB in local](https://mariadb.com/docs/server/server-installation/mariadb-package-repositories/)  
- 
-  Add virtual host to rabbitMQ: sudo rabbitmqctl add_vhost review_dev
-  Create RabbitMq user: sudo rabbitmqctl add_user reviewService 1asd21
-  Grant read and write permission to the user reviewService on the vhost review_dev  
  sudo rabbitmqctl set_permissions -p review_dev reviewService "" ".*" ".*"*"




-  reviewService/deploy/Dockerfile - file to build container of our service - working
- reviewService/.github/workflows/docker-image.yml  - pipeline to build and push a Docker image to GitHub Container Registry (GHCR). - working

_____________________
Design Decisions:
Tech Stack: Golang1.22, RabbitMQ 4.1.3, Redis, MariaDB 15.1


- Incase  , of  failure from Broker / Exchange in publishing msg , the sender will be retying using exponential backoff strategy at max 5 times.
- In case no files are present in the s3 path , api will respond with message :"No files found in s3 path"
- **

_________________

Flow:
- The api checks if files are present in the s3 bucket path , and we launch n  go routines where n is  minm(no of cpus, no of files in s3)
- Each go routine reads and parses the records in .jl file , and does validation of json object 
- If the records is Valid it is published to "reviews"  AMQP exchange. The exchange declared is of direct type.
- In case of errors in publishing to AMQP exchange , there is exponential backoff strategy implemented.
- Once the msg is 
- 



_________________________
Assumptions:







____________________
 Instructions to run the ingestion flow::
curl --location --request POST 'localhost:8080/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"admin",
    "password":"3#%sdf"
}'

sample o/p:

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzU2NjU2NjgyLCJpYXQiOjE3NTY2NTMwODJ9.ifE80l3K_4bRC9S1Dn0Cq2Iy4O7W6cx1m10dELKDT2w"}




Trigger Review Ingestion :

curl --location --request POST 'http://localhost:8080/v1/reviews/injest' \
--header 'Authorization: Bearer {token}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "PathPrefix":"reviews",
    "Force":false

}'

API Parameter:
PathPrefix - path within Bucket where 1 or  multiple *.jl files exists.





______________________

How duplicate records are handled ?
We are doing dedupe on hotelRecordId at subscriber end.

