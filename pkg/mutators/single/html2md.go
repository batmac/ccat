package mutators

import (
	"io"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/batmac/ccat/pkg/log"
)

func init() {
	singleRegister("html2md", h2m, withDescription("html -> markdown"),
		withHintLexer("md"),
		withCategory("convert"),
		withAliases("h2m", "h2md"),
		withConfigBuilder(stdConfigStringWithDefault("github")))
}

func h2m(w io.WriteCloser, r io.ReadCloser, c any) (int64, error) {
	converter := md.NewConverter("", true, nil)
	if c.(string) == "github" {
		log.Debugf("Using GitHub flavored markdown")
		converter.Use(plugin.GitHubFlavored())
	}
	markdown, err := converter.ConvertReader(r)
	if err != nil {
		log.Fatal(err)
	}
	return markdown.WriteTo(w)
}
