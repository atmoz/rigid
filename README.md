# Rigid

Static web sites that makes sense.

This is in heavy development, prepare your soul!

## What?

* Simple
* Fast
* Zero configuration
* Templates
* Markdown

## How?

* Structure your files just like you want the site to be structured.
* html, md, and markdown file extensions are regarded as "pages".
* Add some meta data to your pages with [YAML](https://en.wikipedia.org/wiki/YAML), if you want.
    * title
    * tags
* Expect this to happen to your pages:
    * projects.md --> projects/index.html (pretty URL!)
    * boring.html.md --> boring.html (boring URL)
    * about.html --> about.html (html is boring by default)
* Pages inherits special [go templates](http://golang.org/pkg/text/template/) based on directory location:
    * `_current.template`: only current dir.
    * `_partial.template`: current dir and parent and child dirs.
    * `_final.template`: current dir and only child dirs.
* Template functions (experimental):
    * Sitemap(pattern): Print unordered list of files matching pattern.
    * TaggedPages(pattern): Return array of pages with tags matching pattern.
