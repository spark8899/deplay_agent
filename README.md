# deploy-agent
Cooperate with jenkin to update software and execute restart procedure

# Test example
command api
```
# command not allow
curl -XPOST http://localhost:8000/command -d 'command=../tech/deploy-agent/start.sh sss'
# command path error
curl -XPOST http://localhost:8000/command -d 'command=../tech/deploy-agent/startbb.sh'
# command exec error
curl -XPOST http://localhost:8000/command -d 'command=../tech/deploy-agent/start.sh error'
# command exec success
curl -XPOST http://localhost:8000/command -d 'command=../tech/deploy-agent/start.sh'
curl -XPOST http://localhost:8000/command -d 'command=../tech/deploy-agent/start.sh aaa'
```

upload api
```
# upload filename  not allow
curl http://localhost:8000/upload/file -F "file=@../tech/ccc.txt" -F type=1 -v

# upload file suffix not supported
curl http://localhost:8000/upload/file -F "file=@../tech/aaa.txt" -F type=2 -v

# upload success
curl http://localhost:8000/upload/file -F "file=@../tech/aaa.txt" -F type=1 -v
curl http://localhost:8000/upload/file -F "file=@../tech/bbb.xml" -F type=2 -v
```

# Build app
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
```

# Deploy app
create app directory. 
```
mkdir -p /opt/app/deploy-agent/{config,logs}
```

create configuration file
```
vi /opt/app/deploy-agent/config/config.yaml
Server:
  RunMode: release
  HttpPort: 8881
  ReadTimeout: 60
  WriteTimeout: 60
App:
  DefaultPageSize: 10
  MaxPageSize: 100
  DefaultContextTimeout: 60
  LogSavePath: logs
  LogFileName: deploy-agent
  LogFileExt: .log
  DeployPath: /opt/app/xxx
  DeployFiles:
    - xxx.jar
    - logback-spring.xml
  ExecScripts:
    - /opt/app/xxx/startjar.sh
    - /opt/app/xxx/startjar.sh start
    - /opt/app/xxx/startjar.sh stop
    - /opt/app/xxx/startjar.sh restart
    - /opt/app/xxx/startjar.sh status
    - /opt/app/xxx/startjar.sh update
  UploadMaxSize: 500  # MB
  UploadAllowExts:
    - .jar
    - .xml
    - .rc
```

vi /etc/systemd/system/deploy-agent.service
```
[Unit]
Description=deploy-agent
After=network.target

[Service]
ExecStart=/opt/app/deploy-agent/deploy-agent
WorkingDirectory=/opt/app/deploy-agent
StandardOutput=inherit
StandardError=inherit
Restart=always
User=apprun

[Install]
WantedBy=multi-user.target

systemctl daemon-reload
systemctl start deploy-agent
systemctl enable deploy-agent
```
