## The dotnet-sample sample

The *dotnet-sample* shows how to create a basic .NET application. This sample provides a simple key-value store API using .NET and SQLite. And provides a swagger sample endpoint.
The focal point of the sample is to show how to create a .NET container-based application and deploy it via App Center.

The sample contains the following files.

```
├── build
│   └── common.mak
└── dotnet
    ├── containers
    │   └── dotnet-sample
    │       └── webserver
    │           ├── Dockerfile
    │           └── dotnet-webserver
    │               ├── appsettings.Development.json
    │               ├── appsettings.json
    │               ├── dotnet-webserver.csproj
    │               ├── Program.cs
    │               └── Properties
    │                   └── launchSettings.json
    ├── Makefile
    └── README.md

```

In the `containers` directory, you find the source code for building the container images used by the dotnet-sample. Only one container image is needed in this sample and it is placed in the subdirectory named `dotnet-sample/webserver`. Note that when built, the image name will reflect this directory structure and thus you will see `dotnet-sample/webserver` be part of the container image name.

The other main artifact that goes into building an App is the helm chart, which you find in the `helm-charts` directory. For this sample we have one helm chart named `dotnet-sample` which is placed in a subdirectory of the same name.

<br>

### The App sandbox

With the App sandbox developer option enabled, the system will host both a *container registry* and a *helm chart repository*. These come in handy when you want to run and test your App in a system while still keeping everything local. You can push container images to the sandbox registry by prefixing them with `<system-ip>:5000/sandbox.io/`. These images can be pulled from any helm chart that is installed in the system. Thus, you can now refer to sandbox container images from your local helm chart and install it directly with the *helm* command.

