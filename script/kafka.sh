.\bin\windows\zookeeper-server-start.bat .\config\zookeeper.properties

.\bin\windows\kafka-server-start.bat .\config\server.properties

docker run --ulimit stack=10240:10240 -m 6g -e JAVA_OPTS="-Xms512m -Xmx2g -XX:MaxMetaspaceSize=256m -XX:ReservedCodeCacheSize=128m -XX:+UseContainerSupport -XX:MaxRAMPercentage=80" -p 8083:8083 tchiotludo/akhq