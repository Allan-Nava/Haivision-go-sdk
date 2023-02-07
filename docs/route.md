---
layout: default
title: Managing and Configuring Routes
nav_order: 3
description: "Haivision Go SDK Official Documentation"
permalink: /route
last_modified_date: 2023-02-01T10:00:00+0000
---

# Managing and Configuring Routes

- Get List of Routes
- Get Route Configuration
- Create a Route
- Update a Route
- Delete a Route
- Start or Stop Routes
- Start or Stop a Route's Destination
- Export System Preset
- Import System Preset
- Example: Stopping an Individual Destination of a Route

### Get List of Routes

Use this command to get details of all routes. Retrieve the Device ID via the Get Device Info command.

```
Requests
GET /api/gateway/[Device ID]/routes
cookie: sessionID: [Session ID]
Parameters
N/A

Response
{
  "data": { [
    <Array of Route Objects>
  ]},
  "numPages": 1,
  "numResults": 3,
  "numActiveOutputConnections": 1,
  "pendingRouteCreates": 1
}
```

### Get Route Configuration

Use this command to get details of an individual route. Retrieve the Route ID via the Get Device Configuration command.

```
Requests
GET /api/gateway/[Device ID]/routes/[Route ID]
cookie: sessionID: [Session ID]
Parameters
N/A

Response
Note

See Route Object Model for definition of the response.

{
  <Route Object>
}
```

### Create a Route

Use this command to create an individual route.

```
Requests
POST /api/devices/[Device ID]/updates
cookie: sessionID: [Session ID]

{ 
 "action": "create",
 "deviceID": "[Device ID]", 
 "elementType": "route",
 "fields":
   {
    "name": "[Route name]",
    "startRoute": [true,false],
    "source":
      {
       <Source object>
      },
      "destinations": [
       <Destination object list>
      ]
   }
}

Response
{
  "status": "[success message]"
}
```

### Update a Route

Use this command to update an individual route.

```
Requests
POST /api/devices/[Device ID]/updates
cookie: sessionID: [Session ID]

{ 
 "action": "update",
 "deviceID": "[Device ID]", 
 "elementType": "route",
 "elementID": "[Route ID]",
 "fields":
   {
    "name": "[Route name]",
    "source":
      {
        <Source object>
      },
      "destinations": [
        <Destination object list>
      ]
   }
}

Parameters
Name	Type	Description
action	update  Update the specified element type.
deviceID	string	Device ID retrieved via the Get Device Info command. 
elementType	route Update the route element.
elementID	string	Route ID retrieved via the Get Device Configuration command
name	string	Name of route.
source	object	Source object model. See Source Object Model for definition.
destinations	string	Optional, if no destinatsions are desired in the route. Destinations object model. See Destinations Object Model for definition.
Response
{
  "status": "[success message]"
}
```