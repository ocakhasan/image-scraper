# Parallel Image Scraper

`image-scraper` scrapes images from given website to given folder in parallel.

## USAGE

First you need to clone the project.
```
> git clone https://github.com/ocakhasan/image-scraper.git
> cd image-scraper
```

Then you need to build the project.
```
> go build
```

Then you can use the program as below.
```
> ./image-scraper -w 'google.com' -f google
There are 2 images in google.com
- Image tia.png is downloaded
- Image googlelogo_white_background_color_272x92dp.png is downloaded
```

It will download images from `google.com` into `google` folder.
