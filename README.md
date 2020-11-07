# Reddit Image Downloader
Basic image downloader that uses publicly open reddit endpoint to download images.  Currently it downloads the first request but it can be modified as download recursively to download as much as it can.

# Usage
```
$ go run main.go <subredditname>
```
It downloads the images into "downloads" folder and resets everytime that the command is executed.

# Libraries
It only uses standard libraries.