# Image Scraper

`image-scraper` scrapes images from given website to given folder.

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
Folder google does not exist. Creating one...
Image googlelogo_white_background_color_272x92dp.png is done
Image tia.png is done
```

It will download images from `google.com` into `google` folder.

## TODOS
- [ ] make it parallel
