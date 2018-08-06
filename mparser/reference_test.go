package mparser

import (
	"testing"
)

func TestAnchorFromReference(t *testing.T) {
	ref := []byte(`<reference anchor='ts' target=''>
<front>
 <title>Old Possum's Book of Practical Cats</title>
   <author initials='TS' surname='Stearns' fullname='TS. Stearns'></author>
   <date/>
 </front>
</reference>`)

	got := string(anchorFromReference(ref))
	want := "ts"

	if got != want {
		t.Errorf("want %s, got %s, for input %s...", want, got, ref[:20])
	}
}

func TestReferenceHook(t *testing.T) {
	ref := []byte(`<reference anchor='ts' target=''>
<front>
 <title>Old Possum's Book of Practical Cats</title>
   <author initials='TS' surname='Stearns' fullname='TS. Stearns'></author>
   <date/>
 </front>
</reference>`)

	_, _, read := ReferenceHook(ref)
	if read != len(ref) {
		t.Errorf("want %d, got %d, for input %s...", len(ref), read, ref[:20])
	}
}
