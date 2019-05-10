docker stop $(docker ps -aq)
docker rm c1
docker rm c2
docker rm c3
docker run -d --name=c1 -p 8500:8500 consul agent -dev -client=0.0.0.0 -bind=0.0.0.0
#ifconfig ver (docker0)

IP=$(docker inspect --format '{{ .NetworkSettings.IPAddress }}' c1); echo $IP
docker run -d --name c2 consul agent -dev -bind=0.0.0.0 -join=$IP
docker run -d --name c3 consul agent -dev -bind=0.0.0.0 -join=$IP

docker exec -ti c1 /bin/sh
consul monitor

docker exec -ti c2 /bin/sh
apk add git
apk add go
apk add libc-dev
go get github.com/Mendez9000/consul/discoveryServiceNode
cd ~/go/src/github.com/Mendez9000/consul/discoveryServiceNode
go run discoveryServiceNode.go

docker exec -ti c3 /bin/sh
apk add git
apk add go
apk add libc-dev
go get github.com/Mendez9000/consul/discoveryServiceNode
cd ~/go/src/github.com/Mendez9000/consul/discoveryServiceNode
go run discoveryServiceNode.go

#**************************************************************************************
#Configurar watcher a servicio

vi /consul/config/config.json

{
  "watches": [
    {
        "type": "key",
        "key": "foo",
        "handler": "/scripts/handler.sh"
    },
    {
        "type": "keyprefix",
        "prefix": "foo",
        "handler": "/scripts/handler.sh"
    },
    {
        "type": "key",
        "key": "app1/config/secureModeLevel",
        "handler_type": "http",
        "http_handler_config": {
           "path":"http://172.17.0.1:8091/watch-conf",
           "method": "POST",
           "timeout": "10s"
        }
     }
  ]
}
~

#Recargar configuracion de consul
consul reload

#Conul ui
http://localhost:8500/ui/

#Enviar cambios a consul
curl     --request PUT     --data '26'     http://localhost:8500/v1/kv/app1/config/secureModeLevel
~



