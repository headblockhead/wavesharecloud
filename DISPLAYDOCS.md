# Features

- The device is suitable for use as equipment labels, shelf labels, information storage and more.
- It integrates an SW6106 chip that supports USB-C quick charging protocols like: PD / QC / FCP / PE / SFCP, etc.
- The device allows the user to control the display content via a remote server making it convenient and flexible.
- There is a voltage detection circuit can detect the battery level to make sure the device is not over or under charged.
- The device has support for a user-configured ID string that can be used to identify the device in a network.
- There is an Android app available that allows the user to configure the device over Bluetooth.
- Sadly, there is no backlight. However, the device's display will not change after the display has been written to (until the next write).
- Power consumption is tiny, the battery is only used when refreshing the display.
- It comes with an instruction manual.

# Specifications

| Label                     | Value                |
| ------------------------- | -------------------- |
| Power supply              | Battery              |
| Resolution                | 400x300 (4:3 aspect) |
| Pixel pitch               | 0.212 × 0.212        |
| Display colors            | Black, White         |
| Refresh time              | 4 seconds            |
| Viewing angle             | >170°                |
| Device Dimensions         | 96.5mm x 85mm        |
| Display Dimensions        | 84.8mm x 63.6mm      |
| Refreshes/1000mah battery | 1500+                |

# Setup

