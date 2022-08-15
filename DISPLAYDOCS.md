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
| Operating temperature     | 0°C to 50°C          |
| Operating humidity        | ~35% to ~65%         |
| Storage temperature       | <30°C                |
| Storage humidity          | <55%                 |
| Storage time              | <6 months            |

# Setup

When first setting up, you will need the Android app. You can download it from the [WaveShare Website](https://www.waveshare.net/w/upload/Cloud_app.apk), However if you cannot download it from this link you can download it from [The Wayback Machine](https://web.archive.org/web/20220719161209/https://www.waveshare.net/w/upload/Cloud_app.apk).

# Communication Protocol

Communication is divided into two modes: The Commnd mode and the Data mode. Command mode is used for sending commands (example: Shutdown). Data Mode is used for sending image data to the display.

## Checksum

The checksum used is a simple XOR of the data. For command mode, an implementation in Go would be:

```go
command := "C"

var check uint32
for i := 0; i < len(command); i++ {
    check = check ^ uint32(command[i])
}

display.Connection.Write([]byte(";" + command + "/" + string(rune(check))))
```

And for data mode, an implementation in Go would be:

```go
var check byte

payload := frame.Bytes()

// The first byte is the 0x57 identifier, so we skip it.
for i := 1; i < len(payload[1:])+1; i++ {
    // CheckSum8 Xor
    check ^= payload[i]
}

// The checksum byte is the last byte of the frame. It is stored as BigEndian.
// This function appends the checksum byte to the frame.
err = binary.Write(frame, binary.BigEndian, check)
if err != nil {
    return err
}

display.Connection.Write(frame.Bytes())
```

## Command Format

For the command mode, the command is sent in the following format:  
‘;’ + Command (with optional data) + '/' + Checksum

## Data Format

For the data mode, the data is sent in the following format:  
0x57 + 4 Byte addr + 4 Byte len + 1 Byte num + Data + Checksum

| Label | Value                           |
| ----- | ------------------------------- |
| Addr  | Address of data                 |
| Len   | Length of data                  |
| Num   | Frame number of current section |
| Data  | The data to be transmitted      |

### Recommendations:

- Frames should be the same length
- Frames should not be larger than 1100 Bytes
- Num should be static (at 0x00), due to a software update bug.
- Wait until the device has replied before sending the next frame.

### Information:

- If the addr and len are both 0x00, the display will assume transmission has completed, and the display will be refreshed.

## Response Format

The response is sent in the following format:
'$' + Data + '#'

This format is the same for command and data mode.

# Command List

| Command | Arguments               | Description                                          | 1st response | 2nd response | Requires Unlock                       |
| ------- | ----------------------- | ---------------------------------------------------- | ------------ | ------------ | ------------------------------------- |
| C       |                         | Checks if the device is locked                       | Parity Bit   | Locked?      | <ul><li>[ ] Requires Unlock</li></ul> |
| N       | Password                | Unlock the device                                    | Parity Bit   | Sucsessful?  | <ul><li>[ ] Requires Unlock</li></ul> |
| G       |                         | Gets the custom ID of the device                     | Parity Bit   | ID           | <ul><li>[ ] Requires Unlock</li></ul> |
| r       | Time in seconds (<9999) | Sleep - Only on 2.13 inch                            | Parity Bit   |              | <ul><li>[ ] Requires Unlock</li></ul> |
| 0       | New ID                  | Gives the device a new ID                            | Parity Bit   |              | <ul><li>[x] Requires Unlock</li></ul> |
| 1       | New IP                  | Gives the device a new IP adress on the network      | Parity Bit   |              | <ul><li>[x] Requires Unlock</li></ul> |
| 2       | New WIFI SSID           | Gives the device a new SSID to connect to with WIFI  | Parity Bit   |              | <ul><li>[x] Requires Unlock</li></ul> |
| 3       | New WIFI password       | Gives the device a new password for the WIFI network | Parity Bit   |              | <ul><li>[x] Requires Unlock</li></ul> |
| P       | New device password     | Sets a new password for the device                   | Parity Bit   |              | <ul><li>[x] Requires Unlock</li></ul> |
| L       | Boolean - Lock Device   | Controls the device's locked state                   | Parity Bit   |              | <ul><li>[x] Requires Unlock</li></ul> |
| s       | Boolean - Flag Bit      | I do not know what the flag bit does :P              | Parity Bit   |              | <ul><li>[x] Requires Unlock</li></ul> |
| B       |                         | Open for bluetooth connections                       | Parity Bit   |              | <ul><li>[x] Requires Unlock</li></ul> |
| b       |                         | Check battery voltage                                | Parity Bit   | Level (mv)   | <ul><li>[x] Requires Unlock</li></ul> |
| S       |                         | Shutdown the device                                  | Parity Bit   |              | <ul><li>[x] Requires Unlock</li></ul> |
| R       |                         | Restart the device                                   | Parity Bit   |              | <ul><li>[x] Requires Unlock</li></ul> |

# TODO: Format below text

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
