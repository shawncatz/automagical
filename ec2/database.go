package ec2

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//go:generate counterfeiter . Database
type Database interface {
	Insert(*ec2.Instance) error
	Find(string) (*ec2.Instance, error)
	Remove(string) error
}

type InstanceDatabase struct {
	tableName string
	db        *dynamodb.DynamoDB
}

func NewDatabase(table string) *InstanceDatabase {
	ses := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(ses)

	return &InstanceDatabase{
		tableName: table,
		db:        db,
	}
}

func (d *InstanceDatabase) create() error {
	if d.tableExists() {
		return nil
	}

	createTableInput := dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("InstanceId"),
				AttributeType: aws.String("S"),
			},
		},
		TableName: aws.String(d.tableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("InstanceId"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	_, err := d.db.CreateTable(&createTableInput)
	if err != nil {
		return err
	}
	return d.db.WaitUntilTableExists(&dynamodb.DescribeTableInput{
		TableName: aws.String(d.tableName),
	})
}

func (d *InstanceDatabase) tableExists() bool {
	out, err := d.db.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return false
	}

	for _, t := range out.TableNames {
		if *t == d.tableName {
			return true
		}
	}

	return false
}

func (d *InstanceDatabase) Insert(i *ec2.Instance) error {
	encoder := dynamodbattribute.NewEncoder()
	updated := time.Now().Unix()

	dynamoInstance, err := encoder.Encode(i)
	if err != nil {
		return err
	}

	dynamoItem := map[string]*dynamodb.AttributeValue{
		"InstanceID": &dynamodb.AttributeValue{
			S: i.InstanceId,
		},
		"Instance": dynamoInstance,
		"Updated": &dynamodb.AttributeValue{
			N: aws.String(fmt.Sprintf("%d", updated)),
		},
	}

	_, err = d.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      dynamoItem,
	})

	return err
}

func (d *InstanceDatabase) Find(id string) (*ec2.Instance, error) {
	encoder := dynamodbattribute.NewEncoder()
	av, err := encoder.Encode(id)
	if err != nil {
		return nil, err
	}

	output, err := d.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key:       map[string]*dynamodb.AttributeValue{"InstanceID": av},
	})
	if err != nil {
		return nil, err
	}

	if len(output.Item) == 0 {
		return nil, errors.New("not found")
	}

	decoder := dynamodbattribute.NewDecoder()
	instance := &ec2.Instance{}
	err = decoder.Decode(output.Item["Instance"], instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (d *InstanceDatabase) Remove(id string) error {
	encoder := dynamodbattribute.NewEncoder()
	av, err := encoder.Encode(id)
	if err != nil {
		return err
	}

	_, err = d.db.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"InstanceID": av,
		},
		TableName: aws.String(d.tableName),
	})

	return err
}

func (d *InstanceDatabase) destroy() error {
	deleteTableInput := &dynamodb.DeleteTableInput{
		TableName: aws.String(d.tableName),
	}
	if _, err := d.db.DeleteTable(deleteTableInput); err != nil {
		return err
	}

	waitInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(d.tableName),
	}
	return d.db.WaitUntilTableNotExists(waitInput)
}
