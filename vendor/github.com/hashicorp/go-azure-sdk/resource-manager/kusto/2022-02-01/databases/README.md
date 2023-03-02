
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-02-01/databases` Documentation

The `databases` SDK allows for interaction with the Azure Resource Manager Service `kusto` (API Version `2022-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-02-01/databases"
```


### Client Initialization

```go
client := databases.NewDatabasesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatabasesClient.AddPrincipals`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := databases.DatabasePrincipalListRequest{
	// ...
}


read, err := client.AddPrincipals(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := databases.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := databases.CheckNameRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := databases.Database{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Delete`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Get`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.ListByCluster`

```go
ctx := context.TODO()
id := databases.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

read, err := client.ListByCluster(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.ListPrincipals`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

read, err := client.ListPrincipals(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.RemovePrincipals`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := databases.DatabasePrincipalListRequest{
	// ...
}


read, err := client.RemovePrincipals(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.Update`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := databases.Database{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
