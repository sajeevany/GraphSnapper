# graph-snapper
Tool which calls snapshot APIs on a defined cadence and stores them within the desired data storage platform.

Generate swagger files:

    swag init -g cmd/graph-snapper/main.go -o ./docs

Access swagger API:

    ${HOST}:{PORT}/swagger/index.html
    ie. http://localhost:80/swagger/index.html

Run unit tests:

    go test -short ./...
    
Build image:
    
    ./build/buildImage.sh
    
    or 
    
    docker build --build-arg "GIT_COMMIT=CustomCommit" --build-arg "CONFIG_FILE=deploy/graph-snapper-conf_standalone.json" --tag "portfolio-service:CUSTOM" -f build/Dockerfile .