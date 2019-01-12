package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/joho/godotenv/autoload"

	"github.com/shawncatz/automagical/ec2"
)

func main() {
	lambda.Start(ec2.Handle)
}
