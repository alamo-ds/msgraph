# msgraph

[![license](https://img.shields.io/github/license/alamo-ds/msgraph)](https://github.com/alamo-ds/msgraph/blob/master/LICENSE)
[![CI](https://github.com/alamo-ds/msgraph/actions/workflows/ci.yaml/badge.svg)](https://github.com/alamo-ds/msgraph/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alamo-ds/msgraph)](https://goreportcard.com/report/github.com/alamo-ds/msgraph)

Go library for interacting with the Microsoft Graph API. This project is a work in progress; support for more endpoints will be added over time.

## Overview

Fetch resources through the Microsoft Graph API with the client:

```go
import (
    "os"
    "github.com/ads/msgraph/graph"
)

aadConfig := graph.AzureADConfig{
    TenantID: os.Getenv("TENANT_ID"),
    ClientID: os.Getenv("CLIENT_ID"),
    ClientSecret: os.Getenv("CLIENT_SECRET"),
}

client := graph.NewClient(ctx, aadConfig)
```

Alternatively, set values for your tenant ID and client ID/secret with the command-line utility

```bash
msgraph set \
    --tenant-id <TENANT_ID> \
    --client-id <CLIENT_ID> \
    --secret <CLIENT_SECRET> \
```

With `TENANT_ID`, `CLIENT_ID`, and `CLIENT_SECRET` environment variables set, you can simply run `msgraph set`.

This will create a file at `$HOME/.msgraph/config.json`, simplifying the client creation:

```go
client := graph.NewClient(ctx) // no config object necessary!
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
