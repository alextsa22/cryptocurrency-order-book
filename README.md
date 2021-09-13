# Cryptocurrency Order Book

Task description
----------------

<img align="right" width="30%" src="assets/binance-logo.svg">

It is necessary to obtain the BID and ASK order books using the open API of the BINANCE exchange. The number of orders in each order book is 15. For each order book, calculate the sum of the volumes of orders in it. For each order book, order data and total order quantities are displayed.
Then display the resulting information in a convenient format.

**Data collection method:**
* ✅`Low level`: Using the REST protocol at a rate of 1 request / sec.
* 🔲`High level`: Continuous over the WebSocket protocol.

**Display method:**
* ✅`Low level`: In cmd.
* 🔲`High level`: In a web browser.

Running
-------

To run the application, run the commands:

```
make build & make run
```

Or you can use docker commands directly:

```
docker build -t fetcher .
```

```
docker run --rm -t -i fetcher
```

### Note

You can stop the program with ```Ctrl + C```.

Example
-------

<p align="center">
    <img src="assets/screenshot-v1.0.png">
</p>
