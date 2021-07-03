package reference

import (
	"encoding/xml"
	"testing"
)

func TestReference(t *testing.T) {
	in := []byte(`<reference anchor='IANA' target='https://www.iana.org/assignments/media-types/media-types.xhtml'>
    <front>
        <title abbrev='IANA'>IANA Media Types</title>
        <author>
            <organization>IANA</organization>
        </author>
        <date month='February' year='2019'/>
    </front>
    <refcontent>blah</refcontent>
</reference>
`)
	expect := `<reference anchor="IANA" target="https://www.iana.org/assignments/media-types/media-types.xhtml"><front><title>IANA Media Types</title><author><organization>IANA</organization></author><date year="2019" month="February"></date></front><refcontent>blah</refcontent></reference>`

	var x Reference
	if err := xml.Unmarshal(in, &x); err != nil {
		t.Errorf("failed to unmarshal reference: %s: %s", in, err)
	}
	out, err := xml.Marshal(x)
	if err != nil {
		t.Errorf("failed to marshal reference: %s", err)
	}
	str := string(out)

	if str != expect {
		t.Errorf("expected\n%s\ngot\n%s", expect, str)
	}
}

func TestReferenceDate(t *testing.T) {
	in := []byte(`<reference anchor='IANA' target='https://www.iana.org/assignments/media-types/media-types.xhtml'>
    <front>
        <title abbrev='IANA'>IANA Media Types</title>
        <author>
            <organization>IANA</organization>
        </author>
    </front>
</reference>
`)
	expect := `<reference anchor="IANA" target="https://www.iana.org/assignments/media-types/media-types.xhtml"><front><title>IANA Media Types</title><author><organization>IANA</organization></author></front></reference>`

	var x Reference
	if err := xml.Unmarshal(in, &x); err != nil {
		t.Errorf("failed to unmarshal reference: %s: %s", in, err)
	}
	out, err := xml.Marshal(x)
	if err != nil {
		t.Errorf("failed to marshal reference: %s", err)
	}
	str := string(out)

	if str != expect {
		t.Errorf("expected\n%s\ngot\n%s", expect, str)
	}
}
