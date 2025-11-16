# docker run --rm -itd --name=postgres -p 5432:5432 -e POSTGRES_USER="root" -e POSTGRES_PASSWORD="root" postgres:13.22-alpine3.22
# docker run --rm -itd --name=postgres --network mynet -e POSTGRES_USER="root" -e POSTGRES_PASSWORD="root" postgres:13.22-alpine3.22
docker run -it --network mynet --env-file="./config/.env" --name avito --rm -p 8080:8080 avito