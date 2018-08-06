package xml

import (
	"bytes"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/mmarkdown/mmark/mast"
)

const tomldata = `
Title = "A Standard for the Transmission of IP Datagrams on Avian Carriers"
abbrev = "IP Datagrams on Avian Carriers"
updates = [1034, 1035]
obsoletes = [4094]
category = "info"
docName = "rfc-1149"
ipr= "trust200902"
area = "Internet"
workgroup = "Network Working Group"
keyword = ["a", "b", "c"]

date = 1990-04-01T00:00:00Z

[[author]]
initials="D."
surname="Waitzman"
fullname="David Waitzman"
organization = "BBN STC"

	[author.address]
	email = "dwaitzman@BBN.COM"
	phone = "(617) 873-4323"
	[author.address.postal]
	street = "10 Moulton Street"
	city = "Cambridge"
	code = "MA 02238"

`

func TestTitle(t *testing.T) {
	// TODO: fix test
	node := mast.NewTitle()

	if _, err := toml.Decode(tomldata, node.TitleData); err != nil {
		t.Fatalf("Failure to parsing title block: %s", err.Error())
	}
	r := NewRenderer(RendererOptions{})
	buf := &bytes.Buffer{}
	r.titleBlock(buf, node)
	println(buf.String())
}
