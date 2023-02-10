---
layout: default
title: Statistics
nav_order: 4
description: "Haivision Go SDK Official Documentation"
permalink: /stats
last_modified_date: 2023-02-01T10:00:00+0000
---

# Statistics

- Get Route Statistics
- Get Source Statistics
- Get Destination Statistics
- Get SRT Client Statistics

### Get Route Statistics

Retrieves statistics about a specific route.
```
Requests
GET /api/gateway/[Device ID]/statistics?routeID=[Route ID]
cookie: sessionID: [Session ID]


Response
Response varies depending on which protocol is used.
{
  "collectedAt": [Date/time in Unix time],
  "route": {
    "name": "[Route Name]",
    "elapsedRunningTime": "00:00:14",
    "id": "[Route ID]",
    "state": "running",
    "source": {
      <Source Statistics object>
    },
    "destinations": [
      {
        <Destination Statistics object>
      }
    ]
  }
}
```

### Get Source Statistics
Retrieves statistics about a specific route's source.
```
Requests
GET /api/gateway/[Device ID]/statistics?routeID=[Route ID]&sourceID=
  [Source ID]
cookie: sessionID: [Session ID]

Response
Response varies depending on which protocol is used. See Source Statistics Object Model for parameter definitions of the source statistics object.

{
  "collectedAt": [Date/time in Unix time],
  "source": {
    <Source Statistics Object>
  }
}

```

### Get Destination Statistics

Retrieves statistics about a specific route's destination.

```
Requests
GET /api/gateway/[Device ID]/statistics?routeID=[Route ID]&destinationID=
  [Destination ID]
cookie: sessionID: [Session ID]

or

GET /api/gateway/[Device ID]/statistics?routeID=[Route ID]&destinationName=
  [Destination Name]
cookie: sessionID: [Session ID]

Response
Response varies depending on which protocol is used. See Destinations Statistics Object Model for parameter definitions of the destination statistics object.

{
  "collectedAt": [Date/time in Unix time],
  "destination": {
    <Destination Statistics Object>
  }
}

```

### Get SRT Client Statistics

Retrieves statistics for a specific client connected to an SRT listener destination.

```
Requests
GET /api/gateway/[Device ID]/statistics/client?routeID=[Route ID]&destinationID=
  [Destination ID]&clientAddress=[Client Address]&clientPort=[Client Port]
cookie: sessionID: [Session ID]

Response
See Client Statistics Model in Destinations Statistics Object Model for parameter definitions of the client statistics object.

{
  "collectedAt": [Date/time in Unix time],
  "clientStat": [
    <Client Statistics Object>
  ]
}
```