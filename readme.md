# Docker Container Erstellen
sudo DOCKER_BUILDKIT=1 docker build -t go-plz:v1 .

# BuildKit ContainerD Erstellen
## Start the BuildKit daemon:
sudo buildkitd &


export BUILDKIT_HOST=unix:///run/buildkit/buildkitd.sock
buildctl build \
  --frontend=dockerfile.v0 \
  --local context=. \
  --local dockerfile=. \
  --output type=image,name=go-plz:v1

# Oder mit Docker und Import
docker build -t go-plz:v1 .
sudo ctr images import go-plz.tar

# Ausführen
sudo ctr run --rm --net-host docker.io/library/go-plz:v1 go-plz-instance

# Client Code zum Testen
JWT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibm
FtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3
nMiGM6H9FNFUROf3wh7SmqJp-QV30    RPCTYPE=mTLS  go run cmd/client/main.go loc
alhost:50051 Bochum