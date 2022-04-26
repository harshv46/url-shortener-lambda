GOOS=linux GOARCH=amd64 go build -o shorten main.go
zip deployment.zip shorten

aws lambda update-function-code --function-name ShortenFunction --region ap-south-1 --zip-file fileb://./deployment.zip
