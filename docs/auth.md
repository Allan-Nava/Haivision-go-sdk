---
layout: default
title: Login and Device Info
nav_order: 2
description: "Haivision Go SDK Official Documentation"
permalink: /auth
last_modified_date: 2023-02-01T10:00:00+0000
---

# Login and Device Info

- Initiate Session
- Get Session Information
- Delete Existing Session
- Get Device Info
- Get Device Configuration


### Initiate Session

Initiating a session requires a user name and a password. The configured users are the same as configured for the web application. When you issue the Initiate Session request, a sessionID cookie is returned in the response.
```
Requests
POST /api/session

{
 "username" :  "[user name]",
 "password" :  "[password]"
}
```

```
Parameters
Name	Type	Description
username	string	Valid user.
password	string	Valid password.

Response
{
    "response": {
        "type": "Session",
        "message": "Session successfully started for haiadmin",
        "sessionID": "[Session ID]",
        "lastLoginDate": 1536777877871,
        "numLoginFailures": 0
    }
}
```

### Get Session Information

Get existing session information.

```
Response
{
    "sessionID": "[Session ID]",
    "displayName": "Administrator",
    "email": "haiadmin@localhost",
    "roles": [
        "Administrator"
    ],
    "startAt": 1536937838919,
    "expireAt": 1536938857529,
    "lastLoginDate": null,
    "numLoginFailures": null,
    "isLicensed": true
}
```

### Delete Existing Session

Delete existing session.

```
Requests
DELETE /api/session
cookie: sessionID: [Session ID]
Parameters
N/A

Response
{
 "type" : "Session",
 "message" :  "[success message]"
}
```

### Get Device Info

Retrieves information about the appliance. Use the Device ID value from the _id property to reference the device in subsequent API calls.

```
Requests
GET /api/devices
cookie: sessionID: [Session ID]
Parameters
N/A

Response
[
 { 
  "_id": "[Device ID]", 
  "type": "gateway", 
  "ip": "127.0.0.1", 
  "name": "Haivision Gateway", 
  "lastConnectedAt": [Date/time shown in Unix time], 
  "statusCode": "ok", 
  "status": "Online", 
  "statusDetails": "Connection has been established in the last 1 minutes.", 
  "serialNumber": null, 
  "firmware": "5.0.180611.1530", 
  "hasAdminError": false, 
  "pendingSync": false, 
  "lastConnection": "<1m" 
 } 
]
```

### Get Device Configuration

Retrieves the configuration of the device, including the list of routes. The Device ID used in the command is retrieved via the Get Device Info command.

```
Requests
GET /api/devices/[Device ID]
cookie: sessionID: [Session ID]
Parameters
N/A

Response
Note

Use the config > routes > id properties on subsequent API calls to start/stop the corresponding routes, as shown in Start or Stop Routes.

{
  "_id": "[Device ID]",
  "type": "gateway",
  "ip": "127.0.0.1",
  "name": "Haivision Media Gateway",
  "lastConnectedAt": [Date/time in Unix time],
  "lastConfigReadAt": [Date/time in Unix time],
  "statusCode": "ok",
  "status": "Online",
  "statusDetails": "Connection has been established in the last 1 minutes.",
  "serialNumber": null,
  "firmware": "5.0.180611.1530",
  "hasAdminError": false,
  "config": {
    "routes": [ { <Route Object> } ],
    "settings": {
      "networkAdapters": [
        {
          "name": "eth0",
          "address": "10.67.12.128"
        },
        {
          "name": "idrac",
          "address": "169.254.0.2"
        }
      ]
    },
    "calypsos": [
      {
        "id": "Lh5ZNDm7LspsseJ3pg82jw",
        "gatewayName": "MG 141",
        "address": "dev.haivision.com",
        "httpPort": 80,
        "httpsPort": 443,
        "passcode": "[Pairing passcode]",
        "lastConnectedAt": 1536759437398,
        "status": "accepted",
        "statusDetails": "pairing accepted",
        "proxyAddress": "10.69.12.141",
        "proxyHttpPort": 80,
        "proxyHttpsPort": 443,
        "lastConnection": "<1m"
      }
    ]
  },
  "pendingSync": false,
  "lastAcceptedUpdate": 1530219677119,
  "lastConnection": "<1m",
  "gateway": {
    "port": 1080
  },
  "suggestedPollingInterval": {
    "ms": 10000
  }
}
```