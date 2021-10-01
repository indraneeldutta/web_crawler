# Web crawler

This project is a REST based application to get details of a webpage. It builds and runs on docker for ease of installation and execution. 

# Installation
1. Clone the main branch
2. Run `make run` inside the project dir
3. The project will automatically get deployed on docker and serve on port 9000.

PS: Make sure docker daemon is running before running `make ..` command else it will fail to deploy.

PS: if you do not wish to run it on docker follow the below steps
1. run `go mod vendor` -> this will download all the dependencies required for the project
2. run `go build && ./web_crawler`
3. This will run the application natively and serve on port 9000.

# Usage
1. This is a REST application. The below is the sample request for getting the results. The command can be imported in REST tool like postman or directly run on command line.
```bash
curl --location --request POST 'http://localhost:9000/v1/page/details' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url": "https://google.in"
}'
```
2. Change the url to your desired one inside the body to get results of it.
3. Test cases are written for the apis package with below coverage

```
Running tool: /usr/local/bin/go test -timeout 30s -coverprofile=/var/folders/bd/_twy1rd90bjd1p_b9__yn5vm0000gq/T/vscode-gomWF1BF/go-code-cover github.com/web_crawler/apis

ok  	github.com/web_crawler/apis	3.170s	coverage: 95.9% of statements
```