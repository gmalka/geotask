package router

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"time"
)

const (
	swaggerTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-standalone-preset.js"></script>
    <!-- <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-standalone-preset.js"></script> -->
    <script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
    <!-- <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-bundle.js"></script> -->
    <link rel="stylesheet" href="//unpkg.com/swagger-ui-dist@3/swagger-ui.css" />
    <!-- <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui.css" /> -->
	<style>
		body {
			margin: 0;
		}
	</style>
    <title>Swagger</title>
</head>
<body>
    <div id="swagger-ui"></div>
    <script>
        window.onload = function() {
          SwaggerUIBundle({
            url: "/swagger.json?{{.Time}}",
            dom_id: '#swagger-ui',
            presets: [
              SwaggerUIBundle.presets.apis,
              SwaggerUIStandalonePreset
            ],
            layout: "StandaloneLayout"
          })
        }
    </script>
</body>
</html>
`
)

func swaggerUI(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.New("swagger").Parse(swaggerTemplate)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = tmpl.Execute(c.Writer, struct {
		Time int64
	}{
		Time: time.Now().Unix(),
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
