# Cryptocurrency Order Book

Task description
----------------

<img align="right" width="30%" src="assets/binance-logo.svg">

It is necessary to obtain the BID and ASK order books using the open API of the BINANCE exchange. The number of orders in each order book is 15. For each order book, calculate the sum of the volumes of orders in it. For each order book, order data and total order quantities are displayed.
Then display the resulting information in a convenient format.

**Data collection method:**
* `Low level`: Using the REST protocol at a rate of 1 request / sec.
* `High level`: Continuous over the WebSocket protocol.

**Display method:**
* `Low level`: In cmd.
* `High level`: In a web browser.
