# Parallel Image Scraper

`image-scraper` scrapes images from given website to given folder in parallel.

## USAGE

```
go get github.com/ocakhasan/image-scraper
```

It will generate a binary program in your `$GOBIN` folder.

Then you can use the program as below.
```
> image-scraper -w 'google.com' -f google
There are 2 images in google.com
- Image tia.png is downloaded
- Image googlelogo_white_background_color_272x92dp.png is downloaded
```

It will download images from `google.com` into `google` folder.
