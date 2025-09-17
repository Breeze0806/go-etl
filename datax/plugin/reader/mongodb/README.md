# MongodbReader Plugin Documentation

## Quick Introduction

MongodbReader plugin enables reading data from MongoDB collections. Internally, MongodbReader connects to a remote MongoDB database using the official MongoDB Go driver `go.mongodb.org/mongo-driver` and executes corresponding queries to retrieve data from the MongoDB server.

## Implementation Principles

MongodbReader connects to a remote MongoDB database using the official MongoDB Go driver. Based on user-provided configuration information, it generates queries and sends them to the remote MongoDB server. The returned results from these queries are assembled into an abstract dataset using go-etl's custom data types and passed to downstream Writers for processing.

The plugin splits data reading tasks based on the ObjectId range. It first determines the minimum and maximum ObjectIds in the collection, then divides the ObjectId range into multiple segments based on the number of channels configured, and assigns each segment to a separate task for parallel processing.

## Functionality Description

### Configuration Example

Configuring a job to synchronize data from a MongoDB collection to another MongoDB collection:

```json
{
    "job": {
        "content": [
            {
                "reader": {
                    "name": "mongodbreader",
                    "parameter": {
                        "connection": {
                            "address": "localhost:27017",
                            "username": "root",
                            "password": "123456",
                            "table": {
                                "db": "test_database",
                                "collection": "users"
                            }
                        },
                        "column": [
                            {
                                "name": "_id",
                                "type": "string"
                            },
                            {
                                "name": "name",
                                "type": "string"
                            },
                            {
                                "name": "email",
                                "type": "string"
                            },
                            {
                                "name": "age",
                                "type": "int"
                            },
                            {
                                "name": "created_at",
                                "type": "date"
                            }
                        ]
                    }
                }
            }
        ]
    }
}
```

### Parameter Explanation

#### connection

##### address

- Description: Used to configure the MongoDB server address, in the format `host:port`.
- Required: Yes
- Default: None

##### username

- Description: Used to configure the MongoDB username for authentication.
- Required: No
- Default: None

##### password

- Description: Used to configure the MongoDB password for authentication.
- Required: No
- Default: None

##### table

###### db

- Description: Used to configure the name of the database containing the collection to read from.
- Required: Yes
- Default: None

###### collection

- Description: Used to configure the name of the collection to read from.
- Required: Yes
- Default: None

#### column

- Description: An array describing the fields that need to be synchronized from the configured collection. Users can use the JSON array format to describe field information.

  Each field configuration includes:
  - `name`: The field name in the MongoDB document
  - `type`: The target data type (string, int, date, etc.)
  - `spliter`: For array fields, the separator used when converting to string

  Example:
  ```json
  {"column": [
      {
          "name": "_id",
          "type": "string"
      },
      {
          "name": "tags",
          "type": "string",
          "spliter": ","
      }
  ]}
  ```

  Supported data types:
  - `string`: String type
  - `int`: Integer type
  - `date`: Date/time type
  - `Array`: Array type (converted to string with separator)

- Required: Yes
- Default: None

#### split_key

- Description: Used to configure the field for data splitting. Currently, only "_id" is supported as the split key.
- Required: No
- Default: "_id"

### Type Conversion

Currently, MongodbReader supports most MongoDB types with appropriate conversions to go-etl internal types:

| go-etl Type | MongoDB Data Type           | Notes                                       |
| ----------- | --------------------------- | ------------------------------------------- |
| bigInt      | int, int32, int64           |                                             |
| decimal     | float32, float64            |                                             |
| string      | string                      | Also used for ObjectId and other types      |
| time        | primitive.DateTime          |                                             |
| bool        | bool                        |                                             |
| bytes       | []byte                      |                                             |

Special handling:
1. `primitive.ObjectID`: By default converted to hexadecimal string
2. `[]interface{}` (arrays): Can be converted to strings with custom separators
3. `bson.M` (embedded documents): Serialized as JSON strings
4. Other types: Converted to string representation

## Performance Report

To be tested.

## Constraints and Limitations

### Database Encoding Issues

Currently, only the utf8 character set is supported.

### Split Key Limitation

Currently, only "_id" field is supported as the split key for parallel processing.