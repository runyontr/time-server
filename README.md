# Time Server

Serves up current time on host via json message.  Endpoint hosted at `/v1/time`




## Run Copy
```bash
docker run -it -p 8080:8080 runyonsolutions/time-server
```
 
## Query against Server
```bash
$ curl localhost:8080/v1/time
{"CurrentTimeMillis":1516059077169}
```