When first setting up, you will need the Android app. You can download it from the [WaveShare Website](https://www.waveshare.net/w/upload/Cloud_app.apk), However if you cannot download it from this link you can download it from [The Wayback Machine](https://web.archive.org/web/20220719161209/https://www.waveshare.net/w/upload/Cloud_app.apk).

---

# Below is not a valid markdown file, it has not yet been edited by me. Proceed if you dare! (This is a copy-paste of the original website's text)

You can configure the devices via wifi or Bluetooth, if you are the first time configure the device, only the Bluetooth is available.
The firmware of the 4,2inch cannot be reprogrammed by used because the flash is locked.
You need to wakeup the device by button

Lead Information
Every time the device start, it will do partial refresh and display status icons.
Hereby provide the refrence of icons.
Cloud ESP32 e-Paper Board wait.png Cloud ESP32 e-Paper Board set.png Cloud ESP32 e-Paper Board batter.png Cloud ESP32 e-Paper Board wifi.png Cloud ESP32 e-Paper Board wifi connect.png
Waiting Setting Low Voltage WIFI Host
Waiting: The device is waiting for commands.
Setting: Setting is finished.
Low Voltage: The voltage of battareis is lower than warnning value.
WIFI: The wifi is connected.
Host: The devices is connected to target host by IP address.
Generally, you should press the button to wake up the device and check the icons. The warnning voltage is 3600mv, once the voltage is 150mv lower than warnning voltage (3450mv), the devices will shutdown automatically to protect the stable of the whole system.
Configure Device
First Setup
If you didn't configure the device before, you should configure the device by APP after pressing the Wakeup button to update the display①. Please refer to #1.4 Configure Device by APP to configure the device.②
Note：
①If the device isn't configured, Waiting icon is displaed in the top-right area. If the Low Voltage icon is displayed without Waiting, it means that the batteris is less than 3450mv and it is going to shutdown.
②If the device doesn't connect to Bluetooth, it will shutdown after 90s after booting and refreshing the display.

Reconfiguration
If the device was configured, the device will update and display Setted icon in the top-right area① after pressing wakeup button. The Bluetooth is disabled by default, if required, you should hold the wakeup button for 5s at least to enable the Bluetooth②. The updating process of device will not be interrupted③ while enabling the Bluetooth. After enabling the Bluetooth, you can re-configure the device by referring to #1.4 Configure Device by APP ④.
Note：
①The update time of the device is determined by the speed of WIFI, It should less than 30s as we test.。
②You can hold the wake-up button until your phone scan the device via Bluetooth. Otherwise, the device auto-shutdown if the Bluetooth is disconnected.。
③If you enable the Bluetooth of the device, it will try to connect to master (phone) in 30s if the shutdown command is received via wifi. If the device is connected to the phone by Bluetooth, it will keep waking, otherwise, it will be turned off after 30s.
④You should reboot the device after configuration to make the configuration effect. Please do not reboot the device when transmitting data via WIFI, it will cost data loss.

APP Description
Cloud Epd app 1.png
Bluetooth Connection Button and the information
Cloud Epd app 2.png
Device ID: for distinguish devices
Cloud Epd app 3.png
WIFI_SSID
Cloud Epd app 4.png
Check the current SSID connected
Cloud Epd app 5.png
WIFI_Password
Cloud Epd app 6.png
Host_IP, The IP of host, for example, the IP of Raspberry Pi.
Cloud Epd app 8.png
Device_IP: This is used to set the static IP. If you enable DHCP, the static IP is unavailable
Cloud Epd app 7.png
Device_Password, you should input the device password to vertify if the device is locked.
Note: The default device password of the Raspberry Pi example is 123456. If you lock the device you have to unlock it with password 123456, otherwise, the device cannot work.

Cloud Epd app 9.png
Warning_voltage: If the voltage of batteries is less than the warning voltage, the device will display the warning icon. If the battery is 150mV lower than the warning voltage, the device will shutdown automatically.

Cloud Epd app 10.png
Load the configuration saved
Cloud Epd app 11.png
Save the current configuration, you can save up to four sets of configuration
Cloud Epd app 12.png
Upload the current configuration to device.
Cloud Epd app 13.png
Formating: Clean the current configuration
Configure Device by APP
Ⅰ. Open APP（APP will auto-save the last configuration information）
Cloud ESP32 e-Paper Board manual 1.png
Ⅱ. Click the Bluetooth CONNECTION button, the default Bluetooth device is WaveShare_EPD or the Device ID configured.
Unpaired
Unparied
Paried
Paried
Ⅲ. Choose the device, for example, connect the WaveShare_EPD. If you are the first time to connect the WaveShare_EPD device, it should be paired first.
Cloud ESP32 e-Paper Board manual 4.png
Ⅳ. Modify the configuration information and Upload(If you have configured the password, you need to input the password as well.)
Cloud ESP32 e-Paper Board manual 5.png
Ⅴ. The APP will disconnect and reboot the device if the configuration is uploaded successfully.
Cloud ESP32 e-Paper Board manual 6.png
Note: We recommend you set static IP for the device.
Communicating Protocol
Communicating is divided into two modes: Command mode and the data mode. Command mode is used for sending commands. data Mode is used for sending image data to e-Paper.

Command Format
‘;’+Command（+Data）+'/'+Parity

Data Format
0x57+4Byte addr+ 4Byte len +1Byte num + len Byte data +Verify

Return Format
'$'+Data+'#' The format of response of Command and Data are the same

Note:The Verity is the XOR result of data which is marked in red

Command Mode：
Command Format
‘;’+Comamnd（+Data）+'/'+Verify
Commands (locked)
Comamnd Desctiption Return
'C' Check if the device is locked Parity bit + Flag bit 0 or 1, 0: unclocked, 1: locked.
‘N' + Device password Unclocked the device Parity bit + Flag bit 0 or 1, 0: failed to lock; 1: lock the device successffully.
'G' Get the ID of device ID
'r' + Sleep time (<9999) Set the device to sleep mode Parity bit
PS: The sleep command is only available in 2.13inch e-Paper Cloud Module.

These commands can be used if the device is locked.
Comamnds (unlocked)
Command Description Return
'0' + name Modify the ID Parity bit
'1' + IP address Modify the IP address of Host Parity bit
‘2’+SSID Modify the WIFI SSID Parity bit
‘3’+password Modify the WIFI password Parity bit
‘P’+userpassword Modify the device password Parity bit
'L' +'0' / 'L' + '1' Set the device lock; 1 to lock and 0 to unlock Parity bit
's' + '0' /'s' + '1' Set the flag bit, 0 to disable and 1 to enable Parity bit
'F' Enter data mode Parity bit
'B' Open Bluetooth Parity bit
'b' Check the current voltage of battery Parity bit + voltage of battery (mv)
'S' Shutdown Parity bit
'R' Restart Parity bit
These commands can be used when the device is unlocked.
Data Mode：
Data Format
0x57+4Byte addr+ 4Byte len +1Byte num + len Byte data +Verify
Data Lebs Content
addr 4byte Address if data
len 4byte Length os data
num 1byte The frame number of current sector
data len byte The data transmitted
Note：
Recommend you to transmit the frames with same lenght
The size of frame transmitted should not larger than 1100Byte, otherwise it cose data lose.
num should be static variable because it may be invalid because of version update.
The data frame doesn't have a stop bit, you need to wait for the verity data before sending the next frame, otherwise, it causes failure.
The e-Paper will update automatically and exit from update mode when the addr and len are 0.
For mare detailes, please refer to the python3 examples provided.

Using Guides for RPI
Install Libraries
#python3
sudo apt-get update
sudo apt-get install python3-pip
sudo apt-get install python3-pil
sudo apt-get install python3-tqdm
sudo apt-get install python3-numpy
sudo apt-get install python3-progressbar
Download the demo codes
Open a terminal and runthe following commands：
sudo apt-get install p7zip-full
sudo wget https://www.waveshare.com/w/upload/2/2e/Cloud_RPI.7z
7z x Cloud_RPI.7z
cd Cloud_RPI
python
The demo codes can only support python3.
Please go to the directory of Cloud_RPI and run the command:
#This code is used to climb the picture of the Waveshare website and transmit the image data to the slave device.
sudo python3 ./examples/display_WS.py
#The codes will draw figures and send the image data to the slave device.
sudo python3 ./examples/display_EPD.py
API Description
There are three directories in lib,http_get、tcp_server 和 waveshare_epd, They are used to climb HTTP pictures, TCP service, and the functions of e-Paper.
Cloud ESP32 e-Paper Board manual 10.png

tcp_sver.py
Path:Cloud_RPI/lib/tcp_server
Vreate a tcp_server class in tcp_sver.py file. You need to inherits the class and refactor the handle function when using.

def handle(self)
Every time the new client connected, it should call the handle function.

Receive Message
def Get_msg(self)
Command Return
'$'+Data+'#' Data
Send Command
def Send_cmd(self,cmd)
Parameter cmd is the command sent
Command
cmd ‘;’+cmd+'/'+Parity
Sent data
def Send_data(self,data)
The parameter data is message transmitte (inlcuded addree and lenghta nd so on) data.
Command
data 0x57+data+ Parity
Set size
def set_size(self,w,h)
w: The width of image; h: The height of image.
Take bicolor e-Paper as example, 1 bit stands of one pixel, then you we gets
Len of data=Width of image(w)\*Height of image(h)/8
Refer to 4.2inch e-Paper Module

Update function
def flush_buffer(self,DATA)
DATA; The image data. the image data can be get by the getbuffer function.
Paramter Send times Lenght of every frame（len） Content DATA(Image data) Total lenght of image data/lenght of singal frame（len） 1024 Byte（Configurable） 0x57+4 Byte addr+ 4 Byte len +1 Byte num + len Byte data+Parity
The lenght of singal transmittion should less than 1100 Byte, or it will cause data loss.

Check voltage of battery
Get the current voltage
def check_batter(self)
Power Off Function
Power off or low power state
def Shutdown(self)
http_get.py
Path：Cloud_RPI/lib/http_get
Download picture
def Get_PNG(Url,Name)
This function is used to download the picture from Url and save it to the current directory with Name

epd4in2.py
Path：Cloud_RPI/lib/waveshare_epd
Convert the picture to image data.
def getbuffer(self, image):
waveshare_epd.py
Directory: Cloud_RPI/lib/waveshare_epd
Convert image information to queue
def getbuffer(self, image):
Configure Windows
This guide is made in Windows 10
Note：
Please make sure that you have installed python3 on your Windows PC and the default version is python3 if you installed multiple versions.
You may need to close the firewall to make the python work.
Install libraries
Open a CMD or Powershell to install libraries with the following commands:

#python3
python -m pip install tqdm
python -m pip install pillow
python -m pip install numpy
python -m pip install pypiwin32
python -m pip install progressbar
Download demo codes
Download the Windows demo, unzip, and enter the Cloud_WIN directory.

python
Note that you need to run the CMD or PowerShell under the Cloud_WIN direcoty and turn the following command.

python ./examples/\*\*\*inch_display_EPD
For example：

#If you have 4.2inch e-Paper Cloud Module
python ./examples/4.2inch_display_EPD
#If you have2.13inch e-Paper Cloud Module
python ./examples/2.13inch_display_EPD
Resource
Related Guides
Making BMP file for e-Paper
Androdi app
Android APP
Android APP Souces Codes
You can also scan the below QR Code to install the APP
Cloud APP.png
Raspberry Pi Examples
Raspberry Pi examples
Windows Examples

FAQ
Question:What is the usage environment of the e-ink screen?
Answer:
【Working conditions】Temperature range: 0~50°C; Humidity range: 35%~65%RH
【Storage conditions】: Temperature range: below 30°C; Humidity range: below 55%RH; Maximum storage time: 6 months

Support
If you require technical support, please go to the Support page and open a ticket.
