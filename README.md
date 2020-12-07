# Mimik

Simple application to simulate being a microservice in a mesh. 

## Usage

Any Mimik instance needs to have the following configuration:

### Environment Variables

The following environment variables are needed to create a Mimik instance:

| Variable | Description |
| - | - |
| MIMIK_SERVICE_NAME | The instance name |
| MIMIK_SERVICE_PORT | The instance port |
| MIMIK_ENDPOINTS_FILE | The file with the endpoints configuration and they connections to upstream services |

### Endpoints

The following file describes the endpoints that a Mimik instances listens to and the connections it has to other upstream services:

```json
[
    {
        "name": "Get songs",
        "path": "/",
        "method": "GET",
        "connections": [
            {
                "name": "songs-service",
                "port": "8080",
                "path": "songs",
                "method": "GET"
            }
        ]
    },
    {
        "name": "Get song with id 1",
        "path": "/songs/1",
        "method": "GET",
        "connections": [
            {
                "name": "songs-service",
                "port": "8080",
                "path": "songs/1",
                "method": "GET"
            },
            {
                "name": "hits-service",
                "port": "8080",
                "path": "hits/1",
                "method": "POST"
            }
        ]
    },
    {
        "name": "Health",
        "path": "/health",
        "method": "GET",
        "connections": []
    }
]
```

## Example

The health request from above, doesn't have any connections and its response looks like:

```bash
curl http://localhost:8080/health
{"name":"lyrics-page","version":"v1","path":"/health","statusCode":200,"upstreamResponse":[]}
```