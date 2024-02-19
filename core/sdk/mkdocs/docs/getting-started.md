# Getting Started

## Install Docker Desktop

All we need to run Flare Hotspot development runtime is [Docker](https://www.docker.com/). You can download it from [here](https://www.docker.com/products/docker-desktop).

For a detailed instruction to install Docker Desktop, please refer to the [official documentation](https://docs.docker.com/desktop/) or follow our beginner-friendly tutorial for [windows](./guides/docker-install/windows.md) and for [linux](./guides/docker-install/linux.md).

# Installing Flare Hotspot SDK

## Download Flare Hotspot SDK
To install the plugin SDK, download the latest **devkit-x.x.xzip** file from [sdk-releases](https://github.com/flarehotspot/core/sdk-releases/releases) repository.
After downloading, extract the zip file to your desired location.
```sh
# replace ~/Downloads/devkit-0.0.5.zip with the path to the downloaded zip file
unzip ~/Downloads/devkit-0.0.5.zip -d ~/Documents/devkit-0.0.5
cd ~/Documents/devkit-0.0.5
```

Below is the directory structure of the zip file:
```
|- config
|- core
|- main
|- mock-files
|- plugins
    |-- com.flarego.default-theme
    |-- com.flarego.sample-plugin
```

## Start The SDK Runtime

To start the SDK runtime, you need to run:
```sh
cd ~/Documents/devkit-0.0.5
docker compose up
```

Now you can access the Flare Hotspot web interface:

- Captive Portal: [http://localhost:3000](http://localhost:3000)
- Admin Dashboard: [http://localhost:3000/admin](http://localhost:3000/admin)
- Database Management: [http://localhost:8080](http://localhost:8080)

The default login for the admin dashboard is:
```
username: admin
password: admin
```

To stop the SDK runtime, you need to run:
```sh
cd ~/Documents/devkit-0.0.5
docker compose down
```