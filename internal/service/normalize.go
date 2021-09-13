package service

import "strings"

// normalizeSymbolsList normalizes symbols in a slice using normalizeSymbol.
func normalizeSymbolsList(symbols []string) {
	for i, symbol := range symbols {
		symbols[i] = normalizeSymbol(symbol)
	}
}

// normalizeSymbol converts the string with the symbol to the desired form.
func normalizeSymbol(symbol string) string {
	symbol = strings.TrimSpace(symbol)
	symbol = strings.ToUpper(symbol)
	return symbol
}
