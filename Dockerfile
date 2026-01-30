FROM eclipse-temurin:22-jdk

WORKDIR /app

ENV RAM=4000M

ADD https://piston-data.mojang.com/v1/objects/64bb6d763bed0a9f1d632ec347938594144943ed/server.jar server.jar

RUN echo "eula=true" > eula.txt

EXPOSE 25565

ENTRYPOINT ["sh", "-c", "java -Xmx${RAM} -Xms1024M -jar server.jar nogui"]
