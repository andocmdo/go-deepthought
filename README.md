# gostockd
RESTful JSON (and other formats) API Backend Job Server for gostock simulator


# Example curl commands:
```

curl -H "Content-Type: application/json" -d '{"args":{"command":"gastockagent","endDate":"2017-12-27","maxGen":"10df","mutRate":"0.04","popSize":"235","startDate":"2016-12-26","symbol":"jsasdfon"}}' http://localhost:8080/api/v1/jobs
curl -H "Content-Type: application/json" -d '{"symbol":"192.168.1.45","port":"12365"}' http://localhost:8080/api/v1/workers

```
