package email

import (
	"bytes"
	"log"
	"text/template"
)

func GenerateEmailHTML(data map[string]string) string {
	const emailTemplate = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>{{.Subject}}</title>
		<style>
			body {
				font-family: 'Arial', sans-serif;
				background-color: #f9f9f9;
				margin: 0;
				padding: 0;
				text-align: center;
				color: #333;
			}
			.container {
				max-width: 600px;
				margin: 20px auto;
				background-color: #ffffff;
				border-radius: 8px;
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
				padding: 30px;
				font-size: 16px;
			}
			.header {
				background-color: #007bff;
				color: white;
				padding: 20px 0;
				border-radius: 8px 8px 0 0;
			}
			.header h1 {
				font-size: 28px;
				margin: 0;
			}
			.content {
				padding: 20px 0;
			}
			.cta-button {
				display: inline-block;
				padding: 12px 25px;
				background-color: #28a745;
				color: white;
				text-decoration: none;
				border-radius: 5px;
				font-weight: bold;
				margin-top: 20px;
				font-size: 16px;
			}
			.cta-button:hover {
				background-color: #218838;
			}
			.footer {
				background-color: #f1f1f1;
				color: #777;
				padding: 15px;
				font-size: 12px;
				border-radius: 0 0 8px 8px;
			}
			.footer p {
				margin: 0;
			}
			@media (max-width: 600px) {
				.container {
					width: 100% !important;
					padding: 15px;
				}
				.header h1 {
					font-size: 24px;
				}
				.cta-button {
					padding: 10px 20px;
					font-size: 14px;
				}
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>{{.Title}}</h1>
			</div>
			<div class="content">
				<p>{{.Message}}</p>
				<a href="{{.Link}}" class="cta-button">{{.Button}}</a>
			</div>
			<div class="footer">
				<p>This is an automated email. Please do not reply. For assistance, contact support@sanedge.com</p>
			</div>
		</div>
	</body>
	</html>`

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}

	return buf.String()
}
