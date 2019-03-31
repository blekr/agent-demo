# Requirement
golang 1.8 is required to build

# Features
- <strong>Detached child process: </strong>Forked child process keep running despite the parent's running status, event though the parent process is killed abnormally
- <strong>Graceful shutdown: </strong>sending SIGTERM to process first and then wait for it to terminate. Only kill it when it doest not exit in a specified period of time.

# Usage

- download source:
```
go get github.com/blekr/agent-demo
```
- cd source directory:
```
cd $GOPATH/src/github.com/blekr/agent-demo
```
- build
```
go build -o /tmp/agent-demo
```
- run
```
/tmp/agent-demo
```
- test

Open a new terminal and use below curl commands to start, show status of and stop the process. Replace ```sleep.sh``` with your actual command and ```26085``` your actual pid.
```

curl http://localhost:8080/start -d '{"Path":"sleep.sh","Args":["abc","def"]}'
curl http://localhost:8080/show -d '{"Pid":26085}'
curl http://localhost:8080/stop -d '{"Pid":26085}'

```
# Issues
- Fail to download package ```https://golang.org/x/sys/unix?go-get=1```

If go get fails to download package ```https://golang.org/x/sys/unix?go-get=1``` because of the GFW, please manually download it and place it under correct directory.