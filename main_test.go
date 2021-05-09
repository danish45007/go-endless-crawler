package main

import "testing"

func TestToFixedUrl(t *testing.T) {
	fixedUrl := validUrl("/aboutus.html", "http://test.com/")
	if fixedUrl != "http://test.com/aboutus.html" {
		t.Error("toFixedUrl did not get expected href")
	}

	mailToUrl := validUrl("mailto:test@gmail.com", "http://test.com/")
	if mailToUrl != "http://test.com/" {
		t.Error("expected baseUrl instead of mailto link")
	}

	telephoneUrl := validUrl("tel://9820098200", "http://test.com/")
	if telephoneUrl != "http://test.com/" {
		t.Error("expected baseUrl instead of telephone link")
	}
}