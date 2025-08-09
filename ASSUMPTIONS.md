# Web Crawler Tool Documentation

## Overview

This document explains how my web crawler tool works, what assumptions it makes, and the possible improvements.

## Link categorization 

To keep things simple and straightforward, I've made some decisions about how to classify different types of links on a webpage.

### Innaccesible links

These are links that don't actually take you anywhere useful:

- **Empty links** - Links with no destination (`href=""`)
- **Page anchors** - Links that jump to sections on the same page (`href="#section"`) 
- **Email links** - Links that open your email client (`href="mailto:someone@doruk.com"`)

### Internal links

These are links that stay within the same website:

- **Relative paths** - Simple paths like `/about` or `contact.html`
- **Same domain URLs** - Full URLs that point to the same website

### External links

These are links that take you to completely different websites - any URL starting with `http://` or `https://` that doesn't match the current site's domain.

### Domain matching

When comparing domains, I:
- Ignore uppercase vs lowercase differences
- Remove the `www.` part (so `www.doruk.com` matches `doruk.com`)
- Use a simple "contains" check rather than exact matching

### Current implementation

- Crawler only checks `<a>` tags that have an `href` attribute
- Other types of links like `<link>` tags or `<area>` tags are ignored
- Crawler assumes all relative links are internal to the site

## Dealing with Website Protection

Many websites try to block automated tools like this.

When I first built the tool, some websites (like home24.de) would give me "403 Forbidden" errors. This happens because they most probably detected I was a bot, not a real person browsing with a web browser.

### Solution

**First attempt - Going all out:**
I tried to make my tool look exactly like a real browser by adding lots of headers:

```http
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate
Connection: keep-alive
```

**The backfire:**
This actually made things worse. Some websites (like amazon.com) are smart enough to detect when someone is trying too hard to look like a browser, and they block those requests too.

**Finding the sweet spot:**
After some experimentation, I found that a simpler approach works better:

```http
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Connection: keep-alive
```

This "just enough" approach fools most basic bot detection without triggering the more sophisticated filters.

## Ideas for Future Improvements

### Make It Faster
Right now, I process websites one step at a time. I could use Go's goroutines (think of them as multiple workers) to analyze different parts of a website simultaneously, making the whole process much faster.

### Make It More Accurate
Currently, I just categorize links based on their URLs. I could actually try to visit each link to see if it really works, giving you more accurate information about broken links.

### Handle Modern Websites Better

**The challenge:**
Many modern websites use JavaScript to load content after the page first loads. When this tool fetches a webpage, it only gets the initial HTML structure, not the content that JavaScript adds later.

**What this means:**
If a website heavily relies on JavaScript (like many single-page applications), crawler might miss some links that only appear after the JavaScript runs.

**Potential solution:**
I could add detection to identify when a website uses a lot of JavaScript and warn you that the analysis might be incomplete.

### Stay Ahead of Bot Detection

While current header approach works for many websites, some sites will always be challenging to analyze. Websites that require human verification or use advanced bot detection will continue to be problematic for automated tools like this.
