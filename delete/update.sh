set GOOS=linux

go build -o delete main.go
build-lambda-zip -o deployment.zip delete

aws lambda update-function-code --function-name DeleteFunction --region ap-south-1 --zip-file fileb://./deployment.zip
