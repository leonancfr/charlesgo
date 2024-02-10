
# CharlesGO
This software is responsible by connect the CharlinhOS into Gabriel infrastructure.
## Dependencies
You need install Golang. Is recommended use `asdf` with Golang plugin to install the correct Golang version.

### Install `asdf`
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.13.1

Open `~/.bashrc` and put it at end of file:

    # asdf - language/tools version manager #
    source ~/.asdf/asdf.sh
    source ~/.asdf/completions/asdf.bash

After that, run `. .bashrc`, to test the `asdf` use `asdf info`

### Install Golang
    asdf plugin-add golang
    asdf install golang 1.21.1
    asdf global golang 1.21.1
    asdf shell golang 1.21.1

Test the golang installation with `go version`



## How to Run and Build

1. Create a `.env` file in the root folder of the project. Copy the contents of the `.env.sample` file into the newly created `.env` file, and fill in the empty fields.

2. Export the environment variables using the following command:
    ```bash
    if [ -f .env ]; then
        export $(cat .env | sed 's/#.*//g' | xargs)
    fi
    ```

### Run Without Building

To run the application without building, execute the following command:
```bash
go run main/main.go
```



### Container Build
To build the Docker container, run:
```bash
docker build -o . .
```

### Local Build
To perform a local build, execute:
```bash
./build.sh
```

### Generated Artifacts
The files CharlesGo (mipsle) and LinuxGo (x64) will be generated after the build process.

### The --config flag
You can use the --config flag to specify a file with additional configurations when using the software on x64 architecture. Example: `./LinuxGo --config config.ini`.

##### Set config.ini file
Set the label as the MAC address of the target device and set each flag as `true` or `false` as desired.

`config.ini`
```
[DEVICE_INFO]
LABEL="XX:XX:XX:XX:XX:XX"

[SUPERVISOR]
ENABLE=true

[MQTT]
ENABLE=true

[UPDATE]
ENABLE_STM32=true
ENABLE_HLK7628=true
```
## Upload to device
1. Disable root ssh protection
    1. Log-in with gabriel user
    2. Run `su` to change to root user
    3. Change the field `option RootPasswordAuth` in file `/etc/config/dropbear` to `'off'`
    4. Run `service dropbear restart`

2. Now you can send the software to device with 
    ```
    scp -P 2202 CharlesGo root@172.17.2.1:/opt/gabriel/bin
    ```


