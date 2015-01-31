package vcs

import(
	"testing"
	"encoding/xml"
)

func TestParseSvnXml(t *testing.T) {
	info := Info{}
	data := `
	<?xml version="1.0" encoding="UTF-8"?>
	<info>
		<entry revision="38" kind="dir" path=".">
			<url>https://github.com/MoriTanosuke/Hound/trunk</url>
			<relative-url>^/trunk</relative-url>
			<repository>
				<root>https://github.com/MoriTanosuke/Hound</root>
				<uuid>ca2d1e40-ed62-6735-98bd-ce57b7db7bff</uuid>
			</repository>
			<wc-info>
				<wcroot-abspath>/go/pub/data/hound-svn</wcroot-abspath>
				<schedule>normal</schedule>
				<depth>infinity</depth>
			</wc-info>
			<commit revision="27">
				<author>jonathan.klein</author>
				<date>2015-01-29T22:34:32.000000Z</date>
			</commit>
		</entry>
	</info>
	`

	err := xml.Unmarshal([]byte(data), &info)
	if err != nil {
		t.Error(err)
		return
	}
	
	if (info.Entry.Url != "https://github.com/MoriTanosuke/Hound/trunk") {
		t.Error("Revision was not read correctly, expected 'https://github.com/MoriTanosuke/Hound/trunk', was " + info.Entry.Url)
	}
	if (info.Entry.Revision != "38") {
		t.Error("Revision was not read correctly, expected '38', was " + info.Entry.Revision)
	}
	if(info.Entry.Repository.Uuid != "ca2d1e40-ed62-6735-98bd-ce57b7db7bff") {
		t.Error("Revision was not read correctly, expected 'ca2d1e40-ed62-6735-98bd-ce57b7db7bff', was " + info.Entry.Repository.Uuid)
	}
}
