# GRANDA

- (n) granddaddy informally

Granda is the grand daddy of lambda, an event based serverless compute engine that lets you run your applications without provisioning a server. 

For a better explanation, try [AWS Lambda](https://aws.amazon.com/lambda/)


### How to run

- You need to run the go server for the backend to work. In order to do that,
```go
cd src/granda
go build
./granda
```
- The front end is not powered by golang. Place the code on your apache server path and visit
[http://localhost/granda/](http://localhost/granda/).
- When you click the call to action button, you will be taken to the app page, where you need 
to enter the func name and the image name. The other params are optional. Click on the submit button
and you will see some response


### Things to do

What we know so far
1. We know how to pass in the values into the containers.
2. The conainers need to built, run and removed for each run. 

Here's what we want to do

1. Start a server
2. On specific requests, specific containers are called
3. Track the time spent when the container is running
4. Maybe, use postman to make these requests
5. Work on the front-end to select the image and generated the URL
6. These new URLs when invoked should run fine

This should be good for now. 04/13/2017
