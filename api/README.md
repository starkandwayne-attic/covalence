#Covalence API

The Covalence API can be used to create a new connection in the Covalence system and to view existing
connections in the Covalence system. 

The API can be accessed via JSON over HTTP.

This document specifies the HTTP parameters and JSON formatting necessary for each action that the API supports. 

## Adding New Connection

####POST => /connections

Body's JSON Format:
```
{
    "source_ip":"10.80.76.9",
    "source_port":"5893",
    "desination_ip":"10.80.79.34",
    "desination_port":"6009",
    
    "destination_deployment_name":"asv-sb-cf-postgres",
    "source_deployment_name":"asv-cloud-foundry",
    
}
```

A list of the above JSON objects may be sent in order to create connections in batches.

Response JSON Format:
```
{
    
    "success":"True",
    "error":"Blank if no error, contains error message to present to use otherwise."
    "connection_uuid":"9a304734-9692-4b42-a39b-889bc748ba9c"
    
}
```

If a batch of connection objects are sent to this endpoint, the `connection_uuid` key in the response will contain a list of the UUIDs of the connections created.


## Retrieiving Existing Connections

####GET => /connections

No GET URL params. 

Response JSON Format:
```
[
    {
     
        "source_ip":"10.80.76.9",
        "source_port":"5893",
        "desination_ip":"10.80.79.34",
        "desination_port":"6009",
        "destination_deployment_name":"asv-sb-cf-postgres",
        "source_deployment_name":"asv-cloud-foundry",
        "connection_uuid":"b21e83b5-2a43-4a1b-b592-79f5ae957cd2"
        
    },
    {
        "source_ip":"10.80.35.2",
        "source_port":"9282",
        "desination_ip":"10.80.67.73",
        "desination_port":"4563",
        "destination_deployment_name":"asv-pr-ha-rabbitmq-36",
        "source_deployment_name":"asv-pr-cloud-foundry",
        "connection_uuid":"026459d8-b016-42cf-999c-f48a1cc0c841"
    }
        
    }, ...
]
```