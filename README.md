# deploy-agent
Cooperate with jenkin to update software and execute restart procedure

# version plan
 * v1 Support specific file uploads and specific commands for a single project(Has been completed.)
 * v2 Support specific file uploads and specific commands for multiple projects

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
mkdir -p /opt/app/deploy-agent/{configs,logs}
```

create configuration file
```
vi /opt/app/deploy-agent/configs/config.yaml
Server:
  RunMode: release
  HttpPort: 58881
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
User=xxxx
Group=xxxx

[Install]
WantedBy=multi-user.target

systemctl daemon-reload
systemctl start deploy-agent
systemctl enable deploy-agent
```

# jenkins shell
```
set +x
ENV=`echo ${JOB_NAME} | awk -F'-' '{print $1}'`
PROJECT=`echo ${JOB_NAME} | awk -F"${ENV}-" '{print $NF}'`
DATA="/data/project"
APPPATH="/opt/app/${PROJECT}"

if [ -z $VERSION ];then
  VERSION=${DEPLOY_NUMBER}
fi

for num in 01 02 03 04
do
  server_host=`ansible-inventory -i ${DATA}/${ENV}/inventory_dir/${ENV}-project --host ${PROJECT}${num} -y 2>/dev/null | awk '{print $NF}'`
  if [ `echo ${server_host} | grep ansible_ssh_host | wc -l` -ne 1 ]; then echo "Error: Get inventory host is fault!"; exit 124; fi
  server_url="http://${server_host}:58881"
  echo "deploy ${PROJECT}${num}"
  echo "server_host: $server_host"

  echo "\n### Stop App ${PROJECT}${num} ###"
  curl -sXPOST ${server_url}/command -d 'command='${APPPATH}'/start.sh stop' | jq . | tee tmp01
  if [ $? -ne 0 ]; then echo "Error: curl is fault!"; exit 127; fi
  if [ "`cat tmp01 | grep code`x" != "x" ]; then echo "Error: Stop App is fault!"; /bin/rm tmp01; exit 128; fi

  echo "\n### Clean work space ${PROJECT}${num} ###"
  curl -sXPOST ${server_url}/command -d 'command=/bin/rm -f '${APPPATH}'/'${PROJECT}'.jar' | jq . | tee tmp01
  if [ $? -ne 0 ]; then echo "Error: curl is fault!"; exit 127; fi
  if [ "`cat tmp01 | grep code`x" != "x" ]; then echo "Error: Clean work space is fault!"; /bin/rm tmp01; exit 128; fi

  echo "\n### Rsync jar ${PROJECT}${num} ###"
  curl -s ${server_url}/upload/file -F "file=@${DATA}/env_dir/${PROJECT}/${VERSION}/${PROJECT}.jar" -F type=1 | jq . | tee tmp01
  if [ $? -ne 0 ]; then echo "Error: curl is fault!"; exit 127; fi
  if [ "`cat tmp01 | grep code`x" != "x" ]; then echo "Error: Rsync jar is fault!"; /bin/rm tmp01; exit 128; fi

  echo "\n### Start App ${PROJECT}${num} ###"
  curl -sXPOST ${server_url}/command -d 'command='${APPPATH}'/start.sh start' | jq . | tee tmp01
  if [ $? -ne 0 ]; then echo "Error: curl is fault!"; exit 127; fi
  if [ "`cat tmp01 | grep code`x" != "x" ]; then echo "Error: Start App is fault!"; /bin/rm tmp01; exit 128; fi

  sleep 5
  echo "\n### Check App ${PROJECT}${num} ###"
  curl -sXPOST ${server_url}/command -d 'command='${APPPATH}'/start.sh status' | jq . | tee tmp01
  if [ $? -ne 0 ]; then echo "Error: curl is fault!"; exit 127; fi
  if [ "`cat tmp01 | grep code`x" != "x" ]; then echo "Error: Check App is fault!"; /bin/rm tmp01; exit 128; fi
  /bin/rm tmp01
done

echo ""
echo "##### DEPLOY INFO #####"
echo "PROJECT: ${PROJECT}"
echo "BUILD_VERSION: ${VERSION}"
echo "MD5: `md5sum ${BUILD}/${PROJECT}/${VERSION}/${PROJECT}.jar`"
echo "#######################"

echo ""
if [ -f ${DATA}/${PROJECT}/${VERSION}/version ]; then cat ${DATA}/${PROJECT}/${VERSION}/version; fi
echo ""
echo "##### ALL done. #####"
```