The helm chart itself can also be pushed to the sandbox helm repository which is accessible at `<system-ip>/app-sandbox`. We use [ChartMuseum](https://chartmuseum.com/) as the back-end for this repository and you can read more about how to access the API's [here](https://github.com/helm/chartmuseum?tab=readme-ov-file#api).

Once the helm chart of an App is pushed to the sandbox repository, it will become visible in the App Center and you can then install, upgrade and uninstall your App from here. Note, that for a helm chart to be considered an App, it must use the `app-registration` chart as a *subchart* (see the `Chart.yaml` file for how this is done).

<br>

### Building and running the dotnet-sample App

In the `dotnet-sample` directory you will find a `Makefile` that makes it easy to build and publish the container images and helm charts of the sample. Everything published using the `Makefile` will be to the App sandbox running inside the system.

```make
IMAGES = dotnet-sample/webserver_v0.0.1
CHARTS = dotnet-sample_v0.1.4

include ../build/common.mak
```

The [`Makefile`](./Makefile) located at this same directory provides an convenient way to build and publish both the images and the helm-charts of the sample. Everything published using the `Makefile` will be to the App sandbox running inside the system.

The `Makefile` defines two macros listing respectively the container images and the helm charts to handle.

The last line of the `Makefile` includes the `common.mak` file, which defines a set of useful rules by using the **App Builder** in the background. The most important *targets* defined by these rules are:

```bash
# build container images from source code and helm charts
make build
# push built container images and helm charts to the App sandbox
make push

# install helm chart in system from local file
make install-chart-from-file
# install helm chart in system from sandbox repository
make install-chart-from-repo
# uninstall helm chart from system
make uninstall-chart
```

#### From the Docker view
So, to build the dotnet-sample container image, you must navigate your terminal to the directory of the `Makefile` and then run the command

```bash
make build
```

After a successful build, running `docker images` should now give you output similar to what is shown below

```text
REPOSITORY                                          TAG       IMAGE ID       CREATED         SIZE
10.10.16.34:5000/sandbox.io/dotnet-sample/webserver   0.0.1     a1f80cd4eb00   3 seconds ago   148MB
```

Here you can see that the image is named in accordance with the sandbox container registry hosted inside the system. Also you can see that the image is tagged with version `0.0.1` which is the version used in the `Makefile` (the version used in the `IMAGES` macro).

If you are curious what happens behind the scene when running `make build-image`, then you can add the `-n` command line parameter. This will instruct `make` to not build the image, but instead just show the commands that it would have executed in order to do the build. For the above case you should see output similar to this when running `make -n build-image`

```bash
docker build containers/dotnet-sample/webserver -t 10.10.16.34:5000/sandbox.io/dotnet-sample/webserver:0.0.1
```

You can use the `-n` option for all the other make targets as well.

To push the image to the sandbox registry, run the command

```bash
make push
```

To verify that the push worked, you can first remove the version you already have in the local cache and then try to pull it again from the sandbox. Remember to replace the system IP address below to match the one you are using.

```bash
docker rmi 10.10.16.34:5000/sandbox.io/dotnet-sample/webserver:0.0.1
docker pull 10.10.16.34:5000/sandbox.io/dotnet-sample/webserver:0.0.1
```

#### From the Chart view
The sandbox images can also be pulled from within the system and this is exactly what the dotnet-sample helm chart does. So, let us now build the helm chart and test that it works. You build the helm chart by running

```bash
make build
```

You will notice that there will now be a file named `dotnet-sample-0.1.4.tgz` in the `helm-charts` directory. 

```
dotnet-sample/
└── helm-charts
    ├── dotnet-sample
    └── dotnet-sample-0.1.4.tgz   <---- helm chart package
```

This is a file that contains the entire helm chart as one re-distributable package. The version number `0.1.4` comes from the `Chart.yaml` file which is updated by the `Makefile` from the `CHARTS` macro. To push the package to the sandbox repository, run the command

```bash
make push
```

To confirm that your helm chart has been uploaded successfully, you can run the command

```bash
make list-charts
```

You should get output similar to what is shown below

```yaml
apiVersion: v1
entries:
  dotnet-sample:
  - apiVersion: v2
    appVersion: 0.1.4
    created: "2024-11-08T09:17:40.487168011Z"
    dependencies:
    - name: app-registration
      repository: https://horizonsystem.azurewebsites.net/system
      version: 1.3.0
    description: A dotnet sample application
    digest: 50c17cae1f20d432369defcb664325a57b9a43fd2f4e5df93533193ea57b9a4
    name: dotnet-sample
    type: application
    urls:
    - charts/dotnet-sample-0.1.4.tgz
    version: 0.1.4
generated: "2022-04-09T09:17:40Z"
serverInfo: {}
```

At this point the App should be visible in the App Center and you can install / upgrade and uninstall it from here.

As an alternative you can also install the App directly from the terminal. There are two targets available

```bash
make install-chart-from-file
make install-chart-from-repo
```

The first will install the chart directly from the helm chart file you built locally. In the above case, it would install the chart from the file named `dotnet-sample-0.1.4.tgz` in the `helm-charts` directory.

The second will install the chart from the sandbox repository; so basically the same as what would happen if you install the App from the App Center. For this to work, you of course have to push the chart first, like we did with the command `make push-chart`.

To verify that the dotnet-sample App is running, navigate your browser to `https://<system-ip>/dotnet-sample`. This this show the simple hello world message.

There is one final target available named `uninstall-chart` which can be quite helpful. It does what it says; it simply uninstalls the dotnet-sample chart from the system.

<br>

### HTTP route and API Base Path

The base path for this sample is `dotnet-sample/`. <br>
In the `httproute.yaml`, requests with the path prefix `/dotnet-sample/` are routed to the dotnet-sample service, and the `/dotnet-sample` prefix is removed.

### API Endpoints

- `GET /items`: Returns all stored key-value pairs.
- `POST /items`: Creates a new key-value pair.
- `PUT /items/{key}`: Updates the value of an existing key-value pair.
- `DELETE /items/{key}`: Deletes a key-value pair based on the provided key.

### Swagger UI

Swagger UI is available at `/swagger` for development and `/dotnet-sample/swagger` for production.

### Example Requests

#### Get All Items

```sh
curl -X GET http://0.0.0.0:80/dotnet-sample/items
```

#### Add a New Item

```sh
curl -X POST http://0.0.0.0:80/dotnet-sample/items -H "Content-Type: application/json" -d '{"key": "exampleKey", "value": "exampleValue"}'
```

#### Update an Item

```sh
curl -X PUT http://0.0.0.0:80/items/dotnet-sample/exampleKey -H "Content-Type: application/json" -d '{"value": "newValue"}'
```

#### Delete an Item

```sh
curl -X DELETE http://0.0.0.0:80/dotnet-sample/items/exampleKey
```