# Disgo

Disgo is a simple commenting system for the web, written in Go. It is inspired by [Disqus](http://disqus.com) but does not come with all the bells and whistles. Some (anti-)features are:

- Ajax-based: no Iframes
- no Javascript dependency, especially not jQuery
- works on modern browsers (CORS is required)
- supports MySQL and PostgreSQL
- customizable CSS (in fact, there is currently no default client CSS)

Optional:

- admin approval of new comments
- e-mail notification of new comment
- IP address-based rate limiting
- markdown support

Here's what the admin interface looks like:
![Disg Admin Interface](http://pascalj.github.io/disgo.png)


## Installation

### Binary release

You can download a ready-to-go release for most common platforms (FreeBSD, Linux, OS X, Windows). Just unpack and [configure](#configuration) it.

### Manual build

You need to have the Go environment installed. To build and install disgo, simply run:

```
$ go get github.com/pascalj/disgo
$ go install github.com/pascalj/disgo
```

You can build a release with

```
$ make release
```

`build/` then contains all you need: binaries for the most common platforms, the public files, templates and a sample of the configuration file.

## Configuration

The various switches in the [configuration sample file](config.gcfg.sample) are documented. To get started, just copy the `disgo.cfg.sample` to `disgo.gcfg` and edit it to your needs. Two things are mandatory:
a database must be configured, either MySQL or PostgreSQL and you need to set at least one `origin`. An origin is a URL where you want to allow comments. It may contain `*` as a wildcard. You can tell disgo to load a specific config by providing the `-config` command line parameter. It is strongly suggested to change the `secret`.

The server will listen on `0.0.0.0:3000`. You can change that by setting the `HOST` and/or `PORT` environment variable.

Fire up the server process:

```
$ disgo
```

After that, the admin interface will be available at the host and port you configured (e.g. http://localhost:3000/). The first time you access the admin panel you're asked to create an admin user.

## Embedding comments

Once you got everything up and running, it's easy to embedd comments in your website. Just add the following script to your website, just before the closing `</body>` tag:

```html
<script>
var disgo = {
    base: 'http://example.com/'
}
document.write('<script src="' + disgo.base + '/js/disgo.js">\x3C/script>')
</script>
```

Replace `example.com` with the location of your disgo installation. To display a form and comments, add a `div` with a `data-disgo-url` attribute:

```html
<div data-disgo-url="http://example.com/2014/04/01/facebook-buys-golang"></div>
```

The `data-disgo-url` does not need to be the current URL, it does not even need to be a URL. However, it is used to identify the content that the comments belong to. So if you're hosting a blog and have comments on your posts, it's recommended to use the URL of the posts as `data-disgo-url`.

## Known issues

*Disgo will currently run only at the top-level of a (sub-)domain.* Even if you put it behind a reverse-proxy and serve it under a different directory, certain links will break.

*Only one admin is configurable.* There is no user management interface, yet.

## Contributing

Any feedback is welcome. If you have features/suggestions, please add them to the [Trello board](https://trello.com/b/HU7Vc3NT/disgo) or drop me a mail. If you find a bug, please open an issue on Github. Of course I'm happy to merge pull requests.