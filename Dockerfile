FROM eclipse-temurin:22-jdk

WORKDIR /data

ENV RAM=4000M

ADD https://piston-data.mojang.com/v1/objects/64bb6d763bed0a9f1d632ec347938594144943ed/server.jar /app/server.jar

EXPOSE 25565

ENTRYPOINT ["sh", "-c", "echo 'eula=true' > /data/eula.txt && java -Xmx${RAM} -Xms1024M -jar /app/server.jar nogui"]
