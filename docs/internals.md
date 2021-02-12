# Internals

Any Mimik instance needs to have the following configuration:

#### Environment Variables

The following environment variables are needed to create a Mimik instance manually (without the Helm chart):

| Variable | Description |
| - | - |
| MIMIK_SERVICE_NAME | The instance name |
| MIMIK_SERVICE_PORT | The instance port |
| MIMIK_ENDPOINTS_FILE | A file containing the endpoints configuration and the connections to upstream services |
| MIMIK_LABELS_FILE | A file containing labels, Mimik looks for the version label specifically in the file, if the file does not exists or does not have the version label it defaults to v1 |

#### Endpoints

The following file describes the endpoints that a Mimik instance listens for and the connections it has to other upstream services:

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