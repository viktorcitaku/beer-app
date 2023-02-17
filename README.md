# beer
Simple Golang Application using a public API to pickup beers

### System requirements

```
Docker: 20.10.22
Golang: 1.19.5
Browser: Chrome (latest)
```


### How to build and run the application?

Run the dependent services under `deployments` folder with following command:

`docker compose up --detach --build`

To tear down the above services, run the following:

`docker compose up --detach --build`

Set the environment variable `BEER_STATIC_FILES` with the absolute path of `web` folder.

#### NOTE: WORK IN PROGRESS!!!