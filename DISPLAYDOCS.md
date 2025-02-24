# Waveshare 4.2inch e-Paper Cloud Module

## Features

- The device is suitable for use as equipment labels, shelf labels, information storage and more.
- It integrates an SW6106 chip that supports USB-C quick charging protocols like: PD / QC / FCP / PE / SFCP, etc.
- The device allows the user to control the display content via a remote server making it convenient and flexible.
- There is a voltage detection circuit can detect the battery level to make sure the device is not over or under charged.
- The device has support for a user-configured ID string that can be used to identify the device in a network.
- There is an Android app available that allows the user to configure the device over Bluetooth.
- Sadly, there is no backlight. However, the device's display will not change after the display has been written to (until the next write).
- Power consumption is tiny, the battery is only used when refreshing the display.
- It comes with an instruction manual.

## Specifications

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

## Setup

When first setting up, you will need the Android app, as Bluetooth is the only radio enabled at first start-up. You can download it from the [WaveShare Website](https://www.waveshare.net/w/upload/Cloud_app.apk), However if you cannot download it from this link, you can download it from [The Wayback Machine](https://web.archive.org/web/20220719161209/https://www.waveshare.net/w/upload/Cloud_app.apk).

## Communication Protocol

Communication is divided into two modes: The Commnd mode and the Data mode. Command mode is used for sending commands (example: Shutdown). Data Mode is used for sending image data to the display.

### Checksum

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

### Command Format

For the command mode, the command is sent in the following format:  
‘;’ + Command (with optional data) + '/' + Checksum

### Data Format

For the data mode, the data is sent in the following format:  
0x57 + 4 Byte addr + 4 Byte len + 1 Byte num + Data + Checksum

| Label | Value                           |
| ----- | ------------------------------- |
| Addr  | Address of data                 |
| Len   | Length of data                  |
| Num   | Frame number of current section |
| Data  | The data to be transmitted      |

#### Recommendations:

- Frames should be the same length
- Frames should not be larger than 1100 Bytes
- Num should be static (at 0x00), due to a software update bug.
- Wait until the device has replied before sending the next frame.

#### Information:

- If the addr and len are both 0x00, the display will assume transmission has completed, and the display will be refreshed.

### Response Format

The response is sent in the following format:
'$' + Data + '#'

This format is the same for command and data mode.

## Command List

| Command | Arguments               | Description                                          | 1st response | 2nd response | Requires Unlock                       |
| ------- | ----------------------- | ---------------------------------------------------- | ------------ | ------------ | ------------------------------------- |
| C       |                         | Checks if the device is locked                       | Parity Bit   | Locked?      | <ul><li>[ ] No</li></ul> |
| N       | Password                | Unlock the device                                    | Parity Bit   | Sucsessful?  | <ul><li>[ ] No</li></ul> |
| G       |                         | Gets the custom ID of the device                     | Parity Bit   | ID           | <ul><li>[ ] No</li></ul> |
| r       | Time in seconds (<9999) | Sleep - Only on 2.13 inch                            | Parity Bit   |              | <ul><li>[ ] No</li></ul> |
| 0       | New ID                  | Gives the device a new ID                            | Parity Bit   |              | <ul><li>[x] Yes</li></ul> |
| 1       | New IP                  | Gives the device a new IP adress on the network      | Parity Bit   |              | <ul><li>[x] Yes</li></ul> |
| 2       | New WIFI SSID           | Gives the device a new SSID to connect to with WIFI  | Parity Bit   |              | <ul><li>[x] Yes</li></ul> |
| 3       | New WIFI password       | Gives the device a new password for the WIFI network | Parity Bit   |              | <ul><li>[x] Yes</li></ul> |
| P       | New device password     | Sets a new password for the device                   | Parity Bit   |              | <ul><li>[x] Yes</li></ul> |
| L       | Boolean - Lock Device   | Controls the device's locked state                   | Parity Bit   |              | <ul><li>[x] Yes</li></ul> |
| s       | Boolean - Flag Bit      | I do not know what the flag bit does :P              | Parity Bit   |              | <ul><li>[x] Yes</li></ul> |
| B       |                         | Open for bluetooth connections                       | Parity Bit   |              | <ul><li>[x] Yes</li></ul> |
| b       |                         | Check battery voltage                                | Parity Bit   | Level (mv)   | <ul><li>[x] Yes</li></ul> |
| S       |                         | Shutdown the device                                  | Parity Bit   |              | <ul><li>[x] Yes</li></ul> |
| R       |                         | Restart the device                                   | Parity Bit   |              | <ul><li>[x] Yes</li></ul> |

## Device details

The firmware cannot be reprogrammed, as the flash on the ESP32 chip is locked.
To connect to the device, you need to wake the device up by pressing the button.

When the device starts it will overlay the display's content with the status icons. From left to right, they are:
- Waiting for commands (clock icon)
- Setup complete (gear icon)
- Low voltage (battery with lightning bolt icon)
- WiFi connected (wifi/radio symbol with curved lines)
- Connected to target host (two arrows icon)

It is advised you check the display's icons periodically. The low-voltage warnning appears below 3600mv. If the voltage goes 150mv lower than warnning voltage (3450mv), the device will shutdown automatically to protect its stability.

### Setup

#### If the Setup Complete icon is not shown

To configure the display for the first time (if the 'Setup complete' gear-shaped icon is not shown) the companion Android app must be used. **Press the button**, then connect using Bluetooth via the app to set the display's WiFi credentials. If the display doesn't connect to Bluetooth within 90 seconds of showing the status icons, it will shut down again.

#### If the Setup Complete icon is shown

To re-configure the display (if the 'Setup complete' gear-shaped icon is shown), the companion Android app can be used, but configuration can also be changed through the server the display connects to. If you want to use the app, **hold the button for at least 5 seconds** (or until the device appears on a Bluetooth scan), then connect and configure. The display will try to update like normal over WiFi while the button is held, and if it updates and then recieves a shutdown command from the server it connects to, it will give 30 seconds afterwards for a Bluetooth connection to be initiated. The device should be rebooted after configuration changes are made.
