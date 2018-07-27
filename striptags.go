package htmltags

import (
	"strings"
	"bytes"
	"golang.org/x/net/html"
	)

type Nodes struct{
	Elements *html.Node
}

//Strip tags from provided html string
func Strip(content string, allowedTags []string, stripInlineAttributes bool) (Nodes, error){
	document, err := toNodes(content)
	handleError(err)
	var nodeTree html.Node

	var output func(document *html.Node, nt *html.Node)
	output = func(document *html.Node, nt *html.Node){
		for c := document.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode || ( c.Type == html.ElementNode && inArray(c.Data, allowedTags)) {
				var childNode html.Node
				childNode.Type = c.Type
				childNode.Data = c.Data
				if stripInlineAttributes == true {
					childNode.Attr = []html.Attribute{}
				}else{
					childNode.Attr = c.Attr
				}
				nt.AppendChild(&childNode)
				output(c, nt.LastChild)
			}else{
				output(c, nt)
			}
		}
	}
	output(document, &nodeTree)
	return Nodes{Elements:&nodeTree}, nil
}

//String to nodes helper.
func toNodes(document string) (*html.Node, error){
	nodes, err := html.Parse(strings.NewReader(html.UnescapeString(document)))
	handleError(err)
	return nodes, nil
}

//Nodes method. Converts nodes to string
func (nodes *Nodes) ToString() string{
	var buf bytes.Buffer
	for n := nodes.Elements.FirstChild; n != nil; n = n.NextSibling {
		html.Render(&buf, n)
	}
	return html.UnescapeString(buf.String())
}

//Check if needle is in the array
func inArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//Show error
func handleError(err error){
	if err != nil {
		panic(err)
	}
}