# Infestor Example Tests
As an example we will use `setzer` test.
[Setzer](https://github.com/makerdao/setzer-mcd) is a executable that pulls prices from exchanges.

#### Setting up Setzer
To test it we need first of all to build a Docker container with the `setzer` executable and it's requirements.
NOTE: it's not so easy, so we have prebuilt Docker image with the `setzer` and you have to build it on your machine.

Run `docker build` command from project root folder:
```bash
$ docker build -t setzer -f example/Dockerfile .
```

#### Setting up Smocker

Next we need running [Smocker](https://smocker.dev/), we will run it as separate container as well.

Starting it using `docker run` command:

```bash
$ docker run -d \
  --restart=always \
  -p 8080:8080 \
  -p 8081:8081 \
  --name smocker \
  thiht/smocker
```

#### Running environment

**Yes we will use legacy `--link` docker feature**

We need to start `setzer` container linked with running `smocker` container.

```bash
$ docker run -it --rm -v $(pwd):/app --link smocker --name setzer setzer
```

From there you will be able to send requests to `smocker` container using `http://smocker` URL.

Execution of tests: 

```bash
$ cd example && go test ./...
ok  	github.com/chronicleprotocol/infestor/example	0.129s
```


