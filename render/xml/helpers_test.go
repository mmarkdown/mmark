package xml

import "testing"

func TestAttributesContains(t *testing.T) {
	attrs := []string{`style="symbols"`, `class="boo"`}

	if !AttributesContains("style", attrs) {
		t.Errorf("expected %s to be present in attrs", "style")
	}
	if AttributesContains("stle", attrs) {
		t.Errorf("expected %s to be not present in attrs", "stle")
	}
}
