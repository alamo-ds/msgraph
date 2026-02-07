# msgraph

**(WIP)**

Go library for interacting with the Microsoft Graph API.

## Overview

Fetch resources through the Microsoft Graph API with the client:

```go
import (
    "os"
    "github.com/ads/msgraph/auth"
    "github.com/ads/msgraph/graph"
)

aadConfig := auth.AzureADConfig{
    TenantID: os.Getenv("TENANT_ID"),
    ClientID: os.Getenv("CLIENT_ID"),
    ClientSecret: os.Getenv("CLIENT_SECRET"),
}

client := graph.NewClient(aadConfig)
```

_Coming soon: set values for your tenant ID and client ID/secret with the command-line utility_

```bash
msgraph set \
    --tenant-id <TENANT_ID> \
    --client-id <CLIENT_ID> \
    --secret <CLIENT_SECRET> \
```

This will create a file at `$HOME/.msgraph/config.json`, simplifying the client creation:

```go
client := graph.NewClient() // no config object necessary!
```

Refer to resoruces for the full list of supported endpoints.

## Endpoints

The base URL is `https://graph.microsoft.com/v1.0`

**GET `/groups`**

```go
groups, err := client.Groups().Get(ctx)
if err != nil {
    ...
}
```

**GET `/groups/{group-id}`**

```go
group, err := client.Groups().ById(groupId).Get(ctx)
if err != nil {
    ...
}
```

**GET `/groups/{group-id}/planner/plans`**

```go
group, err := client.Groups().ById(groupId).Plans().Get(ctx)
if err != nil {
    ...
}
```

**GET `/planner/plans/{plan-id}`**

```go
group, err := client.Planner().ById(plannerId).Get(ctx)
if err != nil {
    ...
}
```

**PATCH `/planner/plans/{plan-id}`**

```go
params := graph.PatchPlanParams{
    Title: "Updated Title",
}

plan, err := client.Planner().ById(plannerId).Patch(ctx, params)
if err != nil {
    ...
}
```

**GET `/planner/plans/{plan-id}/buckets`**

```go
buckets, err := client.Planner().ById(plannerId).Buckets().Get(ctx)
if err != nil {
    ...
}
```

**POST `/planner/buckets`**

```go
params := graph.PostBucketParams{
    Name: "New Bucket Name",
    PlanID: plannerId,
    OrderHint: "", // optional
}

bucket, err := client.Planner().Buckets().Post(ctx, params)
if err != nil {
    ...
}
```

**PATCH `/planner/buckets/{plan-id}`**

```go
params := graph.PatchBucketParams{
    Name: "Updated Title",
    OrderHint: "", // optional
}

bucket, err := client.Planner().Buckets().ById(bucketId).Patch(ctx, params)
if err != nil {
    ...
}
```

**GET `/planner/plans/{plan-id}/tasks`**

```go
tasks, err := client.Planner().ById(plannerId).Tasks().Get(ctx)
if err != nil {
    ...
}
```

**GET `/planner/plans/{plan-id}/tasks`**

```go
tasks, err := client.Planner().ById(plannerId).Tasks().Get(ctx)
if err != nil {
    ...
}
```

**GET `/me/planner/tasks`**

_Coming soon_

**GET `/users/{coworker-mail}/planner/tasks`**

_Coming soon_

**GET `/planner/tasks/{task-id}`**

```go
tasks, err := client.Planner().Tasks().ById(taskId).Get(ctx)
if err != nil {
    ...
}
```

**POST `/planner/tasks`**

```go
params := graph.PostTaskParams{
    PlanID: planId,
	BucketID: bucketId, // optional
	Title: "My New Task",
	OrderHint: "", // optional
    ...
}

task, err := client.Planner().Tasks().Post(ctx, params)
if err != nil {
    ...
}
```

**PATCH `/planner/tasks/{task-id}`**

```go
params := graph.PatchTaskParams{
    Title: "Updated Task Title",
    BucketID: "", // optional
    Priority: 0, // optional
    ...
}

task, err := client.Planner().Tasks().ById(taskId).Patch(ctx, params)
if err != nil {
    ...
}
```

**GET `/planner/tasks/{task-id}/details`**

```go
taskDetails, err := client.Planner().Tasks().ById(taskId).Details().Get(ctx)
if err != nil {
    ...
}
```
