package godge

import "html/template"

var scoreboardTmpl = template.Must(template.New("scoreboard").Parse(`
<html>
	<head>
		<style>
			table {
					border-collapse: collapse;
					width: 100%;
			}
			table, th, td {
					border: 1px solid black;
					text-align: center;
			}
		</style>

	</head>

	<body>
		<h1>Scoreboard!</h1>
		<table>
			<tbody>
				{{ range $i1, $row1 :=  $.Scoreboard }}
					<tr>
						{{ range $i2, $row2 := $row1 }}
							<td>
								{{ index $.Scoreboard $i1 $i2 }}
							</td>
						{{ end }}
					</tr>
				{{ end }}
			</tbody>
		</table>
	</body>
</html>

`))
