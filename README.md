# drone-rancher

Drone plugin to deploy or update a project on Rancher.

This is unofficial drone-rancher plugin for drone v0.5

Docker image available at:
https://hub.docker.com/r/majest/drone-rancher/

## Config
The following parameters are used to configure this plugin in .drone.yml:

* **rancher_url** - url to your rancher server, including protocol and port
* **access_key** - rancher api access key
* **secret_key** - rancher api secret key
* **stack** - name of rancher stack to deploy to
* **service** - name of rancher service that's in the given stack to act on
* **docker_image** - new image to assign to service, including tag (`drone/drone:latest`)
* **start_first** - start the new container before stopping the old one, defaults to `true`
* **confirm** - auto confirm the service upgrade if successful, defaults to `false`
* **timeout** - the maximum wait time in seconds for the service to upgrade, default to `30`

The following secret values can be set to configure the plugin by command line.

* **PLUGIN_RANCHER_URL** corresponds to **rancher_url**
* **PLUGIN_ACCESS_KEY** corresponds to **access_key**
* **PLUGIN_SECRET_KEY** corresponds to **secret_key**
* **PLUGIN_STACK** corresponds to **stack**
* **PLUGIN_SERVICE** corresponds to **service**
* **PLUGIN_DOCKER_IMAGE** corresponds to **docker_image**
* **PLUGIN_START_FIRST** corresponds to **start_first**
* **PLUGIN_CONFIRM** corresponds to **confirm**
* **PLUGIN_TIMEOUT** corresponds to **timeout**


The following is a sample Rancher configuration in your `.drone.yml` file:

```yaml
deploy:
  image: majest/drone-rancher
  rancher_url: https://example.rancher.com
  access_key: 1234567abcdefg
  secret_key: abcdefg1234567
  stack: mystack
  service: myservice
  docker_image: drone/drone:latest
```

You can also add secrets via command line **_instead_** of adding them to your ``drone.yml`` by 
```

drone secret add --image=majest/drone-rancher {$YOUR_REPO} PLUGIN_ACCESS_KEY {$YOURACCESSKEY}
drone secret add --image=majest/drone-rancher {$YOUR_REPO} PLUGIN_SECRET_KEY {$YOURSECRETKEY}

# don't forget to sign your .drone.yml.sig after changing the .drone.yml by 

drone sign exocode/gn2016
```
