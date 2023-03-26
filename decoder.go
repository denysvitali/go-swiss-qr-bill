package swiss_qr_code

import "bufio"

type decoder struct {
	r   *bufio.Scanner
	err error
}

func (d *decoder) readLine() string {
	if d.r.Scan() {
		return d.r.Text()
	} else {
		return ""
	}
}
