package parse

import (
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
)

/*
Этот код заменяет кастомные теги принятыми в HTML и возвращает результат в виде дерева.
Содержимое остаётся на своих местах.

Кастомные теги				Чем их заменять
<heading>	 				<strong>
<source url="some url">		<a href="some url">
<link url="some url">		<a href="some url">
<pre>						<p><code>text here</code></p>
<code lang="go">			<pre><code class="go"></code></pre>
*/

func Parse(input string) (string, error) {
	parsed_html, err := ConvertToHTML(input)
	if err != nil {
		return "", err
	}

	markdown, err := ConvertToMarkdown(&parsed_html)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func ConvertToMarkdown(html *string) (string, error) {
	parsed_node, err := ConvertToHTML(*html)
	if err != nil {
		return "", err
	}

	markdown, err := md.ConvertString(parsed_node)
	if err != nil {
		return "", err

	}

	return markdown, nil
}

func ConvertToHTML(input string) (string, error) {
	// Хер знает, почему, но программа отказывается видеть содержимое тегов <source> и <link>
	// Именно поэтому я просто заменил все проблемные теги на <a>
	// Если вы знаете, как решить эту проблему, напишите пожалуйста t.me/Auserum
	input = strings.ReplaceAll(input, "<source", "<a")
	input = strings.ReplaceAll(input, "</source>", "</a>")
	input = strings.ReplaceAll(input, "<link", "<a")
	input = strings.ReplaceAll(input, "</link>", "</a>")
	input = strings.ReplaceAll(input, " url=", " href=")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(input))
	if err != nil {
		return "", err
	}

	doc.Find("heading").Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml("<strong>" + s.Text() + "</strong>")
	})

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if ok {
			s.ReplaceWithHtml("<a href=\"" + url + "\">" + s.Text() + "</a>")
		}
	})

	doc.Find("pre").Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml("<p><code>" + s.Text() + "</code></p>")
	})

	doc.Find("code").Each(func(i int, s *goquery.Selection) {
		lang, ok := s.Attr("lang")
		if ok {
			s.ReplaceWithHtml("<pre><code class=\"" + lang + "\">" + s.Text() + "</code></pre>")
		}
	})

	html, err := doc.Html()
	if err != nil {
		return "", err
	}

	return html, nil
}
