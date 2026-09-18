package inertia

var inertiaEntry = `
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1"/>
            %s
			{{ .inertiaHead }}
		</head>
		<body class="font-sans antialiased">
			{{ .inertia }}
		</body>
	</html>
`
