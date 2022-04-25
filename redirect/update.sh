set GOOS=linux

go build -o redirect main.go
build-lambda-zip -o deployment.zip redirect

aws lambda update-function-code --function-name RedirectFunction --region ap-south-1 --zip-file fileb://./deployment.zip
