# drone-rancher

Drone plugin to deploy or update a project on Rancher.

This is unofficial drone-rancher plugin for drone v0.5

Docker image available at:
https://hub.docker.com/r/majest/drone-rancher/

## Config
The following parameters are used to configure this plugin:

* **url** - url to your rancher server, including protocol and port
* **access_key** - rancher api access key
* **secret_key** - rancher api secret key
* **stack** - name of rancher stack to deploy to
* **service** - name of rancher service that's in the given stack to act on
* **docker_image** - new image to assign to service, including tag (`drone/drone:latest`)
* **start_first** - start the new container before stopping the old one, defaults to `true`
* **confirm** - auto confirm the service upgrade if successful, defaults to `false`
* **timeout** - the maximum wait time in seconds for the service to upgrade, default to `30`

The following secret values can be set to configure the plugin.

* **URL** corresponds to **url**
* **RANCHER_ACCESS_KEY** corresponds to **access_key**
* **RANCHER_SECRET_KEY** corresponds to **secret_key**
* **RANCHER_STACK** corresponds to **stack**
* **RANCHER_SERVICE** corresponds to **service**
* **DOCKER_IMAGE** corresponds to **docker_image**
* **START_FIRST** corresponds to **start_first**
* **CONFIRM** corresponds to **confirm**
* **TIMEOUT** corresponds to **timeout**


The following is a sample Rancher configuration in your `.drone.yml` file:

```yaml
deploy:
  image: majest/drone-rancher
  url: https://example.rancher.com
  access_key: 1234567abcdefg
  secret_key: abcdefg1234567
  stack: mystack
  service: myservice
  docker_image: drone/drone:latest
```
