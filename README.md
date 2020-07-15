# graph-snapper
Tool which takes snapshots of desired dashboards/graphs and stores them in the specified document store. Intended to capture Grafana data and store it in Confluence.

Generate swagger files:

    swag init -g cmd/graph-snapper/main.go -o ./docs

Access swagger API:

    ${HOST}:{PORT}/swagger/index.html
    ie. http://localhost:80/swagger/index.html

Run unit tests:

    go test -short ./...
    
Run integration tests(requires docker):

    go test -run Integration ./...
    
Run unit and integration(requires docker):

    go test ./...
    
Build image:
    
    ./build/buildImage.sh
    
    or 
    
    docker build --build-arg "GIT_COMMIT=CustomCommit" --build-arg "CONFIG_FILE=deploy/graph-snapper-conf_standalone.json" --tag "portfolio-service:CUSTOM" -f build/Dockerfile .