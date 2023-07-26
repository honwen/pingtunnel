# Deploy with docker-compose

**First** edit `.env` file in this directory to your appropriate value.

**Then** run stack with these commands:

- build image

```shell
ocker build -t pingtunnel .
```

- in the server

```shell
docker-compose -f server.yml up -d
```

- in client machine

```shell
docker-compose -f client.yml up -d
```

**Now** use socks5 proxy at port `1080` of your client machine
