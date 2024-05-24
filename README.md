# obs-drops-overlay

A obs-browser-source overlay for playing sounds, music and Videos right in OBS. 
Controlled by a webserver the overlay can be controlled from a tablet or phone.

## Functionality

### Video player
* provides a transparent video player
* when a video ends the overlay becomes transparent
* remote play/pause and seek
* control playbackRate from 0.1 - 2.0
* loads local files and urls

### Audio player
* plays mp3 files
* volume control effects all playing sounds

### Music player
* has separate volume control

## Use

requires a directory tree 👇

```shell
$PWD
  └─assets
      ├─drops  ## audio and video
      │   ├─drop0.mp3
      │   ├─drop1.mp3
      │   └─drop0.mp4
      └─music  ## music
          └─music0.mp3
      
``` 

* audio and video drops go in drops
* I am sure you can figure where music goes

## adding to obs


### Supported Formats

| type  | format  |
|-------|---------|
| audio | mp3     |
| video | mkv mp4 |

## TODO

* control page CSS
