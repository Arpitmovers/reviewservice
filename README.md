# reviewService
microservice  





Setup Instructions:



_____________________
Design Decisions:
Tech Stack: Golang, RabbitMQ, Redis, MariaDB





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

curl --location --request POST 'http://localhost:8080/reviews/injest' \
--header 'Authorization: Bearer {token}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Bucket":"hotelservice",
    "PathPrefix":"reviews",
    "Force":false

}'



PathPrefix - path within Bucket where 1 or  multiple files exists.

______________________

How duplicate records are handled ?
We are doing dedupe on hotelRecordId at subscriber end.

