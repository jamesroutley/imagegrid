# imagegrid

`imagegrid` is a CLI tool which stitches multiple images into a single image, with a margin around each.

## Example

Given two images, `cat-1.jpg`:

![Cat 1](example-images/cat-1.jpg)

And `cat-2.jpg`:

![Cat 2](example-images/cat-2.jpg)

Comibe them with into a new image:

```sh
$ imagegrid cat-1.jpg cat-2.jpg
```

![Combined](example-images/imagegrid-image-2019-09-22.png)


## API

```sh
$ imagegrid -h
Usage of imagegrid:
  -margin float
        margin size (percentage) (default 5)
  -output-filename string
        name of the file to save the image to (default "imagegrid-image-<date>.png")
```

## Why?

At [work](https://monzo.com/), I often need to post a series of screenshots of the Monzo app to Slack - stitching them together into a single image seems to be the easiest way to do this.
