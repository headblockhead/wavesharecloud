# wavesharecloud
A Go Package to control and display data on [Waveshare's 4.2 inch black and white e-paper cloud module](https://www.waveshare.com/4.2inch-e-paper-cloud-module.htm) easily.

## Development

I reverse-engineered the original protocol used to communicate to the display by using Wireshark to monitor TCP traffic to and from the device, then later found official documentation (although not very clear), and created a neat library to display images to the screen, and read data from the device, quickly and easily.

## Features

The library can:

- Display images, resizing and dithering when needed.
- Restart the display
- Get the display's battery level
- Shutdown/Sleep the display
- Get the display's ID
- Unlock displays that are locked, using their PIN/password.

## Examples

- [Blank template](examples/templateExample/)
- [Show an image](examples/showImage/)
- [Get the display's battery level](examples/getBattery/)

All examples support multiple displays connecting simultaneously.

## Device documentation

[The offical waveshare documentation](https://www.waveshare.com/wiki/4.2inch_e-Paper_Cloud_Module) does not cover everything about the device and is very poorly written.

I re-created the most important parts of the documentation in the **[DISPLAYDOCS.md](DISPLAYDOCS.md)** file, adding small code examples in Go, and re-writing the majority of it to hopefully be understood easier than the official version.
