# topics-transformer

[![Circle CI](https://circleci.com/gh/Financial-Times/topics-transformer/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/topics-transformer/tree/master)

Retrieves Topics taxonomy from TME vie the structure service and transforms the topics to the internal UP json model.
The service exposes endpoints for getting all the topics and for getting topic by uuid.

# Usage
`go get github.com/Financial-Times/topics-transformer`

`$GOPATH/bin/topics-transformer --port=8080 -base-url="http://localhost:8080/transformers/topics/" -structure-service-base-url="http://metadata.internal.ft.com:83" -structure-service-username="user" -structure-service-password="pass" -structure-service-principal-header="app-preditor"`
```
export|set PORT=8080
export|set BASE_URL="http://localhost:8080/transformers/topics/"
export|set STRUCTURE_SERVICE_BASE_URL="http://metadata.internal.ft.com:83"
export|set STRUCTURE_SERVICE_USERNAME="user"
export|set STRUCTURE_SERVICE_PASSWORD="pass"
export|set PRINCIPAL_HEADER="app-preditor"
$GOPATH/bin/topics-transformer
```

With Docker:

`docker build -t coco/topics-transformer .`

`docker run -ti --env BASE_URL=<base url> --env STRUCTURE_SERVICE_BASE_URL=<structure service url> --env STRUCTURE_SERVICE_USERNAME=<user> --env STRUCTURE_SERVICE_PASSWORD=<pass> --env PRINCIPAL_HEADER=<header> coco/topics-transformer`
