CanvasStreamTest
================

After reading [Intro to (images in) Go](http://www.pheelicks.com/2013/10/intro-to-images-in-go/) by [@pheelicks](https://twitter.com/@pheeelicks) I started thinking about a better approach to stream the canvas to the browser.
~~Using [data URIs](http://en.wikipedia.org/wiki/Data_URI_scheme), this is quite simple. I hammered together the Nodes example with my approach as a showcase.~~
I don't use base64 data URIs anymore because it adds a lot of overhead. Instead I serve the image in a 2nd handler and just update the img.src attribute with a unused parameter so the browser dosn't cache it.