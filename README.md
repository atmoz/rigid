# Rigid

Static web sites that just makes sense.

This is still in early development, PREPARE YOUR SOUL!

## What?

* Zero configuration
* Default and custom templates
* Markdown support

## How?

* Structure your files just like you want the web site to be structured.
* HTML and markdown files are regarded as web pages.
* Web page paths:
    * projects.md --> projects/index.html (pretty URL!)
    * boring.html.md --> boring.html (boring URL)
    * about.html --> about.html (html is boring by default)
* Web site is rendered with a simple menu, ready to use.
* *Optional:*
    * Add meta data to your pages.
    * Use custom CSS and templates.

## Example page with meta data

    ---
    title: My title
    tags: [ blog/ramblings, blog/example, whatever ]
    ---

    So this is my page, you like!?

## Using custom CSS and templates

If you don't like the default look, you can add your own CSS and/or templates.

### CSS

All you need is to edit `rigid.css` (created on first build) in the root folder.

### Templates

*More info coming later*

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

