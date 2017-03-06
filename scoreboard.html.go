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
			<thead>
				<tr>
					<td> Users </td>
					{{ range .Tasks }}
						<td> {{.}} </td>
					{{ end }}
				</tr>
			</thead>
			<tbody>
				{{ range $i1, $user :=  $.Users }}
					<tr>
						<td> {{ $user }} </td>
						{{ range $i2, $task := $.Tasks }}
							<td>
								{{ call $.Results $user $task }}
							</td>
						{{ end }}
					</tr>
				{{ end }}
			</tbody>
		</table>
	</body>
</html>

`))
