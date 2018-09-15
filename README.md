# LoRa IoT Sensors in Low Power Wide Area Network (LPWAN)

Low-power WAN (LPWAN) is a wireless wide area network specification that interconnects lowbandwidth, battery-powered sensors with low bit rates over long ranges.
An Implementation of a LoRa and LoRaWAN based system of sensors with fully working uplink and downlink and support for sending mac-commands to the end-nodes

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

- The server must run on Unix based system (Linux/macOS)
- the gateway must run on an RPI3 running Raspbian. a LoRa Shield sx1272 is needed in order for it to work as a receiver. follow this site in order to install Raspbian and attach the module to the RPI3:
	(https://www.raspberrypi.org/documentation/installation/installing-images/)
	(https://www.hackster.io/ChrisSamuelson/lora-raspberry-pi-single-channel-gateway-cheap-d57d36) 
- The end node is made of the Nucleo Board Nucleo-L073RZ and a LoRa Shield sx1272.
- Install System Workbench for STM32 in order to flash the program to the nucleo board.
	(http://www.openstm32.org/HomePage)


### Installing the Server
 - Follow loraserver.io to install the server on a Linux machine. Make sure to install both lora-server and lora-app-server.
	https://www.loraserver.io/guides/debian-ubuntu/
	any needed configurations can be found in the site
	This will give you a fully working server that supports uplink and downlink

- In order to support mac commands install gRPC-Golang from https://grpc.io/docs/quickstart/go.html
	copy the scripts from the repository and place them in a new directory in /Documents

### Installing the Gateway
- Download the implementation to the RPI3 from the repository and place it in /Documents
- Compile the software using make.
	- If you wish to work only in uplink mode use make lora_gateway_pi2
	- If you wish to have uplink and downlink use make lora_gateway_pi2_downlink
- open key_MQTT.py and set the variable MQTT_server with the ip of the server (can be found with ifconfig)


### Installing the end node
- Import the downloaded project to System Workbench for STM32 IDE
- Connect the nucleo board with usb to the PC
- Flash the implementation to the board using the green play button in the IDE
- Run the command "cat /dev/serial/by-path/<serial address>" in order to get the log of the end node in the terminal


## Running the System
The end node:
	The end node is ready to use and already transmitting data
	You can choose an AES key in "STM32CubeExpansion_LRWAN_V1.2.0/Projects/STM32L073RZ-Nucleo/Applications/LoRa/End_Node/LoRaWAN/App/inc/commissioning.h" so that the messages will be encrypted

The gateway:
- cd to the directory where you unpaced and compiled the program
- run ./cmd.sh
- A menu will appear in the terminal. choose 1 or 2 depending on how you compiled the program (1 for only downlink, otherwise choose 2)
- The program will check if the LoRa module is available and you should see an incoming message from the end node

The server:
-	Run the follwoig commands:
	mosquitto -d
	sudo systemctl start loraserver.service
	sudo systemctl start lora-app-server.service
- Run the command "sudo journalctl -u <service name>" and make sure no error have accoured 
- open the browser and go to  https://localhost:8080/ (or if you changed the port during the installation and configuration proccess then use you new port)
- follow the on screen discription to add the gateway and end node to the server.


### Deployment
The end node is working on a 50 second cycle, it will send a message to the server automatically as long as it's connected to the computer.  The message will go through the gateway to the server, there you will be able to see in the app server the uplink message and its info. next you'll see the downlink log of the message sent from the server to the end node. By default, if no mac command is sent, an ACK message will be sent.  If you want to send a mac command run one of the gRPC scripts you downloaded, each script builds and send a specific mac command



## Refernces and Repositories

[brocaar]https://github.com/brocaar/loraserver
[CongducPham]https://github.com/CongducPham/LowCostLoRaGw

## Authors

* **Ibraheem Baleek** **Bashir Khayat** **Jeries Nasser**


