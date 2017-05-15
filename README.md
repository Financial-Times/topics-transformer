*DECOMISSIONED*
See [Basic TME Transformer](https://github.com/Financial-Times/basic-tme-transformer) instead

# topics-transformer

[![Circle CI](https://circleci.com/gh/Financial-Times/topics-transformer/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/topics-transformer/tree/master)

Retrieves Topics taxonomy from TME and transforms the topics to the internal UP json model.
The service exposes endpoints for getting all the topics and for getting topic by uuid.

# How to run
`go get github.com/Financial-Times/topics-transformer`

`$GOPATH/bin/ ./topics-transformer.exe  --base-url=http://localhost:8080/transformers/topics/ --tme-base-url=<TME URL> --tme-username=<USER> --tme-password=<PWD> --token=<TOKEN> --port=8080 --maxRecords=1000 --slices=10 `                

```
export|set PORT=8080
export|set BASE_URL="http://localhost:8080/transformers/topics/"
export|set TME_BASE_URL="http://tme.ft.com"
export|set TME_USERNAME="user"
export|set TME_PASSWORD="pass"
export|set TOKEN="token"
export|set MAX_RECORDS="10"
$GOPATH/bin/topics-transformer
```

With Docker:

`docker build -t coco/topics-transformer .`

`docker run -ti --env BASE_URL=<base url> --env TME_BASE_URL=<TME URL> --env TME_USERNAME=<user> --env TME_PASSWORD=<pass> --env TOKEN=<token> -env MAX_RECORDS=<recors> coco/topics-transformer`

#Usage

Get all topics:
`curl -X GET -H "Cache-Control: no-cache" -H "Postman-Token: 4f5f5bec-91ae-d714-3fb5-49b6e09a5a1b" "https://semantic-up.ft.com/__topics-transformer/transformers/topics"`

Get by topic:
`curl -X GET -H "Cache-Control: no-cache" -H "Postman-Token: 4f5f5bec-91ae-d714-3fb5-49b6e09a5a1b" "https://semantic-up.ft.com/__topics-transformer/transformers/topics/0205c4dd-5430-33ac-bb1f-fbfc347b1475"`
