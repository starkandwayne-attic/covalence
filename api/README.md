#Covalence API

The Covalence API can be used to create a new connection in the Covalence system and to view existing
connections in the Covalence system. 

The API can be accessed via JSON over HTTP.

This document specifies the HTTP parameters and JSON formatting necessary for each action that the API supports. 

The fundamental resource created and managed by the Covalence API is the `connection` object. Each `connection` object is made up of a 
`source` and a `destination`. The `source` represents the `connection`'s side that is held by the job that reports the connection to the API. 
The `destination` is the other end of that connection. 

## Connection Object

## Adding New Connection

####POST => /connections

Body's JSON Format:
```
[	
		{
		        "protocol": "tcp",
			"source":{

				"ip": "10.80.76.9",
				"port": "5893",
				"deployment": "asv-sb-cf-postgres",
				"job": "secret-agent",
				"index": 1,
				"user": vcap,
				"group": vcap,
				"pid": 205,
				"process_name": "ruby",
				"age": 2345

		    },
		    "destination":{

				"ip":"10.80.79.34",
				"port":"6009"

		    }

     },

    {

                        "source":{

                                "ip": "10.80.76.9",
                                "port": "5893",
                                "deployment": "asv-sb-cf-postgres",
                                "job": "secret-agent",
                                "index": 1,
                                "user": vcap,
                                "group": vcap,
                                "pid": 205,
                                "process_name": "ruby",
                                "age": 2345

                    },
                    "destination":{

                                "ip":"10.80.79.34",
                                "port":"6009"

                    }

                },		},
		{

			"source": {

				...

			},
			"destination": {

				...	

			}
		
		},...

]
```

Response JSON Format:
```
{
    
    "success":"True",
    "error":"Blank if no error, contains error message to present to use otherwise."
    "connection_uuids":["9a304734-9692-4b42-a39b-889bc748ba9c", ... ]
    
}
```


## Retrieiving Existing Connections

####GET => /connections

No GET URL params. 

Response JSON Format:
(See POST => /connections section for keys available)
```
[
	{	

		"connection_uuid":"b21e83b5-2a43-4a1b-b592-79f5ae957cd2",
		"connection": {
		
				"source": {
     
				...			
	
				},
				"destination": {    

				...

				}			
		}

	},...
]
```