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


### Details
- This concept was presented at [DevopsDays Austin](http://devopsdaysaustin.com). I will update this with a link to the video once it's uploaded.
- You can refer to the slides [here](http://adityarelangi.com/talks/2017/may/)
