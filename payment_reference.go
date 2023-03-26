package swiss_qr_code

type PaymentReferenceType string

const (
	ReferenceQRR  PaymentReferenceType = "QRR"
	ScorReference PaymentReferenceType = "SCOR"
	NonReference  PaymentReferenceType = "NON"
)
