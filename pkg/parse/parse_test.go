package parse

import (
	"testing"
)

func TestConvertMarkdown(t *testing.T) {
	input := `
	<heading>Создание API на Java</heading>

API (Application Programming Interface) на Java может быть создан с использованием различных технологий и фреймворков, таких как Spring Boot, Jakarta EE (ранее Java EE) или фреймворков полегче, вроде Javalin. Выбор конкретного инструментария зависит от требований проекта, его масштаба и предпочтений разработчика. Spring Boot, например, предоставляет упрощенный способ создания веб-приложений и RESTful API, абстрагируя сложности конфигурации и развертывания. <a href="https://spring.io/projects/spring-boot">[1]</a> Jakarta EE, в свою очередь, предлагает более сложные возможности и подходит для enterprise-приложений.

Для создания API на Java обычно требуется определить структуру данных, разработать логику обработки запросов и реализовать методы для обработки HTTP-запросов (GET, POST, PUT, DELETE и т.д.). Это включает в себя обработку входных данных, взаимодействие с базой данных (если необходимо), генерацию ответов в формате JSON или XML и обработку ошибок.<source url="https://jakarta.ee/">[2]</source> При использовании Spring Boot это может быть сделано с помощью аннотаций, таких как <pre>@RestController</pre> и <pre>@RequestMapping</pre>, которые упрощают процесс создания API endpoints.

Важно учитывать аспекты безопасности при разработке API, включая аутентификацию, авторизацию и защиту от распространенных уязвимостей, таких как SQL-инъекции и Cross-Site Scripting (XSS). Также необходимо предусмотреть обработку ошибок и логирование для мониторинга работы API и оперативного решения проблем. Документирование API с помощью таких инструментов, как Swagger/OpenAPI, поможет разработчикам легче понимать и использовать ваш API. <source url="https://swagger.io/">[3]</source>
	`

	parsed_html, err := ConvertToHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	markdown, err := ConvertToMarkdown(&parsed_html)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(markdown)
}
