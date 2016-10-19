package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rancher/go-rancher/client"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

type Rancher struct {
	client     *client.RancherClient
	Url        string `json:"url"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Stack      string `json:"stack"`
	Service    string `json:"service"`
	Image      string `json:"docker_image"`
	StartFirst bool   `json:"start_first"`
	Confirm    bool   `json:"confirm"`
	Timeout    int64  `json:"timeout"`
}

func main() {
	app := cli.NewApp()
	app.Name = "docker rancher"
	app.Usage = ""
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "url",
			Usage:  "Url to your rancher server, including protocol and port",
			EnvVar: "URL",
		},
		cli.StringFlag{
			Name:   "access_key",
			Usage:  "Rancher api access key",
			EnvVar: "RANCHER_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "secret_key",
			Usage:  "Rancher api secret key",
			EnvVar: "RANCHER_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "service",
			Usage:  "Name of rancher service to act on",
			EnvVar: "RANCHER_SERVICE",
		},
		cli.StringFlag{
			Name:   "stack",
			Usage:  "Name of rancher stack to act on",
			EnvVar: "RANCHER_STACK",
		},
		cli.StringFlag{
			Name:   "docker_image",
			Usage:  "New image to assign to service, including tag (drone/drone:latest)",
			EnvVar: "DOCKER_IMAGE",
		},
		cli.BoolFlag{
			Name:   "start_first",
			Usage:  "Start the new container before stopping the old one, defaults to true",
			EnvVar: "START_FIRST",
		},
		cli.BoolFlag{
			Name:   "confirm",
			Usage:  "Auto confirm the service upgrade if successful, defaults to false",
			EnvVar: "CONFIRM",
		},
		cli.Int64Flag{
			Name:   "timeout",
			Usage:  "the maximum wait time in seconds for the service to upgrade, default to 30",
			EnvVar: "TIMEOUT",
			Value:  30,
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Rancher{
		Url:        c.String("url"),
		AccessKey:  c.String("access_key"),
		SecretKey:  c.String("secret_key"),
		Service:    c.String("service"),
		Image:      c.String("docker_image"),
		StartFirst: c.Bool("start_first"),
		Confirm:    c.Bool("confirm"),
		Timeout:    c.Int64("timeout"),
		Stack:      c.String("stack"),
	}

	plugin.Exec()
}

// Exec executes the plugin step
func (r Rancher) Exec() {
	rancher, err := client.NewRancherClient(&client.ClientOpts{
		Url:       r.Url,
		AccessKey: r.AccessKey,
		SecretKey: r.SecretKey,
	})

	r.client = rancher

	if err != nil {
		fmt.Printf("Failed to create rancher client: %s\n", err)
		os.Exit(1)
	}

	service, err := r.getService(r.getStackId())

	// try to create the service
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	r.upgrade(service)
}

// get stack id  by stack name
func (r Rancher) getStackId() string {
	environments, err := r.client.Environment.List(&client.ListOpts{})
	if err != nil {
		fmt.Printf("Failed to list rancher environments: %s\n", err)
		os.Exit(1)
	}

	for _, env := range environments.Data {
		if env.Name == r.Stack {
			return env.Id
		}
	}

	fmt.Printf("Unable to find stack %s\n", r.Stack)
	os.Exit(1)
	return ""
}

// Gets the service by stack Id
func (r Rancher) getService(stackId string) (*client.Service, error) {
	services, err := r.client.Service.List(&client.ListOpts{})
	if err != nil {
		fmt.Printf("Failed to list rancher services: %s\n", err)
		os.Exit(1)
	}

	for _, svc := range services.Data {
		if svc.Name == r.Service && svc.EnvironmentId == stackId {
			return &svc, nil
		}
	}

	return nil, fmt.Errorf("Unable to find service %s\n", r.Service)
}

// does not work
// func (r Rancher) create() {
// 	service := &client.Service{Name: r.Service}
// 	service.LaunchConfig = &client.LaunchConfig{ImageUuid: fmt.Sprintf("docker:%s", r.Image)}
// 	_, err := r.client.Service.Create(service)
// 	if err != nil {
// 		fmt.Printf("Unable to create service %s - %s\n", r.Service, err.Error())
// 		os.Exit(1)
// 	}
// }

// publishes the service
func (r Rancher) upgrade(service *client.Service) {

	service.LaunchConfig.ImageUuid = fmt.Sprintf("docker:%s", r.Image)

	upgrade := &client.ServiceUpgrade{}
	upgrade.InServiceStrategy = &client.InServiceUpgradeStrategy{
		LaunchConfig:           service.LaunchConfig,
		SecondaryLaunchConfigs: service.SecondaryLaunchConfigs,
		StartFirst:             r.StartFirst,
	}

	upgrade.ToServiceStrategy = &client.ToServiceUpgradeStrategy{}
	_, err := r.client.Service.ActionUpgrade(service, upgrade)
	if err != nil {
		fmt.Printf("Unable to upgrade service %s\n", r.Service)
		os.Exit(1)
	}

	if r.Confirm {
		srv, err := retry(func() (interface{}, error) {

			s, e := r.client.Service.ById(service.Id)

			if e != nil {
				return nil, e
			}

			if s.State != "upgraded" {
				return nil, fmt.Errorf("Service state: %s\n", s.State)
			}

			fmt.Printf("Service state: %s\n", s.State)

			return s, nil
		}, time.Duration(r.Timeout)*time.Second, 3*time.Second)

		if err != nil {
			fmt.Printf("Error waiting for service upgrade to complete: %s", err)
			os.Exit(1)
		}

		_, err = r.client.Service.ActionFinishupgrade(srv.(*client.Service))
		if err != nil {
			fmt.Printf("Unable to finish upgrade %s\n", r.Service)
			os.Exit(1)
		}
		fmt.Printf("Finished upgrade %s\n", r.Service)
	}
}

type retryFunc func() (interface{}, error)

func retry(f retryFunc, timeout time.Duration, interval time.Duration) (interface{}, error) {
	finish := time.After(timeout)
	for {
		result, err := f()
		if err == nil {
			return result, nil
		}
		select {
		case <-finish:
			return nil, err
		case <-time.After(interval):
		}
	}
}
