# Rigid

Static web sites that makes sense.

This is in heavy development, PREPARE YOUR SOUL!

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
* Pages inherits these [go templates](http://golang.org/pkg/text/template/) based on directory location:
    * `_current.template`: only in current dir.
    * `_partial.template`: current dir, traverse parent and child dirs.
    * `_final.template`: current dir and only traverse child dirs.
* Template functions (experimental):
    * Sitemap(pattern): Print unordered list of files matching pattern.
    * TaggedPages(pattern): Return array of pages with tags matching pattern.

## Example file structure

    posts/
        _partial.template
        programming/
            _current.template
            some-post.md
        music/
            playlists/
                jazz.md
            _current.template
            another-post.md
    _final.template
    style.css
    index.html

Templates are applied to pages is this order:

* posts/programming/some-post.md
    * posts/programming/\_current.template
    * posts/\_partial.template
    * \_final.template
* posts/music/playlist/jazz.md
    * posts/\_partial.template
    * \_final.template
* posts/music/another-post.md
    * posts/music/\_current.template
    * posts/\_partial.template
    * \_final.template
* index.html
    * \_final.template

## Example page

    ---
    title: My FANCY title, indeed
    tags: [ blog/fancy, blog/example, whatever ]
    ---

    So this is my page, you like!?

## Example template

    <h1>{{.Page.Meta.Title}}</h1>
    <article>
        {{.Content}}

        <p>This page is tagged: {{range .Page.Meta.Tags}}<span>{{.}}</span> {{end}}</p>
    </article>

    <h2>All blog posts:</h2>
    <ul>
    {{range .TaggedPages "blog/*"}}
        <li><a href="{{.PublicPath}}">{{.Meta.Title}}</a></li>
    {{end}}
    </ul>

