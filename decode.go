package swiss_qr_code

import (
	"bufio"
	"bytes"
	"fmt"
)

// Documentation https://www.six-group.com/dam/download/banking-services/standardization/qr-bill/ig-qr-bill-v2.2-en.pdf

type QrCode struct {
	Header                Header
	CreditorInformation   CreditorInformation
	Creditor              Party
	UltimateCreditor      *Party
	PaymentAmount         PaymentAmount
	UltimateDebtor        *Party
	PaymentReference      PaymentReference
	AdditionalInformation AdditionalInformation
	AlternativeSchemes    AlternativeSchemes
}

type QRType string

type VersionMajMin struct {
	Major int
	Minor int
}

type Header struct {
	QRType     QRType
	Version    VersionMajMin
	CodingType CodingType
}

type Party struct {
	AddressType      AddressType
	Name             string
	StrtNmOrAdrLine1 string
	BldgNbOrAdrLine2 string
	PostalCode       string
	Town             string
	CountryCode      string
}

type PaymentAmount struct {
	Amount   *MoneyValue
	Currency Currency
}

type CreditorInformation struct {
	IBAN string
}

type PaymentReference struct {
	Type      PaymentReferenceType
	Reference string
}

type AdditionalInformation struct {
	Unstructured    string
	Trailer         string
	BillInformation string
}

type AlternativeSchemes struct {
	Params []string
}

func Decode(text string) (qrCode *QrCode, err error) {
	d := decoder{
		r: bufio.NewScanner(bytes.NewBufferString(text)),
	}

	qrCode = &QrCode{}
	header := Header{}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// Header
	header.QRType = parseOrFail(d, parseQrType)
	header.Version = parseOrFail(d, parseQrVersion)
	header.CodingType = parseOrFail(d, parseQrCodingType)
	qrCode.Header = header

	qrCode.CreditorInformation = CreditorInformation{}
	qrCode.CreditorInformation.IBAN = parseOrFail(d, parseIBAN)

	// Party
	creditor, err := parseParty(d, false)
	if err != nil {
		return nil, fmt.Errorf("unable to parse creditor")
	}

	qrCode.Creditor = *creditor

	// Ultimate Party (optional)
	ultimateCreditor, err := parseParty(d, true)
	if err != nil {
		return nil, fmt.Errorf("unable to parse ultimate creditor: %v", err)
	}
	qrCode.UltimateCreditor = ultimateCreditor

	paymentAmount := PaymentAmount{}
	paymentAmount.Amount = parseOrFail(d, parseAmount)
	paymentAmount.Currency = parseOrFail(d, parseCurrency)
	qrCode.PaymentAmount = paymentAmount

	// Ultimate Debtor
	ultimateDebtor, err := parseParty(d, true)
	if err != nil {
		return nil, fmt.Errorf("unable to parse ultimate debtor: %v", err)
	}
	qrCode.UltimateDebtor = ultimateDebtor

	// Payment Reference
	paymentReference := PaymentReference{}
	paymentReference.Type = parseOrFail(d, parsePaymentReference)

	switch paymentReference.Type {
	case ReferenceQRR:
		paymentReference.Reference = parseOrFail(d, parseX(27))
	case ScorReference:
		paymentReference.Reference = parseOrFail(d, parseX(25))
	case NonReference:
		paymentReference.Reference = parseOrFail(d, parseX(0))
	}
	qrCode.PaymentReference = paymentReference

	// Additional Information
	additionalInformation := AdditionalInformation{}
	additionalInformation.Unstructured = parseOrFail(d, parseX(140))
	additionalInformation.Trailer = parseOrFail(d, parseTrailer)
	additionalInformation.BillInformation = parseOrSkip(d, parseX(140))
	qrCode.AdditionalInformation = additionalInformation

	// Alternative Schemes (Optional)
	alternativeSchemes := AlternativeSchemes{}
	alternativeSchemes.Params = []string{
		parseOrSkip(d, parseX(100)),
		parseOrSkip(d, parseX(100)),
	}

	qrCode.AlternativeSchemes = alternativeSchemes

	if d.err != nil {
		return nil, d.err
	}
	return qrCode, nil
}
