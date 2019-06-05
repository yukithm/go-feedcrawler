go-feedcrawler
==============

Feed (RSS and Atom) crawler library (an example application included).


Features
--------

- Support RSS and Atom
- Filtering entries
    - Regexp based filter for title, description, content, author and categories
    - Callback function filter
- State management (keep published date and detect new entries)
- Multiple workers


Examples
--------

See `_examples` directory.

- TOML based configuration file
- Fake feed server (dynamic entries feed)


TODO
----

- Support local files (local path and/or file scheme such as "file://")


License
-------

MIT


Author
------

Yuki (@yukithm)
