<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
	<channel>
		<title>{{ with .Title }}{{.}} on {{ end }}{{ .Site.Title }}</title>
		<link>{{ .Permalink }}</link>
		<description>Recent content {{ with .Title }}in {{.}} {{ end }}on {{ .Site.Title }}</description>
		<generator>Hugo -- gohugo.io</generator>{{ with .Site.LanguageCode }}
		<language>{{.}}</language>{{end}}{{ with .Site.Author.email }}
		<managingEditor>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</managingEditor>{{end}}{{ with .Site.Author.email }}
		<webMaster>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</webMaster>{{end}}{{ with .Site.Copyright }}
		<copyright>{{.}}</copyright>{{end}}{{ if not .Date.IsZero }}
		<lastBuildDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 -0700" | safeHTML }}</lastBuildDate>{{ end }}
		<atom:link href="{{.URL}}" rel="self" type="application/rss+xml"/>
		{{ range first 20 .Data.Pages }}
		<item>
			<title>{{ .Title }}</title>
			<link>{{ .Permalink }}</link>
			<pubDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 -0700" | safeHTML }}</pubDate>
			{{ with .Site.Author.email }}<author>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</author>{{end}}
			<guid>{{ .Permalink }}</guid>
			<description>
				<p>{{ .Summary }}...</p>
				<p><a href="{{ .Permalink }}">Read More</a></p>
			</description>
		</item>
		{{ end }}
	</channel>
</rss>
