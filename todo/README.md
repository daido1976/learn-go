```sh
$ curl -X POST http://localhost:8080/todos/ | jq .
$ curl -X POST http://localhost:8080/todos/ -d '{ "title": "Test", "body": "bodyだよ" }' | jq .
```
