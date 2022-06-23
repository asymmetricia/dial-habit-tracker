web/css/bootstrap.min.css: sass.Dockerfile input.scss
	[ -f v5.2.0-beta1.zip ] || wget -O v5.2.0-beta1.zip https://github.com/twbs/bootstrap/archive/v5.2.0-beta1.zip
	docker build -t sass -f sass.Dockerfile .
	docker run \
		-v "$$PWD/web/css:/out" \
		-u $$(id -u) \
		-i sass \
		sass \
			--style compressed \
			input.scss \
			/out/bootstrap.min.css
