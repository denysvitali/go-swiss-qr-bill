package swiss_qr_code

import (
	"fmt"
	"github.com/almerlucke/go-iban/iban"
	"strconv"
	"strings"
)

func parseAmount(line string) (*MoneyValue, error) {
	if line == "" {
		return nil, nil
	}
	sp := strings.Split(line, ".")
	if len(sp) != 2 {
		return nil, fmt.Errorf("invalid amount")
	}

	base, err := strconv.ParseInt(sp[0], 10, 32)
	if err != nil {
		return nil, err
	}

	cents, err := strconv.ParseInt(sp[1], 10, 32)
	if err != nil {
		return nil, err
	}

	return &MoneyValue{Base: uint32(base), Cents: uint32(cents)}, nil
}

func parseParty(d decoder, relaxed bool) (*Party, error) {
	party := Party{}

	parseFunc := parseStringOrFail
	if relaxed {
		parseFunc = parseOrSkip
	}

	party.AddressType = AddressType(parseFunc(d, parseAddressType))
	party.Name = parseFunc(d, parse70)
	party.StrtNmOrAdrLine1 = parseFunc(d, parseX(70))
	party.BldgNbOrAdrLine2 = parseFunc(d, parseX(70))

	switch party.AddressType {
	case Structured:
		if len(party.BldgNbOrAdrLine2) > 16 {
			return nil, fmt.Errorf("invalid BldgNbOrAdrLine2")
		}
	case Combined:
		if len(party.BldgNbOrAdrLine2) == 0 && !relaxed {
			return nil, fmt.Errorf("BldgNbOrAdrLine2 not provided")
		}
	}

	party.PostalCode = parseFunc(d, parseX(16))
	party.Town = parseFunc(d, parseX(35))
	if party.AddressType == Combined {
		if len(party.PostalCode) > 0 {
			return nil, fmt.Errorf("cannot provide postal code in Combined address mode")
		}
		if len(party.Town) > 0 {
			return nil, fmt.Errorf("cannot provide town in Combined address mode")
		}
	}
	party.CountryCode = parseFunc(d, parseX(2))
	return &party, nil
}

func parseX(size int) func(string) (string, error) {
	return func(line string) (string, error) {
		if len(line) > size {
			return "", fmt.Errorf("invalid length")
		}
		return line, nil
	}
}

func parse70(line string) (string, error) {
	return parseX(70)(line)
}

func parseIBAN(line string) (string, error) {
	_, err := iban.NewIBAN(line)
	if err != nil {
		return "", err
	}
	return line, nil
}

func parseOrFail[T any](d decoder, parseFunc func(line string) (v T, err error)) T {
	line := d.readLine()
	res, err := parseFunc(line)
	if err != nil {
		panic(err)
	}
	return res
}

func parseStringOrFail(d decoder, parseFunc func(line string) (v string, err error)) string {
	line := d.readLine()
	res, err := parseFunc(line)
	if err != nil {
		panic(err)
	}
	return res
}

func parseOrSkip(d decoder, parseFunc func(line string) (v string, err error)) string {
	line := d.readLine()
	res, err := parseFunc(line)
	if err != nil {
		return ""
	}
	return res
}

func parseAddressType(addressType string) (string, error) {
	switch addressType {
	case "S":
		return string(Structured), nil
	case "K":
		return string(Combined), nil
	}
	return "", fmt.Errorf("invalid address type \"%s\"", addressType)
}

func parseQrCodingType(codingType string) (CodingType, error) {
	num, err := strconv.ParseInt(codingType, 10, 16)
	if err != nil {
		return -1, err
	}

	return CodingType(num), nil
}

func parseQrVersion(version string) (v VersionMajMin, err error) {
	if len(version) != 4 {
		return v, fmt.Errorf("invalid version length")
	}

	maj := version[0:2]
	min := version[2:4]

	majInt, err := strconv.ParseInt(maj, 10, 16)
	if err != nil {
		return v, err
	}
	minInt, err := strconv.ParseInt(min, 10, 16)
	if err != nil {
		return v, err
	}

	return VersionMajMin{
		Major: int(majInt),
		Minor: int(minInt),
	}, nil
}

func parseQrType(qrType string) (QRType, error) {
	switch qrType {
	case string(SwissPaymentsCodeQrType):
		return SwissPaymentsCodeQrType, nil
	default:
		return "", fmt.Errorf("invalid QR type")
	}
}

func parseCurrency(line string) (Currency, error) {
	switch line {
	case "CHF":
		return CurrencyChf, nil
	case "EUR":
		return CurrencyEur, nil
	}

	return "", nil
}

func parsePaymentReference(line string) (PaymentReferenceType, error) {
	switch line {
	case "QRR":
		return ReferenceQRR, nil
	case "SCOR":
		return ScorReference, nil
	case "NON":
		return NonReference, nil
	}

	return "", fmt.Errorf("invalid reference \"%s\"", line)
}

func parseTrailer(line string) (string, error) {
	if line == "EPD" {
		// End of Payment Data
		return line, nil
	}
	return "", fmt.Errorf("invalid trailer \"%s\"", line)
}
