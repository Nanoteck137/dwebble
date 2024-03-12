.PHONY: gen

gen:
	pyrin gen go api.pyrin -o types/api_gen.go -p types -f
